package main

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"strings"
	"strconv"
	"encoding/json"
)


type Backpack struct {
	Grid  [][]bool   `json:"grid"`
	Items map[string]BackpackSlot `json:"items"`
}

type BackpackSlot struct {
	Attributes 	ItemAttributes 	`json:"itemAttributes"`
	ItemHash 	common.Hash 	`json:"itemHash"`
	Qty        	int64 			`json:"qty"`
}

func (bp *Backpack) updateBackpackPosition(fighter *Fighter, itemHash common.Hash, newPosition Coordinate) error {
	log.Printf("[updateBackpackPosition] itemHash=%v newPosition=%v", itemHash, newPosition)
	// Find the current position of the item in the backpack
	var currentItemSlot BackpackSlot
	var currentCoordKey string
	for coordKey, slot := range bp.Items {
		if slot.ItemHash == itemHash {
			currentItemSlot = slot
			currentCoordKey = coordKey
			break
		}
	}

	if currentItemSlot.ItemHash == (common.Hash{}) {
		return errors.New("[updateBackpackPosition] item not found in backpack")
	}

	coords := strings.Split(currentCoordKey, ",")
	currentX, _ := strconv.Atoi(coords[0])
	currentY, _ := strconv.Atoi(coords[1])

	// Check if the new position is available and has enough space for the item
	width := int(currentItemSlot.Attributes.ItemWidth.Int64())
	height := int(currentItemSlot.Attributes.ItemHeight.Int64())
	if !bp.isSpaceAvailable(int(newPosition.X), int(newPosition.Y), width, height, currentX, currentY) {
		return errors.New("[updateBackpackPosition] not enough space in the new position")
	}

	// Remove the item from its current position in the grid and the Items map
	x, y := parseCoordinates(currentCoordKey)
	bp.clearSpace(x, y, width, height)
	delete(bp.Items, currentCoordKey)

	// Add the item to the new position in the grid and the Items map
	bp.fillSpace(int(newPosition.X), int(newPosition.Y), width, height, itemHash)
	newCoordKey := fmt.Sprintf("%d,%d", newPosition.X, newPosition.Y)
	bp.Items[newCoordKey] = currentItemSlot

	type jsonResponse struct {
		Action string `json:"action"`
        Backpack *Backpack `json:"backpack"`
	}

    jsonResp := jsonResponse{
    	Action: "backpack_update",
        Backpack: fighter.Backpack,
    }

    response, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[updateBackpackPosition] error: ", err)
        return err
    }
    respondFighter(fighter, response)

	return nil
}

func parseCoordinates(coordKey string) (int, int) {
	coords := strings.Split(coordKey, ",")
	x, _ := strconv.Atoi(coords[0])
	y, _ := strconv.Atoi(coords[1])
	return x, y
}

func (bp *Backpack) clearSpace(x, y, width, height int) {
	for row := y; row < y+height; row++ {
		for col := x; col < x+width; col++ {
			bp.Grid[row][col] = false
		}
	}
}


func handleItemPickedEvent(itemHash common.Hash, logEntry *types.Log, fighter *Fighter) {

	// Parse the contract ABI
	parsedABI := loadABI("Backpack")

	// Iterate through logs and unpack the event data

	event := ItemPickedEvent{}

	log.Printf("[handleItemPickedEvent] logEntry: %v", logEntry)

	err := parsedABI.UnpackIntoInterface(&event, "ItemPicked", logEntry.Data)
	if err != nil {
		log.Printf("[handleItemPickedEvent] Failed to unpack log data: %v", err)
		return
	}

	fmt.Printf("[handleItemPickedEvent] event: %+v\n", event)
    DroppedItemsMutex.Lock()
	item := DroppedItems[itemHash].Item
	item.TokenId = event.TokenId
    DroppedItemsMutex.Unlock()
	saveItemAttributesToDB(item)


    DroppedItemsMutex.Lock()
	delete(DroppedItems, itemHash)
    DroppedItemsMutex.Unlock()

	broadcastPickupMessage(fighter, item, event.Qty)
}


func NewBackpack(width, height int) *Backpack {
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}
	return &Backpack{Grid: grid, Items: make(map[string]BackpackSlot)}
}


func (bp *Backpack) AddItem(item ItemAttributes, qty int64, itemHash common.Hash) (int, int, error) {
	gridHeight := len(bp.Grid)
	gridWidth := len(bp.Grid[0])
	for y := 0; y < gridHeight-int(item.ItemHeight.Int64())+1; y++ {
    	for x := 0; x < gridWidth-int(item.ItemWidth.Int64())+1; x++ {

			if bp.isSpaceAvailable(x, y, int(item.ItemWidth.Int64()), int(item.ItemHeight.Int64()), -10, -10) {
				bp.fillSpace(x, y, int(item.ItemWidth.Int64()), int(item.ItemHeight.Int64()), itemHash)

				// Store the item and quantity in the Items map
				coordKey := fmt.Sprintf("%d,%d", x, y)
				bp.Items[coordKey] = BackpackSlot{Attributes: item, Qty: qty, ItemHash: itemHash}

				return x, y, nil
			}
		}
	}
	return -1, -1, errors.New("not enough space in backpack")
}


func (bp *Backpack) isSpaceAvailable(x, y, width, height, currentX, currentY int) bool {
	gridHeight := len(bp.Grid)
	gridWidth := len(bp.Grid[0])

	if x < 0 || x+width > gridWidth || y < 0 || y+height > gridHeight {
		return false
	}

	for row := y; row < y+height; row++ {
		for col := x; col < x+width; col++ {
			if bp.Grid[row][col] {
				// Check if row and col are outside the range of the item's current position
				if row < currentY || row >= currentY+height || col < currentX || col >= currentX+width {
					return false
				}
			}
		}
	}
	return true
}


func (bp *Backpack) fillSpace(x, y, width, height int, itemHash common.Hash) {
    for row := y; row < y+height; row++ {
        for col := x; col < x+width; col++ {
            bp.Grid[row][col] = true
        }
    }
    coordKey := fmt.Sprintf("%d,%d", x, y)
    bp.Items[coordKey] = BackpackSlot{Attributes: bp.Items[coordKey].Attributes, Qty: bp.Items[coordKey].Qty, ItemHash: itemHash}
}




