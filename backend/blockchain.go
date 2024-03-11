// blockchain.go

package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

    "github.com/gorilla/websocket"

	"context"
    "log"
    "math/big"
    "encoding/json"

    "time"
    "sync"
    "fmt"

    "github.com/mriusd/game-contracts/db"
    "github.com/mriusd/game-contracts/maps"
    "github.com/mriusd/game-contracts/items"
    "github.com/mriusd/game-contracts/inventory"
    "github.com/mriusd/game-contracts/blockchain"
    "github.com/mriusd/game-contracts/fighters"
)

var NONCE_MUTEX = sync.RWMutex{}

func getNonce() (uint64, error) {
    nonce, err := getRpcClient().PendingNonceAt(context.Background(), ADMIN_ADDRESS)
    if err != nil {
        log.Printf("[getNonce] Failed to get nonce: %v", err)
        return 0, fmt.Errorf("[getNonce] Failed to get nonce, err: %v", err)
    }
    return nonce, nil
}


func sendBlockchainTransaction(
    fighter *fighters.Fighter, 
    contractName string, 
    contractAdr string, 
    data []byte, 
    eventContractName string,
    eventContractAddress string,
    eventTolisten string, 
    coords maps.Coordinate, 
    someHash common.Hash,
    conn *websocket.Conn,
) {
    log.Printf("[sendBlockchainTransaction] eventTolisten=%v", eventTolisten)
    // Connect to the Ethereum network
    client      := getRpcClient();


    // Set contract address
    contractAddress := common.HexToAddress(contractAdr)

    

    gasLimit := uint64(6000000)
    gasPrice := big.NewInt(GAS_PRICE)

    // auth := bind.NewKeyedTransactor(PRIVATE_KEY)
    // auth.Nonce = big.NewInt(int64(nonce))
    // auth.Value = big.NewInt(0)
    // auth.GasLimit = gasLimit 
    // auth.GasPrice = gasPrice

    NONCE_MUTEX.Lock()
    nonce, err := getNonce()
    if err != nil {
        log.Printf("[sendBlockchainTransaction] Failed to get nonce: %v", err)
        return
    }

    log.Printf("[sendBlockchainTransaction] nonce=%v gasPrice=%v", nonce, gasPrice)

    // Create transaction and sign it
    tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)
    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(RPCNetworkID), PRIVATE_KEY)
    if err != nil {
        log.Printf("[sendBlockchainTransaction] Failed to sign transaction: %v", err)
        NONCE_MUTEX.Unlock()
        return
    }

    // Send the transaction
    err = client.SendTransaction(context.Background(), signedTx)
    NONCE_MUTEX.Unlock()
    if err != nil {
        log.Printf("[sendBlockchainTransaction] Failed to send transaction: %v gasPrice=%v", err, gasPrice)
    } else {
        log.Printf("[sendBlockchainTransaction] Transaction hash: %v", signedTx.Hash().Hex())

        receiptChan := make(chan *types.Receipt)
        errChan := make(chan error)

        go waitForReceiptP(client, signedTx.Hash(), common.HexToAddress(eventContractAddress), receiptChan, errChan)

        select {
        case receipt := <-receiptChan:
            // Process the receipt
            //log.Printf("[PickupDroppedItem] Logs %+v:", receipt.Logs[0])
            handleBlockchainEvent(eventTolisten, eventContractName, receipt, fighter, coords, someHash, conn)
        case err := <-errChan:
            log.Printf("[sendBlockchainTransaction] Failed to get transaction receipt: %v", err)
        }
    }

    log.Printf("[sendBlockchainTransaction] contractName=%v eventTolisten=%v TX: %v", contractName, eventTolisten, signedTx.Hash().Hex());
}

