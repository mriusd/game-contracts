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

    for _, opponentFighter := range targets {
        opponentId := opponentFighter.TokenID

        def2 := opponentFighter.Defence

        dmg1 := randomValueWithinRange(playerFighter.Damage, 0.25)


       	var damage float64;
       	var oppNewHealth int64;
        npcHealth := getNpcHealth(opponentFighter)

       	// Update battle 
    	damage = float64(min(npcHealth, max(0, dmg1 - def2)));
    	oppNewHealth = max(0, npcHealth - int64(damage));    	

        opponentFighter.ConnMutex.Lock()
       	if opponentFighter.IsNpc && oppNewHealth == 0 {
            opponentFighter.IsDead = true
       	}

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