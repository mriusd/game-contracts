// inventory.go

package inventory

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/mriusd/game-contracts/items"
	"github.com/mriusd/game-contracts/maps"
	"github.com/mriusd/game-contracts/db"
)


type Inventory struct {
	Grid  [][]bool   					`json:"grid" bson:"grid"`
	Items map[string]*InventorySlot 	`json:"items" bson:"items"`
	Gold int 							`json:"gold" bson:"gold"`
	Type string 						`json:"-" bson:"type"`
	OwnerId int 						`json:"-" bson:"owner_id"`
	sync.RWMutex						`json:"-" bson:"-"`
}

type InventorySlot struct {
	Attributes 	*items.TokenAttributes	`json:"itemAttributes" bson:"item_attributes"`
	ItemHash 	string					`json:"itemHash" bson:"item_hash"`
	Qty        	int 					`json:"qty" bson:"qty"`

	sync.RWMutex						`json:"-" bson:"-"`
}


func (i *Inventory) Set (v *Inventory)  {
	i.Lock()
	defer i.Unlock()
	
	i = v
}

func (i *Inventory) FindByHash (hash string) *InventorySlot {
	i.RLock()
	defer i.RUnlock()
	for _, InventorySlot := range i.Items {
		if InventorySlot.GetItemHash() == hash {
			return InventorySlot
		}
	}
	return nil
}

func NewInventory(ownerId int, inventoryType string) *Inventory {
	return &Inventory{
		Grid: GetGrid(inventoryType), 
		Items: make(map[string]*InventorySlot),
		Type: inventoryType,
		OwnerId: ownerId,
	}
}


func GetGrid(inventoryType string) [][]bool {
	height := 8
	width := 8

	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}

	return grid
}

func (i *Inventory) GetGrid() [][]bool {
	i.RLock()
	defer i.RUnlock()

	return i.Grid
}

func (i *Inventory) GetItems() map[string]*InventorySlot {
	i.RLock()
	i.RUnlock()

	return i.Items
}

func (i *Inventory) GetGold() int {
	i.RLock()
	defer i.RUnlock()

	return i.Gold
}



func (i *InventorySlot) GetItemHash() string {
	i.RLock()
	defer i.RUnlock()

	return i.ItemHash
}

func (i *InventorySlot) GetAttributes() *items.TokenAttributes {
	i.RLock()
	defer i.RUnlock()

	return i.Attributes
}

func (b *Inventory) Consume (itemHash string) error {
    slot := b.FindByHash(itemHash)

    if slot.Qty <= 0 {
        return fmt.Errorf("no items left to consume for itemHash %s", itemHash)
    }

    slot.Qty--

    if slot.Qty == 0 {
    	b.RemoveItemByHash(itemHash);
    	//BurnConsumable(fighter, slot.Attributes);
    }

    //applyConsumable(fighter, slot.Attributes)

    return nil
}

func (i *Inventory) SetGold(v int) {
	i.Lock()
	i.Gold = v
	i.Unlock()	

	i.RecordToDB()
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


// func (i *Inventory) RecordToDB() error {
//     i.RLock()
//     dbId := fmt.Sprintf("%s_%s", i.Type, i.OwnerId)
//     copyOfInventory := *i // Make a copy of the Inventory
//     i.RUnlock()

//     collection := db.Client.Database("game").Collection("inventory")
//     _, err := collection.InsertOne(context.Background(), copyOfInventory) // Use the copy
//     if err != nil {
//     	log.Printf("[Inventory: RecordToDB]: %w", err)
//         return fmt.Errorf("[Inventory: RecordToDB]: %w", err)
//     }

//     log.Printf("[Inventory: RecordToDB] Inventory Recorded %v", dbId)
//     return nil
// }

func (i *Inventory) RecordToDB() error {
    i.RLock()
    copyOfInventory := *i 
    i.RUnlock()

    filter := bson.M{"owner_id": i.OwnerId, "type": i.Type}
    update := bson.M{"$set": copyOfInventory}
    options := options.Update().SetUpsert(true)

    collection := db.Client.Database("game").Collection("inventory")
    _, err := collection.UpdateOne(context.Background(), filter, update, options)
    if err != nil {
        log.Printf("[Inventory: RecordToDB]: %w", err)
        return fmt.Errorf("[Inventory: RecordToDB]: %w", err)
    }

    log.Printf("[Inventory: RecordToDB] Inventory Recorded or Updated")
    return nil
}

func GetInventoryFromDB(ownerId int, inventoryType string) (*Inventory, error) {
    filter := bson.M{"owner_id": ownerId, "type": inventoryType}
    collection := db.Client.Database("game").Collection("inventory")

    var inventory Inventory
    err := collection.FindOne(context.Background(), filter).Decode(&inventory)
    if err != nil {
        if err == mongo.ErrNoDocuments {
        	log.Printf("[Inventory: GetFromDB] Inventory not found in db")
            return NewInventory(ownerId, inventoryType), nil // Return an empty Inventory struct if no document is found
        }
        return nil, fmt.Errorf("[GetFromDB]: %w", err)
    }

    log.Printf("[Inventory: GetFromDB] Inventory found=%v", inventory)
    inventory.PopulateAttributes()
    return &inventory, nil
}

func (i *Inventory) PopulateAttributes() {
	for itemHash, _ := range i.Items {
		i.Items[itemHash].Attributes.ItemAttributes = items.BaseItemAttributes[i.Items[itemHash].Attributes.Name]
		i.Items[itemHash].Attributes.ItemParameters = items.BaseItemParameters[i.Items[itemHash].Attributes.Name]
	}
}




func (i *Inventory) RemoveItemByHash(itemHash string) bool {
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
            i.clearSpace(int(x), int(y), width, height)

            i.Lock()
            delete(i.Items, key)
            i.Unlock()
            i.RecordToDB()
            return true
        }
    }

    
    return false
}