func handleBlockchainEvent(eventName, contractName string, receipt *types.Receipt, fighter *fighters.Fighter, coords maps.Coordinate, someHash common.Hash, conn *websocket.Conn) {
    parsedABI := blockchain.LoadABI(contractName)
    switch eventName {
        case "Faucet":
            getUserCredits(conn)
            break

        case "FighterCreated":
            // event := FighterCreatedEvent{}

            // err := parsedABI.UnpackIntoInterface(&event, eventName, receipt.Logs[0].Data)
            // if err != nil {
            //     log.Printf("[handleBlockchainEvent:FighterCreated] Failed to unpack log err=%v data=%v", err, receipt.Logs[0].Data)
            //     return
            // }

            // log.Printf("[handleBlockchainEvent:FighterCreated] event=%v", event)

            getUserFighters(conn)
            break

        case "BackpackItemDropped":
            event := items.ItemDroppedEventSolidity{}

            err := parsedABI.UnpackIntoInterface(&event, eventName, receipt.Logs[0].Data)
            if err != nil {
                log.Printf("[handleBlockchainEvent:InventoryItemDropped] Failed to unpack log data: %v", err)
                return
            }

            log.Printf("[handleBlockchainEvent:InventoryItemDropped] ItemHash: %v", event.ItemHash)

            event.BlockNumber = receipt.BlockNumber
            event.Coords = coords
            event.OwnerId = big.NewInt(fighter.GetTokenID())

            // Add a self-destruct timer to remove the item from the map after 30 seconds
            time.AfterFunc(30*time.Second, func() {
                // DroppedItemsMutex.Lock() // Use a mutex if needed to protect access to the map
                // delete(DroppedItems, event.ItemHash)
                // DroppedItemsMutex.Unlock()
                items.DroppedItems.Remove(event.ItemHash)
                log.Printf("[handleBlockchainEvent:InventoryItemDropped] Item with hash %v has been removed after 30 seconds", event.ItemHash)
                broadcastDropMessage()
            })

            // DroppedItemsMutex.Lock()
            // DroppedItems[event.ItemHash] = &event
            // DroppedItemsMutex.Unlock()
            items.DroppedItems.Add(event.ItemHash, &event)


            if fighter.GetBackpack().FindByHash(someHash) != nil {
                fighter.GetBackpack().RemoveItemByHash(someHash)
                WsSendBackpack(fighter)
            } else {
                //removeItemFromEquipmentSlotByHash(fighter, someHash)
                fighter.GetEquipment().RemoveByHash(someHash)
            }

            db.RemoveItemFromDB(event.TokenId.Int64())

            broadcastDropMessage()

            break

        case "ItemDropped":
            event := items.ItemDroppedEventSolidity{}
            log.Printf("[handleBlockchainEvent:ItemDropped] receipt.Logs[0].Data: %v", receipt.Logs[0].Data)
            err := parsedABI.UnpackIntoInterface(&event, eventName, receipt.Logs[0].Data)
            if err != nil {
                log.Printf("[handleBlockchainEvent:ItemDropped] Failed to unpack log data: %v", err)
                return
            }

            log.Printf("[handleBlockchainEvent:ItemDropped] ItemHash: %v", event.ItemHash)

            event.BlockNumber = receipt.BlockNumber
            event.Coords = maps.Coordinate{X: event.CX.Int64(), Y: event.CY.Int64()}
            event.OwnerId = big.NewInt(fighter.TokenID)

            // Add a self-destruct timer to remove the item from the map after 30 seconds
            time.AfterFunc(30*time.Second, func() {
                // DroppedItemsMutex.Lock() // Use a mutex if needed to protect access to the map
                // delete(DroppedItems, event.ItemHash)
                // DroppedItemsMutex.Unlock()
                items.DroppedItems.Remove(event.ItemHash)
                log.Printf("[handleBlockchainEvent:ItemDropped] Item with hash %v has been removed after 30 seconds", event.ItemHash)
                broadcastDropMessage()
            })

            // DroppedItemsMutex.Lock()
            // DroppedItems[event.ItemHash] = &event
            // DroppedItemsMutex.Unlock()
            items.DroppedItems.Add(event.ItemHash, &event)

            log.Printf("[handleBlockchainEvent:ItemDropped] DroppedItems: %v", items.DroppedItems)

            if fighter.GetBackpack().FindByHash(someHash) != nil {
                fighter.GetBackpack().RemoveItemByHash(someHash)
                WsSendBackpack(fighter)
            } else {
                //removeItemFromEquipmentSlotByHash(fighter, someHash)
                fighter.GetEquipment().RemoveByHash(someHash)
            }

            db.RemoveItemFromDB(event.TokenId.Int64())

            broadcastDropMessage()
            break

        case "ItemPicked":
            event := items.ItemPickedEvent{}

            log.Printf("[handleBlockchainEvent:ItemPicked] logEntry: %v", receipt.Logs[0])

            err := parsedABI.UnpackIntoInterface(&event, eventName, receipt.Logs[0].Data)
            if err != nil {
                log.Printf("[handleBlockchainEvent:ItemPicked] Failed to unpack log data: %v", err)
                return
            }


            // DroppedItemsMutex.Lock()
            // dropEvent := DroppedItems[someHash]
            // DroppedItemsMutex.Unlock()
            dropEvent := items.DroppedItems.Find(someHash)
            

            log.Printf("[handleBlockchainEvent:ItemPicked] someHash: %v", someHash)
            log.Printf("[handleBlockchainEvent:ItemPicked] DroppedItems: %v", items.DroppedItems)
            log.Printf("[handleBlockchainEvent:ItemPicked] dropEvent: %v", dropEvent)

            item := dropEvent.GetItem()
            item.TokenId = event.TokenId
            tokenAtts := items.ConvertSolidityItemToGoItem(item);
            items.SaveItemAttributesToDB(tokenAtts)

            // DroppedItemsMutex.Lock()
            // delete(DroppedItems, someHash)
            // DroppedItemsMutex.Unlock()
            //DroppedItems.Remove(someHash)

            log.Printf("[ItemPicked] tokenAtts=%v", tokenAtts)

            if item.Name != "Gold" {
                _, _, err := fighter.GetBackpack().AddItem(tokenAtts, dropEvent.Qty.Int64(), someHash)
                saveBackpackToDB(fighter)
                if err != nil {
                    log.Printf("[handleBlockchainEvent:ItemPicked] Inventory full: %v", someHash)
                    sendErrorMessage(fighter, "Inventory full")
                }
            }

            log.Printf("[handleBlockchainEvent:ItemPicked] event: %+v\n", event)  
            WsSendBackpack(fighter) 
            broadcastPickupMessage(fighter, tokenAtts, event.Qty)
            sendChatMessageToFighter(fighter, "SYSTEM", "Picked "+item.Name, "system")
            break 


        default:
            log.Printf("[handleBlockchainEvent] Uknown eventName=%v receipt=%v", eventName, receipt);
            break
    }
}


