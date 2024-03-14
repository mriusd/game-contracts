package main

import (
    "context"
    "github.com/gorilla/websocket"	
    "log"
    "time"
    "fmt"
    "math"
    "sort"
    "strconv"

    "encoding/json"

    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"

    "github.com/mriusd/game-contracts/db"    
    "github.com/mriusd/game-contracts/maps" 
    "github.com/mriusd/game-contracts/battle" 
    "github.com/mriusd/game-contracts/items" 
    "github.com/mriusd/game-contracts/inventory" 
    "github.com/mriusd/game-contracts/fighters" 
    "github.com/mriusd/game-contracts/skill" 
)



// func getHealthRegenerationRate(atts *fighters.FighterAttributes) (float64) {
//     vitality := atts.GetVitality()
//     healthRegenBonus := 0

//     regenRate := (float64(vitality)/float64(HealthRegenerationDivider) + float64(healthRegenBonus))/5000;
//     log.Printf("[getHealthRegenerationRate] vitality=%v regenRate=%v", vitality, regenRate)
//     return regenRate
// }

func getNpcHealth(fighter *fighters.Fighter) int {
	if !fighter.IsDead {
		return getHealth(fighter);
	} else {
		now := int(time.Now().UnixNano()) / 1e6
		elapsedTimeMs := now - fighter.GetLastDmgTimestamp()

		if elapsedTimeMs >= 5000 {
            maxHealth := fighter.GetMaxHealth()
			log.Printf("[getNpcHealth] At least 5 seconds have passed since TimeOfDeath.")
			fighter.SetIsDead(false)
			fighter.SetHealthAfterLastDmg(maxHealth)
			return maxHealth;
		} else {
			return 0;
		}	
	}
	
}

func getHealth(fighter *fighters.Fighter) int {
	maxHealth := fighter.GetMaxHealth();
	lastDmgTimestamp := fighter.GetLastDmgTimestamp();
	healthAfterLastDmg := fighter.GetHealthAfterLastDmg();

    healthRegenRate := fighter.GetHpRegenerationRate();
    currentTime := int(time.Now().UnixNano()) / int(time.Millisecond);

    health := float64(healthAfterLastDmg) + (float64((currentTime - lastDmgTimestamp)) * healthRegenRate)

    //log.Printf("[getHealth] currentTime=", currentTime," maxHealth=", maxHealth," lastDmgTimestamp=",lastDmgTimestamp," healthAfterLastDmg=",healthAfterLastDmg," healthRegenRate=", healthRegenRate, " health=", health)

    currHealth := min(maxHealth, int(health))
    fighter.SetCurrentHealth(currHealth)

    return currHealth
}

func initiateFighterRoutine(fighter *fighters.Fighter) {
    log.Printf("[initiateFighterRoutine] fighter=%v", fighter.ID)
    speed := fighter.GetMovementSpeed()

    msPerHit := 60000 / speed
    delay := time.Duration(msPerHit) * time.Millisecond

    for {
        // Check if LastChatMsg is not empty and if 30 seconds passed from the LastChatMsgTimestamp
        currentTimeMillis := int(time.Now().UnixNano()) / int(time.Millisecond)
        if fighter.LastChatMsg != "" && (currentTimeMillis - fighter.GetLastChatMsgTimestamp() > 30 * 1000) {
            // Lock the fighter struct to prevent concurrent write
            fighter.SetLastChatMsg("")
        }

        err := pingFighter(fighter)
        if err != nil {
            log.Printf("[initiateFighterRoutine] Error pinging fighter %v", err)
            
            break
        }

        //WsSendBackpack(fighter)

        time.Sleep(delay)
    }
}



