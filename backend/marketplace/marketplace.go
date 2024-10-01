// marketplace.go

package marketplace

import (
	"context"
	"log"
	"sync"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mriusd/game-contracts/items"
	"github.com/mriusd/game-contracts/fighters"
	"github.com/mriusd/game-contracts/db"
)

type MarketplaceItem struct {
	Item            *items.TokenAttributes      `json:"item" bson:"item"`
	ItemHash    	string 						`json:"item_hash" bson:"item_hash"`
	Qty       		int 						`json:"qty" bson:"qty"`
	OwnerAccountID  primitive.ObjectID          `json:"owner_account_id" bson:"owner_account_id"`
	FighterTokenId  int                         `json:"fighter_token_id" bson:"fighter_token_id"`
	PriceGold       int                    		`json:"price_gold" bson:"price_gold"`
	PriceNefesh     int                         `json:"price_nefesh" bson:"price_nefesh"`
	PriceRuach      int                         `json:"price_ruach" bson:"price_ruach"`
	PriceNeshamah   int                         `json:"price_neshamah" bson:"price_neshamah"`
	PriceHaya      	int                         `json:"price_haya" bson:"price_haya"`

	CreatedAt 		time.Time 					`json:"created_at" bson:"created_at"`

	sync.RWMutex
}

type SafeMarketplaceMap struct {
	Items []*MarketplaceItem
	sync.RWMutex
}

var MarketplaceItems SafeMarketplaceMap


func (i *MarketplaceItem) GetItem() *items.TokenAttributes {
	i.RLock()
	defer i.RUnlock()

	return i.Item
}

func (i *MarketplaceItem) GetQty() int {
	i.RLock()
	defer i.RUnlock()

	return i.Qty
}

func (i *MarketplaceItem) GetOwnerAccountID() primitive.ObjectID {
	i.RLock()
	defer i.RUnlock()

	return i.OwnerAccountID
}

func (i *MarketplaceItem) GetFighterTokenId() int {
	i.RLock()
	defer i.RUnlock()

	return i.FighterTokenId
}

func (i *MarketplaceItem) GetItemHash() string {
	i.RLock()
	defer i.RUnlock()

	return i.ItemHash
}

func (i *SafeMarketplaceMap) FindItem (itemHash string) *MarketplaceItem {
	MarketplaceItems.RLock()
	defer MarketplaceItems.RUnlock()
	
	for _, marketplaceItem := range MarketplaceItems.Items {
		if marketplaceItem.ItemHash == itemHash {
			return marketplaceItem
		}
	}	

	return nil
}

func (i *SafeMarketplaceMap) AddItem(item *MarketplaceItem) {
	i.Lock()
	defer i.Unlock()

	i.Items = append(i.Items, item)
}

func (i *SafeMarketplaceMap) AddItemToMarketplace(fighter *fighters.Fighter, itemHash string, priceGold, priceNefesh, priceRuach, priceNeshamah, priceHaya int) error {
	wh := fighter.GetWarehouse()

	itemSlot := wh.Find(itemHash)

	// Create a new MarketplaceItem
	newMarketplaceItem := &MarketplaceItem{
		Item:           itemSlot.GetAttributes(),
		ItemHash: 		itemHash,
		Qty:			itemSlot.GetQty(),
		OwnerAccountID: fighter.GetAccountID(),  
		FighterTokenId: wh.GetOwnerID(),              
		PriceGold:      priceGold,
		PriceNefesh:    priceNefesh,
		PriceRuach:     priceRuach,
		PriceNeshamah:  priceNeshamah,
		PriceHaya:      priceHaya,

		CreatedAt: 		time.Now(),

	}

	// Add the new marketplace item to the database
	err := AddMarketplaceItemToDB(newMarketplaceItem)
	if err != nil {
		log.Printf("[AddItemToMarketplace] Error adding item to marketplace: %v", err)
		return err
	}

	i.AddItem(newMarketplaceItem)

	itemSlot.SetInTrade(true)

	log.Printf("[AddItemToMarketplace] Item successfully added to marketplace: %s", itemHash)
	return nil
}

