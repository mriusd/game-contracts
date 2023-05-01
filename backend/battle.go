package main

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"log"
	"time"
    "fmt"  
    "math/big"
)

type Damage struct {
    FighterId        *big.Int
    Damage           *big.Int
}


type RecordHitMsg struct {
    OpponentID  string  `json:"opponentID`
    PlayerID    string  `json:"playerID"`    
    Skill       int64   `json:"skill"`
    Direction   Direction   `json:"direction"`
}


func ProcessHit(conn *websocket.Conn, data json.RawMessage) {
    var hitData RecordHitMsg
    err := json.Unmarshal(data, &hitData)
    if err != nil {
        log.Printf("websocket unmarshal error: %v", err)
        return
    }

    playerFighter := getFighterSafely(hitData.PlayerID)
    opponentFighter := getFighterSafely(hitData.OpponentID)

    if playerFighter == nil || opponentFighter == nil {
        log.Printf("[ProcessHit] Error: Player or opponent fighter not found, player=%v, opponent=%v", playerFighter, opponentFighter)
        return
    }


    playerId := playerFighter.TokenID
    opponentId := opponentFighter.TokenID
    

    stats1, err := getFighterAttributes(playerFighter);
    stats2, err := getFighterAttributes(opponentFighter);

    //def1 := stats1.Agility.Int64()/4;
    def2 := stats2.Agility.Int64()/4;

    dmg1 := randomValueWithinRange((stats1.Strength.Int64()/4 + stats1.Energy.Int64()/4), 0.25);
    //dmg2 := randomValueWithinRange((stats2.Strength.Int64()/4 + stats2.Energy.Int64()/4), 0.25);


    //items1 := getEquippedItems(stats1);
    items2 := getEquippedItems(stats2);

    //itemDefence1 := getTotalItemsDefence(items1)
    itemDefence2 := getTotalItemsDefence(items2)

   	var damage float64;
   	var oppNewHealth int64;
    npcHealth := getNpcHealth(opponentFighter)

   	// Update battle 
	damage = float64(min(npcHealth, max(0, dmg1 - def2 - itemDefence2)));
	oppNewHealth = max(0, npcHealth - int64(damage));    	


   	if (opponentFighter.IsNpc) {
   		if opponentFighter.IsDead {
			now := time.Now().UnixNano() / 1e6
			elapsedTimeMs := now - opponentFighter.LastDmgTimestamp

			if elapsedTimeMs >= 5000 {
				fmt.Println("[ProcessHit] At least 5 seconds have passed since TimeOfDeath.")
				opponentFighter.IsDead = false;
                opponentFighter.DamageReceived = []Damage{};
				opponentFighter.HealthAfterLastDmg = opponentFighter.MaxHealth;
			} else {
				log.Printf("[ProcessHit] NPC Dead playerId=", playerId, "opponentId=", opponentId)
				return
			}
   			
   		} else if oppNewHealth == 0 {
            opponentFighter.IsDead = true;
   			
   		}
   	}
   	opponentFighter.LastDmgTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
   	opponentFighter.HealthAfterLastDmg = oppNewHealth

    if damage > 0 {
        addDamageToFighter(hitData.OpponentID, big.NewInt(playerId), big.NewInt(int64(damage)))
    }	

   	if (oppNewHealth == 0) {
   		recordBattleOnChain(opponentFighter)
   	}

    
   	
   	type jsonResponse struct {
		Action string `json:"action"`
    	Damage int64 `json:"damage"`
    	Opponent string `json:"opponent"`
    	Player string `json:"player"`
    	OpponentNewHealth int64 `json:"opponentHealth"`
    	LastDmgTimestamp int64 `json:"lastDmgTimestamp"`
    	HealthAfterLastDmg int64 `json:"healthAfterLastDmg"`
    }

    jsonResp := jsonResponse{
    	Action: "damage_dealt",
    	Damage: int64(damage),
		Opponent: hitData.OpponentID,
		Player: hitData.PlayerID,
		OpponentNewHealth: oppNewHealth,
		LastDmgTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
    }

    // Convert the struct to JSON
    response, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[ProcessHit] error: ", err)
        return
    }

    broadcastWsMessage(playerFighter.Location, response)


    log.Println("[ProcessHit] damage=", damage, "opponentId=", opponentId, "playerId=", playerId);
}