func authFighter(playerId int, ownerAddess string, locationKey string) (*fighters.Fighter, error) {
    log.Printf("[authFighter] playerId=%v ownerAddess=%v locationKey=%v", playerId, ownerAddess, locationKey)

    if playerId == 0 {
        log.Printf("[authFighter] Player id cannot be zero")
        //sendErrorMsgToConn(conn, "SYSTEM", "Invalid player id")
        return nil, fmt.Errorf("Player id cannot be zero")
    }

    fighter, err := fighters.GetFromDB(playerId)
    if err != nil {
        return nil, fmt.Errorf("Failed to load character: %v", err)
    }  

    maxHealth := fighter.CalcMaxHealth()

    log.Printf("[authFighter] fighter=%v", fighter)

    backpack, err := inventory.GetInventoryFromDB(playerId, "backpack")
    if err != nil {
        return nil, fmt.Errorf("Failed to load backpack: %v", err)
    } 

    equipment, err := inventory.GetEquipmentFromDB(playerId)
    if err != nil {
        return nil, fmt.Errorf("Failed to load equipment: %v", err)
    } 

    fighter.ID = strconv.Itoa(fighter.TokenID)
    fighter.MovementSpeed = 270
    fighter.Backpack = backpack
    //fighter.Vault = inventory.GetFromDB(playerId, "vault") 
    fighter.Equipment = equipment
    fighter.HealthAfterLastDmg = maxHealth
    fighter.MaxHealth = maxHealth
    fighter.CurrentHealth = maxHealth
    fighter.Skills = skill.Skills
    fighter.Level = fighter.CalcLevel()

    log.Printf("[authFighter] props set=%v", fighter)
    PopulationMap.Add("lorencia", fighter) 
    log.Printf("[authFighter] added to population=%v", fighter)

    updateFighterParams(fighter) 
    go initiateFighterRoutine(fighter)
    
    //getFighterItems(fighter)
    return fighter, nil
}

func WsSendBackpack(fighter *fighters.Fighter) {
    //log.Printf("[wsSendBackpack] fighter: %v backpack: %v", fighter.GetName(), fighter.GetBackpack())
    type jsonResponse struct {
        Action string `json:"action"`
        Backpack inventory.Inventory `json:"backpack"`
        Equipment map[int]inventory.InventorySlot `json:"equipment"`
    }

    jsonResp := jsonResponse{
        Action: "backpack_update",
        Backpack: *fighter.GetBackpack(),
        Equipment: fighter.GetEquipment().GetMap(),
    }

    response, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[wsSendInventory] error: ", err)
        return
    }
    respondFighter(fighter, response)
}

// func authFighter(conn *websocket.Conn, playerId int, ownerAddess string, locationKey string) {
//     log.Printf("[authFighter] playerId=%v ownerAddess=%v locationKey=%v conn=%v", playerId, ownerAddess, locationKey)

//     if playerId == 0 {
//         log.Printf("[authFighter] Player id cannot be zero")
//         sendErrorMsgToConn(conn, "SYSTEM", "Invalid player id")
//         return
//     }

//     strId := strconv.Itoa(int(playerId))
//     //stats := getFighterStats(playerId)
//     log.Printf("[authFighter] playerId=%v stats=%v", playerId, stats)

//     atts, _ := getFighterAttributes(playerId)

//     location := maps.DecodeLocation(locationKey)
//     town := location[0]  

//     connection := ConnectionsMap.Find(conn)
//     if connection != nil {
//         //removeFighterFromPopulation(connection.gFighter())
//         PopulationMap.Remove(connection.gFighter())
//     }
    
//     fighter := fighters.FightersMap.Find(strId)
//     if fighter != nil {
//         log.Printf("[authFighter] Fighter already exists, only update the Conn value")

//         oldConn, _ := findConnectionByFighter(fighter)
//         if oldConn != nil {
//             ConnectionsMap.Remove(oldConn)
//         }

//         // PopulationMutex.Lock()
//         // if Population[town] == nil {
//         //     Population[town] = make([]*Fighter, 0)
//         // }
//         // Population[town] = append(Population[town], fighter)
//         // PopulationMutex.Unlock() 

//         PopulationMap.Add(town, fighter)

