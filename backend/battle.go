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

type DamageType struct {
    IsCritical          bool `json:"isCritical"`
    IsExcellent         bool `json:"isExcellent"`
    IsDouble            bool `json:"isDouble"`
    IsIgnoreDefence     bool `json:"isIgnoreDefence"`
}

/*
    Damage colors:

    if isIgnoreDefence { light yelow }
    else if isExcellent { light green }
    else if isCritical { light blue }

    if double { display twice damage/2 }

*/

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

var ExcellentDamageBonus = float64(0.15)
var CriticalDamageBonus = float64(0.05)

func getSkillSafely(skill int64) *Skill {
    SkillsMutex.RLock()
    defer SkillsMutex.RUnlock()

    return Skills[skill]
}

func ProcessHit(conn *websocket.Conn, data json.RawMessage) {
    var hitData Hit
    err := json.Unmarshal(data, &hitData)
    if err != nil {
        log.Printf("[ProcessHit] websocket unmarshal error: %v", err)
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
        var damage float64;
        var oppNewHealth int64;

        npcHealth := getNpcHealth(opponentFighter)

        damageType := DamageType{}

        def2 := opponentFighter.Defence

        dmg1 := randomValueWithinRange(playerFighter.Damage, 0.25)

        damage = float64(min(npcHealth, max(0, dmg1 - def2)));

        if /* playerFighter.IgnoreDefRate > 0 &&   */ randomValueWithinRange(100, 1) <= playerFighter.IgnoreDefRate + 20 {
            damageType.IsIgnoreDefence = true
            damage = float64(min(npcHealth, max(0, dmg1)))
        }

        if /*  playerFighter.ExcellentDmgRate > 0 &&  */ randomValueWithinRange(100, 1) <= playerFighter.ExcellentDmgRate + 20 {
            damageType.IsExcellent = true
            damage = damage * (1 + ExcellentDamageBonus)
        }

        if /* playerFighter.CriticalDmgRate > 0 &&  */ randomValueWithinRange(100, 1) <= playerFighter.CriticalDmgRate + 20 {
            damageType.IsCritical = true
            damage = damage * (1 + CriticalDamageBonus)
        }

        if /* playerFighter.DoubleDmgRate > 0 && */ randomValueWithinRange(100, 1) <= playerFighter.DoubleDmgRate + 20 {
            damageType.IsDouble = true
            damage *= 2
        }   	
       	
        

       	// Update battle 
    	oppNewHealth = max(0, npcHealth - int64(damage));    	

        opponentFighter.Mutex.Lock()
       	if opponentFighter.IsNpc && oppNewHealth == 0 {
            opponentFighter.IsDead = true
       	}

       	opponentFighter.LastDmgTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
        opponentFighter.HealthAfterLastDmg = oppNewHealth
       	opponentFighter.CurrentHealth = oppNewHealth
        opponentFighter.Mutex.Unlock()

        if damage > 0 {
            addDamageToFighter(opponentFighter.ID, big.NewInt(playerId), big.NewInt(int64(damage)))
        }	

       	if (oppNewHealth == 0) {
       		recordBattleOnChain(opponentFighter)
       	}

        
       	
       	type jsonResponse struct {
    		Action string `json:"action"`
        	Damage int64 `json:"damage"`
            Type DamageType `json:"type"`
        	Opponent string `json:"opponent"`
        	Player string `json:"player"`
        	OpponentNewHealth int64 `json:"opponentHealth"`
        	LastDmgTimestamp int64 `json:"lastDmgTimestamp"`
        	HealthAfterLastDmg int64 `json:"healthAfterLastDmg"`
            PlayerFighter *Fighter `json:"playerFighter"`
            OpponentFighter *Fighter `json:"opponentFighter"`
            Skill *Skill `json:"skill"`
        }

        jsonResp := jsonResponse{
        	Action: "damage_dealt",
        	Damage: int64(damage),
            Type: damageType,
    		Opponent: opponentFighter.ID,
    		Player: hitData.PlayerID,
    		OpponentNewHealth: oppNewHealth,
    		LastDmgTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
            PlayerFighter: playerFighter,
            OpponentFighter: opponentFighter,
            Skill: Skills[hitData.Skill],
        }

        // Convert the struct to JSON
        response, err := json.Marshal(jsonResp)
        if err != nil {
            log.Print("[ProcessHit] error: ", err)
            return
        }



        broadcastWsMessage(opponentFighter.Location, response)

        


        log.Println("[ProcessHit] damage=", damage, "opponentId=", opponentFighter.ID, "playerId=", playerFighter.ID);
    }
}