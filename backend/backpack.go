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
	Items map[string]*BackpackSlot `json:"items"`
}

type BackpackSlot struct {
	Attributes 	ItemAttributes 	`json:"itemAttributes"`
	ItemHash 	common.Hash 	`json:"itemHash"`
	Qty        	int64 			`json:"qty"`
}

func (b *Backpack) ConsumeBackpackItem(fighter *Fighter, itemHash common.Hash) error {
    slot := getBackpackSlotByHash(fighter, itemHash)

    if slot.Qty <= 0 {
        return fmt.Errorf("no items left to consume for itemHash %s", itemHash.String())
    }

    slot.Qty--

    if slot.Qty == 0 {
    	b.removeItemByHash(fighter, itemHash);
    	BurnItem(slot.Attributes);
    }

    return nil
}


func removeItemFromEquipmentSlotByHash(fighter *Fighter, itemHash common.Hash) bool {
	fighter.Mutex.Lock()
	// Iterate through the equipment slots in the fighter
	for slotID, slot := range fighter.Equipment {
		// If the itemHash matches the current slot's itemHash, remove the item from the slot
		if slot.ItemHash == itemHash {
			// Set the equipment slot to nil or delete the slot from the map, depending on your requirements
			//fighter.Equipment[slotID] = nil
			delete(fighter.Equipment, slotID)
			// Return true to indicate that the item was successfully removed
			fighter.Mutex.Unlock()
			wsSendBackpack(fighter)
			return true
		}
	}
	fighter.Mutex.Unlock()
	// If no matching equipment slot is found, return false
	return false
}

func UnequipBackpackItem (fighter *Fighter, itemHash common.Hash, coords Coordinate) {
	log.Printf("[UnequipBackpackItem] itemHash=%v, coords=%v ", itemHash, coords)
	
	slot := getEquipmentSlotByHash(fighter, itemHash)
	if slot == nil {
		log.Printf("[UnequipBackpackItem] slot empty=%v", itemHash)
		return
	}

	atts := slot.Attributes

	fighter.Mutex.Lock()
	_, _, error := fighter.Backpack.AddItem(atts, 1, itemHash)
	fighter.Mutex.Unlock()
	if error != nil {
		log.Printf("[UnequipBackpackItem] Not enough space=%v", itemHash)
		return
	}
	
	removeItemFromEquipmentSlotByHash(fighter, itemHash)
	updateFighterParams(fighter)

}


func EquipBackpackItem (fighter *Fighter, itemHash common.Hash, slotId int64) {
	
	slot := getBackpackSlotByHash(fighter, itemHash)
	log.Printf("[EquipBackpackItem] itemHash=%v, slotId=%v slot=%v", itemHash, slotId, slot)
	if slot == nil {
		log.Printf("[EquipBackpackItem] slot not found=%v", itemHash)
		return
	}

	atts := slot.Attributes
	if atts.AcceptableSlot1.Int64() != slotId && atts.AcceptableSlot2.Int64() != slotId {
		log.Printf("[EquipBackpackItem] Invalid slot for slotId=%v AcceptableSlot1=%v AcceptableSlot2=%v", slotId, atts.AcceptableSlot1, atts.AcceptableSlot2)
		return
	}

	fighter.Mutex.RLock()
	currSlot, ok := fighter.Equipment[slotId]
	fighter.Mutex.RUnlock()
	if ok && currSlot != nil {
		log.Printf("[EquipBackpackItem] Slot not empty %v", slotId)		
		return
	}


	fighter.Mutex.Lock()
	fighter.Equipment[slotId] = slot
	fighter.Mutex.Unlock()
	fighter.Backpack.removeItemByHash(fighter, itemHash)
	//wsSendBackpack(fighter)

	updateFighterParams(fighter)

	return
}

func getBackpackSlotByHash(fighter *Fighter, itemHash common.Hash) *BackpackSlot {
	fighter.Mutex.RLock()
	defer fighter.Mutex.RUnlock()
	for _, backpackSlot := range fighter.Backpack.Items {
		if backpackSlot.ItemHash == itemHash {
			return backpackSlot
		}
	}
	return nil
}

