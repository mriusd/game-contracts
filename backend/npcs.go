package main

import (
	"strconv"
	"log"
	"encoding/json"
	"time"
	"io/ioutil"
	"os"
    "math/rand"
    "strings"
)



type NPC struct {
    ID               int64      `json:"id"`
    Name             string     `json:"name"`
    Level            int64      `json:"level"`
    Strength         int64      `json:"strength"`
    Agility          int64      `json:"agility"`
    Energy           int64      `json:"energy"`
    Vitality         int64      `json:"vitality"`
    AttackSpeed      int64      `json:"attackSpeed"`
    DropRarityLevel  int64      `json:"dropRarityLevel"`
    RespawnLocations [][]string `json:"respawnLocations"`
    CanFight         bool       `json:"canFight"`
    MaxHealth        int64      `json:"maxHealth"`
    Skill            int64      `json:"skill"`
    MovementSpeed    int64      `json:"movementSpeed"`
}

var npcs []NPC;
var uniqueNpcIdCounter int64 = 1000
var npcVissionDistance int64 = 10

func initiateNpcRoutine(fighter *Fighter) {
    speed := fighter.gMovementSpeed()

    msPerHit := 60000 / speed
    delay := time.Duration(msPerHit) * time.Millisecond

    location := decodeLocation(fighter.gLocation())
    town := location[0]

    for {
        time.Sleep(delay)       

        now := time.Now().UnixNano() / 1e6
        elapsedTimeMs := now - fighter.gLastDmgTimestamp()

        if fighter.gIsDead() && elapsedTimeMs >= 5000 {
            log.Printf("[initiateNpcRoutine] At least 5 seconds have passed since TimeOfDeath.")

            emptySquares := getEmptySquares(fighter.gSpawnCoords(), 5, town)

            if len(emptySquares) == 0 {
                continue // No empty squares available to spawn the NPC
            }

            rand.Seed(time.Now().UnixNano())
            spawnCoord := emptySquares[rand.Intn(len(emptySquares))]

            maxHealth := fighter.gMaxHealth()

            fighter.Lock()
            fighter.IsDead = false
            fighter.HealthAfterLastDmg = maxHealth
            fighter.DamageReceived = []Damage{}
            fighter.Coordinates = spawnCoord
            fighter.CurrentHealth = maxHealth
            fighter.Unlock()

            //emitNpcSpawnMessage(fighter)
        } else if fighter.gIsDead() {
            continue
        } else {
            nonNpcFighters := findNearbyFighters(fighter.gLocation(), fighter.gCoordinates(), npcVissionDistance, false)

            if len(nonNpcFighters) > 0 {
                closestFighter := nonNpcFighters[0]

                skill := getSkillSafely(fighter.gSkill());

                distance := euclideanDistance(fighter.gCoordinates(), closestFighter.gCoordinates())
                //log.Printf("[initiateNpcRoutine] id=%v distance=%v ActiveDistance=%v Npc coords=%v fighterCoords=%v", fighter.ID, distance, skill.ActiveDistance, fighter.Coordinates, closestFighter.Coordinates)
                if distance <= float64(skill.ActiveDistance)+0.5 {
                    data := Hit{
                        OpponentID: closestFighter.gID(),
                        PlayerID:   fighter.gID(),
                        Skill:      fighter.gSkill(),
                        Direction:  getDirection(fighter.gCoordinates(), closestFighter.gCoordinates()),
                    }

                    rawMessage, err := json.Marshal(data)
                    if err != nil {
                        log.Printf("[initiateNpcRoutine] Error marshaling data:", err)
                        return
                    }
                    //log.Printf("[initiateNpcRoutine] ProcessHit data=%v", data )
                    direction := getDirection(fighter.gCoordinates(), closestFighter.gCoordinates())
                    fighter.Lock()
                    fighter.Direction = direction
                    fighter.Unlock()

                    conn, _ := findConnectionByFighter(closestFighter)

                    if conn != nil {
                        ProcessHit(conn, rawMessage)
                    }
                    
                    
                } else {
                    nextSquare := findNearestEmptySquareToPlayer(fighter.gCoordinates(), closestFighter.gCoordinates())
                    if fighter.gCoordinates() != nextSquare {
                        direction := getDirection(fighter.gCoordinates(), nextSquare)
                        // fighter.Lock()
                        // fighter.Direction = direction
                        // fighter.Coordinates = nextSquare
                        // fighter.Unlock()

                        fighter.sDirection(direction)
                        fighter.sCoordinates(nextSquare)
                        //broadcastNpcMove(fighter, nextSquare)
                    }                        
                }                    
            }
        }
    }
}


