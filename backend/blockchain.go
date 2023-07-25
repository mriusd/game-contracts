package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/crypto"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/gorilla/websocket"

	"context"
    "log"
    "math/big"
    "encoding/json"

    "io/ioutil"
    "fmt"
    "strings"
    "time"
    "strconv"
)


// This function is vor development only
// Will not work in production
func MakeItem(fighter *Fighter, item *ItemAttributes) {
    log.Printf("[MakeItem] ItemAttributes=%v", item);


    // Load contract ABI from file
    contractABI := loadABI("Items");

    data, err := contractABI.Pack("makeItem", item)
    if err != nil {
        log.Printf("[MakeItem] Failed to encode function arguments: %v", err)
    }

    sendBlockchainTransactionShort(
        fighter, 
        "Items", 
        ItemsContract, 
        data, 
        "ItemDropped", 
        fighter.Coordinates, 
        common.Hash{},
    )
}

func sendBlockchainTransactionShort(
    fighter *Fighter, 
    contractName string, 
    contractAdr string, 
    data []byte, 
    eventTolisten string, 
    coords Coordinate, 
    someHash common.Hash,
) {

    sendBlockchainTransaction(
        fighter, 
        contractName,
        contractAdr,
        data,
        contractName,
        contractAdr,
        eventTolisten,
        coords,
        someHash,
        nil,
    )
}