// This function is vor development only
// Will not work in production
func MakeItem(fighter *fighters.Fighter, item *items.SolidityItemAtts) {
    log.Printf("[MakeItem] ItemAttributes=%v", item);


    // Load contract ABI from file
    contractABI := blockchain.LoadABI("DropHelper");

    data, err := contractABI.Pack("makeItem", item, fighter.GetLocation(), big.NewInt(fighter.GetCoordinates().X), big.NewInt(fighter.GetCoordinates().Y))
    if err != nil {
        log.Printf("[MakeItem] Failed to encode function arguments: %v", err)
    }

    go sendBlockchainTransaction(
        fighter, 
        "DropHelper", 
        DropHelperContract, 
        data, 
        "Drop",
        DropContract,
        "ItemDropped", 
        fighter.GetCoordinates(), 
        common.Hash{},
        nil,
    );


}


func getUserFighters(conn *websocket.Conn)  {
    log.Printf("[getUserFighters]")

    connection := GetConnection(conn)

    if (connection.gOwnerAddress() == common.Address{}) {
        log.Fatalf("[getUserFighters] Zero address")
        return
    }

    // Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(FightersContract)
    contractABI := blockchain.LoadABI("FightersHelper")

    // Prepare the call to the getTokenAttributes function
    callData, err := contractABI.Pack("getUserFighters", connection.gOwnerAddress())
    if err != nil {
        log.Fatalf("[getUserFighters] Failed to pack call data: %v", err)
    }


    //log.Printf("[getUserFighters] FreqData: %v", reqData)

    // Call the contract using the Ethereum client
    result, err := client.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
    }, nil)
    if err != nil {
        log.Fatalf("[getUserFighters] Failed to call contract: %v", err)
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var response []interface{};

    //log.Printf("[getUserFighters] Result", result)
    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    response, err = contractABI.Unpack("getUserFighters", result)
    if err != nil {
        log.Printf("[getUserFighters] Failed to unpack error: %v", err)
    }

    var fighterList []*fighters.FighterAttributes

    log.Printf("[getUserFighters] response: %v", response)

    for _, v := range response {
        attrs, ok := v.([]*big.Int)
        if !ok {
            // handle invalid attribute format
            log.Printf("[getUserFighters] Error iterating attributes")
            continue
        }

        for i := 0; i < len(attrs); i += 1 {

            itemId := attrs[i]

            // get item attributes
            fighterAttributes, _ := getFighterAttributes(itemId.Int64());

            //recordItemToDB(itemAttributes);

            fighterList = append(fighterList, fighterAttributes);
        }
    }


    log.Printf("[getUserFighters] fighters: %v", fighterList)


    type jsonResponse struct {
        Action string `json:"action"`
        Fighters []*fighters.FighterAttributes `json:"fighters"`        
    }

    jsonResp := jsonResponse{
        Action: "user_fighters",
        Fighters: fighterList,
    }


    jr, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[getUserFighters] error: ", err)
        return
    }
    respondConn(conn, jr)
}


