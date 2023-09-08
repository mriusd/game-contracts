package main

import (
	"log"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Coordinate struct {
	X int64 `json:"x"`
	Y int64 `json:"z"`
}

var Population = make(map[string][]*Fighter)
var PopulationMutex sync.RWMutex

type Direction struct {
    Dx int64 `json:"dx"`
    Dy int64 `json:"dz"`
}

var directions = []Direction{
    {-1, 0}, {1, 0}, {0, -1}, {0, 1},
    {-1, -1}, {-1, 1}, {1, -1}, {1, 1},
}

type Location struct {
	X float64 `json:"x"`
	Z float64 `json:"z"`
}

type Rotation struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Scale struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type MapObject struct {
	Type     string   `json:"type"`
	Location Location `json:"location"`
	Rotation Rotation `json:"rotation"`
	Scale    Scale    `json:"scale"`
}

func updateFighterDirection(fighter *Fighter, dir Direction) {
	if fighter == nil {
        log.Println("[updateFighterDirection] Error: Received nil fighter in updateFighterDirection")
        return
    }


	fighter.Mutex.Lock()
	fighter.Direction = dir
	fighter.Mutex.Unlock()
}

var MapObjects = make(map[string][]MapObject)
var MapObjectsMutex sync.RWMutex

func distance(x1, z1, x2, z2 float64) float64 {
	dx := x1 - x2
	dz := z1 - z2
	return math.Sqrt(dx*dx + dz*dz)
}

func removeFighterFromPopulation(fighter *Fighter) {
    PopulationMutex.Lock()
    defer PopulationMutex.Unlock()
    
    for key, fighters := range Population {
        for i, f := range fighters {
            if f == fighter {
                // Remove the fighter from the slice.
                Population[key] = append(fighters[:i], fighters[i+1:]...)
                break
            }
        }
    }
}


func getMapObjectsInRadius(mapName string, radius, x, z float64) []MapObject {
	MapObjectsMutex.RLock()
	objects, found := MapObjects[mapName]
	MapObjectsMutex.RUnlock()

	if !found {
		log.Printf("[getMapObjectsInRadius] mapName=%s not found", mapName)
		return nil
	}

	var objectsInRadius []MapObject
	for _, obj := range objects {
		if distance(obj.Location.X, obj.Location.Z, x, z) <= radius {
			objectsInRadius = append(objectsInRadius, obj)
		}
	}

	return objectsInRadius
}

func loadMap(mapName string) {
	log.Printf("[loadMap] mapName=%v", mapName)
	data, err := ioutil.ReadFile("./maps/"+mapName+".json")
	if err != nil {
		log.Printf("[loadMap] err1=%v", err)
	}
	
	var objects []MapObject
	err = json.Unmarshal(data, &objects)
	if err != nil {
		log.Printf("[loadMap] err2=%v", err)
	}

	MapObjectsMutex.Lock()
	MapObjects[mapName] = objects
	MapObjectsMutex.Unlock()
}

func loadMaps() {
	loadMap("lorencia")
}


func getDirection(coord1, coord2 Coordinate) Direction {
	deltaX := coord2.X - coord1.X
	deltaY := coord2.Y - coord1.Y

	if deltaX > 0 {
		deltaX = 1
	} else if deltaX < 0 {
		deltaX = -1
	}

	if deltaY > 0 {
		deltaY = 1
	} else if deltaY < 0 {
		deltaY = -1
	}

	return Direction{Dx: deltaX, Dy: deltaY}
}


func findTargetsByDirection(fighter *Fighter, dir Direction, skill *Skill, targetId string) []*Fighter {
	PopulationMutex.RLock()
	defer PopulationMutex.RUnlock()

	targets := []*Fighter{}
	
	for _, candidate := range Population[fighter.Location] {
		if !fighter.IsNpc && !candidate.IsNpc { continue }
		if fighter.IsNpc && candidate.IsNpc { continue }
		if candidate.IsNpc && candidate.IsDead { continue }
		if fighter == candidate { continue }
		distance := euclideanDistance(fighter.Coordinates, candidate.Coordinates)
		if distance <= float64(skill.ActiveDistance)+0.5 {
			angle := math.Atan2(float64(dir.Dx), float64(dir.Dy)) * 180 / math.Pi
			targetAngle := math.Atan2(float64(candidate.Coordinates.Y-fighter.Coordinates.Y), float64(candidate.Coordinates.X-fighter.Coordinates.X)) * 180 / math.Pi
			angleDifference := math.Abs(angle - targetAngle)

			// Handle angle difference greater than 180 degrees
			if angleDifference > 180 {
				angleDifference = 360 - angleDifference
			}

			//log.Printf("[findTargetsByDirection] candidate=%v angleDifference=%v compAngle=%v", candidate, angleDifference, float64(skill.HitAngle) )

			if angleDifference <= float64(skill.HitAngle) {
				// If the skill is not multihit, return the list with a single target
				if !skill.Multihit && candidate.ID == targetId {
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



func findNearestEmptySquareToPlayer(npcCoord, playerCoord Coordinate) Coordinate {
    bestSquare := npcCoord
    minDistance := euclideanDistance(npcCoord, playerCoord)

    for _, dir := range directions {
        candidateSquare := moveInDirection(npcCoord, dir, 1)
        if !isSquareOccupied(candidateSquare) {
            distance := euclideanDistance(candidateSquare, playerCoord)
            for _, nextDir := range directions {
                nextSquare := moveInDirection(candidateSquare, nextDir, 1)
                if !isSquareOccupied(nextSquare) {
                    nextDistance := euclideanDistance(nextSquare, playerCoord)
                    for _, finalDir := range directions {
                        finalSquare := moveInDirection(nextSquare, finalDir, 1)
                        if !isSquareOccupied(finalSquare) {
                            finalDistance := euclideanDistance(finalSquare, playerCoord)
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
    npcDistance := euclideanDistance(npcCoord, playerCoord)
    if npcDistance < minDistance {
        return npcCoord
    }

    return bestSquare
}

func moveInDirection(coord Coordinate, dir Direction, steps int64) Coordinate {
    return Coordinate{
        X: coord.X + dir.Dx * steps,
        Y: coord.Y + dir.Dy * steps,
    }
}

func isSquareOccupied(coord Coordinate) bool {
    FightersMutex.Lock()
    defer FightersMutex.Unlock()
    for _, fighter := range Fighters {
        if fighter.Coordinates.X == coord.X && fighter.Coordinates.Y == coord.Y {
            return true
        }
    }
    



    return false
}

func decodeLocation(locationHash string) []string {
	return strings.Split(locationHash, "_")
}

func getEmptySquares(center Coordinate, radius int64, town string) []Coordinate {
    emptySquares := []Coordinate{}

    PopulationMutex.RLock()
    defer PopulationMutex.RUnlock()

    for x := center.X - radius; x <= center.X + radius; x++ {
        for y := center.Y - radius; y <= center.Y + radius; y++ {
            if euclideanDistance(center, Coordinate{X: x, Y: y}) > float64(radius) {
                continue
            }

            occupied := false
            for _, fighter := range Population[town] {
                if fighter.Coordinates.X == x && fighter.Coordinates.Y == y {
                    occupied = true
                    break
                }
            }

            if !occupied {
                emptySquares = append(emptySquares, Coordinate{X: x, Y: y})
            }
        }
    }

    return emptySquares
}


func getNextSquare(position Coordinate, destination Coordinate) Coordinate {
	dx := destination.X - position.X
	dy := destination.Y - position.Y

	nextX := position.X
	nextY := position.Y

	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		nextX += sign(dx)
	} else if math.Abs(float64(dy)) > math.Abs(float64(dx)) {
		nextY += sign(dy)
	} else {
		nextX += sign(dx)
		nextY += sign(dy)
	}

	return Coordinate{X: nextX, Y: nextY}
}

func sign(x int64) int64 {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	} else {
		return 0
	}
}

func moveFighter(fighter *Fighter, coords Coordinate) {
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

	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
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

	fighter.Mutex.Lock()
	fighter.Coordinates = coords
	fighter.LastMoveTimestamp = currentTime
	fighter.Mutex.Unlock()


	broadcastNpcMove(fighter, coords)
	pingFighter(fighter)
}

func euclideanDistance(coord1, coord2 Coordinate) float64 {
	deltaX := float64(coord1.X - coord2.X)
	deltaY := float64(coord1.Y - coord2.Y)
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

func findNearbyFighters(coords Coordinate, distance int64, isNpc bool) []*Fighter {
	nearbyFighters := []*Fighter{}

	PopulationMutex.RLock()
	defer PopulationMutex.RUnlock()

	for _, fighters := range Population {
		for _, fighter := range fighters {
			// Calculate the Euclidean distance between the given coordinates and the fighter's coordinates
			dist := euclideanDistance(coords, fighter.Coordinates)

			// Check if the distance is within the given range
			if dist <= float64(distance) && fighter.IsNpc == isNpc {
				nearbyFighters = append(nearbyFighters, fighter)
			}
		}
	}

	// Sort the nearbyFighters by their distance to the coords
	sort.Slice(nearbyFighters, func(i, j int) bool {
		distI := euclideanDistance(coords, nearbyFighters[i].Coordinates)
		distJ := euclideanDistance(coords, nearbyFighters[j].Coordinates)
		return distI < distJ
	})

	return nearbyFighters
}