func (bp *Inventory) UpdateInventoryPosition(itemHash string, newPosition maps.Coordinate) error {
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

	if currentItemSlot.ItemHash == "" {
		return errors.New("[updateInventoryPosition] item not found in Inventory")
	}

	coords := strings.Split(currentCoordKey, ",")
	currentX, _ := strconv.Atoi(coords[0])
	currentY, _ := strconv.Atoi(coords[1])

	// Check if the new position is available and has enough space for the item
	width := currentItemSlot.Attributes.ItemParameters.ItemWidth
	height := currentItemSlot.Attributes.ItemParameters.ItemHeight
	if !bp.isSpaceAvailable(newPosition.X, newPosition.Y, width, height, int(currentX), int(currentY)) {
		return errors.New("[updateInventoryPosition] not enough space in the new position")
	}

	// Remove the item from its current position in the grid and the Items map
	x, y := ParseCoordinates(currentCoordKey)
	bp.clearSpace(x, y, width, height)
	delete(bp.Items, currentCoordKey)

	// Add the item to the new position in the grid and the Items map
	bp.fillSpace(newPosition.X, newPosition.Y, width, height, itemHash)
	newCoordKey := fmt.Sprintf("%d,%d", newPosition.X, newPosition.Y)
	bp.Items[newCoordKey] = currentItemSlot

	bp.RecordToDB()
	return nil
}



func ParseCoordinates(coordKey string) (int,  int) {
	coords := strings.Split(coordKey, ",")
	x, _ := strconv.Atoi(coords[0])
	y, _ := strconv.Atoi(coords[1])
	return int(x), int(y)
}

func (bp *Inventory) clearSpace(x, y, width, height int) {
	for row := y; row < y+height; row++ {
		for col := x; col < x+width; col++ {
			bp.Grid[row][col] = false
		}
	}
}






func (i *Inventory) AddItem(item *items.TokenAttributes, qty int, itemHash string) (int,  int,  error) {
	log.Printf("[Inventory: AddItem] item: %v", item)
	grid := i.GetGrid()

	itemParameters := item.GetItemParameters()

	gridHeight := len(grid)
	gridWidth := len(grid[0])

	ih := itemParameters.ItemHeight
	iw := itemParameters.ItemWidth

	for y := 0; y < gridHeight - int(ih) + 1; y++ {
    	for x := 0; x < gridWidth - int(iw) + 1; x++ {

			if i.isSpaceAvailable(int(x), int(y), iw, ih, -10, -10) {
				i.fillSpace(int(x), int(y), iw, ih, itemHash)

				// Store the item and quantity in the Items map
				coordKey := fmt.Sprintf("%d,%d", x, y)
				i.Items[coordKey] = &InventorySlot{Attributes: item, Qty: qty, ItemHash: itemHash}
				i.RecordToDB()
				return x, y, nil
			}
		}
	}

	
	return -1, -1, errors.New("not enough space in Inventory")
}

func (bp *Inventory) AddItemToPosition(item *items.TokenAttributes, qty int, itemHash string, x,y int) (int,  int,  error) {
	log.Printf("[AddItemToPosition] item: %v", item)
	itemParameters := item.GetItemParameters()

	ih := itemParameters.ItemHeight
	iw := itemParameters.ItemWidth

	if bp.isSpaceAvailable(x, y, iw, ih, -10, -10) {
		bp.fillSpace(x, y, iw, ih, itemHash)

		// Store the item and quantity in the Items map
		coordKey := fmt.Sprintf("%d,%d", x, y)
		bp.Items[coordKey] = &InventorySlot{Attributes: item, Qty: qty, ItemHash: itemHash}
		bp.RecordToDB()
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


func (bp *Inventory) fillSpace(x, y, width, height int,  itemHash string) {
    for row := y; row < y+height; row++ {
        for col := x; col < x+width; col++ {
            bp.Grid[row][col] = true
        }
    }
    // coordKey := fmt.Sprintf("%d,%d", x, y)
    // bp.Items[coordKey] = &InventorySlot{Attributes: bp.Items[coordKey].Attributes, Qty: bp.Items[coordKey].Qty, ItemHash: itemHash}
}