func sendBlockchainTransaction(
    fighter *Fighter, 
    contractName string, 
    contractAdr string, 
    data []byte, 
    eventContractName string,
    eventContractAddress string,
    eventTolisten string, 
    coords Coordinate, 
    someHash common.Hash,
    conn *websocket.Conn,
) {
    // Connect to the Ethereum network
    client      := getRpcClient();

    // Load your private key
    privateKey, err := crypto.HexToECDSA(PrivateKey)
    if err != nil {
        log.Fatalf("[sendBlockchainTransaction] Failed to load private key: %v", err)
    }

    

    // Set contract address
    contractAddress := common.HexToAddress(contractAdr)

    // Prepare transaction options
    nonce, err := client.NonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey), nil)
    if err != nil {
        log.Printf("[sendBlockchainTransaction] Failed to retrieve nonce: %v", err)
    }
    gasLimit := uint64(3000000)
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        log.Printf("[sendBlockchainTransaction] Failed to retrieve gas price: %v", err)
    }
    auth := bind.NewKeyedTransactor(privateKey)
    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)
    auth.GasLimit = gasLimit
    auth.GasPrice = gasPrice

    // Create transaction and sign it
    tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)
    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(RPCNetworkID), privateKey)
    if err != nil {
        log.Printf("[sendBlockchainTransaction] Failed to sign transaction: %v", err)
    }

    // Send the transaction
    err = client.SendTransaction(context.Background(), signedTx)
    if err != nil {
        log.Printf("[sendBlockchainTransaction] Failed to send transaction: %v", err)
    } else {
        fmt.Println("[sendBlockchainTransaction] Transaction hash:", signedTx.Hash().Hex())

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

func handleBlockchainEvent(eventName, contractName string, receipt *types.Receipt, fighter *Fighter, coords Coordinate, someHash common.Hash, conn *websocket.Conn) {
    parsedABI := loadABI(contractName)
    switch eventName {
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
            event := ItemDroppedEvent{}

            err := parsedABI.UnpackIntoInterface(&event, eventName, receipt.Logs[0].Data)
            if err != nil {
                log.Printf("[handleBlockchainEvent] Failed to unpack log data: %v", err)
                return
            }

            log.Printf("[handleBlockchainEvent] event=%v", event)

            event.BlockNumber = receipt.BlockNumber
            event.Coords = coords
            event.OwnerId = big.NewInt(fighter.TokenID)

            // Add a self-destruct timer to remove the item from the map after 30 seconds
            time.AfterFunc(30*time.Second, func() {
                DroppedItemsMutex.Lock() // Use a mutex if needed to protect access to the map
                delete(DroppedItems, event.ItemHash)
                DroppedItemsMutex.Unlock()
                log.Printf("[handleBlockchainEvent] Item with hash %v has been removed after 30 seconds", event.ItemHash)
                broadcastDropMessage()
            })

            DroppedItemsMutex.Lock()
            DroppedItems[event.ItemHash] = &event
            DroppedItemsMutex.Unlock()


            if getBackpackSlotByHash(fighter, someHash) != nil {
                fighter.Backpack.removeItemByHash(fighter, someHash)
            } else {
                removeItemFromEquipmentSlotByHash(fighter, someHash)
            }

            removeItemFromDB(event.TokenId.Int64())

            broadcastDropMessage()

        break

        case "ItemDropped":
            event := ItemDroppedEvent{}
            log.Printf("[handleBlockchainEvent:ItemDropped] receipt: %v", receipt)
            err := parsedABI.UnpackIntoInterface(&event, "ItemDropped", receipt.Logs[0].Data)
            if err != nil {
                log.Printf("[handleBlockchainEvent:ItemDropped] Failed to unpack log data: %v", err)
                return
            }

            log.Printf("[handleBlockchainEvent:ItemDropped] ItemHash: %v", event.ItemHash)

            event.BlockNumber = receipt.BlockNumber
            event.Coords = coords
            event.OwnerId = big.NewInt(fighter.TokenID)

            // Add a self-destruct timer to remove the item from the map after 30 seconds
            time.AfterFunc(30*time.Second, func() {
                DroppedItemsMutex.Lock() // Use a mutex if needed to protect access to the map
                delete(DroppedItems, event.ItemHash)
                DroppedItemsMutex.Unlock()
                log.Printf("[handleBlockchainEvent:ItemDropped] Item with hash %v has been removed after 30 seconds", event.ItemHash)
                broadcastDropMessage()
            })

            DroppedItemsMutex.Lock()
            DroppedItems[event.ItemHash] = &event
            DroppedItemsMutex.Unlock()

            if getBackpackSlotByHash(fighter, someHash) != nil {
                fighter.Backpack.removeItemByHash(fighter, someHash)
            } else {
                removeItemFromEquipmentSlotByHash(fighter, someHash)
            }

            removeItemFromDB(event.TokenId.Int64())

            broadcastDropMessage()
        break

        case "ItemPicked":
            event := ItemPickedEvent{}

            log.Printf("[handleBlockchainEvent:ItemPicked] logEntry: %v", receipt.Logs[0])

            err := parsedABI.UnpackIntoInterface(&event, "ItemPicked", receipt.Logs[0].Data)
            if err != nil {
                log.Printf("[handleBlockchainEvent:ItemPicked] Failed to unpack log data: %v", err)
                return
            }

            DroppedItemsMutex.Lock()
            dropEvent := DroppedItems[someHash]
            DroppedItemsMutex.Unlock()
            
            item := dropEvent.Item
            item.TokenId = event.TokenId
            saveItemAttributesToDB(item)

            DroppedItemsMutex.Lock()
            delete(DroppedItems, someHash)
            DroppedItemsMutex.Unlock()

            if (item.ItemAttributesId.Int64() != GoldItemId) {
                fighter.Mutex.Lock()
                _, _, err := fighter.Backpack.AddItem(item, dropEvent.Qty.Int64(), someHash)
                fighter.Mutex.Unlock()
                saveBackpackToDB(fighter)
                if err != nil {
                    log.Printf("[handleBlockchainEvent:ItemPicked] Backpack full: %v", someHash)
                    sendErrorMessage(fighter, "Backpack full")
                }
            }

            fmt.Printf("[handleBlockchainEvent:ItemPicked] event: %+v\n", event)   
            broadcastPickupMessage(fighter, item, event.Qty)

        break 


        default:
            log.Printf("[handleBlockchainEvent] Uknown eventName=%v receipt=%v", eventName, receipt);
        break
    }
}


func getUserFighters(conn *websocket.Conn)  {
    log.Printf("[getUserFighters]")

    if (Connections[conn].OwnerAddress == common.Address{}) {
        log.Fatalf("[getUserFighters] Zero address")
        return
    }

    // Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(FighterAttributesContract)
    contractABI := loadABI("FighterAttributes")

    // Prepare the call to the getTokenAttributes function
    Connections[conn].Mutex.RLock()
    callData, err := contractABI.Pack("getUserFighters", Connections[conn].OwnerAddress)
    Connections[conn].Mutex.RUnlock()
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

    var fighters []FighterAttributes

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

            fighters = append(fighters, fighterAttributes);
        }
    }


    log.Printf("[getUserFighters] fighters: %v", fighters)


    type jsonResponse struct {
        Action string `json:"action"`
        Fighters []FighterAttributes `json:"fighters"`        
    }

    jsonResp := jsonResponse{
        Action: "user_fighters",
        Fighters: fighters,
    }


    jr, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[getUserFighters] error: ", err)
        return
    }
    respondConn(conn, jr)
}


