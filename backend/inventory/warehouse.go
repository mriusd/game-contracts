// warehouse.go

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

type Warehouse struct {
	Map map[string]*InventorySlot	`json:"items" bson:"items"`
	OwnerID int `json:"owner_id" bson:"owner_id"`
	
	sync.RWMutex `json:"-" bson:"-"`
}

func (i *Warehouse) GetOwnerID() int {
	i.RLock()
	defer i.RUnlock()

	return i.OwnerID
}

func (i *Warehouse) GetMap() map[string]*InventorySlot {
	i.RLock()
	defer i.RUnlock()

	return i.Map
}

func (i *Warehouse) SetMap(v map[string]*InventorySlot)  {
	i.Lock()
	defer i.Unlock()

	i.Map = v
}

func (i *Warehouse) Find(itemHash string) *InventorySlot {
	i.RLock()
	defer i.RUnlock()

	return i.Map[itemHash]
}

func (i *Warehouse) RemoveByHash(itemHash string) {
	i.Lock()
	delete(i.Map, itemHash)
	i.Unlock()

	i.RecordToDB()
}

func (i *Warehouse) AddItem(item *InventorySlot) {
	i.Lock()
	i.Map[item.ItemHash] = item
	i.Unlock()

	i.RecordToDB()
}


func NewWarehouse(ownerId int) *Warehouse {
    warehouse := &Warehouse{
        Map: make(map[string]*InventorySlot),
		OwnerID: ownerId,
    }

    return warehouse
}


func GetWarehouseFromDB(ownerId int) (*Warehouse, error) {
    filter := bson.M{"owner_id": ownerId}
    collection := db.Client.Database("game").Collection("warehouses")

    var warehouse Warehouse
    err := collection.FindOne(context.Background(), filter).Decode(&warehouse)
    if err != nil {
        if err == mongo.ErrNoDocuments {
        	log.Printf("[GetWarehouseFromDB] Warehouse not found in db")
            return NewWarehouse(ownerId), nil 
        }
        return nil, fmt.Errorf("[GetWarehouseFromDB]: %w", err)
    }

    log.Printf("[GetWarehouseFromDB] Warehouse found=%v", warehouse)
    warehouse.PopulateAttributes()
    return &warehouse, nil
}


func (i *Warehouse) RecordToDB() error {
    i.RLock()
    copy := *i 
    i.RUnlock()

    filter := bson.M{"owner_id": i.OwnerID}
    update := bson.M{"$set": copy}
    options := options.Update().SetUpsert(true)

    collection := db.Client.Database("game").Collection("warehouse")
    _, err := collection.UpdateOne(context.Background(), filter, update, options)
    if err != nil {
        log.Printf("[Warehouse: RecordToDB]: %w", err)
        return fmt.Errorf("[Warehouse: RecordToDB]: %w", err)
    }

    log.Printf("[Warehouse: RecordToDB] Warehouse Recorded or Updated")
    return nil
}

func (i *Warehouse) PopulateAttributes() {
	for index, slot := range i.Map {
		i.Map[index].Attributes.ItemAttributes = items.BaseItemAttributes[slot.Attributes.Name]
		i.Map[index].Attributes.ItemParameters = items.BaseItemParameters[slot.Attributes.Name]
	}
}









