// broadcast.go

package main 

import (
	"encoding/json"
    "github.com/gorilla/websocket"
    "log"
    "fmt"

    "github.com/mriusd/game-contracts/maps"
    "github.com/mriusd/game-contracts/items"
    "github.com/mriusd/game-contracts/fighters"
    "github.com/mriusd/game-contracts/drop"
)

func pingFighter(fighter *fighters.Fighter) error {
    //log.Printf("[pingFighter] fighter: %v", fighter)
    type jsonResponse struct {
        Action          string      `json:"action"`
        Fighter         *fighters.Fighter    `json:"fighter"`
        MapObjects      []maps.MapObject `json:"mapObjects"`
        Npcs            []*fighters.Fighter  `json:"npcs"`
        Players         []*fighters.Fighter  `json:"players"`
    }

    mapObjects := maps.GetMapObjectsInRadius("lorencia", float64(20), float64(fighter.Coordinates.X), float64(fighter.Coordinates.Y))

    jsonResp := jsonResponse{
        Action: "ping",
        Fighter: fighter,
        MapObjects: mapObjects,
        Npcs: findNearbyFighters(fighter.GetLocation(), fighter.GetCoordinates(), 20, true),
        Players: findNearbyFighters(fighter.GetLocation(), fighter.GetCoordinates(), 20, false), 
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[pingFighter] %v %v", fighter, err)
    }

    return respondFighter(fighter, messageJSON)
}

func sendErrorMessage(fighter *fighters.Fighter, msg string) {
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
        DroppedItems map[string]drop.ItemDroppedEvent `json:"droppedItems"`
    }

    jsonResp := jsonResponse{
        Action: "dropped_items",
        DroppedItems: drop.DroppedItems.GetMap(),
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[broadcastDropMessage] Error=%v", err)
    }

    broadcastWsMessage("lorencia", messageJSON)
}


func broadcastPickupMessage(fighter *fighters.Fighter, item *items.TokenAttributes, qty int) {
    //log.Printf("[broadcastPickupMessage] item: %v fighter: %v", item, fighter)
    type jsonResponse struct {
        Action      string          `json:"action"`
        Item        *items.TokenAttributes  `json:"item"`
        Fighter     *fighters.Fighter        `json:"fighter"`
        Qty         int           `json:"qty"`
    }

    jsonResp := jsonResponse{
        Action: "item_picked",
        Item: item,
        Fighter: fighter,
        Qty: qty,
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

func broadcastNpcMove(npc *fighters.Fighter, coords maps.Coordinate) {
    //log.Printf("[broadcastNpcMove] npc=%v coords=%v", npc, coords)
    type jsonResponse struct {
        Action string `json:"action"`
        Npc *fighters.Fighter `json:"npc"`
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
    for _, fighter := range PopulationMap.GetTownMap(locationHash) {
        if !fighter.GetIsNpc() {

            conn, connection := findConnectionByFighter(fighter)          

            if conn != nil {
                connection.Lock()
                err := conn.WriteMessage(websocket.TextMessage, messageJSON)
                connection.Unlock()
                if err != nil {
                    log.Printf("[broadcastWsMessage] Error broadcasting to %s: %v", fighter.GetID(), err)
                    ConnectionsMap.Remove(conn)
                    PopulationMap.Remove(fighter)
                }
            }
            
        }
    }
}

func respondFighter(fighter *fighters.Fighter, response json.RawMessage) error {
    conn, connection := findConnectionByFighter(fighter)

    if conn == nil {
        log.Println("[respondFighter] Connection not found") 
        PopulationMap.Remove(fighter)       
        return fmt.Errorf("Connection not found")
    }

    connection.Lock()
    err := conn.WriteMessage(websocket.TextMessage, response)
    connection.Unlock()

    if err != nil {
        log.Printf("[respondFighter] Error: %v", err)
        PopulationMap.Remove(fighter)
        ConnectionsMap.Remove(conn)
        return fmt.Errorf("respondFighter] Connection Error %v", err)
    }

    return nil
}

func respondConn(conn *Connection, response json.RawMessage) error {
    conn.Lock()
    err := conn.WSConn.WriteMessage(websocket.TextMessage, response)
    conn.Unlock()

    if err != nil {
        log.Println("[respondConn] Error: %v conn=%v", err, conn)
        fighter := conn.GetFighter()       
        log.Println("[respondConn] fighter: %v", fighter)
        if fighter != nil {
            PopulationMap.Remove(fighter)
        }

        ConnectionsMap.Remove(conn.WSConn)

        return fmt.Errorf("Connection lost")
    }

    return nil
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


func sendChatMessageToConn(conn *Connection, author, message, msgType string) {
    messageJSON := prepareChatMessage(author, message, msgType)
    respondConn(conn, messageJSON)
}

func sendErrorMsgToConn(conn *Connection, author, message string) {
    sendChatMessageToConn(conn, author, message, "error")
}

func sendLocalMsgToConn(conn *Connection, author, message string) {
    sendChatMessageToConn(conn, author, message, "local")
}



func sendChatMessageToFighter(fighter *fighters.Fighter, author, message, msgType string) {
    messageJSON := prepareChatMessage(author, message, msgType)
    respondFighter(fighter, messageJSON)
}

func sendErrorMsgToFighter(fighter *fighters.Fighter, author, message string) {
    sendChatMessageToFighter(fighter, author, message, "error")
}

func sendLocalMsgToFighter(fighter *fighters.Fighter, author, message string) {
    sendChatMessageToFighter(fighter, author, message, "local")
}