func DropBackpackItem(conn *websocket.Conn, itemHash common.Hash, coords maps.Coordinate) {
    log.Printf("[DropBackpackItem] ItemHash=%v", itemHash);

    fighter, err := findFighterByConn(conn)

    if err != nil {
        return
    }

    bacpackSlot := fighter.GetBackpack().FindByHash(itemHash)

    if bacpackSlot == nil {
        log.Printf("[DropBackpackItem] Item not found in Inventory: %v", itemHash)

        return
    }    

    item := bacpackSlot.Attributes   

    log.Printf("[DropBackpackItem] item=%v", item);

    if item.ItemAttributes.IsBox {
        DropBox(conn, itemHash, coords, fighter, item)
        return
    }

    //tokenId := big.NewInt(0).SetInt64(item.TokenId)
    qty := big.NewInt(bacpackSlot.Qty)

    log.Printf("[DropBackpackItem] TokenId=%v Qty=%v", item.TokenId, qty);
    

    // Load contract ABI from file
    contractABI := blockchain.LoadABI("BackpackHelper");

    data, err := contractABI.Pack("dropBackpackItem", item.TokenId, qty)
    if err != nil {
        log.Printf("[DropBackpackItem] Failed to encode function arguments: %v", err)
    }

    go sendBlockchainTransaction(
        fighter, 
        "BackpackHelper", 
        BackpackHelperContract, 
        data, 
        "Backpack",
        BackpackContract,
        "BackpackItemDropped", 
        coords, 
        itemHash,
        nil,
    )
}

// func BurnConsumable(fighter *Fighter, item *TokenAttributes) {
//     log.Printf("[BurnConsumable] item=%v", item.gTokenId());
    
//     // Load contract ABI from file
//     contractABI := loadABI("ItemsHelper");

//     data, err := contractABI.Pack("burnConsumable", item.gTokenId())
//     if err != nil {
//         log.Printf("[BurnConsumable] Failed to encode function arguments: %v", err)
//     }

//     go sendBlockchainTransaction(
//         fighter, 
//         "ItemsHelper", 
//         ItemsHelperContract, 
//         data, 
//         "Items",
//         ItemsContract,
//         "ConsumableItemBurnt", 
//         maps.Coordinate{X: 0, Y: 0}, 
//         common.Hash{},
//         nil,
//     )
// }

// func CreateFighter(conn *websocket.Conn, ownerAddress, name string, class string) {
//     log.Printf("[CreateFighter] ownerAddress=%v, class=%v", ownerAddress, class);


//     err := validateFighterName(name)
//     if err != nil {
//         log.Printf("[CreateFighter] Invalid fighter name=%v error=%v", name, err)
//         sendErrorMsgToConn(conn, "SYSTEM", "Invalid character name. Onlye letters a to Z and numbers 0 to 9 allowed. Max length 13 characters.")
//         return
//     }


//     // Load contract ABI from file
//     contractABI := blockchain.LoadABI("FightersHelper")


//     data, err := contractABI.Pack("createFighter", common.HexToAddress(ownerAddress), name, class)
//     if err != nil {
//         log.Printf("[CreateFighter] Failed to encode function arguments: %v", err)
//     }

//     go sendBlockchainTransaction(
//         nil, 
//         "FightersHelper", 
//         FightersContract, 
//         data, 
//         "Fighters",
//         FightersContract,
//         "FighterCreated", 
//         maps.Coordinate{X: 0, Y: 0}, 
//         common.Hash{},
//         conn,
//     )
// }

func DropBox(conn *websocket.Conn, itemHash common.Hash, coords maps.Coordinate, fighter *fighters.Fighter, item *items.TokenAttributes) {
    log.Printf("[DropBox] ItemHash=%v", itemHash);

    // Load contract ABI from file
    contractABI := blockchain.LoadABI("DropHelper");

    data, err := contractABI.Pack("openBox", item.GetTokenId(), fighter.GetLocation(), big.NewInt(fighter.GetCoordinates().X), big.NewInt(fighter.GetCoordinates().Y))
    if err != nil {
        log.Printf("[DropBox] Failed to encode function arguments: %v", err)
    }

    go sendBlockchainTransaction(
        fighter, 
        "DropHelper", 
        DropHelperContract, 
        data, 
        "Drop",
        DropContract,
        "ItemDropped", 
        coords, 
        itemHash,
        nil,
    );

}


