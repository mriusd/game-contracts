// broadcast.go

package main 

import (
	"encoding/json"
    "github.com/gorilla/websocket"
    "math/big"
    "log"
    "github.com/ethereum/go-ethereum/common"
)

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
        Npcs: findNearbyFighters(fighter.gLocation(), fighter.gCoordinates(), 20, true),
        Players: findNearbyFighters(fighter.gLocation(), fighter.gCoordinates(), 20, false), 
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
        DroppedItems map[common.Hash]*ItemDroppedEventGo `json:"droppedItems"`
    }

    jsonResp := jsonResponse{
        Action: "dropped_items",
        DroppedItems: getDroppedItemsInGo(),
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[broadcastDropMessage] Error=%v", err)
    }

    broadcastWsMessage("lorencia", messageJSON)
}


func broadcastPickupMessage(fighter *Fighter, item *TokenAttributes, qty *big.Int) {
    //log.Printf("[broadcastPickupMessage] item: %v fighter: %v", item, fighter)
    type jsonResponse struct {
        Action      string          `json:"action"`
        Item        *TokenAttributes  `json:"item"`
        Fighter     *Fighter        `json:"fighter"`
        Qty         int64           `json:"qty"`
    }

    jsonResp := jsonResponse{
        Action: "item_picked",
        Item: item,
        Fighter: fighter,
        Qty: qty.Int64(),
    }

    fighter.RLock()
    item.RLock()
    messageJSON, err := json.Marshal(jsonResp)
    item.RUnlock()
    fighter.RUnlock()

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
    for _, fighter := range FightersMap.gMap() {
        if !fighter.gIsNpc() && fighter.gLocation() == locationHash {

            conn, connection := findConnectionByFighter(fighter)          

            if conn != nil {
                connection.Lock()
                err := conn.WriteMessage(websocket.TextMessage, messageJSON)
                connection.Unlock()
                if err != nil {
                    log.Printf("[broadcastWsMessage] Error broadcasting to %s: %v", fighter.gID(), err)
                    ConnectionsMap.Remove(conn)
                }
            }
            
        }
    }
}

func respondFighter(fighter *Fighter, response json.RawMessage) {
    conn, connection := findConnectionByFighter(fighter)

    if conn == nil {
        log.Println("[respondFighter] Connection not found")
        return
    }

    connection.Lock()
    defer connection.Unlock()

    err := conn.WriteMessage(websocket.TextMessage, response)

    if err != nil {
        log.Println("[respondFighter] Error: ", err)
        return
    }
}

func respondConn(conn *websocket.Conn, response json.RawMessage) {
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