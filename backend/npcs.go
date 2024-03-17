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

    "github.com/mriusd/game-contracts/battle"   
    "github.com/mriusd/game-contracts/maps"   
    "github.com/mriusd/game-contracts/fighters"   
    "github.com/mriusd/game-contracts/skill"   
)

type NPC struct {
    ID               int      `json:"id"`
    Name             string     `json:"name"`
    Level            int      `json:"level"`
    Strength         int      `json:"strength"`
    Agility          int      `json:"agility"`
    Energy           int      `json:"energy"`
    Vitality         int      `json:"vitality"`
    AttackSpeed      int      `json:"attackSpeed"`
    DropRarityLevel  int      `json:"dropRarityLevel"`
    RespawnLocations [][]string `json:"respawnLocations"`
    CanFight         bool       `json:"canFight"`
    MaxHealth        int      `json:"maxHealth"`
    Skill            int      `json:"skill"`
    MovementSpeed    int      `json:"movementSpeed"`
}

var npcs []NPC;
var uniqueNpcIdCounter int = 1000
var npcVissionDistance int = 5
var npcAllowedDistanceFromSpawn int = 10

func initiateNpcRoutine(npc *fighters.Fighter) {
    speed := npc.GetMovementSpeed()

    msPerHit := 60000 / speed
    delay := time.Duration(msPerHit) * time.Millisecond

    //location := maps.DecodeLocation(npc.GetLocation())
    town := npc.GetLocation()

    rand.Seed(time.Now().UnixNano())

    for {
        time.Sleep(delay)       

        now := int(time.Now().UnixNano()) / 1e6
        elapsedTimeMs := now - npc.GetLastDmgTimestamp()

        if npc.GetIsDead() && elapsedTimeMs >= 5000 {
            log.Printf("[initiateNpcRoutine] At least 5 seconds have passed since TimeOfDeath.")

            emptySquares := GetEmptySquares(npc.GetSpawnCoords(), 5, town)

            if len(emptySquares) == 0 {
                continue // No empty squares available to spawn the NPC
            }

            
            spawnCoord := emptySquares[rand.Intn(len(emptySquares))]
            maxHealth := npc.GetMaxHealth()

            npc.Lock()
            npc.IsDead = false
            npc.HealthAfterLastDmg = maxHealth
            npc.DamageReceived = []battle.Damage{}
            npc.Coordinates = spawnCoord
            npc.CurrentHealth = maxHealth
            npc.Unlock()

            //emitNpcSpawnMessage(npc)
        } else if npc.GetIsDead() {
            continue
        } else {
            nonNpcFighters := findNearbyFighters(npc.GetLocation(), npc.GetCoordinates(), npcVissionDistance, false)

            if len(nonNpcFighters) > 0 {
                closestFighter := nonNpcFighters[0]

                skill := skill.Get(npc.GetSkill());

                distance := maps.EuclideanDistance(npc.GetCoordinates(), closestFighter.GetCoordinates())
                //log.Printf("[initiateNpcRoutine] id=%v distance=%v ActiveDistance=%v Npc coords=%v fighterCoords=%v", fighter.ID, distance, skill.ActiveDistance, fighter.Coordinates, closestFighter.Coordinates)
                if distance <= float64(skill.ActiveDistance)+0.5 {
                    data := Hit{
                        OpponentID: closestFighter.GetID(),
                        PlayerID:   npc.GetID(),
                        Skill:      npc.GetSkill(),
                        Direction:  maps.GetDirection(npc.GetCoordinates(), closestFighter.GetCoordinates()),
                    }

                    rawMessage, err := json.Marshal(data)
                    if err != nil {
                        log.Printf("[initiateNpcRoutine] Error marshaling data: %v", err)
                        return
                    }
                    //log.Printf("[initiateNpcRoutine] ProcessHit data=%v", data )
                    direction := maps.GetDirection(npc.GetCoordinates(), closestFighter.GetCoordinates())
                    npc.SetDirection(direction)

                    //conn, _ := findConnectionByFighter(closestFighter)

                    // c, conn := findConnectionByFighter(closestFighter)
                    // log.Printf("closestFighter=%v conn=%v", closestFighter, conn)
                    // if c != nil {
                        
                    // } else {
                    //     log.Printf("[initiateNpcRoutine] Could not find connection")
                    // }
                    ProcessHit(npc, rawMessage)
                    
                } else {
                    nextSquare := findNearestEmptySquareToPlayer(npc.GetCoordinates(), closestFighter.GetCoordinates())
                    if npc.GetCoordinates() != nextSquare {
                        direction := maps.GetDirection(npc.GetCoordinates(), nextSquare)
                        // fighter.Lock()
                        // fighter.Direction = direction
                        // fighter.Coordinates = nextSquare
                        // fighter.Unlock()

                        npc.SetDirection(direction)
                        npc.SetCoordinates(nextSquare)
                        //broadcastNpcMove(fighter, nextSquare)
                    }                        
                }                    
            } else {
                // move randomly
                if rand.Intn(2) == 0 { 
                    // NPC decides to stop this time
                    continue
                }

                // Fetch current position and direction
                currentCoord := npc.GetCoordinates()
                currentDirection := npc.GetDirection()
                spawnCoordinate := npc.GetSpawnCoords()

                // Generate a new coordinate based on the current direction
                newCoord := maps.Coordinate{X: currentCoord.X + currentDirection.Dx, Y: currentCoord.Y + currentDirection.Dy}

                // Check if new coordinate is within the allowed radius from the spawn point
                if !isSquareOccupied(newCoord) && isWithinRadius(spawnCoordinate, newCoord, npcAllowedDistanceFromSpawn) {
                    npc.SetCoordinates(newCoord)
                    continue
                }

                // If the new coordinate is occupied or outside the radius, select a new direction randomly
                newDirectionIndex := rand.Intn(len(maps.Directions))
                newDirection := maps.Directions[newDirectionIndex]
                newCoord = maps.Coordinate{X: currentCoord.X + newDirection.Dx, Y: currentCoord.Y + newDirection.Dy}

                if !isSquareOccupied(newCoord) && isWithinRadius(spawnCoordinate, newCoord, npcAllowedDistanceFromSpawn) {
                    npc.SetDirection(newDirection)
                    npc.SetCoordinates(newCoord)
                    continue
                }
                
            }
        }
    }
}