func RecordKill(opponent *fighters.Fighter) {
    coords := opponent.Coordinates

    type DamageTuple struct {
        FighterId        *big.Int
        Damage           *big.Int
    }

    killedFighter := big.NewInt(opponent.GetTokenID())
    damageDealt := opponent.GetDamageReceived()
    battleNonce := big.NewInt(time.Now().UnixNano() / int64(time.Millisecond))

    log.Printf("[recordBattleOnChain] killedFighter=%v damageDealt=%v", killedFighter, damageDealt)
    if len(damageDealt) == 0 {
        return;
    }

    damageDealtTuples := make([]DamageTuple, len(damageDealt))
    for i, d := range damageDealt {
        damageDealtTuples[i] = DamageTuple{
            FighterId:        d.FighterId,
            Damage:           d.Damage,
        }
    }


    log.Printf("[recordBattleOnChain] damageDealt 2 %v", damageDealt)
    //killer := fighters.FightersMap.Find( strconv.FormatInt(damageDealt[0].FighterId.Int64(), 10) )
    killer := PopulationMap.Find("lorencia", damageDealt[0].FighterId.Int64())

    contractABI := blockchain.LoadABI("BattleHelper")

    // Encode function arguments
    data, err := contractABI.Pack("recordKill", killedFighter, damageDealtTuples, battleNonce, opponent.GetLocation(), big.NewInt(opponent.GetCoordinates().X), big.NewInt(opponent.GetCoordinates().Y))
    if err != nil {
        log.Printf("[recordBattleOnChain] Failed to encode function arguments: %v", err)
        return
    }

    log.Printf("[recordBattleOnChain] damageDealt 3 %v", damageDealt)

    go sendBlockchainTransaction(
        killer, 
        "BattleHelper", 
        BattleHelperContract, 
        data, 
        "Drop",
        DropContract,
        "ItemDropped", 
        coords, 
        common.Hash{},
        nil,
    );
}

func lastBlockNumber() (uint64, error) {
    client := getRpcClient()
    header, err := client.HeaderByNumber(context.Background(), nil)
    if err != nil {
        return 0, err
    }

    var blockNumber = header.Number.Uint64()
    log.Printf("[lastBlockNumber] %v", blockNumber)

    return blockNumber, nil
}

func PickupDroppedItem(fighter *fighters.Fighter, itemHash common.Hash) {
    log.Printf("[PickupDroppedItem] ItemHash=%v", itemHash);

    dropEvent := items.DroppedItems.Find(itemHash)

    if dropEvent == nil {
        log.Printf("[PickupDroppedItem] DroppedItems=%v", items.DroppedItems)
        log.Printf("[PickupDroppedItem] Dropped item not found: %v", itemHash)
        return
    }    

    item := dropEvent.GetItem()

    blockNumber := dropEvent.GetBlockNumber()

    contractABI := blockchain.LoadABI("BackpackHelper");

    fighterID := big.NewInt(fighter.GetTokenID())
    log.Printf("[PickupDroppedItem] itemHash=%v item=%+v blockNumber=%v fighter=%v", itemHash, item, blockNumber, fighterID)

    data, err := contractABI.Pack("pickupItem", itemHash, item, blockNumber, fighterID)
    if err != nil {
        log.Printf("[PickupDroppedItem] Failed to encode function arguments: %v", err)
    }

    go sendBlockchainTransaction(
        fighter, 
        "BackpackHelper", 
        BackpackHelperContract, 
        data, 
        "Backpack",
        BackpackContract,
        "ItemPicked", 
        fighter.GetCoordinates(), 
        itemHash,
        nil,
    );
}

func getTokenAttributes(itemId int64) *items.TokenAttributes {
	//log.Printf("[getTokenAttributes] itemId: %v", itemId)
    if itemId == 0 {
        return &items.TokenAttributes{};
    }

    atts := items.ItemAttributesCache.Find(itemId)
    if atts != nil {
        return atts
    } else {
        dbItem, ok := items.GetItemAttributesFromDB(itemId)
        if ok {
            items.ItemAttributesCache.Add(itemId, dbItem)
            return dbItem;
        }
    }
	

	// Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(ItemsHelperContract)
    contractABI := blockchain.LoadABI("ItemsHelper")

    // Prepare the call to the getTokenAttributes function
    tokenID := big.NewInt(itemId)
    callData, err := contractABI.Pack("getTokenAttributes", tokenID)
    if err != nil {
        log.Fatalf("[getTokenAttributes] Failed to pack call data: %v", err)
    }

    // Call the contract using the Ethereum client
    result, err := client.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
    }, nil)
    if err != nil {
        log.Fatalf("[getTokenAttributes] Failed to call contract: %v", err)
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var attributes []interface{};


    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    attributes, err = contractABI.Unpack("getTokenAttributes", result)
    if err != nil {
        log.Printf("[getTokenAttributes] Failed to unpack error: %v", err)
    }

    jsonatts, err := json.Marshal(attributes[0])

    var item items.SolidityItemAtts
    json.Unmarshal(jsonatts, &item)
    if err != nil {
        log.Fatalf("[getTokenAttributes] Failed to call contract: %v", err)
    }

    tokenAtts := items.ConvertSolidityItemToGoItem(item);

   	items.ItemAttributesCache.Add(itemId, tokenAtts);
    items.SaveItemAttributesToDB(tokenAtts);
   	return tokenAtts;
}

