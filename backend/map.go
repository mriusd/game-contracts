package main

import (
	"log"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
	"fmt"

	"github.com/gorilla/websocket"
)

type Coordinate struct {
	X int64 `json:"x"`
	Y int64 `json:"z"`
}

type Decoration struct {
	Coords Coordinate `json:"coords"`
	Type string `json:"type"`
}

var Population = make(map[string][]*Fighter)
var PopulationMutex sync.RWMutex

var Decorations = make(map[string][]*Decoration)
var DecorationsMytex sync.RWMutex


type Direction struct {
    Dx int64 `json:"dx"`
    Dy int64 `json:"dz"`
}

var directions = []Direction{
    {-1, 0}, {1, 0}, {0, -1}, {0, 1},
    {-1, -1}, {-1, 1}, {1, -1}, {1, 1},
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
		if fighter.IsNpc && candidate.IsNpc { continue }
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
    minDistance := float64(100) //euclideanDistance(npcCoord, playerCoord)

    for _, dir := range directions {
        candidateSquare := Coordinate{
            X: npcCoord.X + dir.Dx,
            Y: npcCoord.Y + dir.Dy,
        }

        if !isSquareOccupied(candidateSquare) {
            distance := euclideanDistance(candidateSquare, playerCoord)
            if distance < minDistance {
                minDistance = distance
                bestSquare = candidateSquare
            }
        }
    }

    return bestSquare
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

func moveFighter(conn *websocket.Conn, coords Coordinate) {
	log.Printf("[moveFighter] coords=%v", coords)
	fighter := findFighterByConn(conn)
	if fighter.Coordinates == coords {
		log.Printf("[moveFighter] Fighter already in the spot coords=%v", coords)
		sendErrorMessage(fighter, fmt.Sprintf("Already in spot coords=%v", coords))
		return
	}

	if isSquareOccupied(coords) {
		log.Printf("[moveFighter] Square occupiedt coords=%v", coords)
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

	fighter.ConnMutex.Lock()
	fighter.Coordinates = coords
	fighter.LastMoveTimestamp = currentTime
	fighter.ConnMutex.Unlock()
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
