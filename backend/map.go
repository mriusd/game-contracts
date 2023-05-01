package main

import (
	"log"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Coordinate struct {
	X int64 `json:"x"`
	Y int64 `json:"z"`
}

var Population = make(map[string][]*Fighter)
var PopulationMutex sync.RWMutex




type Direction struct {
    dx int64
    dy int64
}

var directions = []Direction{
    {-1, 0}, {1, 0}, {0, -1}, {0, 1},
    {-1, -1}, {-1, 1}, {1, -1}, {1, 1},
}

func findNearestEmptySquareToPlayer(npcCoord, playerCoord Coordinate) Coordinate {
    bestSquare := npcCoord
    minDistance := euclideanDistance(npcCoord, playerCoord)

    for _, dir := range directions {
        candidateSquare := Coordinate{
            X: npcCoord.X + dir.dx,
            Y: npcCoord.Y + dir.dy,
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
		log.Printf("[moveFighter] Fighter already in the spot")
		return
	}

	fighter.Coordinates = coords
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