func getNextUniqueNpcId() string {
    uniqueNpcIdCounter++
    return "npc_" + strconv.Itoa(int(uniqueNpcIdCounter))
}

func findNpcById(id int64) *NPC {
    for _, npc := range npcs {
        if npc.ID == id {
            return &npc
        }
    }
    return nil
}

func findNpcByName(name string) *NPC {
    lowerName := strings.ToLower(name)
    for _, npc := range npcs {
        if strings.ToLower(npc.Name) == lowerName {
            return &npc
        }
    }
    return nil
}

func emitNpcSpawnMessage(npc *Fighter) {
    sendSpawnNpcMessage(npc)
}

func sendSpawnNpcMessage(npc *Fighter)  {
    //log.Printf("[sendSpawnNpcMessage] ", npc)

    type jsonResponse struct {
        Action string `json:"action"`
        Npc *Fighter `json:"npc"`
    }

    jsonResp := jsonResponse{
        Action: "spawn_npc",
        Npc: npc,
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        log.Printf("[sendSpawnNpcMessage] %v %v", npc.ID, err)
    }

    broadcastWsMessage(npc.Location, messageJSON)
}

func spawnNPC(npcId int64, location []string) {
    
    npc := findNpcById(npcId)
    //log.Printf("[spawnNPC] %v %v", npcId, npc)

    uniqueNpcId := getNextUniqueNpcId()

    town := location[0]
    x, _ := strconv.ParseInt(location[1], 10, 64)
    y, _ := strconv.ParseInt(location[2], 10, 64)

    centerCoord := Coordinate{X: x, Y: y}
    emptySquares := getEmptySquares(centerCoord, 5, town)

    if len(emptySquares) == 0 {
        return // No empty squares available to spawn the NPC
    }

    rand.Seed(time.Now().UnixNano())
    spawnCoord := emptySquares[rand.Intn(len(emptySquares))]
    

    fighter := &Fighter{
        ID: uniqueNpcId,
        MaxHealth: npc.MaxHealth, 
        CurrentHealth: npc.MaxHealth, 
        Name: npc.Name,
        IsNpc: true,
        CanFight: npc.CanFight,
        HpRegenerationRate: 0,
        HpRegenerationBonus: 0,
        LastDmgTimestamp: 0,
        HealthAfterLastDmg: npc.MaxHealth,
        TokenID: npcId,
        Location: town,
        AttackSpeed: npc.AttackSpeed,
        Coordinates: spawnCoord,
        Skill: npc.Skill,
        SpawnCoords: centerCoord,
        Strength: npc.Strength,
        Agility: npc.Agility,
        Energy: npc.Energy,
        Vitality: npc.Vitality,
        MovementSpeed: npc.MovementSpeed,
        Direction: Direction{Dx: 0, Dy: 1},
    }

    FightersMap.Add(uniqueNpcId, fighter)

    emitNpcSpawnMessage(fighter);

    

    

    // if _, exists := Population[town]; !exists {
    //     Population[town] = make([]*Fighter, 0)
    // }

    // Population[town] = append(Population[town], fighter)

    PopulationMap.Add(town, fighter)
    go initiateNpcRoutine(fighter)
}

func loadNPCs() {
    // Open the JSON file
    file, err := os.Open("../npcList.json")
    if err != nil {
        log.Printf("[loadNPCs] error= %v", err)
    }
    defer file.Close()

    // Read the JSON data
    data, err := ioutil.ReadAll(file)
    if err != nil {
        log.Printf("[loadNPCs] error= %v", err)
    }

    // Unmarshal the JSON data into a slice of NPCs
    err = json.Unmarshal(data, &npcs)
    if err != nil {
        log.Printf("[loadNPCs] error= %v", err)
    }

    log.Printf("[loadNPCs] %v", npcs)

    // Set default values and initiate NPC routines
    for i, npc := range npcs {
        npcs[i].CanFight = true
        npcs[i].MaxHealth = npc.Vitality

        // Iterate through respawn locations
        for _, location := range npc.RespawnLocations {
            mobCount, err := strconv.Atoi(location[3])
            if err != nil {
                log.Printf("Error converting mob count to integer: %v", err)
                continue
            }
            for i := 0; i < mobCount; i++ {
                spawnNPC(npc.ID, location)
                
            }
        }      
    }

    log.Printf("NPCs Loaded %v", npcs )
}

func getNPCs(locationHash string) []*Fighter {
    location := decodeLocation(locationHash);
    zone := location[0]

    // coord := Coordinate{X: x, Y: y}
    npcFighters := []*Fighter{}
    for _, fighter := range PopulationMap.gTownMap(zone) {
        if fighter.gIsNpc() {
            npcFighters = append(npcFighters, fighter)
        }
    }

    return npcFighters
}