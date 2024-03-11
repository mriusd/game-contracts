package main

import (
	"encoding/json"
	"log"
	"time"
    "math/big"

    "github.com/mriusd/game-contracts/battle"
    "github.com/mriusd/game-contracts/maps"
    "github.com/mriusd/game-contracts/skill"
    "github.com/mriusd/game-contracts/fighters"
)

type Hit struct {
    OpponentID  string  `json:"opponentID"`
    PlayerID    string  `json:"playerID"`    
    Skill       int64   `json:"skill"`
    Direction   maps.Direction   `json:"direction"`
}







var ExcellentDamageBonus = float64(0.15)
var CriticalDamageBonus = float64(0.05)



func ProcessHit(conn *Connection, data json.RawMessage) {
    var hitData Hit
    err := json.Unmarshal(data, &hitData)
    if err != nil {
        log.Printf("[ProcessHit] websocket unmarshal error: %v", err)
        return
    }

    if hitData.PlayerID == hitData.OpponentID { return }

    //playerFighter := fighters.FightersMap.Find(hitData.PlayerID)    
    playerFighter := conn.gFighter()    
    targets := findTargetsByDirection(playerFighter, hitData.Direction, skill.Get(hitData.Skill), hitData.OpponentID)

    if playerFighter == nil {
        log.Printf("[ProcessHit] Error: Player fighter not found, player=%v", hitData.PlayerID)
        return
    }

    type jsonResponse struct {
        Action string `json:"action"`
        Fighter *fighters.Fighter `json:"fighter"`
        Skill skill.Skill `json:"skill"`
    }

    jsonResp := jsonResponse{
        Action: "fire_skill",
        Fighter: playerFighter,
        Skill: skill.Get(hitData.Skill),
    }

    // Convert the struct to JSON
    skilresp, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[ProcessHit] failed broadcasting skill: ", err)
        return
    }

    broadcastWsMessage(playerFighter.GetLocation(), skilresp)


    playerId := playerFighter.GetTokenID()

    for _, opponentFighter := range targets {
        var damage float64;
        var oppNewHealth int64;

        npcHealth := getNpcHealth(opponentFighter)

        damageType := battle.DamageType{}

        def2 := opponentFighter.GetDefence()

        dmg1 := randomValueWithinRange(playerFighter.GetDamage(), 0.25)

        damage = float64(min(npcHealth, max(0, dmg1 - def2)));

        if /* playerFighter.IgnoreDefRate > 0 &&   */ randomValueWithinRange(100, 1) <= playerFighter.GetIgnoreDefRate() + 20 {
            damageType.IsIgnoreDefence = true
            damage = float64(min(npcHealth, max(0, dmg1)))
        }

        if /*  playerFighter.ExcellentDmgRate > 0 &&  */ randomValueWithinRange(100, 1) <= playerFighter.GetExcellentDmgRate() + 20 {
            damageType.IsExcellent = true
            damage = damage * (1 + ExcellentDamageBonus)
        }

        if /* playerFighter.CriticalDmgRate > 0 &&  */ randomValueWithinRange(100, 1) <= playerFighter.GetCriticalDmgRate() + 20 {
            damageType.IsCritical = true
            damage = damage * (1 + CriticalDamageBonus)
        }

        if /* playerFighter.DoubleDmgRate > 0 && */ randomValueWithinRange(100, 1) <= playerFighter.GetDoubleDmgRate() + 20 {
            damageType.IsDouble = true
            damage *= 2
        }   	
       	
        
    	oppNewHealth = max(0, npcHealth - int64(damage));    	

        
       	if opponentFighter.GetIsNpc() && oppNewHealth == 0 {
            opponentFighter.SetIsDead(true)
       	}

       	opponentFighter.SetLastDmgTimestamp(time.Now().UnixNano() / int64(time.Millisecond))
        opponentFighter.SetHealthAfterLastDmg(oppNewHealth)
       	opponentFighter.SetCurrentHealth(oppNewHealth)

        if damage > 0 {
            addDamageToFighter(opponentFighter.GetTokenID(), big.NewInt(playerId), big.NewInt(int64(damage)))
        }	

       	if (oppNewHealth == 0) {
       		RecordKill(opponentFighter)
       	}

        
       	
       	type jsonResponse struct {
    		Action string `json:"action"`
        	Damage int64 `json:"damage"`
            Type battle.DamageType `json:"type"`
        	Opponent string `json:"opponent"`
        	Player string `json:"player"`
        	OpponentNewHealth int64 `json:"opponentHealth"`
        	LastDmgTimestamp int64 `json:"lastDmgTimestamp"`
        	HealthAfterLastDmg int64 `json:"healthAfterLastDmg"`
            PlayerFighter *fighters.Fighter `json:"playerFighter"`
            OpponentFighter *fighters.Fighter `json:"opponentFighter"`
            Skill skill.Skill `json:"skill"`
        }

        jsonResp := jsonResponse{
        	Action: "damage_dealt",
        	Damage: int64(damage),
            Type: damageType,
    		Opponent: opponentFighter.GetID(),
    		Player: hitData.PlayerID,
    		OpponentNewHealth: oppNewHealth,
    		LastDmgTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
            PlayerFighter: playerFighter,
            OpponentFighter: opponentFighter,
            Skill: skill.Get(hitData.Skill),
        }

        // Convert the struct to JSON
        response, err := json.Marshal(jsonResp)
        if err != nil {
            log.Print("[ProcessHit] error: ", err)
            return
        }



        broadcastWsMessage(opponentFighter.GetLocation(), response)

        


        log.Println("[ProcessHit] damage=", damage, "opponentId=", opponentFighter.GetID(), "playerId=", playerFighter.GetID());
    }
}