func DropBackpackItem(conn *websocket.Conn, itemHash common.Hash, coords Coordinate) {
    log.Printf("[DropBackpackItem] ItemHash=%v", itemHash);

    fighter     := findFighterByConn(conn)
    bacpackSlot := getBackpackSlotByHash(fighter, itemHash)

    if bacpackSlot == nil {
        log.Printf("[DropBackpackItem] Item not found in backpack: %v", itemHash)

        bacpackSlot = getEquipmentSlotByHash(fighter, itemHash)

        if bacpackSlot == nil {
            log.Printf("[DropBackpackItem] Item not found in equipment: %v", itemHash)
            return
        }
    }    

    item := bacpackSlot.Attributes   

    if item.IsBox {
        DropBox(conn, itemHash, coords, fighter, item)
        return
    }

    log.Printf("[DropBackpackItem] TokenId=%v item=%v", item.TokenId, item);
    

    // Load contract ABI from file
    contractABI := loadABI("Backpack");

    data, err := contractABI.Pack("dropBackpackItem", item.TokenId)
    if err != nil {
        log.Printf("[DropBackpackItem] Failed to encode function arguments: %v", err)
    }

    sendBlockchainTransactionShort(fighter, "Backpack", BackpackContract, data, "BackpackItemDropped", coords, itemHash);

    
}

func CreateFighter(conn *websocket.Conn, ownerAddress, name string, class uint8) {
    log.Printf("[CreateFighter] ownerAddress=%v, class=%v", ownerAddress, class);

    // Load contract ABI from file
    contractABI := loadABI("FighterAttributes");

    data, err := contractABI.Pack("createFighter", common.HexToAddress(ownerAddress), name, class)
    if err != nil {
        log.Printf("[CreateFighter] Failed to encode function arguments: %v", err)
    }

    sendBlockchainTransaction(
        nil, 
        "FighterAttributes", 
        FighterAttributesContract, 
        data, 
        "FighterAttributes",
        FighterAttributesContract,
        "FighterCreated", 
        Coordinate{X: 0, Y: 0}, 
        common.Hash{},
        conn,
    )
}

func DropBox(conn *websocket.Conn, itemHash common.Hash, coords Coordinate, fighter *Fighter, item ItemAttributes) {
    log.Printf("[DropBox] ItemHash=%v", itemHash);

    // Load contract ABI from file
    contractABI := loadABI("Items");

    data, err := contractABI.Pack("openBox", item.TokenId)
    if err != nil {
        log.Printf("[DropBackpackItem] Failed to encode function arguments: %v", err)
    }

    sendBlockchainTransactionShort(fighter, "Items", ItemsContract, data, "ItemDropped", coords, itemHash);
}


