package main

import (
    "log"
    "encoding/json"

    "github.com/gorilla/websocket"
    "net/http"
    "github.com/ethereum/go-ethereum/common"

    "math/big"
  
)

type WsMessage struct {
    Type string  `json:"type"`
    Data Fighter `json:"data"`
}

func pingFighter(fighter *Fighter) {
    //log.Printf("[pingFighter] fighter: %v", fighter)
    type jsonResponse struct {
        Action          string      `json:"action"`
        Fighter         *Fighter    `json:"fighter"`
    }

    jsonResp := jsonResponse{
        Action: "ping",
        Fighter: fighter,
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
        if !fighter.IsNpc && fighter.Location == locationHash && !fighter.IsClosed {
            fighter.ConnMutex.Lock()

            err := fighter.Conn.WriteMessage(websocket.TextMessage, messageJSON)
            if err != nil {
                log.Printf("[broadcastWsMessage] Error broadcasting to %s: %v", fighter.ID, err)
            }
            fighter.ConnMutex.Unlock()
        }
    }
}

func respondFighter(fighter *Fighter, response json.RawMessage) {
    fighter.ConnMutex.Lock()
    defer fighter.ConnMutex.Unlock()

    if fighter.IsClosed {
        log.Println("[respondFighter] Connection is already closed.")
        return
    }

    err := fighter.Conn.WriteMessage(websocket.TextMessage, response)

    if err != nil {
        log.Println("[respondFighter] Error: ", err)
        fighter.IsClosed = true
        return
    }
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
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

    for {
        _, message, err := conn.ReadMessage()


        if err != nil {
            if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
                log.Printf("[handleWebSocket] WebSocket closed err=%v message=%v", err, message)
            } else {
                log.Printf("[handleWebSocket] Failed to read message from WebSocket err=%v message=%v", err, message)
            }
            fighter := findFighterByConn(conn)
            if fighter != nil {
                fighter.ConnMutex.Lock()
                fighter.IsClosed = true
                fighter.ConnMutex.Unlock()
            }
            break
        }

        //log.Printf("message: ", message)

        err = json.Unmarshal(message, &msg)
        if err != nil {
            log.Printf("[handleWebSocket] websocket unmarshal message=%v error=%v", message, err)
            continue
        }

        // log.Printf("Type: ", msg.Type)
        // log.Printf("Data: ", msg.Data)
        switch msg.Type {
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
                    return
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
                    log.Printf("[handleWebSocket] websocket unmarshal error: %v", err)
                    return
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
                    log.Printf("[handleWebSocket] websocket unmarshal error: %v", err)
                    return
                }

                PickupDroppedItem(conn, common.HexToHash(reqData.ItemHash))
            continue

            case "move_fighter":
                var reqData Coordinate
                err := json.Unmarshal(msg.Data, &reqData)
                if err != nil {
                    log.Printf("[handleWebSocket] websocket unmarshal error: %v", err)
                    return
                }

                moveFighter(conn, reqData)
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

        // Handle the message here and send a response back to the client
        // response := "Hello, client!"
        // conn.WriteMessage(websocket.TextMessage, []byte(response))
    }
}