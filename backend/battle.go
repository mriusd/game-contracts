package main

import (
	"encoding/json"
	"log"
	"time"

    "github.com/mriusd/game-contracts/battle"
    "github.com/mriusd/game-contracts/maps"
    "github.com/mriusd/game-contracts/skill"
    "github.com/mriusd/game-contracts/fighters"
    "github.com/mriusd/game-contracts/drop"
)

type Hit struct {
    OpponentID  string  `json:"opponentID"`
    PlayerID    string  `json:"playerID"`    
    Skill       int   `json:"skill"`
    Direction   maps.Direction   `json:"direction"`
}







var ExcellentDamageBonus = float64(0.15)
var CriticalDamageBonus = float64(0.05)



func ProcessHit(playerFighter *fighters.Fighter, data json.RawMessage) {
    log.Printf("[ProcessHit]  playerFighter=%v data=%v", playerFighter, data)
    var hitData Hit
    err := json.Unmarshal(data, &hitData)
    if err != nil {
        log.Printf("[ProcessHit] websocket unmarshal error: %v", err)
        return
    }

    if hitData.PlayerID == hitData.OpponentID { return }

    //playerFighter := fighters.FightersMap.Find(hitData.PlayerID)    
    //playerFighter := conn.gFighter()    
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
        log.Printf("[ProcessHit] failed broadcasting skill: %v", err)
        return
    }

    broadcastWsMessage(playerFighter.GetLocation(), skilresp)


    for _, opponentFighter := range targets {
        log.Printf("[ProcessHit] opponentFighter: %v", opponentFighter)
        var damage float64;
        var oppNewHealth int;

        npcHealth := getNpcHealth(opponentFighter)

        log.Printf("[ProcessHit] npcHealth: %v", npcHealth)

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
       	
        
    	oppNewHealth = int(max(0, npcHealth - int(damage)));    	

        
       	if opponentFighter.GetIsNpc() && oppNewHealth == 0 {
            opponentFighter.SetIsDead(true)
       	}

       	opponentFighter.SetLastDmgTimestamp(int(time.Now().UnixNano()) / int(time.Millisecond))
        opponentFighter.SetHealthAfterLastDmg(oppNewHealth)
       	opponentFighter.SetCurrentHealth(oppNewHealth)

        if damage > 0 {
            addDamageToFighter(opponentFighter.GetID(), playerFighter.GetID(), int(damage))
        }	

       	if (oppNewHealth == 0) {
       		ProcessKill(opponentFighter)
       	}

        
       	
       	type jsonResponse struct {
    		Action string `json:"action"`
        	Damage int `json:"damage"`
            Type battle.DamageType `json:"type"`
        	Opponent string `json:"opponent"`
        	Player string `json:"player"`
        	OpponentNewHealth int `json:"opponentHealth"`
        	LastDmgTimestamp int `json:"lastDmgTimestamp"`
        	HealthAfterLastDmg int `json:"healthAfterLastDmg"`
            PlayerFighter *fighters.Fighter `json:"playerFighter"`
            OpponentFighter *fighters.Fighter `json:"opponentFighter"`
            Skill skill.Skill `json:"skill"`
        }

        jsonResp := jsonResponse{
        	Action: "damage_dealt",
        	Damage: int(damage),
            Type: damageType,
    		Opponent: opponentFighter.GetID(),
    		Player: hitData.PlayerID,
    		OpponentNewHealth: oppNewHealth,
    		LastDmgTimestamp: int(time.Now().UnixNano()) / int(time.Millisecond),
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


func ProcessKill(opponent *fighters.Fighter) {
    log.Printf("[ProcessKill] opponent=%v", opponent)

    coords := opponent.Coordinates

    type DamageTuple struct {
        FighterId        string
        Damage           int
    }

    damageDealt := opponent.GetDamageReceived()
    //battleNonce := big.NewInt(time.Now().UnixNano() / int(time.Millisecond))

    log.Printf("[ProcessKill] opponent=%v damageDealt=%v", opponent.GetTokenID(), damageDealt)
    if len(damageDealt) == 0 {
        return;
    }

    damageDealtTuples := make([]DamageTuple, len(damageDealt))
    for i, d := range damageDealt {
        damageDealtTuples[i] = DamageTuple{
            FighterId:        d.FighterId,
            Damage:           d.Damage,
        }
    }


    log.Printf("[ProcessKill] damageDealt 2 %v", damageDealt)
    //killer := fighters.FightersMap.Find( strconv.FormatInt(damageDealt[0].FighterId.int(), 10) )
    hunter := PopulationMap.Find("lorencia", damageDealt[0].FighterId)

    


    // Drop item
    drop.DropNewItem(opponent.GetLevel(), hunter, "lorencia", coords)
    broadcastDropMessage()
}