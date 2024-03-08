package main

import (
	"math/big"

    "github.com/gorilla/websocket"	
    "log"
    "strconv"

    "time"
    "fmt"
    "math/rand"

    "github.com/ethereum/go-ethereum/common"

    "errors"
    "regexp"
    "unicode/utf8"
)

func validateFighterName(name string) error {
    if utf8.RuneCountInString(name) > 13 {
        return errors.New("Name too long")
    }

    // Check if name contains only A-Z, a-z, 0-9
    matched, err := regexp.MatchString(`^[a-zA-Z0-9]*$`, name)
    if err != nil {
        return err
    }

    if !matched {
        return errors.New("Name contains invalid characters")
    }

    return nil
}

func getHealthRegenerationRate(atts *FighterAttributes) (float64) {

    vitality := atts.gVitality().Int64()
    healthRegenBonus := 0

    regenRate := (float64(vitality)/float64(HealthRegenerationDivider) + float64(healthRegenBonus))/5000;
    log.Printf("[getHealthRegenerationRate] vitality=%v regenRate=%v", vitality, regenRate)
    return regenRate
}

func getNpcHealth(fighter *Fighter) int64 {
	if !fighter.IsDead {
		return getHealth(fighter);
	} else {
		now := time.Now().UnixNano() / 1e6
		elapsedTimeMs := now - fighter.gLastDmgTimestamp()

		if elapsedTimeMs >= 5000 {
			log.Printf("[getNpcHealth] At least 5 seconds have passed since TimeOfDeath.")
            fighter.Lock()
			fighter.IsDead = false;
			fighter.HealthAfterLastDmg = fighter.MaxHealth;
            fighter.Unlock()
			return fighter.MaxHealth;
		} else {
			return 0;
		}	
	}
	
}

func getHealth(fighter *Fighter) int64 {
	maxHealth := fighter.gMaxHealth();
	lastDmgTimestamp := fighter.gLastDmgTimestamp();
	healthAfterLastDmg := fighter.gHealthAfterLastDmg();

    healthRegenRate := fighter.gHpRegenerationRate();
    currentTime := time.Now().UnixNano() / int64(time.Millisecond);

    health := float64(healthAfterLastDmg) + (float64((currentTime - lastDmgTimestamp)) * healthRegenRate)

    //log.Printf("[getHealth] currentTime=", currentTime," maxHealth=", maxHealth," lastDmgTimestamp=",lastDmgTimestamp," healthAfterLastDmg=",healthAfterLastDmg," healthRegenRate=", healthRegenRate, " health=", health)

    currHealth := min(maxHealth, int64(health))
    fighter.sCurrentHealth(currHealth)

    return currHealth
}

func initiateFighterRoutine(conn *websocket.Conn, fighter *Fighter) {
    log.Printf("[initiateFighterRoutine] fighter=%v", fighter.ID)
    speed := fighter.MovementSpeed

    msPerHit := 60000 / speed
    delay := time.Duration(msPerHit) * time.Millisecond

    for {
        conn, _ := findConnectionByFighter(fighter)

        if conn == nil {
            log.Printf("[initiateFighterRoutine] Connection closed, stopping the loop for fighter: %v", fighter.ID)
            //removeFighterFromPopulation(fighter)
            PopulationMap.Remove(fighter)
            return;
        }

       

        // Check if LastChatMsg is not empty and if 30 seconds passed from the LastChatMsgTimestamp
        currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)
        if fighter.LastChatMsg != "" && (currentTimeMillis - fighter.LastChatMsgTimestamp > 30 * 1000) {
            // Lock the fighter struct to prevent concurrent write
            fighter.Lock()
            fighter.LastChatMsg = ""
            fighter.Unlock()
        }

        pingFighter(fighter)
        time.Sleep(delay)
    }
}



