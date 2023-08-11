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


type Inventory struct {
	Grid  [][]bool   `json:"grid"`
	Items map[string]*InventorySlot `json:"items"`
}

type InventorySlot struct {
	Attributes 	TokenAttributes	`json:"itemAttributes"`
	ItemHash 	common.Hash 	`json:"itemHash"`
	Qty        	int64 			`json:"qty"`
}

func (b *Inventory) ConsumeBackpackItem(fighter *Fighter, itemHash common.Hash) error {
    slot := getBackpackSlotByHash(fighter, itemHash)

    if slot.Qty <= 0 {
        return fmt.Errorf("no items left to consume for itemHash %s", itemHash.String())
    }

    slot.Qty--

    if slot.Qty == 0 {
    	b.removeItemByHash(fighter, itemHash);
    	BurnConsumable(fighter, slot.Attributes);
    }

    applyConsumable(fighter, slot.Attributes)

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
	log.Printf("[UnequipInventoryItem] itemHash=%v, coords=%v ", itemHash, coords)
	
	slot := getEquipmentSlotByHash(fighter, itemHash)
	if slot == nil {
		log.Printf("[UnequipInventoryItem] slot empty=%v", itemHash)
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

	fighter.Mutex.RLock()
	currSlot, ok := fighter.Equipment[slotId]
	fighter.Mutex.RUnlock()
	if ok && currSlot != nil {
		log.Printf("[EquipInventoryItem] Slot not empty %v", slotId)		
		return
	}


	fighter.Mutex.Lock()
	fighter.Equipment[slotId] = slot
	fighter.Mutex.Unlock()
	fighter.Backpack.removeItemByHash(fighter, itemHash)
	//wsSendInventory(fighter)

	updateFighterParams(fighter)

	return
}

func getBackpackSlotByHash(fighter *Fighter, itemHash common.Hash) *InventorySlot {
	fighter.Mutex.RLock()
	defer fighter.Mutex.RUnlock()
	for _, InventorySlot := range fighter.Backpack.Items {
		if InventorySlot.ItemHash == itemHash {
			return InventorySlot
		}
	}
	return nil
}

func getVaultSlotByHash(fighter *Fighter, itemHash common.Hash) *InventorySlot {
	fighter.Mutex.RLock()
	defer fighter.Mutex.RUnlock()
	for _, vaultSlot := range fighter.Vault.Items {
		if vaultSlot.ItemHash == itemHash {
			return vaultSlot
		}
	}
	return nil
}

func getEquipmentSlotByHash(fighter *Fighter, itemHash common.Hash) *InventorySlot {
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


func (bp *Inventory) removeItemByHash(fighter *Fighter, itemHash common.Hash) bool {
    log.Printf("[removeItemByHash] itemHash=%v", itemHash)
    fighter.Mutex.Lock()
    

    for key, InventorySlot := range bp.Items {
        if InventorySlot.ItemHash == itemHash {
            // Parse the coordinate string to get x, y, width, and height
            coords := strings.Split(key, ",")
            x, _ := strconv.Atoi(coords[0])
            y, _ := strconv.Atoi(coords[1])

            width := InventorySlot.Attributes.ItemParameters.ItemWidth
            height := InventorySlot.Attributes.ItemParameters.ItemHeight

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

func (bp *Inventory) updateInventoryPosition(fighter *Fighter, itemHash common.Hash, newPosition Coordinate) error {
	log.Printf("[updateInventoryPosition] itemHash=%v newPosition=%v", itemHash, newPosition)
	// Find the current position of the item in the Inventory
	var currentItemSlot *InventorySlot
	var currentCoordKey string
	for coordKey, slot := range bp.Items {
		if slot.ItemHash == itemHash {
			currentItemSlot = slot
			currentCoordKey = coordKey
			break
		}
	}

	if currentItemSlot.ItemHash == (common.Hash{}) {
		return errors.New("[updateInventoryPosition] item not found in Inventory")
	}

	coords := strings.Split(currentCoordKey, ",")
	currentX, _ := strconv.Atoi(coords[0])
	currentY, _ := strconv.Atoi(coords[1])

	// Check if the new position is available and has enough space for the item
	width := int(currentItemSlot.Attributes.ItemParameters.ItemWidth)
	height := int(currentItemSlot.Attributes.ItemParameters.ItemHeight)
	if !bp.isSpaceAvailable(int(newPosition.X), int(newPosition.Y), width, height, currentX, currentY) {
		return errors.New("[updateInventoryPosition] not enough space in the new position")
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
        Backpack *Inventory `json:"backpack"`
        Equipment map[int64]*InventorySlot `json:"equipment"`
	}

    jsonResp := jsonResponse{
    	Action: "backpack_update",
        Backpack: fighter.Backpack,
        Equipment: fighter.Equipment,
    }

    response, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[wsSendInventory] error: ", err)
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

func (bp *Inventory) clearSpace(x, y, width, height int) {
	for row := y; row < y+height; row++ {
		for col := x; col < x+width; col++ {
			bp.Grid[row][col] = false
		}
	}
}


func handleItemPickedEvent(itemHash common.Hash, logEntry *types.Log, fighter *Fighter) {

	// Parse the contract ABI
	parsedABI := loadABI("Inventory")

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
	tokenAtts := convertSolidityItemToGoItem(item);
	saveItemAttributesToDB(tokenAtts)

	DroppedItemsMutex.Lock()
	delete(DroppedItems, itemHash)
    DroppedItemsMutex.Unlock()

	if item.Name != "Gold" {
        fighter.Mutex.Lock()
        _, _, err := fighter.Backpack.AddItem(tokenAtts, dropEvent.Qty.Int64(), itemHash)
        fighter.Mutex.Unlock()
        saveBackpackToDB(fighter)
        if err != nil {
            log.Printf("[handleItemPickedEvent] Inventory full: %v", itemHash)
            sendErrorMessage(fighter, "Inventory full")
        }
    }

	fmt.Printf("[handleItemPickedEvent] event: %+v\n", event)   
	broadcastPickupMessage(fighter, tokenAtts, event.Qty)
}


func NewInventory(width, height int) *Inventory {
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}
	return &Inventory{Grid: grid, Items: make(map[string]*InventorySlot)}
}


func (bp *Inventory) AddItem(item TokenAttributes, qty int64, itemHash common.Hash) (int, int, error) {
	log.Printf("[AddItem] item: %v", item)
	gridHeight := len(bp.Grid)
	gridWidth := len(bp.Grid[0])
	for y := 0; y < gridHeight-int(item.ItemParameters.ItemHeight)+1; y++ {
    	for x := 0; x < gridWidth-int(item.ItemParameters.ItemWidth)+1; x++ {

			if bp.isSpaceAvailable(x, y, int(item.ItemParameters.ItemWidth), int(item.ItemParameters.ItemHeight), -10, -10) {
				bp.fillSpace(x, y, int(item.ItemParameters.ItemWidth), int(item.ItemParameters.ItemHeight), itemHash)

				// Store the item and quantity in the Items map
				coordKey := fmt.Sprintf("%d,%d", x, y)
				bp.Items[coordKey] = &InventorySlot{Attributes: item, Qty: qty, ItemHash: itemHash}
				return x, y, nil
			}
		}
	}
	return -1, -1, errors.New("not enough space in Inventory")
}

func (bp *Inventory) AddItemToPosition(item TokenAttributes, qty int64, itemHash common.Hash, x,y int) (int, int, error) {
	log.Printf("[AddItemToPosition] item: %v", item)
	if bp.isSpaceAvailable(x, y, int(item.ItemParameters.ItemWidth), int(item.ItemParameters.ItemHeight), -10, -10) {
		bp.fillSpace(x, y, int(item.ItemParameters.ItemWidth), int(item.ItemParameters.ItemHeight), itemHash)

		// Store the item and quantity in the Items map
		coordKey := fmt.Sprintf("%d,%d", x, y)
		bp.Items[coordKey] = &InventorySlot{Attributes: item, Qty: qty, ItemHash: itemHash}
		return x, y, nil
	}
	return -1, -1, errors.New("not enough space in Inventory")
}


func (bp *Inventory) isSpaceAvailable(x, y, width, height, currentX, currentY int) bool {
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


func (bp *Inventory) fillSpace(x, y, width, height int, itemHash common.Hash) {
    for row := y; row < y+height; row++ {
        for col := x; col < x+width; col++ {
            bp.Grid[row][col] = true
        }
    }
    // coordKey := fmt.Sprintf("%d,%d", x, y)
    // bp.Items[coordKey] = &InventorySlot{Attributes: bp.Items[coordKey].Attributes, Qty: bp.Items[coordKey].Qty, ItemHash: itemHash}
}




