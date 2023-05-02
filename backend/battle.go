package main

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"log"
	"time"
    "math/big"
    "sync"
)

type Damage struct {
    FighterId        *big.Int
    Damage           *big.Int
}


type Hit struct {
    OpponentID  string  `json:"opponentID`
    PlayerID    string  `json:"playerID"`    
    Skill       int64   `json:"skill"`
    Direction   Direction   `json:"direction"`
}

type Skill struct {
    SkillId             int     `json:"skillId"`
    Name                string  `json:"name"`
    ActiveDistance      int     `json:"activeDistance"`
    Multihit            bool    `json:"multihit"`
    AttackSuccessRate   int     `json:"attackSuccessRate"`
    HitAngle            int     `json:"hitAngle"`
    Disabled            bool    `json:"disabled"`
}

var Skills = map[int64]*Skill{
    0: {
        SkillId:           0,
        Name:              "Malee", 
        Multihit:          false,
        ActiveDistance:    1,
        AttackSuccessRate: 100,
        HitAngle:          180,
    },
    1: {
        SkillId:           1,
        Name:              "Slash", 
        Multihit:          false,
        ActiveDistance:    1,
        AttackSuccessRate: 100,
        HitAngle:          180,
    },
    2: {
        SkillId:           2,
        Name:              "Arrow", 
        Multihit:          false,
        ActiveDistance:    5,
        AttackSuccessRate: 100,
        HitAngle:          180,
    },
    3: {
        SkillId:           3,
        Name:              "Tripple Shot", 
        Multihit:          true,
        ActiveDistance:    4,
        AttackSuccessRate: 100,
        HitAngle:          180,
    },
    4: {
        SkillId:           4,
        Name:              "Dark Spirits", 
        Multihit:          true,
        ActiveDistance:    20,
        AttackSuccessRate: 100,
        HitAngle:          360,
    },
}


var SkillsMutex sync.RWMutex

func getSkillSafely(skill int64) *Skill {
    SkillsMutex.RLock()
    defer SkillsMutex.RUnlock()

    return Skills[skill]
}

func ProcessHit(conn *websocket.Conn, data json.RawMessage) {
    var hitData Hit
    err := json.Unmarshal(data, &hitData)
    if err != nil {
        log.Printf("websocket unmarshal error: %v", err)
        return
    }

    if hitData.PlayerID == hitData.OpponentID { return }

    playerFighter := getFighterSafely(hitData.PlayerID)    
    targets := findTargetsByDirection(playerFighter, hitData.Direction, Skills[hitData.Skill], hitData.OpponentID)

    if playerFighter == nil {
        log.Printf("[ProcessHit] Error: Player fighter not found, player=%v", playerFighter)
        return
    }


    playerId := playerFighter.TokenID
    stats1, err := getFighterAttributes(playerFighter)

    for _, opponentFighter := range targets {
        opponentId := opponentFighter.TokenID
        stats2, err := getFighterAttributes(opponentFighter)

        //def1 := stats1.Agility.Int64()/4;
        def2 := stats2.Agility.Int64()/4

        dmg1 := randomValueWithinRange((stats1.Strength.Int64()/4 + stats1.Energy.Int64()/4), 0.25)
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
    			// now := time.Now().UnixNano() / 1e6
    			// elapsedTimeMs := now - opponentFighter.LastDmgTimestamp

    			// if elapsedTimeMs >= 5000 {
    			// 	fmt.Println("[ProcessHit] At least 5 seconds have passed since TimeOfDeath.")
    			// 	opponentFighter.IsDead = false;
                //     opponentFighter.DamageReceived = []Damage{};
    			// 	opponentFighter.HealthAfterLastDmg = opponentFighter.MaxHealth;
    			// } else {
    				//log.Printf("[ProcessHit] NPC Dead playerId=", playerId, "opponentId=", opponentId)
    				continue
    			// }
       			
       		} else if oppNewHealth == 0 {
                opponentFighter.ConnMutex.Lock()
                opponentFighter.IsDead = true
                opponentFighter.ConnMutex.Unlock()
       			
       		}
       	}
        opponentFighter.ConnMutex.Lock()
       	opponentFighter.LastDmgTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
        opponentFighter.HealthAfterLastDmg = oppNewHealth
       	opponentFighter.CurrentHealth = oppNewHealth
        opponentFighter.ConnMutex.Unlock()

        if damage > 0 {
            addDamageToFighter(opponentFighter.ID, big.NewInt(playerId), big.NewInt(int64(damage)))
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
            Fighter *Fighter `json:"fighter"`
        }

        jsonResp := jsonResponse{
        	Action: "damage_dealt",
        	Damage: int64(damage),
    		Opponent: opponentFighter.ID,
    		Player: hitData.PlayerID,
    		OpponentNewHealth: oppNewHealth,
    		LastDmgTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
            Fighter: opponentFighter,
        }

        // Convert the struct to JSON
        response, err := json.Marshal(jsonResp)
        if err != nil {
            log.Print("[ProcessHit] error: ", err)
            return
        }

        broadcastWsMessage(opponentFighter.Location, response)


        log.Println("[ProcessHit] damage=", damage, "opponentId=", opponentId, "playerId=", playerFighter.ID);
    }
}