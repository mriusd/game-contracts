package maps

import (
	"log"
	"math"
	"strings"
	"sync"
	"encoding/json"
	"io/ioutil"
)

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"z"`
}



type Direction struct {
    Dx int `json:"dx"`
    Dy int `json:"dz"`
}

var Directions = []Direction{
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


var MapObjects = make(map[string][]MapObject)
var MapObjectsMutex sync.RWMutex

func distance(x1, z1, x2, z2 float64) float64 {
	dx := x1 - x2
	dz := z1 - z2
	return math.Sqrt(dx*dx + dz*dz)
}


func GetMapObjectsInRadius(mapName string, radius, x, z float64) []MapObject {
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
	data, err := ioutil.ReadFile("./maps/locations/"+mapName+".json")
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

func Load() {
	loadMap("lorencia")
}


func GetDirection(coord1, coord2 Coordinate) Direction {
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

func MoveInDirection(coord Coordinate, dir Direction, steps int) Coordinate {
    return Coordinate{
        X: coord.X + dir.Dx * steps,
        Y: coord.Y + dir.Dy * steps,
    }
}

func DecodeLocation(locationHash string) []string {
	return strings.Split(locationHash, "_")
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

func sign(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	} else {
		return 0
	}
}

func EuclideanDistance(coord1, coord2 Coordinate) float64 {
	deltaX := float64(coord1.X - coord2.X)
	deltaY := float64(coord1.Y - coord2.Y)
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}


