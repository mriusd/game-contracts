package main

import (
    "log"
    "encoding/json"

    "github.com/gorilla/websocket"
    "net/http"
    "github.com/ethereum/go-ethereum/common"

    "math/big"
    "sync"
    "runtime/debug"
    "fmt"
)

type Connection struct {
    Fighter *Fighter
    OwnerAddress common.Address
    Mutex   sync.RWMutex 
}

var Connections = make(map[*websocket.Conn]*Connection)

var ConnectionsMutex sync.RWMutex

type WsMessage struct {
    Type string  `json:"type"`
    Data Fighter `json:"data"`
}

func getOwnerAddressByConn(conn *websocket.Conn) (common.Address, error) {
    ConnectionsMutex.RLock()
    defer ConnectionsMutex.RUnlock()

    connection, ok := Connections[conn]
    if !ok {
        return common.Address{}, fmt.Errorf("connection not found")
    }

    return connection.OwnerAddress, nil
}


func findConnectionByFighter(fighter *Fighter) (*websocket.Conn, *Connection) {
    ConnectionsMutex.RLock()
    defer ConnectionsMutex.RUnlock()
    
    for conn, connection := range Connections {
        if connection.Fighter == fighter {
            return conn, connection
        }
    }

    return nil, &Connection{}
}


func pingFighter(fighter *Fighter) {
    //log.Printf("[pingFighter] fighter: %v", fighter)
    type jsonResponse struct {
        Action          string      `json:"action"`
        Fighter         *Fighter    `json:"fighter"`
        MapObjects      []MapObject `json:"mapObjects"`
        Npcs            []*Fighter  `json:"npcs"`
        Players         []*Fighter  `json:"players"`
    }

    mapObjects := getMapObjectsInRadius("lorencia", float64(20), float64(fighter.Coordinates.X), float64(fighter.Coordinates.Y))

    jsonResp := jsonResponse{
        Action: "ping",
        Fighter: fighter,
        MapObjects: mapObjects,
        Npcs: findNearbyFighters(fighter.Coordinates, 20, true),
        Players: findNearbyFighters(fighter.Coordinates, 20, false), 
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[pingFighter] %v %v", fighter, err)
    }

    respondFighter(fighter, messageJSON)
}

func sendErrorMessage(fighter *Fighter, msg string) {
    //log.Printf("[pingFighter] fighter: %v", fighter)
    type jsonResponse struct {
        Action          string      `json:"action"`
        Msg         string   `json:"msg"`
    }

    jsonResp := jsonResponse{
        Action: "error_message",
        Msg: msg,
    }
        

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[sendErrorMessage] %v %v %v", fighter, msg, err)
    }

    respondFighter(fighter, messageJSON)
}



func broadcastDropMessage() {
    //log.Printf("[broadcastDropMessage] ")
    type jsonResponse struct {
        Action string `json:"action"`
        DroppedItems map[common.Hash]*ItemDroppedEvent `json:"droppedItems"`
    }

    jsonResp := jsonResponse{
        Action: "dropped_items",
        DroppedItems: DroppedItems,
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[broadcastDropMessage] Error=%v", err)
    }

    broadcastWsMessage("lorencia", messageJSON)
}


func broadcastPickupMessage(fighter *Fighter, item ItemAttributes, qty *big.Int) {
    //log.Printf("[broadcastPickupMessage] item: %v fighter: %v", item, fighter)
    type jsonResponse struct {
        Action      string          `json:"action"`
        Item        ItemAttributes  `json:"item"`
        Fighter     *Fighter        `json:"fighter"`
        Qty         int64           `json:"qty"`
    }

    jsonResp := jsonResponse{
        Action: "item_picked",
        Item: item,
        Fighter: fighter,
        Qty: qty.Int64(),
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[broadcastDropMessage] %v %v", fighter, err)
    }

    broadcastWsMessage("lorencia", messageJSON)
}