func RemoveItemFromMarketplace(fighter *fighters.Fighter, itemHash string) error {
	itemToRemove := MarketplaceItems.FindItem(itemHash)

	fighterTokenId := fighter.GetTokenID()

	// If no item found, return an error
	if itemToRemove == nil {
		log.Printf("[RemoveItemFromMarketplace] Item not found in marketplace for FighterTokenId: %d, ItemHash: %s", fighterTokenId, itemHash)
		return fmt.Errorf("[RemoveItemFromMarketplace] item not found in marketplace")
	}

	if itemToRemove.GetFighterTokenId() != fighter.GetTokenID() {
		log.Printf("[RemoveItemFromMarketplace] Fighter not the item owner: %d, ItemHash: %s", fighterTokenId, itemHash)
		return fmt.Errorf("[RemoveItemFromMarketplace] You are not the item owner")
	}



	// Call the function to remove the item from DB and from the in-memory list
	err := RemoveMarketplaceItemFromDB(itemToRemove)
	if err != nil {
		log.Printf("[RemoveItemFromMarketplace] Error removing item from marketplace: %v", err)
		return err
	}

	itemSlot := fighter.GetWarehouse().Find(itemHash)
	itemSlot.SetInTrade(false)

	log.Printf("[RemoveItemFromMarketplace] Item successfully removed from marketplace: FighterTokenId: %d, ItemHash: %s", fighterTokenId, itemHash)

	return nil
}


func LoadMarketplaceItems() {
	collection := db.Client.Database("game").Collection("marketplace")
	ctx := context.Background()

	// Query for all items in the marketplace
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("[LoadMarketplaceItems] Error fetching items from marketplace: %v", err)
	}
	defer cursor.Close(ctx)

	var items []*MarketplaceItem

	for cursor.Next(ctx) {
		var item MarketplaceItem
		if err := cursor.Decode(&item); err != nil {
			log.Printf("[LoadMarketplaceItems] Error decoding marketplace item: %v", err)
			continue
		}
		items = append(items, &item)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("[LoadMarketplaceItems] Error iterating over marketplace items: %v", err)
	}

	MarketplaceItems.Lock()
	MarketplaceItems.Items = items
	MarketplaceItems.Unlock()

	log.Printf("[LoadMarketplaceItems] Loaded %d items into marketplace", len(items))
}

func AddMarketplaceItemToDB(item *MarketplaceItem) error {
	collection := db.Client.Database("game").Collection("marketplace")
	ctx := context.Background()

	// Insert the MarketplaceItem into the database
	_, err := collection.InsertOne(ctx, item)
	if err != nil {
		log.Printf("[AddMarketplaceItemToDB] Error inserting marketplace item to DB: %v", err)
		return err
	}

	log.Printf("[AddMarketplaceItemToDB] Marketplace item successfully added to DB, FighterTokenId: %d", item.FighterTokenId)

	return nil
}

func RemoveMarketplaceItemFromDB(item *MarketplaceItem) error {
	collection := db.Client.Database("game").Collection("marketplace")
	ctx := context.Background()

	itemHash := item.GetItemHash()

	// Remove the MarketplaceItem from the database
	_, err := collection.DeleteOne(ctx, bson.M{
		"item_hash":   item.GetItemHash(),  // Assuming item has a unique identifier (ItemHash)
	})
	if err != nil {
		log.Printf("[RemoveMarketplaceItemFromDB] Error removing marketplace item from DB: %v", err)
		return err
	}

	log.Printf("[RemoveMarketplaceItemFromDB] Marketplace item successfully removed from DB, FighterTokenId: %d", item.FighterTokenId)

	// Remove the MarketplaceItem from the in-memory list
	MarketplaceItems.Lock()
	defer MarketplaceItems.Unlock()

	for i, marketplaceItem := range MarketplaceItems.Items {
		if marketplaceItem.GetItemHash() == itemHash {
			// Remove item from the slice
			MarketplaceItems.Items = append(MarketplaceItems.Items[:i], MarketplaceItems.Items[i+1:]...)
			log.Printf("[RemoveMarketplaceItemFromDB] Marketplace item successfully removed from in-memory map")
			break
		}
	}

	return nil
}