//         ConnectionsMap.AddWithValues(conn, fighter, common.HexToAddress(ownerAddess))

//         go initiateFighterRoutine(conn, fighter)
//         getFighterItems(fighter)
//     } else {        
//         centerCoord := maps.Coordinate{X: 5, Y: 5}
//         emptySquares := getEmptySquares(centerCoord, 5, town)

//         rand.Seed(time.Now().UnixNano())
//         spawnCoord := emptySquares[rand.Intn(len(emptySquares))]


//         fighter, err := retrieveFighterFromDB(strId)

//         if err == nil {
//             fighter.Backpack = inventory.NewInventory(8, 8) 
//             fighter.Vault = inventory.NewInventory(8, 16) 
//             fighter.Equipment = inventory.NewEquipment()

//         // } else {
//         //     log.Printf("[authFighter] err=%v", err)

//         //     fighter = &fighters.Fighter{
//         //         ID: strId,
//         //         TokenID: playerId,
//         //         MaxHealth: fighter.GetMaxHealth(),
//         //         CurrentHealth: stats.MaxHealth.int(),
//         //         Name: atts.Name,
//         //         IsNpc: false,
//         //         CanFight: true,
//         //         LastDmgTimestamp: 0,
//         //         HealthAfterLastDmg: 0,
//         //         OwnerAddress: ownerAddess,
//         //         MovementSpeed: 270,
//         //         Coordinates: spawnCoord,
//         //         Backpack: inventory.NewInventory(8, 8),
//         //         Vault: inventory.NewInventory(8, 16), 
//         //         Location: town,
//         //         Strength: atts.Strength.int(),
//         //         Agility: atts.Agility.int(),
//         //         Energy: atts.Energy.int(),
//         //         Vitality: atts.Vitality.int(),
//         //         HpRegenerationRate: getHealthRegenerationRate(atts),
//         //         Level: stats.Level.int(),
//         //         Experience: atts.Experience.int(),
//         //         Direction: maps.Direction{Dx: 0, Dy: 1},
//         //         Skills: skill.Skills,
//         //         Equipment: inventory.NewEquipment(),
//         //     }
//         }        


//         fighters.FightersMap.Add(strId, fighter)

//         getBackpackFromDB(fighter)
//         updateFighterParams(fighter) 

//         PopulationMap.Add(town, fighter)    

//         ConnectionsMap.AddWithValues(conn, fighter, common.HexToAddress(ownerAddess))        
        
//         go initiateFighterRoutine(conn, fighter)
//         getFighterItems(fighter)
//     }

//     //FaucetCredits(conn)

    

    
// }



func findFighterByConn(conn *websocket.Conn) (*fighters.Fighter, error) {
    connection := ConnectionsMap.Find(conn)

    if connection == nil {
        return nil, fmt.Errorf("connection not found")
    }

    fighter := connection.GetFighter()
    if fighter == nil {
        return nil, fmt.Errorf("connection has no fighter")
    }

    return fighter, nil
}


