// inventory.go

package inventory

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"strings"
	"strconv"
	"sync"

	"github.com/mriusd/game-contracts/items"
	"github.com/mriusd/game-contracts/maps"
)


type Inventory struct {
	Grid  [][]bool   `json:"grid"`
	Items map[string]*InventorySlot `json:"items"`
	sync.RWMutex
}

type InventorySlot struct {
	Attributes 	*items.TokenAttributes	`json:"itemAttributes"`
	ItemHash 	common.Hash 	`json:"itemHash"`
	Qty        	int64 			`json:"qty"`

	sync.RWMutex
}


func (i *Inventory) Set (v *Inventory)  {
	i.Lock()
	defer i.Unlock()
	
	i = v
}

func (i *Inventory) FindByHash (hash common.Hash) *InventorySlot {
	i.RLock()
	defer i.RUnlock()
	for _, InventorySlot := range i.Items {
		if InventorySlot.GetItemHash() == hash {
			return InventorySlot
		}
	}
	return nil
}

func NewInventory(width, height int) *Inventory {
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}
	return &Inventory{Grid: grid, Items: make(map[string]*InventorySlot)}
}


func (i *Inventory) GetGrid() [][]bool {
	i.RLock()
	i.RUnlock()

	return i.Grid
}

func (i *Inventory) GetItems() map[string]*InventorySlot {
	i.RLock()
	i.RUnlock()

	return i.Items
}



func (i *InventorySlot) GetItemHash() common.Hash {
	i.RLock()
	i.RUnlock()

	return i.ItemHash
}

func (i *InventorySlot) GetAttributes() *items.TokenAttributes {
	i.RLock()
	i.RUnlock()

	return i.Attributes
}

func (b *Inventory) Consume (itemHash common.Hash) error {
    slot := b.FindByHash(itemHash)

    if slot.Qty <= 0 {
        return fmt.Errorf("no items left to consume for itemHash %s", itemHash.String())
    }

    slot.Qty--

    if slot.Qty == 0 {
    	b.RemoveItemByHash(itemHash);
    	//BurnConsumable(fighter, slot.Attributes);
    }

    //applyConsumable(fighter, slot.Attributes)

    return nil
}


// func removeItemFromEquipmentSlotByHash(fighter *Fighter, itemHash common.Hash) bool {
	
// 	// Iterate through the equipment slots in the fighter
// 	for slotID, slot := range fighter.gEquipment().gMap() {
// 		// If the itemHash matches the current slot's itemHash, remove the item from the slot
// 		if slot.ItemHash == itemHash {
// 			// Set the equipment slot to nil or delete the slot from the map, depending on your requirements
// 			//fighter.Equipment[slotID] = nil
// 			fighter.Lock()
// 			delete(fighter.Equipment, slotID)
// 			// Return true to indicate that the item was successfully removed
// 			fighter.Unlock()
// 			wsSendBackpack(fighter)
// 			return true
// 		}
// 	}
// 	fighter.Unlock()
// 	// If no matching equipment slot is found, return false
// 	return false
// }






// func getBackpackSlotByHash(fighter *Fighter, itemHash common.Hash) *InventorySlot {
// 	fighter.RLock()
// 	defer fighter.RUnlock()
// 	for _, InventorySlot := range fighter.Backpack.Items {
// 		if InventorySlot.ItemHash == itemHash {
// 			return InventorySlot
// 		}
// 	}
// 	return nil
// }

// func GetVaultSlotByHash(fighter *Fighter, itemHash common.Hash) *InventorySlot {
// 	fighter.RLock()
// 	defer fighter.RUnlock()
// 	for _, vaultSlot := range fighter.Vault.Items {
// 		if vaultSlot.ItemHash == itemHash {
// 			return vaultSlot
// 		}
// 	}
// 	return nil
// }

// func getEquipmentSlotByHash(fighter *Fighter, itemHash common.Hash) *InventorySlot {
// 	fighter.RLock()
// 	defer fighter.RUnlock()
// 	// Iterate through the equipment slots in the fighter
// 	for _, slot := range fighter.Equipment {
// 		// If the itemHash matches the current slot's itemHash, return that slot
// 		if slot.ItemHash == itemHash {
// 			return slot
// 		}
// 	}

// 	// If no matching equipment slot is found, return nil
// 	return nil
// }


func (i *Inventory) RemoveItemByHash(itemHash common.Hash) bool {
    log.Printf("[removeItemByHash] itemHash=%v", itemHash)

    for key, InventorySlot := range i.GetItems() {
        if InventorySlot.GetItemHash() == itemHash {
            // Parse the coordinate string to get x, y, width, and height
            coords := strings.Split(key, ",")
            x, _ := strconv.Atoi(coords[0])
            y, _ := strconv.Atoi(coords[1])

            width := InventorySlot.GetAttributes().ItemParameters.ItemWidth
            height := InventorySlot.GetAttributes().ItemParameters.ItemHeight

            // Call clearSpace
            i.clearSpace(x, y, int(width), int(height))

            i.Lock()
            delete(i.Items, key)
            i.Unlock()
            return true
        }
    }

    return false
}

func (bp *Inventory) UpdateInventoryPosition(itemHash common.Hash, newPosition maps.Coordinate) error {
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
	x, y := ParseCoordinates(currentCoordKey)
	bp.clearSpace(x, y, width, height)
	delete(bp.Items, currentCoordKey)

	// Add the item to the new position in the grid and the Items map
	bp.fillSpace(int(newPosition.X), int(newPosition.Y), width, height, itemHash)
	newCoordKey := fmt.Sprintf("%d,%d", newPosition.X, newPosition.Y)
	bp.Items[newCoordKey] = currentItemSlot

	return nil
}



func ParseCoordinates(coordKey string) (int, int) {
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






func (i *Inventory) AddItem(item *items.TokenAttributes, qty int64, itemHash common.Hash) (int, int, error) {
	log.Printf("[AddItem] item: %v", item.GetName())
	grid := i.GetGrid()

	itemParameters := item.GetItemParameters()

	gridHeight := len(grid)
	gridWidth := len(grid[0])

	ih := int(itemParameters.GetItemHeight())
	iw := int(itemParameters.GetItemWidth())

	for y := 0; y < gridHeight - ih + 1; y++ {
    	for x := 0; x < gridWidth - iw + 1; x++ {

			if i.isSpaceAvailable(x, y, iw, ih, -10, -10) {
				i.fillSpace(x, y, iw, ih, itemHash)

				// Store the item and quantity in the Items map
				coordKey := fmt.Sprintf("%d,%d", x, y)
				i.Items[coordKey] = &InventorySlot{Attributes: item, Qty: qty, ItemHash: itemHash}
				return x, y, nil
			}
		}
	}
	return -1, -1, errors.New("not enough space in Inventory")
}

func (bp *Inventory) AddItemToPosition(item *items.TokenAttributes, qty int64, itemHash common.Hash, x,y int) (int, int, error) {
	log.Printf("[AddItemToPosition] item: %v", item)
	itemParameters := item.GetItemParameters()

	ih := int(itemParameters.GetItemHeight())
	iw := int(itemParameters.GetItemWidth())

	if bp.isSpaceAvailable(x, y, iw, ih, -10, -10) {
		bp.fillSpace(x, y, iw, ih, itemHash)

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




