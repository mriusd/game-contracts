package main

import (
    "log"
    "fmt"
    "encoding/json"

    "github.com/gorilla/websocket"
    "net/http"
    "github.com/ethereum/go-ethereum/common"

    "runtime/debug"

    "github.com/mriusd/game-contracts/maps" 
    "github.com/mriusd/game-contracts/fighters"
)

type WsMessage struct {
    Type string  `json:"type"`
    Data fighters.Fighter `json:"data"`
}


func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    log.Println("[handleWebSocket] handleWebSocket start")
    var msg struct {
        Type string `json:"type"`
        Data json.RawMessage `json:"data"`
    }


    var upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool {
            // allow all connections by default
            return true
        },
    }

    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("[handleWebSocket] Failed to upgrade to WebSocket: %v", err)
        return
    }
    defer c.Close()

    conn := ConnectionsMap.Add(c)

    for {
        // Use defer/recover to catch any panic inside the loop
        defer func() {
            if r := recover(); r != nil {
                log.Printf("[handleWebSocket] Recovered from %v ", r)
                debug.PrintStack()
                
            }
        }()

        _, message, err := c.ReadMessage()


        if err != nil {
            if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
                log.Printf("[handleWebSocket] WebSocket closed err=%v message=%v", err, message)
            } else {
                log.Printf("[handleWebSocket] Failed to read message from WebSocket err=%v message=%v", err, message)
            }
            ConnectionsMap.Remove(c)
            break
        }

        log.Printf("message: %v", message)

        err = json.Unmarshal(message, &msg)
        if err != nil {
            log.Printf("[handleWebSocket] websocket unmarshal message=%v error=%v", message, err)
            continue
        }

        log.Printf("Type: %v", msg.Type)
        switch msg.Type {
            case "create_fighter":
                type CreateFighterData struct {
                    OwnerAddress string `json:"ownerAddress"`
                    FighterClass string `json:"fighterClass"`
                    Name string `json:"name"`
                }

                var reqData CreateFighterData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:create_fighter] websocket unmarshal error: %v", err)
                    continue
                }

                _, err = fighters.CreateFighter(reqData.OwnerAddress, reqData.Name, reqData.FighterClass)
                if err != nil {
                    sendErrorMsgToConn(c, "SYSTEM", fmt.Sprintf("Failed to create fighter. Error: %v", err))
                }

                //fighters.GetUserFighters(reqData.OwnerAddress)
                //fighters.PushUserFighters(conn, reqData.OwnerAddress)
                serializedFighterList, err := fighters.GetJsonSerializedFighters(reqData.OwnerAddress)
                if err != nil {
                    sendErrorMsgToConn(c, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }
                respondConn(c, serializedFighterList)
            continue

            case "get_user_fighters":
                type UserAddressData struct {
                    OwnerAddress string `json:"ownerAddress"`
                }

                var reqData UserAddressData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:get_user_fighters] websocket unmarshal error: %v", err)
                    continue
                }

                ConnectionsMap.SetConnectionOwnerAddress(c, common.HexToAddress(reqData.OwnerAddress))
                
                //fighters.GetUserFighters(reqData.OwnerAddress)
                //fighters.PushUserFighters(conn, reqData.OwnerAddress)
                serializedFighterList, err := fighters.GetJsonSerializedFighters(reqData.OwnerAddress)
                if err != nil {
                    sendErrorMsgToConn(c, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }
                respondConn(c, serializedFighterList)
            continue

            case "update_fighter_direction":
                type FighterDirection struct {
                    Direction maps.Direction `json:"direction"`
                }

                var reqData FighterDirection
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] websocket unmarshal error: %v", err)
                    continue
                }
                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                } else {
                    go updateFighterDirection(fighter, reqData.Direction)
                }
                
            continue

            case "auth":
                //log.Printf("[handleWebSocket] auth: %v", msg.Data)
                type AuthData struct {
                    PlayerID     int64  `json:"playerID"`
                    UserAddress  string `json:"userAddress"`
                    LocationHash string `json:"locationHash"` 
                }

                var reqData AuthData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket] websocket unmarshal error: %v", err)
                    continue
                }

                //log.Printf("[handleWebSocket] reqData: %v", reqData)
                fighter, err := authFighter(conn, reqData.PlayerID, reqData.UserAddress, reqData.LocationHash);
                if err != nil {
                    sendErrorMsgToConn(c, "SYSTEM", fmt.Sprintf("Auth failed. Error: %v", err))
                }

                ConnectionsMap.AddWithValues(c, fighter, common.HexToAddress(reqData.UserAddress))
            case "submit_attack":
                ProcessHit(conn, msg.Data)

            case "get_fighter_items":
                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                } else {
                    go getFighterItems(fighter)
                }
                
            continue

            case "pickup_dropped_item":
                type PickUpData struct {
                    ItemHash     string  `json:"itemHash"` 
                }
                var reqData PickUpData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:pickup_dropped_item] websocket unmarshal error: %v", err)
                    continue
                }
                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                }

                PickupDroppedItem(fighter, common.HexToHash(reqData.ItemHash))
            continue

            case "move_fighter":

                var reqData maps.Coordinate
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:move_fighter] websocket unmarshal error: %v", err)
                    continue
                }
                
                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                }
                log.Printf("[handleWebSocket] move_fighter: %v fighter=%v", reqData, fighter)
                moveFighter(fighter, reqData)
            continue

            case "update_backpack_item_position":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:update_backpack_item_position] websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                }
                fighter.Backpack.UpdateInventoryPosition(common.HexToHash(reqData.ItemHash), reqData.Position)
            continue

            case "update_vault_item_position":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:update_backpack_item_position] websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                }
                fighter.Vault.UpdateInventoryPosition(common.HexToHash(reqData.ItemHash), reqData.Position)
            continue

            case "move_item_from_backpack_to_vault":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:drop_backpack_item] websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                }



                tokenAtts := fighter.GetVault().FindByHash(common.HexToHash(reqData.ItemHash))
                fighter.GetBackpack().RemoveItemByHash(common.HexToHash(reqData.ItemHash));

                fighter.GetVault().AddItemToPosition(tokenAtts.Attributes, tokenAtts.Qty, common.HexToHash(reqData.ItemHash), int(reqData.Position.X), int(reqData.Position.Y));
            continue


            case "move_item_from_vault_to_backpack":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:drop_backpack_item] websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                }
                tokenAtts := fighter.GetVault().FindByHash(common.HexToHash(reqData.ItemHash))
                fighter.GetVault().RemoveItemByHash(common.HexToHash(reqData.ItemHash))

                fighter.GetBackpack().AddItemToPosition(tokenAtts.Attributes, tokenAtts.Qty, common.HexToHash(reqData.ItemHash), int(reqData.Position.X), int(reqData.Position.Y));
            continue

            case "drop_backpack_item":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:drop_backpack_item] websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                }

                DropBackpackItem(c, common.HexToHash(reqData.ItemHash), fighter.Coordinates)
            continue

            case "equip_backpack_item":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Slot int64 `json:"slot"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:equip_backpack_item]  websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                }

                EquipBackpackItem(fighter, common.HexToHash(reqData.ItemHash), reqData.Slot)
            continue

            case "unequip_backpack_item":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"slot"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:unequip_backpack_item]  websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                }

                UnequipBackpackItem(fighter, common.HexToHash(reqData.ItemHash), reqData.Position)
            continue

            case "message":
                type ReqData struct {
                    Text  string `json:"text"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:message]  websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] fighter not found: %v", err)
                    sendErrorMsgToConn(c, "SYSTEM", "Fighter not found")
                    continue
                }

                handleCommand(fighter, reqData.Text)
            continue

        

            // case "getFighter":
            //     getFighter(conn, msg.Data, w, r)
            //     continue
                
            default:
                log.Printf("[handleWebSocket] unknown message type: %s", msg.Type)
        }

        // data := decodeJson(message);
        // log.Printf("[handleWebSocket] Received message: %s\n", data["action"].(string))
    }
}