func authFighter(conn *websocket.Conn, playerId int64, ownerAddess string, locationKey string) {
    log.Printf("[authFighter] playerId=%v ownerAddess=%v locationKey=%v conn=%v", playerId, ownerAddess, locationKey, conn)

    if playerId == 0 {
        log.Printf("[authFighter] Player id cannot be zero")
        sendErrorMsgToConn(conn, "SYSTEM", "Invalid player id")
        return
    }

    strId := strconv.Itoa(int(playerId))
    stats := getFighterStats(playerId)
    atts, _ := getFighterAttributes(playerId)

    location := decodeLocation(locationKey)
    town := location[0]  

    connection := ConnectionsMap.Find(conn)
    if connection != nil {
        //removeFighterFromPopulation(connection.gFighter())
        PopulationMap.Remove(connection.gFighter())
    }
    
    fighter := FightersMap.Find(strId)
    if fighter != nil {
        log.Printf("[authFighter] Fighter already exists, only update the Conn value")

        oldConn, _ := findConnectionByFighter(fighter)
        if oldConn != nil {
            ConnectionsMap.Remove(oldConn)
        }

        // PopulationMutex.Lock()
        // if Population[town] == nil {
        //     Population[town] = make([]*Fighter, 0)
        // }
        // Population[town] = append(Population[town], fighter)
        // PopulationMutex.Unlock() 

        PopulationMap.Add(town, fighter)

        ConnectionsMap.AddWithValues(conn, fighter, common.HexToAddress(ownerAddess))

        go initiateFighterRoutine(conn, fighter)
        getFighterItems(fighter)
    } else {        
        centerCoord := Coordinate{X: 5, Y: 5}
        emptySquares := getEmptySquares(centerCoord, 5, town)

        rand.Seed(time.Now().UnixNano())
        spawnCoord := emptySquares[rand.Intn(len(emptySquares))]


        fighter, err := retrieveFighterFromDB(strId)

        if err == nil {
            fighter.Backpack = NewInventory(8, 8) 
            fighter.Vault = NewInventory(8, 16) 
            fighter.Equipment = make(map[int64]*InventorySlot)

        } else {
            log.Printf("[authFighter] err=%v", err)

            fighter = &Fighter{
                ID: strId,
                TokenID: playerId,
                BirthBlock: atts.BirthBlock.Int64(),
                MaxHealth: stats.MaxHealth.Int64(),
                CurrentHealth: stats.MaxHealth.Int64(),
                Name: atts.Name,
                IsNpc: false,
                CanFight: true,
                LastDmgTimestamp: 0,
                HealthAfterLastDmg: 0,
                OwnerAddress: ownerAddess,
                MovementSpeed: 270,
                Coordinates: spawnCoord,
                Backpack: NewInventory(8, 8),
                Vault: NewInventory(8, 16), 
                Location: town,
                Strength: atts.Strength.Int64(),
                Agility: atts.Agility.Int64(),
                Energy: atts.Energy.Int64(),
                Vitality: atts.Vitality.Int64(),
                HpRegenerationRate: getHealthRegenerationRate(atts),
                Level: stats.Level.Int64(),
                Experience: atts.Experience.Int64(),
                Direction: Direction{Dx: 0, Dy: 1},
                Skills: Skills,
                Equipment: make(map[int64]*InventorySlot),
            }
        }        


        FightersMap.Add(strId, fighter)

        getBackpackFromDB(fighter)
        updateFighterParams(fighter) 

        PopulationMap.Add(town, fighter)    

        ConnectionsMap.AddWithValues(conn, fighter, common.HexToAddress(ownerAddess))        
        
        go initiateFighterRoutine(conn, fighter)
        getFighterItems(fighter)
    }

    //FaucetCredits(conn)

    

    
}


func findFighterByConn(conn *websocket.Conn) (*Fighter, error) {
    connection := ConnectionsMap.Find(conn)

    if connection == nil {
        return nil, fmt.Errorf("connection not found")
    }

    fighter := connection.gFighter()
    if fighter == nil {
        return nil, fmt.Errorf("connection has no fighter")
    }

    return fighter, nil
}