// func getItemAttributes(itemName string) *ItemAttributes {
//     return BaseItemAttributes[itemName];
//     // //log.Printf("[getItemAttributes] itemId: %v", itemId)
//     // // if itemId == 0 {
//     // //     return ItemAttributes{};
//     // // }

//     // // Connect to the Ethereum network using an Ethereum client
//     // client := getRpcClient();

//     // // Define the contract address and ABI
//     // contractAddress := common.HexToAddress(ItemsHelperContract)
//     // contractABI := loadABI("ItemsHelper")

//     // // Prepare the call to the getTokenAttributes function
//     // // tokenID := big.NewInt(itemId)
//     // callData, err := contractABI.Pack("getItemAttributes", itemName)
//     // if err != nil {
//     //     log.Fatalf("[getItemAttributes] Failed to pack call data: %v", err)
//     // }

//     // // Call the contract using the Ethereum client
//     // result, err := client.CallContract(context.Background(), ethereum.CallMsg{
//     //     To:   &contractAddress,
//     //     Data: callData,
//     // }, nil)
//     // if err != nil {
//     //     log.Fatalf("[getItemAttributes] Failed to call contract: %v", err)
//     // }

//     // // Unpack the result into the attributes struct
//     // //var attributes []FighterAttributes
//     // var attributes []interface{};


//     // //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
//     // //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
//     // attributes, err = contractABI.Unpack("getItemAttributes", result)
//     // if err != nil {
//     //     log.Printf("[getItemAttributes] Failed to unpack error: %v", err)
//     // }

//     // jsonatts, err := json.Marshal(attributes[0])

//     // var item ItemAttributes
//     // json.Unmarshal(jsonatts, &item)
//     // if err != nil {
//     //     log.Fatalf("[getItemAttributes] Failed to call contract: %v", err)
//     // }

//     // //log.Printf("[getItemAttributes] item: %v", item)
//     // //ItemAttributesCache[itemName] = item;
//     // //saveItemAttributesToDB(item);
//     // return item;
// }

func getFighterAttributes(TokenID int64) (*fighters.FighterAttributes, error) {
    // FighterAttributesCacheMutex.RLock()
    // atts, ok := FighterAttributesCache[TokenID];
    // FighterAttributesCacheMutex.RUnlock()

    atts := fighters.FighterAttributesCache.Find(TokenID)
    if atts != nil {
        return atts, nil
    }

	// Connect to the Ethereum network using an Ethereum client
    rpcClient := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(FightersHelperContract)
    contractABI := blockchain.LoadABI("FightersHelper")

    //log.Printf("contractABI: ", contractABI.Methods[method.Name]);

    // Prepare the call to the getTokenAttributes function
    tokenId := big.NewInt(TokenID) 
    callData, err := contractABI.Pack("getTokenAttributes", tokenId)
    if err != nil {
        log.Fatalf("[getFighterAttributes] Failed to pack call data: %v", err)
    }

    // log.Printf("callData: %v ", callData);
    // log.Printf("fighterID: %v ", fighterID);

    // Call the contract using the Ethereum client
    result, err := rpcClient.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
        Gas: 3000000,
    }, nil)
    if err != nil {
        if err.Error()[:36] == "VM Exception while processing transaction" {
            reason, err := abi.UnpackRevert(result)
            if err != nil {
                log.Printf("[getFighterAttributes] Failed to decode revert reason: %v", err)
            }
            log.Printf("[getFighterAttributes] Revert reason: %v", reason)
        } else {
            log.Printf("[getFighterAttributes] Failed to call contract 1: %v err: %v", TokenID, err)

        }
        return nil, err
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var attributes []interface{};

    //log.Printf("result: %v ", result);

    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    attributes, err = contractABI.Unpack("getTokenAttributes", result)
    if err != nil {
        log.Printf("[getFighterAttributes] Failed to unpack error: %v", err)
    }

    jsonatts, err := json.Marshal(attributes[0])

    var fighterAtts fighters.FighterAttributes
    json.Unmarshal(jsonatts, &fighterAtts)
    if err != nil {
        log.Printf("[getFighterAttributes] Failed to call contract 2: %v", err)
    }
   	//log.Printf("[getFighterAttributes] fighter: %v", fighterAtts)

    // FighterAttributesCacheMutex.Lock()
    // FighterAttributesCache[TokenID] = fighterAtts
    // FighterAttributesCacheMutex.Unlock()

    fighters.FighterAttributesCache.Add(TokenID, &fighterAtts)

   	return &fighterAtts, nil;
}

