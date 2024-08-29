package main

import (
    "log"
    "fmt"
    "encoding/json"

    "github.com/gorilla/websocket"
    "net/http"

    "runtime/debug"

    "github.com/mriusd/game-contracts/account" 
    "github.com/mriusd/game-contracts/maps" 
    "github.com/mriusd/game-contracts/fighters"
    "github.com/mriusd/game-contracts/drop"
    "github.com/mriusd/game-contracts/shop"
    "github.com/mriusd/game-contracts/trade"
)

type WsMessage struct {
    Type string  `json:"type"`
    Data fighters.Fighter `json:"data"`
}


func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    log.Println("[handleWebSocket] handleWebSocket start")
    sessionId := r.URL.Query().Get("sessionId")
    if sessionId == "" {
        log.Printf("[handleWebSocket] No session id")
        http.Error(w, "Session ID required", http.StatusUnauthorized)
        return
    }

    session, err := account.ValidateSession(sessionId)
    if err != nil {
        log.Printf("[handleWebSocket] Invalid or expired session: %v", sessionId)
        http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
        return
    }

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



    conn := ConnectionsMap.Add(c, session.AccountID, session)
    log.Printf("[handleWebSocket] ws connection established accountId: %v", session.AccountID)

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

            // fighter, err := findFighterByConn(c)

            // if err != nil {
            //     PopulationMap.Remove(fighter)
            // }

            ConnectionsMap.Remove(c)
            break
        }

        //log.Printf("message: %v", message)

        err = json.Unmarshal(message, &msg)
        if err != nil {
            log.Printf("[handleWebSocket] websocket unmarshal message=%v error=%v", message, err)
            continue
        }

        log.Printf("Type: %v", msg.Type)


        if account.IsSessionExpired(session) {
            sendErrorMsgToConn(conn, "SYSTEM", "Session expired. Please login.")
            ConnectionsMap.Remove(c)
            break
        }

        switch msg.Type {
            
            case "create_fighter":
                type CreateFighterData struct {
                    FighterClass string `json:"fighterClass"`
                    Name string `json:"name"`
                }

                var reqData CreateFighterData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:create_fighter] websocket unmarshal error: %v", err)
                    continue
                }

                _, err = fighters.CreateFighter(conn.GetAccountID(), reqData.Name, reqData.FighterClass)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Failed to create fighter. Error: %v", err))
                }

                serializedFighterList, err := fighters.GetJsonSerializedFighters(conn.GetAccountID())
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }
                respondConn(conn, serializedFighterList)
            continue

            case "get_user_fighters":
                log.Printf("[handleWebSocket:get_user_fighters] accountId: ", conn.GetAccountID())
                serializedFighterList, err := fighters.GetJsonSerializedFighters(conn.GetAccountID())
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }
                respondConn(conn, serializedFighterList)
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
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                } else {
                    go updateFighterDirection(fighter, reqData.Direction)
                }
                
            continue

            case "auth":
                //log.Printf("[handleWebSocket] auth: %v", msg.Data)
                type AuthData struct {
                    PlayerID     int  `json:"playerID"`
                }

                var reqData AuthData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket] websocket unmarshal error: %v", err)
                    continue
                }

                //log.Printf("[handleWebSocket] reqData: %v", reqData)
                fighter, err := authFighter(reqData.PlayerID);
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Auth failed. Error: %v", err))
                }

                conn = ConnectionsMap.AddWithValues(c, fighter)


                WsSendBackpack(fighter)
            case "submit_attack":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:submit_attack] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }
                ProcessHit(fighter, msg.Data)


            case "get_shop": 
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket: get_shop] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    ShopName string  `json:"shopName"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket: get_shop] websocket unmarshal error: %v", err)
                    continue
                }


                shopObj, err := shop.GetShop(reqData.ShopName)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Shop not found")
                }

                WsSendShop(fighter, shopObj, reqData.ShopName)

            case "buy_item": 
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket: buy_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    ShopName string  `json:"shopName"`
                    ItemHash string `json:"itemHash"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket: buy_item] websocket unmarshal error: %v", err)
                    continue
                }

                err = shop.BuyItem(fighter, reqData.ShopName, reqData.ItemHash)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                assignConsumables(fighter)
                WsSendBackpack(fighter)

            case "consumable_bind": 
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket: consumable_bind] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    Binding string  `json:"binding"`
                    Key string `json:"key"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket: consumable_bind] websocket unmarshal error: %v", err)
                    continue
                }

                err = fighter.ConsumableBind(reqData.Binding, reqData.Key)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                assignConsumables(fighter)
                WsSendBackpack(fighter)
                pingFighter(fighter)


            case "consume_backpack_item":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket: consume_backpack_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    ItemHash string  `json:"itemHash"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket: consume_backpack_item] websocket unmarshal error: %v", err)
                    continue
                }

                err = fighter.GetBackpack().Consume(reqData.ItemHash)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                assignConsumables(fighter)
                WsSendBackpack(fighter)
                pingFighter(fighter)


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
                    log.Printf("[handleWebSocket:pickup_dropped_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                err = drop.PickItem(fighter, reqData.ItemHash)
                if err != nil {
                    log.Printf("[handleWebSocket:pickup_dropped_item] error picking item: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Item not found %v", err))
                    continue
                }

                assignConsumables(fighter)
                broadcastDropMessage()
                WsSendBackpack(fighter)


            

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
                    log.Printf("[handleWebSocket:update_backpack_item_position] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }
                fighter.GetBackpack().UpdateInventoryPosition(reqData.ItemHash, reqData.Position)
                WsSendBackpack(fighter)
        
            case "get_vault":
                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:get_vault] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                WsSendVault(fighter)



            case "update_vault_item_position":
                log.Printf("[handleWebSocket:update_vault_item_position] msg.Data: %v", msg.Data)
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:update_vault_item_position] websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:update_vault_item_position] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }
                fighter.GetVault().UpdateInventoryPosition(reqData.ItemHash, reqData.Position)
                WsSendVault(fighter)




            case "drop_backpack_item":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"coordinates"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:drop_backpack_item] websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:drop_backpack_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                drop.DropItem(fighter.GetBackpack(), fighter, reqData.ItemHash, reqData.Position)
                assignConsumables(fighter)
                WsSendBackpack(fighter)
                broadcastDropMessage()


            case "drop_vault_item":

                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"coordinates"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:drop_vault_item] websocket unmarshal error: %v", err)
                    continue
                }

                log.Printf("[drop_vault_item] reqData=%v", reqData)

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:drop_vault_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                drop.DropItem(fighter.GetVault(), fighter, reqData.ItemHash, reqData.Position)
                WsSendVault(fighter)
                broadcastDropMessage()


            case "drop_equipped_item":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"coordinates"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:drop_equipped_item] websocket unmarshal error: %v", err)
                    continue
                }

                log.Printf("[drop_equipped_item] reqData=%v", reqData)

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:drop_equipped_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                drop.DropEquippedItem(fighter, reqData.ItemHash, reqData.Position)
                WsSendBackpack(fighter)
                broadcastDropMessage()

            case "equip_backpack_item":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Slot int `json:"slot"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:equip_backpack_item]  websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:equip_backpack_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                EquipBackpackItem(fighter, reqData.ItemHash, reqData.Slot)
                WsSendBackpack(fighter)

            case "unequip_backpack_item":
                log.Printf("[handleWebSocket:unequip_backpack_item]  msg.Data: %v", msg.Data)
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:unequip_backpack_item]  websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:unequip_backpack_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                UnequipBackpackItem(fighter, reqData.ItemHash, reqData.Position)
                WsSendBackpack(fighter)

            case "equip_vault_item":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Slot int `json:"slot"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:equip_backpack_item]  websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:equip_backpack_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                EquipVaultItem(fighter, reqData.ItemHash, reqData.Slot)
                WsSendBackpack(fighter)
                WsSendVault(fighter)


            case "unequip_vault_item":
                log.Printf("[handleWebSocket:unequip_vault_item]  msg.Data: %v", msg.Data)
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:unequip_vault_item]  websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:unequip_vault_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                UnequipVaultItem(fighter, reqData.ItemHash, reqData.Position)
                WsSendVault(fighter)
                WsSendBackpack(fighter)
                


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
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                handleCommand(fighter, reqData.Text)

            
            case "move_fighter":
                var reqData maps.Coordinate
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:move_fighter] websocket unmarshal error: %v", err)
                    continue
                }
                
                fighter, err := findFighterByConn(c)

                if err != nil {
                    log.Printf("[handleWebSocket:move_fighter] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }
                log.Printf("[handleWebSocket] move_fighter: %v fighter=%v", reqData, fighter)
                moveFighter(fighter, reqData)
                WsSendFighter(fighter)
                pingFighter(fighter)



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
                    log.Printf("[handleWebSocket:drop_backpack_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                tokenAtts := fighter.GetBackpack().FindByHash(reqData.ItemHash)
                fighter.GetBackpack().RemoveItemByHash(reqData.ItemHash);
                assignConsumables(fighter)
                WsSendBackpack(fighter)

                fighter.GetVault().AddItemToPosition(tokenAtts.Attributes, tokenAtts.Qty, reqData.ItemHash, int(reqData.Position.X), int(reqData.Position.Y));
                WsSendVault(fighter)


            case "move_gold_from_backpack_to_vault":
                type ReqData struct {
                    Amount  int `json:"amount"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:move_gold_from_backpack_to_vault] websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:move_gold_from_backpack_to_vault] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                backpack := fighter.GetBackpack()

                availableGold := backpack.GetGold()

                if reqData.Amount > availableGold {
                    sendErrorMsgToConn(conn, "SYSTEM", "Not enough gold")
                    continue
                }

                backpack.SetGold(availableGold - reqData.Amount)
                WsSendBackpack(fighter)

                vault := fighter.GetVault()

                vault.SetGold(vault.GetGold() + reqData.Amount)
                WsSendVault(fighter)


            case "move_gold_from_vault_to_backpack":
                type ReqData struct {
                    Amount  int `json:"amount"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:move_gold_from_backpack_to_vault] websocket unmarshal error: %v", err)
                    continue
                }

                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:move_gold_from_backpack_to_vault] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                vault := fighter.GetVault()

                availableGold := vault.GetGold()

                if reqData.Amount > availableGold {
                    sendErrorMsgToConn(conn, "SYSTEM", "Not enough gold")
                    continue
                }

                vault.SetGold(availableGold - reqData.Amount)
                WsSendVault(fighter)
                

                backpack := fighter.GetBackpack()

                backpack.SetGold(backpack.GetGold() + reqData.Amount)
                WsSendBackpack(fighter)
                


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
                    log.Printf("[handleWebSocket:drop_backpack_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }
                tokenAtts := fighter.GetVault().FindByHash(reqData.ItemHash)
                fighter.GetVault().RemoveItemByHash(reqData.ItemHash)
                WsSendVault(fighter)
                

                fighter.GetBackpack().AddItemToPosition(tokenAtts.Attributes, tokenAtts.Qty, reqData.ItemHash, int(reqData.Position.X), int(reqData.Position.Y));
                assignConsumables(fighter)
                WsSendBackpack(fighter)


            case "skill_bind":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:skill_bind] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    Skill   int `json:"skill"`
                    Key     int `json:"key"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:skill_bind]  websocket unmarshal error: %v", err)
                    continue
                }
        

                err = fighter.BindSkill(reqData.Skill, reqData.Key)
                if err != nil {
                    log.Printf("[handleWebSocket:skill_bind] failed to add stat: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                WsSendFighter(fighter)
                pingFighter(fighter)


            case "stats_add_strength":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_strength] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    Amount  int `json:"amount"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_strength]  websocket unmarshal error: %v", err)
                    continue
                }
        

                err = fighter.AddStrength(reqData.Amount)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_strength] failed to add stat: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                WsSendFighter(fighter)
                pingFighter(fighter)

            case "stats_add_agility":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_agility] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    Amount  int `json:"amount"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_agility]  websocket unmarshal error: %v", err)
                    continue
                }
        

                err = fighter.AddAgility(reqData.Amount)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_agility] failed to add stat: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                WsSendFighter(fighter)
                pingFighter(fighter)


            case "stats_add_energy":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_energy] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    Amount  int `json:"amount"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_energy]  websocket unmarshal error: %v", err)
                    continue
                }
        

                err = fighter.AddEnergy(reqData.Amount)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_energy] failed to add stat: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                WsSendFighter(fighter)
                pingFighter(fighter)


            case "stats_add_vitality":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_vitality] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    Amount  int `json:"amount"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_vitality]  websocket unmarshal error: %v", err)
                    continue
                }
        

                err = fighter.AddVitality(reqData.Amount)
                if err != nil {
                    log.Printf("[handleWebSocket:stats_add_vitality] failed to add stat: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                WsSendFighter(fighter)
                pingFighter(fighter)


            case "trade_initiate":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_initiate] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    PlayerId  string `json:"player_id"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_initiate]  websocket unmarshal error: %v", err)
                    continue
                }

                counterparty := findNearbyFighterById(fighter, reqData.PlayerId, false)
                if counterparty == nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found or not in range")
                    continue
                }

                _, err = trade.Initiate(fighter, counterparty)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                WsSendTrade(fighter)
                WsSendTrade(counterparty)

            case "trade_set_gold":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_set_gold] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    Amount  int `json:"amount"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_set_gold]  websocket unmarshal error: %v", err)
                    continue
                }
        

                err = trade.SetGold(fighter, reqData.Amount)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                fighterTrade := trade.TradesMap.FindByFighter(fighter)

                if fighterTrade == nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Trade not found")
                    continue
                }


                WsSendTrade(fighterTrade.GetFighter1())
                WsSendTrade(fighterTrade.GetFighter2())


            case "trade_add_item":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_add_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    ItemHash  string `json:"item_hash"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_add_item]  websocket unmarshal error: %v", err)
                    continue
                }

                fighterTrade := trade.TradesMap.FindByFighter(fighter)
                if fighterTrade == nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Trade not found")
                    continue
                }        

                err = fighterTrade.AddItem(fighter, reqData.ItemHash)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                


                fighter1 := fighterTrade.GetFighter1()
                fighter2 := fighterTrade.GetFighter2()
                assignConsumables(fighter)
                WsSendTrade(fighter1)
                WsSendTrade(fighter2)

            case "trade_add_item_to_position":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_add_item_to_position] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    ItemHash  string `json:"item_hash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_add_item_to_position]  websocket unmarshal error: %v", err)
                    continue
                }

                fighterTrade := trade.TradesMap.FindByFighter(fighter)
                if fighterTrade == nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Trade not found")
                    continue
                }        

                err = fighterTrade.AddItemToPosition(fighter, reqData.ItemHash, reqData.Position)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                


                fighter1 := fighterTrade.GetFighter1()
                fighter2 := fighterTrade.GetFighter2()
                assignConsumables(fighter)
                WsSendTrade(fighter1)
                WsSendTrade(fighter2)

            case "trade_move_item":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_move_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    ItemHash  string `json:"item_hash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_move_item]  websocket unmarshal error: %v", err)
                    continue
                }

                fighterTrade := trade.TradesMap.FindByFighter(fighter)
                if fighterTrade == nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Trade not found")
                    continue
                }
        

                err = fighterTrade.MoveItem(fighter, reqData.ItemHash, reqData.Position)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }                


                fighter1 := fighterTrade.GetFighter1()
                fighter2 := fighterTrade.GetFighter2()

                WsSendTrade(fighter1)
                WsSendTrade(fighter2)


            case "trade_remove_item":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_remove_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    ItemHash  string `json:"item_hash"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_remove_item]  websocket unmarshal error: %v", err)
                    continue
                }
        

                err = trade.RemoveItem(fighter, reqData.ItemHash)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                fighterTrade := trade.TradesMap.FindByFighter(fighter)

                if fighterTrade == nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Trade not found")
                    continue
                }


                fighter1 := fighterTrade.GetFighter1()
                fighter2 := fighterTrade.GetFighter2()
                assignConsumables(fighter)
                WsSendTrade(fighter1)
                WsSendTrade(fighter2)


            case "trade_remove_item_to_backpack":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_remove_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    ItemHash  string `json:"item_hash"`
                    Position maps.Coordinate `json:"position"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_remove_item]  websocket unmarshal error: %v", err)
                    continue
                }
        

                err = trade.RemoveItemToBackpack(fighter, reqData.ItemHash, reqData.Position)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                fighterTrade := trade.TradesMap.FindByFighter(fighter)

                if fighterTrade == nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Trade not found")
                    continue
                }


                fighter1 := fighterTrade.GetFighter1()
                fighter2 := fighterTrade.GetFighter2()
                assignConsumables(fighter)
                WsSendTrade(fighter1)
                WsSendTrade(fighter2)


            case "trade_remove_item_to_equipment":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_remove_item_to_equipment] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                type ReqData struct {
                    ItemHash  string `json:"item_hash"`
                    Slot int `json:"slot"`
                }

                var reqData ReqData
                err = json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_remove_item_to_equipment]  websocket unmarshal error: %v", err)
                    continue
                }
        

                err = trade.RemoveItemToEquipment(fighter, reqData.ItemHash, reqData.Slot)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }

                fighterTrade := trade.TradesMap.FindByFighter(fighter)

                if fighterTrade == nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Trade not found")
                    continue
                }


                fighter1 := fighterTrade.GetFighter1()
                fighter2 := fighterTrade.GetFighter2()
                assignConsumables(fighter)
                WsSendTrade(fighter1)
                WsSendTrade(fighter2)


            case "trade_approve":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_remove_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                fighterTrade := trade.TradesMap.FindByFighter(fighter)
                if fighterTrade == nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Trade not found")
                    continue
                }

                fighter1 := fighterTrade.GetFighter1()
                fighter2 := fighterTrade.GetFighter2()

                err = trade.Approve(fighter)
                if err != nil {
                    sendErrorMsgToConn(conn, "SYSTEM", fmt.Sprintf("Error: %v", err))
                    continue
                }                

                WsSendTrade(fighter1)
                WsSendTrade(fighter2)


            case "trade_cancel":
                fighter, err := findFighterByConn(c)
                if err != nil {
                    log.Printf("[handleWebSocket:trade_remove_item] fighter not found: %v", err)
                    sendErrorMsgToConn(conn, "SYSTEM", "Fighter not found")
                    continue
                }

                fighterTrade := trade.TradesMap.FindByFighter(fighter)

                if fighterTrade == nil {
                    sendErrorMsgToConn(conn, "SYSTEM", "Trade not found")
                    continue
                }


                fighterTrade.Cancel()

                fighter1 := fighterTrade.GetFighter1()
                fighter2 := fighterTrade.GetFighter2()

                WsSendTrade(fighter1)
                WsSendTrade(fighter2)

                
            default:
                log.Printf("[handleWebSocket] unknown message type: %s", msg.Type)
        }

        // data := decodeJson(message);
        // log.Printf("[handleWebSocket] Received message: %s\n", data["action"].(string))
    }
}