func getEquipmentSlotByHash(fighter *Fighter, itemHash common.Hash) *BackpackSlot {
	fighter.Mutex.RLock()
	defer fighter.Mutex.RUnlock()
	// Iterate through the equipment slots in the fighter
	for _, slot := range fighter.Equipment {
		// If the itemHash matches the current slot's itemHash, return that slot
		if slot.ItemHash == itemHash {
			return slot
		}
	}

	// If no matching equipment slot is found, return nil
	return nil
}


func (bp *Backpack) removeItemByHash(fighter *Fighter, itemHash common.Hash) bool {
    log.Printf("[removeItemByHash] itemHash=%v", itemHash)
    fighter.Mutex.Lock()
    

    for key, backpackSlot := range bp.Items {
        if backpackSlot.ItemHash == itemHash {
            // Parse the coordinate string to get x, y, width, and height
            coords := strings.Split(key, ",")
            x, _ := strconv.Atoi(coords[0])
            y, _ := strconv.Atoi(coords[1])

            width := backpackSlot.Attributes.ItemWidth.Int64()
            height := backpackSlot.Attributes.ItemHeight.Int64()

            // Call clearSpace
            bp.clearSpace(x, y, int(width), int(height))

            delete(bp.Items, key)
            fighter.Mutex.Unlock()
            wsSendBackpack(fighter)
            return true
        }
    }

    fighter.Mutex.Unlock()
    return false
}

func (bp *Backpack) updateBackpackPosition(fighter *Fighter, itemHash common.Hash, newPosition Coordinate) error {
	log.Printf("[updateBackpackPosition] itemHash=%v newPosition=%v", itemHash, newPosition)
	// Find the current position of the item in the backpack
	var currentItemSlot *BackpackSlot
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

	wsSendBackpack(fighter)

	return nil
}

func wsSendBackpack(fighter *Fighter) {
	type jsonResponse struct {
		Action string `json:"action"`
        Backpack *Backpack `json:"backpack"`
        Equipment map[int64]*BackpackSlot `json:"equipment"`
	}

    jsonResp := jsonResponse{
    	Action: "backpack_update",
        Backpack: fighter.Backpack,
        Equipment: fighter.Equipment,
    }

    response, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[wsSendBackpack] error: ", err)
        return
    }
    saveBackpackToDB(fighter)
    pingFighter(fighter)

    respondFighter(fighter, response)
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

	DroppedItemsMutex.Lock()
	dropEvent := DroppedItems[itemHash]
    DroppedItemsMutex.Unlock()
	
	item := dropEvent.Item
	item.TokenId = event.TokenId
	saveItemAttributesToDB(item)

	DroppedItemsMutex.Lock()
	delete(DroppedItems, itemHash)
    DroppedItemsMutex.Unlock()

	if (item.ItemAttributesId.Int64() != GoldItemId) {
        fighter.Mutex.Lock()
        _, _, err := fighter.Backpack.AddItem(item, dropEvent.Qty.Int64(), itemHash)
        fighter.Mutex.Unlock()
        saveBackpackToDB(fighter)
        if err != nil {
            log.Printf("[PickupDroppedItem] Backpack full: %v", itemHash)
            sendErrorMessage(fighter, "Backpack full")
        }
    }

	fmt.Printf("[handleItemPickedEvent] event: %+v\n", event)   
	broadcastPickupMessage(fighter, item, event.Qty)
}


func NewBackpack(width, height int) *Backpack {
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}
	return &Backpack{Grid: grid, Items: make(map[string]*BackpackSlot)}
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
				bp.Items[coordKey] = &BackpackSlot{Attributes: item, Qty: qty, ItemHash: itemHash}
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
    // coordKey := fmt.Sprintf("%d,%d", x, y)
    // bp.Items[coordKey] = &BackpackSlot{Attributes: bp.Items[coordKey].Attributes, Qty: bp.Items[coordKey].Qty, ItemHash: itemHash}
}




