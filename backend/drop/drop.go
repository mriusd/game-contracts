// drop.go

package drop

import (
	"context"
	"math/rand"
    "time"
    "log"
    "sync"
    "fmt"

    // "errors"

	"github.com/mriusd/game-contracts/maps" 
    "github.com/mriusd/game-contracts/fighters" 
    "github.com/mriusd/game-contracts/items" 
    "github.com/mriusd/game-contracts/db" 
)

type ItemDroppedEvent struct {
	ItemHash    string    				`json:"itemHash"`
	Item        *items.TokenAttributes 	`json:"item"`
	Qty         int       				`json:"qty"`
	Town  		string 					`json:"town"`
	Coords      maps.Coordinate     	`json:"coords"`
    OwnerId     int       				`json:"ownerId"`
    Timestamp   int  					`json:"timestamp"`
}
type SafeDroppedItemsMap struct {
	Map map[string]ItemDroppedEvent
	sync.RWMutex
}

var DroppedItems = &SafeDroppedItemsMap{Map: make(map[string]ItemDroppedEvent, 0)}


func (i *SafeDroppedItemsMap) CleanDroppedItems() {
	droppedItems := i.GetMap()
	thirtySecondsInNanoSeconds := int(10 * time.Second)

	for hash, droppedItem := range droppedItems {
		if int(time.Now().UnixNano()) > (droppedItem.Timestamp + thirtySecondsInNanoSeconds) {
			i.Remove(hash)
		}
	}
}

func (i *SafeDroppedItemsMap) GetMap() map[string]ItemDroppedEvent {
    i.RLock()
    defer i.RUnlock()

    copy := make(map[string]ItemDroppedEvent, len(i.Map)) // Initialize with the same size for efficiency
    for key, val := range i.Map {
        copy[key] = val
    }
    
    return copy
}


func (i *SafeDroppedItemsMap) Remove(k string) {
    i.Lock()
    defer i.Unlock()

    delete(i.Map, k)
}


func (i *SafeDroppedItemsMap) Add(k string, item ItemDroppedEvent) {
	i.Lock()
	defer i.Unlock()

	i.Map[k] = item
	log.Printf("[DroppedItems:Add] dropped item added %v", k)
}

func (i *SafeDroppedItemsMap) Find(k string) (ItemDroppedEvent, bool) {
	i.RLock()
	defer i.RUnlock()

	droppedItem, exists := i.Map[k]

	if !exists {
		return ItemDroppedEvent{}, false
	}

	return droppedItem, true
}



func getRandomNumberMax(min, max int) int {
    rand.Seed(time.Now().UnixNano())
    return int(rand.Intn(int(max-min+1))) + min
}

func returnRandomItemFromDropList(dummyParam uint,  its []items.ItemAttributes) items.ItemAttributes {
    if len(its) == 0 {
        log.Printf("[returnRandomItemFromDropList] empty item list")
        goldItem, _ := items.GetItemAttributesByName("Gold")
        return goldItem
    }

    index := getRandomNumberMax(0, int(len(its)-1))
    return its[index]
}


func DropItem (fighter *fighters.Fighter, itemHash string) {
	log.Printf("[DropItem] itemHash=%v fighter=%v", itemHash, fighter)

	backpack := fighter.GetBackpack()
	itemSlot := backpack.FindByHash(itemHash)

	if itemSlot == nil {
		return
	}

	item 	 := itemSlot.GetAttributes()

	dropEvent := ItemDroppedEvent{
		ItemHash: itemHash,
		Item: item,
		Qty: 1,
		Town: fighter.GetLocation(),
		Coords: fighter.GetCoordinates(),
		OwnerId: item.FighterId,
		Timestamp: int(time.Now().UnixNano()),
	}

	backpack.RemoveItemByHash(itemHash)
	DroppedItems.Add(itemHash, dropEvent)
}

func DropNewItem(rarityLevel int, hunter *fighters.Fighter, town string, coords maps.Coordinate, exp int) ItemDroppedEvent {	
	item := getDropItem(rarityLevel)

	log.Printf("[DropNewItem] rarityLevel=%v hunter=%v", rarityLevel, hunter)
	item.FighterId = hunter.GetTokenID()
	item.CreatedAt = int(time.Now().UnixNano())

	itemHash, err := items.HashItemAttributes(item)
	if err != nil {
		log.Fatalf("[DropItem] Failed to hash item=%v err=%v", item, err)
	}

	qty := 1
	if item.Name == "Gold" {
		qty = max(1, exp)
	}

	dropEvent := ItemDroppedEvent{
		ItemHash: itemHash,
		Item: item,
		Qty: qty,
		Town: town,
		Coords: coords,
		OwnerId: item.FighterId,
		Timestamp: int(time.Now().UnixNano()),
	}

	log.Printf("[DropNewItem] itemName=%v", dropEvent.Item.Name)

	DroppedItems.Add(itemHash, dropEvent)
	return dropEvent
}