func recordBattleOnChain(opponent *Fighter) {
    coords := opponent.Coordinates

    type DamageTuple struct {
        FighterId        *big.Int
        Damage           *big.Int
    }

    killedFighter := big.NewInt(opponent.TokenID)
    damageDealt := opponent.DamageReceived
    battleNonce := big.NewInt(time.Now().UnixNano() / int64(time.Millisecond))

    log.Printf("[recordBattleOnChain] damageDealt %v", damageDealt)

    damageDealtTuples := make([]DamageTuple, len(damageDealt))
    for i, d := range damageDealt {
        damageDealtTuples[i] = DamageTuple{
            FighterId:        d.FighterId,
            Damage:           d.Damage,
        }
    }

    killer := getFighterSafely( strconv.FormatInt(damageDealt[0].FighterId.Int64(), 10) )

    contractABI := loadABI("Battle")

    // Encode function arguments
    data, err := contractABI.Pack("recordKill", killedFighter, damageDealtTuples, battleNonce)
    if err != nil {
        log.Printf("[recordBattleOnChain]Failed to encode function arguments: %v", err)
    }

    sendBlockchainTransaction(
        killer, 
        "Battle", 
        BattleContract, 
        data, 
        "Items",
        ItemsContract,
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

func PickupDroppedItem(fighter *Fighter, itemHash common.Hash) {
    log.Printf("[PickupDroppedItem] ItemHash=%v", itemHash);

    dropEvent, ok := DroppedItems[itemHash]

    if !ok {
        log.Printf("[PickupDroppedItem] Dropped item not found: %v", itemHash)
        return
    }    

    item := dropEvent.Item
    blockNumber := dropEvent.BlockNumber

    contractABI := loadABI("Backpack");

    fighterID := big.NewInt(fighter.TokenID)
    log.Printf("[PickupDroppedItem] itemHash=%v item=%+v blockNumber=%v fighter=%v", itemHash, item, blockNumber, fighterID)

    data, err := contractABI.Pack("pickupItem", itemHash, item, blockNumber, fighterID)
    if err != nil {
        log.Printf("[PickupDroppedItem] Failed to encode function arguments: %v", err)
    }

    sendBlockchainTransactionShort(
        fighter, 
        "Backpack", 
        BackpackContract, 
        data, 
        "ItemPicked", 
        fighter.Coordinates, 
        itemHash,
    )
}

func getTokenAttributes(itemId int64) ItemAttributes {
	//log.Printf("[getItemAttributes] itemId: %v", itemId)
    if itemId == 0 {
        return ItemAttributes{};
    }

    atts, ok := ItemAttributesCache[itemId]
    if ok {
        return atts
    } else {
        dbItem, ok := getItemAttributesFromDB(itemId)
        if ok {
            return dbItem;
        }
    }
	

	// Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(ItemsContract)
    contractABI := loadABI("Items")

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

    var item ItemAttributes
    json.Unmarshal(jsonatts, &item)
    if err != nil {
        log.Fatalf("[getTokenAttributes] Failed to call contract: %v", err)
    }

   	ItemAttributesCache[itemId] = item;
    saveItemAttributesToDB(item);
   	return item;
}

func getItemAttributes(itemId int64) ItemAttributes {
    //log.Printf("[getItemAttributes] itemId: %v", itemId)
    if itemId == 0 {
        return ItemAttributes{};
    }

    // Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(ItemsContract)
    contractABI := loadABI("Items")

    // Prepare the call to the getTokenAttributes function
    tokenID := big.NewInt(itemId)
    callData, err := contractABI.Pack("getItemAttributes", tokenID)
    if err != nil {
        log.Fatalf("[getItemAttributes] Failed to pack call data: %v", err)
    }

    // Call the contract using the Ethereum client
    result, err := client.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
    }, nil)
    if err != nil {
        log.Fatalf("[getItemAttributes] Failed to call contract: %v", err)
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var attributes []interface{};


    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    attributes, err = contractABI.Unpack("getItemAttributes", result)
    if err != nil {
        log.Printf("[getItemAttributes] Failed to unpack error: %v", err)
    }

    jsonatts, err := json.Marshal(attributes[0])

    var item ItemAttributes
    json.Unmarshal(jsonatts, &item)
    if err != nil {
        log.Fatalf("[getItemAttributes] Failed to call contract: %v", err)
    }

    //log.Printf("[getItemAttributes] item: %v", item)
    ItemAttributesCache[itemId] = item;
    saveItemAttributesToDB(item);
    return item;
}

func getFighterAttributes(TokenID int64) (FighterAttributes, error) {
    FighterAttributesCacheMutex.RLock()
    atts, ok := FighterAttributesCache[TokenID];
    FighterAttributesCacheMutex.RUnlock()
    if ok {
        return atts, nil
    }

	// Connect to the Ethereum network using an Ethereum client
    rpcClient := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(FighterAttributesContract)
    contractABI := loadABI("FighterAttributes")

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
            log.Printf("[getFighterAttributes] Failed to call contract 1: %v", TokenID, err)

        }
        return FighterAttributes{}, err
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

    var fighterAtts FighterAttributes
    json.Unmarshal(jsonatts, &fighterAtts)
    if err != nil {
        log.Printf("[getFighterAttributes] Failed to call contract 2: %v", err)
    }
   	//log.Printf("[getFighterAttributes] fighter: %v", fighter)

    FighterAttributesCacheMutex.Lock()
    FighterAttributesCache[TokenID] = fighterAtts
    FighterAttributesCacheMutex.Unlock()
   	return fighterAtts, nil;
}