func addDamageToFighter(fighterID string, hitterID *big.Int, damage *big.Int) {
    found := false
    fighter := FightersMap.Find(fighterID);

    log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v", fighterID, hitterID, damage)

    damageReceived := fighter.gDamageReceived()

    // Check if there's already damage from the hitter
    for i, dmg := range damageReceived {
        if dmg.FighterId.Cmp(hitterID) == 0 {
            found = true
            //log.Printf("[addDamageToFighter] Damage found ")
            // Add the new damage to the existing damage
            damageReceived[i].Damage = new(big.Int).Add(damageReceived[i].Damage, damage)
            break
        }
    }

    // If no existing damage from the hitter, create a new Damage object and append it to DamageReceived
    if !found {
        newDamage := Damage{
            FighterId: hitterID,
            Damage:    damage,
        }
        damageReceived = append(damageReceived, newDamage)
    }

    fighter.Lock()
    fighter.DamageReceived = damageReceived
    fighter.Unlock()

    //log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v fighter=%v", fighterID, hitterID, damage, fighter)
}

func updateFighterParams(fighter *Fighter) {

    fighter.RLock()
    equipment := fighter.Equipment
    fighter.RUnlock()


    defence := fighter.Agility/4
    damage := fighter.Strength/4 + fighter.Energy/4

    criticalDmg     := int64(0)
    excellentDmg    := int64(0)
    doubleDmg       := int64(0)
    ignoreDef       := int64(0)

    for _, item := range equipment {
        // Perform your logic with the current item and slot
        defence += item.Attributes.ItemParameters.Defense + item.Attributes.AdditionalDefense.Int64()
        damage += item.Attributes.ItemParameters.MinPhysicalDamage + item.Attributes.AdditionalDamage.Int64()

        if item.Attributes.Luck {
            criticalDmg += 5
        }

        if item.Attributes.ExcellentItemAttributes.ExcellentDamageProbabilityIncrease.Int64() > 0 {
            excellentDmg += 10
        }

        if item.Attributes.ExcellentItemAttributes.DoubleDamageChance.Int64() > 0 {
            doubleDmg += 10
        }

        if item.Attributes.ExcellentItemAttributes.IgnoreOpponentDefenseChance.Int64() > 0 {
            ignoreDef += item.Attributes.ExcellentItemAttributes.IgnoreOpponentDefenseChance.Int64()
        }
    }


    fighter.Lock()
    fighter.Damage = damage
    fighter.Defence = defence

    fighter.CriticalDmgRate = criticalDmg
    fighter.ExcellentDmgRate = excellentDmg
    fighter.DoubleDmgRate = doubleDmg
    fighter.IgnoreDefRate = ignoreDef

    fighter.Unlock()

    //pingFighter(fighter)

    log.Printf("[updateFighterParams] fighter=%v", fighter)
}


func applyConsumable(fighter *Fighter, item *TokenAttributes) {
    switch item.gName() {
        case "Small Healing Potion":
            go graduallyIncreaseHp(fighter, 100, 5)
            break

        case "Small Mana Potion":
            go graduallyIncreaseMana(fighter, 100, 5)
            break

        default:
            log.Printf("[applyConsumable] Unknown consumable=%v", item.gName())
            break
    }
}


func graduallyIncreaseHp(fighter *Fighter, hp int64, chunks int64) {
    // Calculate how much to increase HP by each chunk
    hpIncrease := hp / chunks

    for i := int64(0); i < chunks; i++ {
        // Lock the mutex before updating fighter's HP
        fighter.Lock()
        fighter.CurrentHealth += hpIncrease
        fighter.Unlock()

        // Print HP for debugging purposes, remove in production code
        log.Printf("[graduallyIncreaseHp] HP after chunk %v:%v", i+1, fighter.CurrentHealth)

        // Sleep for one second
        time.Sleep(1 * time.Second)
    }
}

func graduallyIncreaseMana(fighter *Fighter, mana int64, chunks int64) {
    // Calculate how much to increase HP by each chunk
    manaIncrease := mana / chunks

    for i := int64(0); i < chunks; i++ {
        // Lock the mutex before updating fighter's HP
        fighter.Lock()
        fighter.CurrentMana += manaIncrease
        fighter.Unlock()

        // Print HP for debugging purposes, remove in production code
        log.Printf("[graduallyIncreaseMana] MP after chunk %v:%v", i+1, fighter.CurrentMana)

        // Sleep for one second
        time.Sleep(1 * time.Second)
    }
}