func PickItem (fighter *fighters.Fighter, itemHash string) error {
	log.Printf("[PickItem] itemHash=%v fighter=%v", itemHash, fighter)

	droppedItem, exists := DroppedItems.Find(itemHash)
	if !exists {		
		log.Printf("[PickItem] DroppedItems=%v", DroppedItems.GetMap())
		return fmt.Errorf("[PickItem] Item not found=%v", itemHash)
	}

	log.Printf("[PickItem] droppedItem=%v fighter=%v", droppedItem.Item, fighter)

	backpack := fighter.GetBackpack()

	if droppedItem.Item.Name == "Gold" {
		DroppedItems.Remove(droppedItem.ItemHash)
		currGold := backpack.GetGold()
		backpack.SetGold(currGold + droppedItem.Qty)
		return nil
	}


	_, _, err := backpack.AddItem(droppedItem.Item, 1, itemHash)
	if err != nil {
		return err
	}

	DroppedItems.Remove(droppedItem.ItemHash)

	if droppedItem.Item.TokenId == 0 {
		RecordItemToDB(droppedItem.Item)
	}		

	return nil
}

func RecordItemToDB(item *items.TokenAttributes) error {
    // Assuming db.GetNextSequenceValue generates the next sequence value for Fighter ID
    nextID, err := db.GetNextSequenceValue("item")
    if err != nil {
        return err
    }

    item.TokenId = nextID

    collection := db.Client.Database("game").Collection("item")
    _, err = collection.InsertOne(context.Background(), item)
    if err != nil {
        return fmt.Errorf("RecordItemToDB: %w", err)
    }

    log.Printf("[RecordItemToDB] New item recorded with TokenID: %v", nextID)
    return nil
}



func getDropItem(rarityLevel int) *items.TokenAttributes {
	log.Printf("[getDropItem] rarityLevel=%v", rarityLevel)
	params := DropParamsMobMap[rarityLevel]

    randomNumber := getRandomNumberMax(0, 100)
    var itemAtts items.ItemAttributes

    item := items.TokenAttributes{}

    cumulativeRate := params.JewelsDropRate
    if randomNumber < cumulativeRate {
        itemAtts = returnRandomItemFromDropList(100, items.GetDropItems(rarityLevel, "jewel"))
    } else {
        cumulativeRate += params.ArmoursDropRate
        if randomNumber < cumulativeRate {
            itemAtts = returnRandomItemFromDropList(101, items.GetDropItems(rarityLevel, "armour"))
        } else {
            cumulativeRate += params.WeaponsDropRate
            if randomNumber < cumulativeRate {
                itemAtts = returnRandomItemFromDropList(102, items.GetDropItems(rarityLevel, "weapon"))
            } else {
                cumulativeRate += params.MiscDropRate
                if randomNumber < cumulativeRate {
                    itemAtts = returnRandomItemFromDropList(103, items.GetDropItems(rarityLevel, "misc"))
                } else {
                    cumulativeRate += params.BoxDropRate
                    if randomNumber < cumulativeRate {
                        itemAtts = returnRandomItemFromDropList(104, items.GetDropItems(rarityLevel, "box"))
                    } else {
                        itemAtts, _ = items.GetItemAttributesByName("Gold")
                    }
                }
            }
        }
    }

    if 	!itemAtts.IsJewel && 
    	!itemAtts.IsMisc && 
    	!itemAtts.IsBox && 
    	getRandomNumberMax(1, 100) <= params.LuckDropRate && 
    	itemAtts.Name != "Empty item" && 
    	itemAtts.Name != "Gold" {

        item.Luck = true
    }

    if 	itemAtts.IsWeapon && 
    	getRandomNumberMax(1, 100) <= params.SkillDropRate {

        item.Skill = true
    }

    if itemAtts.IsBox {
        item.ItemLevel = rarityLevel
    } else if !itemAtts.IsMisc && 
    	!itemAtts.IsJewel && 
    	itemAtts.Name !=  "Empty item" && 
    	itemAtts.Name != "Gold" {

        item.ItemLevel = params.MinItemLevel + getRandomNumberMax(0, params.MaxItemLevel-params.MinItemLevel)

        if itemAtts.IsWeapon {
            item.AdditionalDamage = 4 * getRandomNumberMax(0, params.MaxAddPoints/4)
        } else if itemAtts.IsArmour {
            item.AdditionalDefense = 4 * getRandomNumberMax(0, params.MaxAddPoints/4)
        }
    }

    exceAtts := items.ExcellentItemAttributes{}

    // if (itemAtts.IsWeapon || itemAtts.IsArmour) && getRandomNumberMax(1, 1000) <= int(params.ExcDropRate) {
    //     exceAtts = AddExcellentOption(itemAtts)
    // }
    
	item.Name = itemAtts.Name
	item.TokenId = 0
	item.PackSize = 1

	item.ItemAttributes = itemAtts
	item.ItemParameters = items.BaseItemParameters[itemAtts.Name]
	item.ExcellentItemAttributes = exceAtts

    return &item
}


func MakeItem(item *items.TokenAttributes, hunter *fighters.Fighter, town string, coords maps.Coordinate) ItemDroppedEvent {	

	log.Printf("[MakeItem] item.Name=%v", item.Name)
	item.FighterId = hunter.GetTokenID()
	item.CreatedAt = int(time.Now().UnixNano())

	itemHash, err := items.HashItemAttributes(item)
	if err != nil {
		log.Fatalf("[MakeItem] Failed to hash item=%v err=%v", item, err)
	}

	dropEvent := ItemDroppedEvent{
		ItemHash: itemHash,
		Item: item,
		Qty: 1,
		Town: town,
		Coords: coords,
		OwnerId: item.FighterId,
		Timestamp: int(time.Now().UnixNano()),
	}

	log.Printf("[MakeItem] itemName=%v", dropEvent.Item.Name)

	DroppedItems.Add(itemHash, dropEvent)

	
	return dropEvent
}


func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}