func getFighterMoney(fighter *fighters.Fighter) int64 {

    //log.Printf("[getFighterAttributes] id: %v", id)
    fighterID := fighter.TokenID
    
    // Connect to the Ethereum network using an Ethereum client
    rpcClient := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(MoneyContract)
    contractABI := blockchain.LoadABI("MoneyHelper")

    //log.Printf("contractABI: ", contractABI);

    // Prepare the call to the getTokenAttributes function
    tokenId := common.HexToAddress(fighter.OwnerAddress)
    callData, err := contractABI.Pack("balanceOf", tokenId)
    if err != nil {
        log.Fatalf("[getFighterMoney] Failed to pack call data: %v", err)
    }

    // log.Printf("callData: %v ", callData);
    // log.Printf("fighterID: %v ", fighterID);

    // Call the contract using the Ethereum client
    result, err := rpcClient.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
        Gas: 3000000,
    }, nil)
    if err != nil {
        if err.Error()[:36] == "VM Exception while processing transaction" {
            reason, err := abi.UnpackRevert(result)
            if err != nil {
                log.Fatalf("[getFighterMoney] Failed to decode revert reason: %v", err)
            }
            log.Fatalf("[getFighterMoney] Revert reason: %v", reason)
        } else {
            log.Fatalf("[getFighterMoney] Failed to call contract: %v, err: %v", fighterID, err)
        }
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var moneyResp  []interface{};

    //log.Printf("result: %v ", result);



    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    moneyResp, err = contractABI.Unpack("balanceOf", result)
    if err != nil {
        log.Printf("[getFighterMoney] Failed to unpack error: %v", err)
    }

    jsonatts, err := json.Marshal(moneyResp[0])

    var money big.Int
    json.Unmarshal(jsonatts, &money)
    if err != nil {
        log.Fatalf("[getFighterAttributes] Failed to call contract: %v", err)
    }

    divisor := new(big.Float).SetInt(big.NewInt(1e18))
    moneyFloat := new(big.Float).SetInt(&money)
    rounded := new(big.Float).Quo(moneyFloat, divisor)

    // Round to 0 decimal places
    roundedInt64, _ := rounded.Int64()
    //log.Printf("[getFighterMoney] fighter:%v money: %v", id, roundedInt64)
    return roundedInt64
}

func waitForReceiptP(client *ethclient.Client, txHash common.Hash, contractAddress common.Address, receiptChan chan *types.Receipt, errChan chan error) {
    for {
        receipt, err := client.TransactionReceipt(context.Background(), txHash)
        if err == nil {
            filteredReceipt := new(types.Receipt)
            *filteredReceipt = *receipt
            filteredReceipt.Logs = nil
            for _, log1 := range receipt.Logs {
                log.Printf("[waitForReceiptP] log=%v", log1)
                if log1.Address == contractAddress {
                    filteredReceipt.Logs = append(filteredReceipt.Logs, log1)
                }
            }
            receiptChan <- filteredReceipt
            return
        }
        if err != ethereum.NotFound {
            errChan <- err
            return
        }
        time.Sleep(1 * time.Second)
    }
}


func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
    for {
        receipt, err := client.TransactionReceipt(context.Background(), txHash)
        if err == nil {
            return receipt, nil
        }
        if err != ethereum.NotFound {
            return nil, err
        }
        time.Sleep(1 * time.Second)
    }
}

func getRpcClient() *ethclient.Client {
	// Connect to the Ethereum network using an Ethereum client
	client, err := ethclient.Dial(RpcClinetAddress)
	if err != nil {
		log.Fatalf("[getRpcClient] Failed to connect to Ethereum network: %v", err)
	}

	return client
}

// func getFighterStats(fighterID int64) fighters.FighterStats {

// 	// Connect to the Ethereum network using an Ethereum client
//     client := getRpcClient();

//     // Define the contract address and ABI
//     contractAddress := common.HexToAddress(FightersHelperContract)
//     contractABI := blockchain.LoadABI("FightersHelper")

//     // Prepare the call to the getTokenAttributes function
//     tokenID := big.NewInt(fighterID)
//     callData, err := contractABI.Pack("getFighterStats", tokenID)
//     if err != nil {
//         log.Fatalf("[getFighterStats] Failed to pack call data: %v", err)
//     }

//     // Call the contract using the Ethereum client
//     result, err := client.CallContract(context.Background(), ethereum.CallMsg{
//         To:   &contractAddress,
//         Data: callData,
//     }, nil)
//     if err != nil {
//         log.Fatalf("[getFighterStats] Failed to call contract fighterID=%v error=%v", fighterID, err)
//     }