func broadcastNpcMove(npc *Fighter, coords Coordinate) {
    //log.Printf("[broadcastNpcMove] npc=%v coords=%v", npc, coords)
    type jsonResponse struct {
        Action string `json:"action"`
        Npc *Fighter `json:"npc"`
    }

    jsonResp := jsonResponse{
        Action: "update_npc",
        Npc: npc,
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[broadcastDropMessage] %v %v %v", npc, coords, err)
    }

    broadcastWsMessage("lorencia", messageJSON)
}

func broadcastWsMessage(locationHash string, messageJSON json.RawMessage) {
    for _, fighter := range Fighters {
        if !fighter.IsNpc && fighter.Location == locationHash {

            conn, connection := findConnectionByFighter(fighter)
            connection.Mutex.Lock()

            if conn != nil {
                err := conn.WriteMessage(websocket.TextMessage, messageJSON)
                if err != nil {
                    log.Printf("[broadcastWsMessage] Error broadcasting to %s: %v", fighter.ID, err)
                }
            }
            
            connection.Mutex.Unlock()
        }
    }
}

func respondFighter(fighter *Fighter, response json.RawMessage) {
    conn, connection := findConnectionByFighter(fighter)

    if conn == nil {
        log.Println("[respondFighter] Connection not found")
        return
    }

    connection.Mutex.Lock()
    defer connection.Mutex.Unlock()

    err := conn.WriteMessage(websocket.TextMessage, response)

    if err != nil {
        log.Println("[respondFighter] Error: ", err)
        return
    }
}

func respondConn(conn *websocket.Conn, response json.RawMessage) {
    Connections[conn].Mutex.Lock()
    defer Connections[conn].Mutex.Unlock()

    err := conn.WriteMessage(websocket.TextMessage, response)

    if err != nil {
        log.Println("[respondConn] Error: ", err)
    }
}


func prepareChatMessage(author, message, msgType string) (response json.RawMessage) {
        type jsonResponse struct {
        Action      string   `json:"action"`
        Author      string   `json:"author"`
        Msg         string   `json:"msg"`
        MsgType        string   `json:"msgType"`
    }

    jsonResp := jsonResponse{
        Action: "chat_message",
        Author: author,
        Msg: message,
        MsgType: msgType,
    }
        

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[prepareChatMessage] Error marshaling JSON: %v %v %v", message, msgType, err)
    }

    return messageJSON
}

func broadcastChatMsg(location, author, message, msgType string) {
    messageJSON := prepareChatMessage(author, message, msgType)
    broadcastWsMessage(location, messageJSON)
}


func sendChatMessageToConn(conn *websocket.Conn, author, message, msgType string) {
    messageJSON := prepareChatMessage(author, message, msgType)
    respondConn(conn, messageJSON)
}

func sendErrorMsgToConn(conn *websocket.Conn, author, message string) {
    sendChatMessageToConn(conn, author, message, "error")
}

func sendLocalMsgToConn(conn *websocket.Conn, author, message string) {
    sendChatMessageToConn(conn, author, message, "local")
}



func sendChatMessageToFighter(fighter *Fighter, author, message, msgType string) {
    messageJSON := prepareChatMessage(author, message, msgType)
    respondFighter(fighter, messageJSON)
}

func sendErrorMsgToFighter(fighter *Fighter, author, message string) {
    sendChatMessageToFighter(fighter, author, message, "error")
}

