package main

func broadcastDropMessage(eventData ItemDroppedEvent) {
    log.Printf("[broadcastDropMessage] eventData: %v", eventData)
    type jsonResponse struct {
        Action string `json:"action"`
        EventData ItemDroppedEvent `json:"eventData"`
    }

    jsonResp := jsonResponse{
        Action: "item_dropped",
        EventData: eventData,
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[broadcastDropMessage] %v %v", eventData, err)
    }

    broadcastWsMessage("lorencia", messageJSON)
}

func broadcastPickupMessage(fighter Fighter, item ItemAttributes) {
    log.Printf("[broadcastPickupMessage] item: %v fighter: %v", item, fighter)
    type jsonResponse struct {
        Action string `json:"action"`
        Item ItemAttributes `json:"item"`
        Fighter Fighter `json:"fighter"`
    }

    jsonResp := jsonResponse{
        Action: "item_picked",
        Item: item,
        Fighter: fighter,
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[broadcastDropMessage] %v %v", eventData, err)
    }

    broadcastWsMessage("lorencia", messageJSON)
}

func broadcastWsMessage(locationHash string, messageJSON json.RawMessage) {
    for _, fighter := range Fighters {
        if !fighter.IsNpc && fighter.Location == locationHash {
            fighter.ConnMutex.Lock()

            err := fighter.Conn.WriteMessage(websocket.TextMessage, messageJSON)
            if err != nil {
                log.Printf("[broadcastWsMessage] Error broadcasting to %s: %v", fighter.ID, err)
            }
            fighter.ConnMutex.Unlock()
        }
    }
}