//     // Unpack the result into the attributes struct
//     //var attributes []FighterAttributes
//     var attributes []interface{};


//     //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
//     //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
//     attributes, err = contractABI.Unpack("getFighterStats", result)
//     if err != nil {
//         log.Printf("[getFighterStats] Failed to unpack error: %v", err)
//     }

//     jsonatts, err := json.Marshal(attributes[0])

//     var fighter fighters.FighterStats
//     json.Unmarshal(jsonatts, &fighter)
//     if err != nil {
//         log.Fatalf("[getFighterStats] Failed to call contract fighterID %v err: %v", fighterID, err)
//     }

//    	return fighter;
// }

func getFighterItems(fighter *fighters.Fighter)  {
    //fighter := getFighterSafely(convertIdToString(FighterId))

    //log.Printf("[getFighterItems] FighterId: %v", FighterId)
    

	// Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(ItemsHelperContract)
    contractABI := blockchain.LoadABI("ItemsHelper")

    // Prepare the call to the getTokenAttributes function
    tokenID := big.NewInt(fighter.GetTokenID())
    callData, err := contractABI.Pack("getFighterItems", common.HexToAddress(fighter.GetOwnerAddress()), tokenID)
    if err != nil {
        log.Fatalf("[getFighterItems] Failed to pack call data: %v", err)
    }


    //log.Printf("[getFighterItems] FreqData: %v", reqData)

    // Call the contract using the Ethereum client
    result, err := client.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
    }, nil)
    if err != nil {
        log.Fatalf("[getFighterItems] Failed to call contract: %v", err)
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var attributes []interface{};

	//log.Printf("[getFighterItems] Result", result)
    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    attributes, err = contractABI.Unpack("getFighterItems", result)
    if err != nil {
        log.Printf("[getFighterItems] Failed to unpack error: %v", err)
    }

    var its []*items.TokenAttributes

	//log.Printf("[getFighterItems] attributes: %v", attributes)

	for _, v := range attributes {
		attrs, ok := v.([][2]*big.Int)
		if !ok {
	        // handle invalid attribute format
	        log.Printf("[getFighterItems] Error iterating attributes")
	    }

	    for _, attr := range attrs {

	    	itemId := attr[0]
	    	// lastUpdBlock := attr[1]

	    	// log.Printf("[getFighterItems] itemId=",itemId ," lastUpdBlock=", lastUpdBlock)

	    	// get item attributes
	    	itemAttributes := getTokenAttributes(itemId.Int64());

	    	//recordItemToDB(itemAttributes);

	    	its = append(its, itemAttributes);

	    }
	}

    //log.Printf("[getFighterItems] items: %v", items)

	

	//stats := getFighterStats(fighter.GetTokenID());

    //jsonstats, err := json.Marshal(stats)
    //log.Print("[getFighter] jsonstats: %s", stats)


    fighterAttributes, err := getFighterAttributes(fighter.GetTokenID());
    jsonfighteratts, err := json.Marshal(fighterAttributes)

    //equipment := fighter.Equipment;

	//jsonequip, err := json.Marshal(equipment)
    //log.Print("[getFighter] jsonstats: %s", equipment)


    

    npcatts := getNPCs("lorencia_0_0");
    //log.Print("[getFighterItems] npcs: ", npcs)
    jsonnpcs, err := json.Marshal(npcatts)


    jsonfighter, err := json.Marshal(fighter)

    type jsonResponse struct {
		Action string `json:"action"`
		Items []*items.TokenAttributes `json:"items"`
		Attributes string `json:"attributes"`
		Equipment map[int64]*inventory.InventorySlot `json:"equipment"`
		//Stats string `json:"stats"`
		NPCs string `json:"npcs"`
		Fighter string `json:"fighter"`
        Money int64 `json:"money"`
        DroppedItems map[common.Hash]*items.ItemDroppedEventGo `json:"droppedItems"`
        Backpack *inventory.Inventory `json:"backpack"`
	}

    jsonResp := jsonResponse{
    	Action: "fighter_items",
    	Items: its,
    	Attributes: string(jsonfighteratts),
    	Equipment: fighter.GetEquipment().GetMap(),
    	//Stats: string(jsonstats),
    	NPCs: string(jsonnpcs),
    	Fighter: string(jsonfighter),
        Money: getFighterMoney(fighter),
        DroppedItems: items.GetDroppedItemsInGo(),
        Backpack: fighter.GetBackpack(),
    }


    response, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[getFighterItems] error: ", err)
        return
    }
    respondFighter(fighter, response)
}

