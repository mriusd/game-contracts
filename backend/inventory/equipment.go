// equipment.go

package inventory

import (
	"context"
	"sync"
	//"encoding/json"
	"log"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/mriusd/game-contracts/items"
	"github.com/mriusd/game-contracts/db"
)

type Equipment struct {
	Map map[int]*InventorySlot	`json:"items" bson:"items"`
	IsEquipped map[int]bool		`json:"is_equipped" bson:"is_equipped"`
	OwnerId int 				`json:"-" bson:"owner_id"`
	sync.RWMutex				`json:"-" bson:"-"`
}

func (i *Equipment) GetMap() map[int]*InventorySlot {
	i.RLock()
	defer i.RUnlock()

	return i.Map
}

func (i *Equipment) SetMap(v map[int]*InventorySlot)  {
	i.Lock()
	defer i.Unlock()

	i.Map = v
}

func (i *Equipment) Find (k int) *InventorySlot {
	i.RLock()
	defer i.RUnlock()

	return i.Map[k]
}

func (i *Equipment) Dress (slotId int, item *InventorySlot) {
	i.Lock()
	i.Map[slotId] = item
	i.IsEquipped[slotId] = true
	i.Unlock()

	i.RecordToDB()
}

func (i *Equipment) RemoveByHash (hash string) {
	i.Lock()

	for slotId, slot := range i.Map {
		if slot.GetItemHash() == hash {
			i.IsEquipped[slotId] = false
			delete(i.Map, slotId)
		}
	}
	i.Unlock()

	i.RecordToDB()
}

func (i *Equipment) FindByHash (hash string) *InventorySlot {
	i.RLock()
	defer i.RUnlock()

	for _, slot := range i.Map {
		if slot.GetItemHash() == hash {
			return slot
		}
	}

	return nil
}

// func (e *Equipment) MarshalJSON() ([]byte, error) {
//     e.RLock()
//     defer e.RUnlock()

//     // Marshal the Map field directly
//     return json.Marshal(e.Map)
// }

// func NewEquipment(ownerId int) *Equipment {
// 	return &Equipment{
// 		Map: make(map[int]*InventorySlot),
// 		IsEquipped: make(map[int]bool),
// 		OwnerId: ownerId,
// 	}


// }

func NewEquipment(ownerId int) *Equipment {
    equipment := &Equipment{
        Map: make(map[int]*InventorySlot),
		IsEquipped: make(map[int]bool),
		OwnerId: ownerId,
    }

    for i := 1; i <= 11; i++ {
        equipment.IsEquipped[i] = false
    }

    return equipment
}


func GetEquipmentFromDB(ownerId int) (*Equipment, error) {
    filter := bson.M{"owner_id": ownerId}
    collection := db.Client.Database("game").Collection("equipment")

    var equipment Equipment
    err := collection.FindOne(context.Background(), filter).Decode(&equipment)
    if err != nil {
        if err == mongo.ErrNoDocuments {
        	log.Printf("[Equipment: GetFromDB] Equipment not found in db")
            return NewEquipment(ownerId), nil 
        }
        return nil, fmt.Errorf("[GetFromDB]: %w", err)
    }

    log.Printf("[Equipment: GetFromDB] Equipment found=%v", equipment)
    equipment.PopulateAttributes()
    return &equipment, nil
}


func (i *Equipment) RecordToDB() error {
    i.RLock()
    copy := *i 
    i.RUnlock()

    filter := bson.M{"owner_id": i.OwnerId}
    update := bson.M{"$set": copy}
    options := options.Update().SetUpsert(true)

    collection := db.Client.Database("game").Collection("equipment")
    _, err := collection.UpdateOne(context.Background(), filter, update, options)
    if err != nil {
        log.Printf("[Equipment: RecordToDB]: %w", err)
        return fmt.Errorf("[Equipment: RecordToDB]: %w", err)
    }

    log.Printf("[Equipment: RecordToDB] Equipment Recorded or Updated")
    return nil
}

func (i *Equipment) PopulateAttributes() {
	for index, slot := range i.Map {
		i.Map[index].Attributes.ItemAttributes = items.BaseItemAttributes[slot.Attributes.Name]
		i.Map[index].Attributes.ItemParameters = items.BaseItemParameters[slot.Attributes.Name]
	}
}