func getFighterMoney(fighter *Fighter) int64 {

    //log.Printf("[getFighterAttributes] id: %v", id)
    fighterID := fighter.TokenID
    
    // Connect to the Ethereum network using an Ethereum client
    rpcClient := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(MoneyContract)
    contractABI := loadABI("Money")

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
            log.Fatalf("[getFighterMoney] Failed to call contract: %v", fighterID, err)
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

func getFighterStats(fighterID int64) FighterStats {

	// Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(FighterAttributesContract)
    contractABI := loadABI("FighterAttributes")

    // Prepare the call to the getTokenAttributes function
    tokenID := big.NewInt(fighterID)
    callData, err := contractABI.Pack("getFighterStats", tokenID)
    if err != nil {
        log.Fatalf("[getFighterStats] Failed to pack call data: %v", err)
    }

    // Call the contract using the Ethereum client
    result, err := client.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
    }, nil)
    if err != nil {
        log.Fatalf("[getFighterStats] Failed to call contract fighterID=",fighterID," error=", err)
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var attributes []interface{};


    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    attributes, err = contractABI.Unpack("getFighterStats", result)
    if err != nil {
        log.Printf("[getFighterStats] Failed to unpack error: %v", err)
    }

    jsonatts, err := json.Marshal(attributes[0])

    var fighter FighterStats
    json.Unmarshal(jsonatts, &fighter)
    if err != nil {
        log.Fatalf("[getFighterStats] Failed to call contract fighterID ", fighterID, "error=", err)
    }

   	return fighter;
}

func getFighterItems(FighterId int64)  {
    fighter := getFighterSafely(convertIdToString(FighterId))

    //log.Printf("[getFighterItems] FighterId: %v", FighterId)
    

	// Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(ItemsContract)
    contractABI := loadABI("Items")

    // Prepare the call to the getTokenAttributes function
    tokenID := big.NewInt(FighterId)
    callData, err := contractABI.Pack("getFighterItems", common.HexToAddress(fighter.OwnerAddress), tokenID)
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

    var items []ItemAttributes

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

	    	items = append(items, itemAttributes);

	    }
	}

    //log.Printf("[getFighterItems] items: %v", items)

	

	stats := getFighterStats(FighterId);

    jsonstats, err := json.Marshal(stats)
    //log.Print("[getFighter] jsonstats: %s", stats)


    fighterAttributes, err := getFighterAttributes(fighter.TokenID);
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
		Items []ItemAttributes `json:"items"`
		Attributes string `json:"attributes"`
		Equipment map[int64]*BackpackSlot `json:"equipment"`
		Stats string `json:"stats"`
		NPCs string `json:"npcs"`
		Fighter string `json:"fighter"`
        Money int64 `json:"money"`
        DroppedItems map[common.Hash]*ItemDroppedEvent `json:"droppedItems"`
        Backpack *Backpack `json:"backpack"`
	}

    jsonResp := jsonResponse{
    	Action: "fighter_items",
    	Items: items,
    	Attributes: string(jsonfighteratts),
    	Equipment: fighter.Equipment,
    	Stats: string(jsonstats),
    	NPCs: string(jsonnpcs),
    	Fighter: string(jsonfighter),
        Money: getFighterMoney(fighter),
        DroppedItems: getDroppedItemsSafely(fighter),
        Backpack: fighter.Backpack,
    }


    response, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[getFighterItems] error: ", err)
        return
    }
    respondFighter(fighter, response)
}

func loadABI(contract string) (abi.ABI) {
    // Read the contract ABI file
    raw, err := ioutil.ReadFile("../build/contracts/" + contract + ".json")
    if err != nil {
        panic(fmt.Sprintf("Error reading ABI file: %v", err))
    }

    // Unmarshal the ABI JSON into the contractABI object
    var contractABIContent struct {
        ABI json.RawMessage `json:"abi"`
    }

    err = json.Unmarshal(raw, &contractABIContent)
    if err != nil {
        panic(fmt.Sprintf("Error unmarshalling ABI JSON: %v", err))
    }

    // Use the abi.JSON function to parse the ABI directly
    parsedABI, err := abi.JSON(strings.NewReader(string(contractABIContent.ABI)))
    if err != nil {
        panic(fmt.Sprintf("Error parsing ABI JSON: %v", err))
    }

    return parsedABI
}