// isWithinRadius checks if a point is within a specified radius of another point
func isWithinRadius(center, point maps.Coordinate, radius int) bool {
    dx := point.X - center.X
    dy := point.Y - center.Y
    distanceSquared := dx*dx + dy*dy
    radiusSquared := radius * radius
    return distanceSquared <= radiusSquared
}


func getNextUniqueNpcId() string {
    uniqueNpcIdCounter++
    return "npc_" + strconv.Itoa(int(uniqueNpcIdCounter))
}

func findNpcById(id int) *NPC {
    for _, npc := range npcs {
        if npc.ID == id {
            return &npc
        }
    }
    return nil
}

func FindNpcByName(name string) *NPC {
    lowerName := strings.ToLower(name)
    for _, npc := range npcs {
        if strings.ToLower(npc.Name) == lowerName {
            return &npc
        }
    }
    return nil
}

func emitNpcSpawnMessage(npc *fighters.Fighter) {
    sendSpawnNpcMessage(npc)
}

func sendSpawnNpcMessage(npc *fighters.Fighter)  {
    //log.Printf("[sendSpawnNpcMessage] ", npc)

    type jsonResponse struct {
        Action string `json:"action"`
        Npc *fighters.Fighter `json:"npc"`
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

func SpawnNPC(npcId int, location []string) {
    
    npc := findNpcById(npcId)
    //log.Printf("[spawnNPC] %v %v", npcId, npc)

    uniqueNpcId := getNextUniqueNpcId()

    town := location[0]
    x, _ := strconv.Atoi(location[1])
    y, _ := strconv.Atoi(location[2])

    centerCoord := maps.Coordinate{X: x, Y: y}
    emptySquares := GetEmptySquares(centerCoord, 5, town)

    if len(emptySquares) == 0 {
        return // No empty squares available to spawn the NPC
    }

    rand.Seed(time.Now().UnixNano())
    spawnCoord := emptySquares[rand.Intn(len(emptySquares))]
    

    fighter := &fighters.Fighter{
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
        Level: npc.Level,
        AttackSpeed: npc.AttackSpeed,
        Coordinates: spawnCoord,
        Skill: npc.Skill,
        SpawnCoords: centerCoord,
        Strength: npc.Strength,
        Agility: npc.Agility,
        Energy: npc.Energy,
        Vitality: npc.Vitality,
        MovementSpeed: npc.MovementSpeed,
        Direction: maps.Direction{Dx: 0, Dy: 1},
    }

    //fighters.FightersMap.Add(uniqueNpcId, fighter)

    emitNpcSpawnMessage(fighter);

    

    

    // if _, exists := Population[town]; !exists {
    //     Population[town] = make([]*Fighter, 0)
    // }

    // Population[town] = append(Population[town], fighter)

    PopulationMap.Add(town, fighter)
    go initiateNpcRoutine(fighter)
}

func LoadNPCs() {
    // Open the JSON file
    file, err := os.Open("./npcs.json")
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
                SpawnNPC(npc.ID, location)
                
            }
        }      
    }

    log.Printf("NPCs Loaded %v", npcs )
}

func getNPCs(locationHash string) []*fighters.Fighter {
    location := maps.DecodeLocation(locationHash);
    zone := location[0]

    // coord := Coordinate{X: x, Y: y}
    npcFighters := []*fighters.Fighter{}
    for _, fighter := range PopulationMap.GetTownMap(zone) {
        if fighter.GetIsNpc() {
            npcFighters = append(npcFighters, fighter)
        }
    }

    return npcFighters
}