func sendLocalMsgToFighter(fighter *Fighter, author, message string) {
    sendChatMessageToFighter(fighter, author, message, "local")
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

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("[handleWebSocket] Failed to upgrade to WebSocket:", err)
        return
    }
    defer conn.Close()

    Connections[conn] = &Connection{}

    for {
        // Use defer/recover to catch any panic inside the loop
        defer func() {
            if r := recover(); r != nil {
                log.Printf("[handleWebSocket] Recovered from ", r)
                debug.PrintStack()
            }
        }()

        _, message, err := conn.ReadMessage()


        if err != nil {
            if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
                log.Printf("[handleWebSocket] WebSocket closed err=%v message=%v", err, message)
            } else {
                log.Printf("[handleWebSocket] Failed to read message from WebSocket err=%v message=%v", err, message)
            }
            ConnectionsMutex.Lock()
            delete(Connections, conn)
            ConnectionsMutex.Unlock()
            break
        }

        log.Printf("message: ", message)

        err = json.Unmarshal(message, &msg)
        if err != nil {
            log.Printf("[handleWebSocket] websocket unmarshal message=%v error=%v", message, err)
            continue
        }

        log.Printf("Type: ", msg.Type)
        log.Printf("Data: ", msg.Data)
        switch msg.Type {
            case "create_fighter":
                type CreateFighterData struct {
                    OwnerAddress string `json:"ownerAddress"`
                    FighterClass int64 `json:"fighterClass"`
                    Name string `json:"name"`
                }

                var reqData CreateFighterData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:create_fighter] websocket unmarshal error: %v", err)
                    continue
                }

                go CreateFighter(conn, reqData.OwnerAddress, reqData.Name, uint8(reqData.FighterClass))
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

                Connections[conn].Mutex.Lock()
                Connections[conn].OwnerAddress = common.HexToAddress(reqData.OwnerAddress)
                Connections[conn].Mutex.Unlock()
    
                go getUserFighters(conn)
            continue

            case "update_fighter_direction":
                type FighterDirection struct {
                    Direction Direction `json:"direction"`
                }

                var reqData FighterDirection
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:update_fighter_direction] websocket unmarshal error: %v", err)
                    continue
                }
                fighter := findFighterByConn(conn)
                go updateFighterDirection(fighter, reqData.Direction)
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
                authFighter(conn, reqData.PlayerID, reqData.UserAddress, reqData.LocationHash);
            continue
                
            case "submit_attack":
                ProcessHit(conn, msg.Data)
            continue

            case "get_fighter_items":
                type ItemReqData struct {
                    FighterId int64 
                }

                var reqData ItemReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:get_fighter_items] websocket unmarshal error: %v", err)
                    continue
                }
                go getFighterItems(reqData.FighterId)
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
                fighter := findFighterByConn(conn)
                PickupDroppedItem(fighter, common.HexToHash(reqData.ItemHash))
            continue

            case "move_fighter":

                var reqData Coordinate
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:move_fighter] websocket unmarshal error: %v", err)
                    continue
                }
                
                fighter := findFighterByConn(conn)
                log.Printf("[handleWebSocket] move_fighter: %v fighter=%v", reqData, fighter)
                moveFighter(fighter, reqData)
            continue

            case "update_backpack_item_position":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position Coordinate `json:"position"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:update_backpack_item_position] websocket unmarshal error: %v", err)
                    continue
                }

                fighter := findFighterByConn(conn)
                fighter.Backpack.updateBackpackPosition(fighter, common.HexToHash(reqData.ItemHash), reqData.Position)
            continue

            case "drop_backpack_item":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position Coordinate `json:"position"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:drop_backpack_item] websocket unmarshal error: %v", err)
                    continue
                }

                fighter := findFighterByConn(conn)

                DropBackpackItem(conn, common.HexToHash(reqData.ItemHash), fighter.Coordinates)
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

                fighter := findFighterByConn(conn)

                EquipBackpackItem(fighter, common.HexToHash(reqData.ItemHash), reqData.Slot)
            continue

            case "unequip_backpack_item":
                type ReqData struct {
                    ItemHash  string `json:"itemHash"`
                    Position Coordinate `json:"slot"`
                }

                var reqData ReqData
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket:unequip_backpack_item]  websocket unmarshal error: %v", err)
                    continue
                }

                fighter := findFighterByConn(conn)

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

                fighter := findFighterByConn(conn)

                handleCommand(fighter, reqData.Text)
            continue

        

            // case "getFighter":
            //     getFighter(conn, msg.Data, w, r)
            //     continue
                
            default:
                log.Printf("[handleWebSocket] unknown message type: %s", msg.Type)
            continue
        }

        data := decodeJson(message);
        log.Printf("[handleWebSocket] Received message: %s\n", data["action"].(string))
    }
}