func addDamageToFighter(fighterID, hitterID string, damage int) {
    found := false
    //fighter := fighters.FightersMap.Find(fighterID);
    fighter := PopulationMap.Find("lorencia", fighterID)

    log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v", fighterID, hitterID, damage)

    damageReceived := fighter.GetDamageReceived()

    // Check if there's already damage from the hitter
    for i, dmg := range damageReceived {
        if dmg.FighterId == hitterID {
            found = true
            //log.Printf("[addDamageToFighter] Damage found ")
            // Add the new damage to the existing damage
            damageReceived[i].Damage += damage
            break
        }
    }

    // If no existing damage from the hitter, create a new Damage object and append it to DamageReceived
    if !found {
        newDamage := battle.Damage{
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

func updateFighterParams(fighter *fighters.Fighter) {
    log.Printf("[updateFighterParams]  1 fighter=%v", fighter)
    defence := fighter.GetAgility()/4
    damage := fighter.GetStrength()/4 + fighter.GetEnergy()/4

    criticalDmg     := int(0)
    excellentDmg    := int(0)
    doubleDmg       := int(0)
    ignoreDef       := int(0)



    for _, item := range fighter.GetEquipment().GetMap() {
        itemAttributes := item.GetAttributes()
        // Perform your logic with the current item and slot
        defence += itemAttributes.ItemParameters.Defense + itemAttributes.AdditionalDefense
        damage += itemAttributes.ItemParameters.MinPhysicalDamage + itemAttributes.AdditionalDamage

        if itemAttributes.Luck {
            criticalDmg += 5
        }

        if itemAttributes.ExcellentItemAttributes.ExcellentDamageProbabilityIncrease > 0 {
            excellentDmg += 10
        }

        if itemAttributes.ExcellentItemAttributes.DoubleDamageChance > 0 {
            doubleDmg += 10
        }

        if itemAttributes.ExcellentItemAttributes.IgnoreOpponentDefenseChance > 0 {
            ignoreDef += itemAttributes.ExcellentItemAttributes.IgnoreOpponentDefenseChance
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




func applyConsumable(fighter *fighters.Fighter, item *items.TokenAttributes) {
    switch item.GetName() {
        case "Small Healing Potion":
            go graduallyIncreaseHp(fighter, 100, 5)
            break

        case "Small Mana Potion":
            go graduallyIncreaseMana(fighter, 100, 5)
            break

        default:
            log.Printf("[applyConsumable] Unknown consumable=%v", item.GetName())
            break
    }
}


func graduallyIncreaseHp(fighter *fighters.Fighter, hp int, chunks int) {
    // Calculate how much to increase HP by each chunk
    hpIncrease := hp / chunks

    for i := int(0); i < chunks; i++ {
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

func graduallyIncreaseMana(fighter *fighters.Fighter, mana int, chunks int) {
    // Calculate how much to increase HP by each chunk
    manaIncrease := mana / chunks

    for i := int(0); i < chunks; i++ {
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

func updateFighterDB(fighter *fighters.Fighter) {
    collection := db.Client.Database("game").Collection("fighters")

    // Marshalling the fighter object to JSON
    jsonFighter, err := json.Marshal(fighter)
    if err != nil {
        log.Fatal(err)
    }

    filter := bson.D{{Key: "fighterID", Value: fighter.ID}}
    update := bson.D{
        {Key: "$set", Value: bson.D{
            {Key: "fighterID", Value: fighter.ID},
            {Key: "atts", Value: string(jsonFighter)},
        }},
    }

    upsert := true
    opt := options.UpdateOptions{
        Upsert: &upsert,
    }

    _, err = collection.UpdateOne(context.Background(), filter, update, &opt)
    if err != nil {
        log.Fatal(err)
    }
}

func retrieveFighterFromDB(fighterID string) (*fighters.Fighter, error) {
    log.Printf("[retrieveFighterFromDB] fighterID=%v", fighterID)
    collection := db.Client.Database("game").Collection("fighters")

    filter := bson.D{{Key: "fighterID", Value: fighterID}}
    var result struct {
        FighterID string `bson:"fighterID"`
        Atts      string `bson:"atts"`
    }

    err := collection.FindOne(context.Background(), filter).Decode(&result)
    if err != nil {
        return nil, err
    }

    var fighter fighters.Fighter
    err = json.Unmarshal([]byte(result.Atts), &fighter)
    if err != nil {
        return nil, err
    }
    log.Printf("[retrieveFighterFromDB] fighter=%v", fighter.GetName())
    return &fighter, nil
}

// func getBackpackFromDB(fighter *fighters.Fighter) (bool) {
//     collection := db.Client.Database("game").Collection("Backpacks")

//     filter := bson.M{"fighterId": fighter.GetTokenID()}

//     var result bson.M
//     err := collection.FindOne(context.Background(), filter).Decode(&result)

//     if err != nil {
//         log.Printf("[getBackpackFromDB] Error getting Inventory from database: %v", err)
//         return false
//     }

//     InventoryStr, ok := result["backpack"].(string)
//     if !ok {
//         log.Printf("[getBackpackFromDB] Error asserting Inventory as string")
//         return false
//     }

//     var backpack *inventory.Inventory
//     err = json.Unmarshal([]byte(InventoryStr), &backpack)
//     if err != nil {
//         log.Printf("[getBackpackFromDB] Error unmarshaling Backpack: %v", err)
//         return false
//     }

//     equipmentStr, ok := result["equipment"].(string)
//     if !ok {
//         log.Printf("[getBackpackFromDB] Error asserting equipment as string")
//         return false
//     }

//     var equipment map[int]*inventory.InventorySlot
//     err = json.Unmarshal([]byte(equipmentStr), &equipment)
//     if err != nil {
//         log.Printf("[getBackpackFromDB] Error unmarshaling equipment: %v", err)
//         return false
//     }

//     log.Printf("[getBackpackFromDB] backpack=%v equipment=%v", backpack, equipment)
//     fighter.GetBackpack().Set(backpack)
//     fighter.GetEquipment().SetMap(equipment)
    

//     return true;
// }

// func saveBackpackToDB(fighter *fighters.Fighter) error {
//     log.Printf("[saveInventoryToDB] fighter=%v", fighter)
//     collection := db.Client.Database("game").Collection("Backpacks")

//     InventoryJSON, err := json.Marshal(fighter.GetBackpack())
//     if err != nil {
//         log.Printf("[saveInventoryToDB] Error marshaling Inventory: %v", err)
//         return err
//     }

//     InventoryStr := string(InventoryJSON)

//     equipmentJSON, err := json.Marshal(fighter.GetEquipment().GetMap())
//     if err != nil {
//         log.Printf("[saveInventoryToDB] Error marshaling Inventory: %v", err)
//         return err
//     }
//     filter := bson.M{"fighterId": fighter.TokenID}
//     equipmentStr := string(equipmentJSON)

    
//     update := bson.M{"$set": bson.M{"backpack": InventoryStr, "equipment": equipmentStr}}
//     opts := options.Update().SetUpsert(true)

//     _, err = collection.UpdateOne(context.Background(), filter, update, opts)
//     if err != nil {
//         log.Printf("[saveInventoryToDB] Error updating Inventory in database: %v", err)
//         return err
//     }

//     return nil
// }


func findTargetsByDirection(fighter *fighters.Fighter, dir maps.Direction, skill skill.Skill, targetId string) []*fighters.Fighter {
    targets := []*fighters.Fighter{}

    candidates := PopulationMap.GetTownMap(fighter.GetLocation())

    for _, candidate := range candidates {
        if fighter == candidate { continue }
        if !fighter.GetIsNpc() && !candidate.GetIsNpc() { continue }
        if fighter.GetIsNpc() && candidate.GetIsNpc() { continue }
        if candidate.GetIsNpc() && candidate.GetIsDead() { continue }

        //log.Printf("[findTargetsByDirection] candidate=%v", candidate)
        distance := maps.EuclideanDistance(fighter.GetCoordinates(), candidate.GetCoordinates())
        if distance <= float64(skill.ActiveDistance)+0.5 {
            angle := math.Atan2(float64(dir.Dx), float64(dir.Dy)) * 180 / math.Pi
            targetAngle := math.Atan2(float64(candidate.GetCoordinates().Y-fighter.GetCoordinates().Y), float64(candidate.GetCoordinates().X-fighter.GetCoordinates().X)) * 180 / math.Pi
            angleDifference := math.Abs(angle - targetAngle)

            // Handle angle difference greater than 180 degrees
            if angleDifference > 180 {
                angleDifference = 360 - angleDifference
            }

            //log.Printf("[findTargetsByDirection] candidate=%v angleDifference=%v compAngle=%v", candidate, angleDifference, float64(skill.HitAngle) )

            if angleDifference <= float64(skill.HitAngle) {
                // If the skill is not multihit, return the list with a single target
                if !skill.Multihit && candidate.GetID() == targetId {
                    targets = append(targets, candidate)
                    return targets
                } else if skill.Multihit {
                    targets = append(targets, candidate)
                }
            }
        }
    }

    return targets
}


func moveFighter(fighter *fighters.Fighter, coords maps.Coordinate) {
    log.Printf("[moveFighter] coords=%v", coords)
    if fighter.Coordinates == coords {
        log.Printf("[moveFighter] Fighter already in the spot coords=%v", coords)
        sendErrorMessage(fighter, fmt.Sprintf("Already in spot coords=%v", coords))
        return
    }


    if isSquareOccupied(coords) {
        log.Printf("[moveFighter] Square occupied coords=%v", coords)
        sendErrorMessage(fighter, fmt.Sprintf("Square occupied coords=%v", coords))
        return
    }

    currentTime := int(time.Now().UnixNano()) / int(time.Millisecond)
    elapsedTime := currentTime - fighter.LastMoveTimestamp
    

    if elapsedTime < 60000 / fighter.MovementSpeed {
        
        speed := float64(99999)
        if elapsedTime > 0 {
            speed = float64(60000) / float64(elapsedTime)
        }

        log.Printf("[moveFighter] Moving too fast=%v", speed)
        sendErrorMessage(fighter, fmt.Sprintf("Moving too fast speed=%v", speed))
        return
    }

    fighter.Lock()
    fighter.Coordinates = coords
    fighter.LastMoveTimestamp = currentTime
    fighter.Unlock()


    broadcastNpcMove(fighter, coords)
    //pingFighter(fighter)
}


func updateFighterDirection(fighter *fighters.Fighter, dir maps.Direction) {
    if fighter == nil {
        log.Println("[updateFighterDirection] Error: Received nil fighter in updateFighterDirection")
        return
    }


    fighter.SetDirection(dir)
}


func isSquareOccupied(coord maps.Coordinate) bool {
    for _, fighter := range PopulationMap.GetTownMap("lorencia") {
        fighterCoords := fighter.GetCoordinates()
        if fighterCoords.X == coord.X && fighterCoords.Y == coord.Y {
            return true
        }
    }

    return false
}
func getEmptySquares(center maps.Coordinate, radius int, town string) []maps.Coordinate {
    emptySquares := []maps.Coordinate{}

    for x := center.X - radius; x <= center.X + radius; x++ {
        for y := center.Y - radius; y <= center.Y + radius; y++ {
            if maps.EuclideanDistance(center, maps.Coordinate{X: x, Y: y}) > float64(radius) {
                continue
            }

            occupied := false
            for _, fighter := range PopulationMap.GetTownMap(town) {
                coords := fighter.GetCoordinates()
                if coords.X == x && coords.Y == y {
                    occupied = true
                    break
                }
            }

            if !occupied {
                emptySquares = append(emptySquares, maps.Coordinate{X: x, Y: y})
            }
        }
    }

    return emptySquares
}

func findNearbyFighters(town string, coords maps.Coordinate, distance int, isNpc bool) []*fighters.Fighter {
    nearbyFighters := []*fighters.Fighter{}

    for _, fighter := range PopulationMap.GetTownMap(town) {  
        // Calculate the Euclidean distance between the given coordinates and the fighter's coordinates
        dist := maps.EuclideanDistance(coords, fighter.GetCoordinates())

        // Check if the distance is within the given range
        if dist <= float64(distance) && fighter.GetIsNpc() == isNpc && !fighter.GetIsDead() {
            nearbyFighters = append(nearbyFighters, fighter)
        }
    }

    // Sort the nearbyFighters by their distance to the coords
    sort.Slice(nearbyFighters, func(i, j int) bool {
        distI := maps.EuclideanDistance(coords, nearbyFighters[i].GetCoordinates())
        distJ := maps.EuclideanDistance(coords, nearbyFighters[j].GetCoordinates())
        return distI < distJ
    })

    return nearbyFighters
}


func findNearestEmptySquareToPlayer(npcCoord, playerCoord maps.Coordinate) maps.Coordinate {
    bestSquare := npcCoord
    minDistance := maps.EuclideanDistance(npcCoord, playerCoord)

    for _, dir := range maps.Directions {
        candidateSquare := maps.MoveInDirection(npcCoord, dir, 1)
        if !isSquareOccupied(candidateSquare) {
            distance := maps.EuclideanDistance(candidateSquare, playerCoord)
            for _, nextDir := range maps.Directions {
                nextSquare := maps.MoveInDirection(candidateSquare, nextDir, 1)
                if !isSquareOccupied(nextSquare) {
                    nextDistance := maps.EuclideanDistance(nextSquare, playerCoord)
                    for _, finalDir := range maps.Directions {
                        finalSquare := maps.MoveInDirection(nextSquare, finalDir, 1)
                        if !isSquareOccupied(finalSquare) {
                            finalDistance := maps.EuclideanDistance(finalSquare, playerCoord)
                            averageDistance := (distance + nextDistance + finalDistance) / 3
                            if averageDistance < minDistance {
                                minDistance = averageDistance
                                bestSquare = candidateSquare
                            }
                        }
                    }
                }
            }
        }
    }

    // Check if npcCoord is better than the bestSquare
    npcDistance := maps.EuclideanDistance(npcCoord, playerCoord)
    if npcDistance < minDistance {
        return npcCoord
    }

    return bestSquare
}

func EquipBackpackItem (fighter *fighters.Fighter, itemHash string, slotId int) {
    
    slot := fighter.GetBackpack().FindByHash(itemHash)
    log.Printf("[EquipInventoryItem] itemHash=%v, slotId=%v slot=%v", itemHash, slotId, slot)
    if slot == nil {
        log.Printf("[EquipInventoryItem] slot not found=%v", itemHash)
        return
    }

    atts := slot.Attributes
    if atts.ItemParameters.AcceptableSlot1 != slotId && atts.ItemParameters.AcceptableSlot2 != slotId {
        log.Printf("[EquipInventoryItem] Invalid slot for slotId=%v AcceptableSlot1=%v AcceptableSlot2=%v", slotId, atts.ItemParameters.AcceptableSlot1, atts.ItemParameters.AcceptableSlot1)
        return
    }

    // fighter.RLock()
    // currSlot, ok := fighter.Equipment[slotId]
    // fighter.RUnlock()
    _, exists := fighter.GetEquipment().Find(slotId)
    if exists {
        log.Printf("[EquipInventoryItem] Slot not empty %v", slotId)        
        return
    }


    fighter.GetEquipment().Dress(slotId, *slot)
    fighter.GetBackpack().RemoveItemByHash(itemHash)
    //wsSendInventory(fighter)

    updateFighterParams(fighter)

    return
}



func UnequipBackpackItem (fighter *fighters.Fighter, itemHash string, coords maps.Coordinate) {
    log.Printf("[UnequipInventoryItem] itemHash=%v, coords=%v ", itemHash, coords)
    
    //slot := getEquipmentSlotByHash(fighter, itemHash)
    slot, exists := fighter.GetEquipment().FindByHash(itemHash)
    if !exists {
        log.Printf("[UnequipInventoryItem] slot empty=%v", itemHash)
        return
    }

    atts := slot.Attributes

    _, _, error := fighter.GetBackpack().AddItem(atts, 1, itemHash)
    if error != nil {
        log.Printf("[UnequipBackpackItem] Not enough space=%v", itemHash)
        return
    }

    
    // removeItemFromEquipmentSlotByHash(fighter, itemHash)
    fighter.GetEquipment().RemoveByHash(itemHash)
    //
    
    updateFighterParams(fighter)

}








