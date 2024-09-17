// marketplace.go

package marketplace

import (
	"context"
	"log"
	"sync"

	"github.com/mriusd/game-contracts/items"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"

	"github.com/mriusd/game-contracts/inventory"
	"github.com/mriusd/game-contracts/fighter"
)

type MarketplaceItem struct {
	Item            *items.TokenAttributes      `json:"item" bson:"item"`
	ItemHash    	string 						`json:"item_hash" bson:"item_hash"`
	Qty       		int 						`json:"qty" bson:"qty"`
	OwnerAccountID  primitive.ObjectID          `json:"owner_account_id" bson:"owner_account_id"`
	FighterTokenId  int                         `json:"fighter_token_id" bson:"fighter_token_id"`
	PriceGold       int                    		`json:"price_gold" bson:"price_gold"`
	PriceCredits    int                    		`json:"price_credits" bson:"price_credits"`
	PriceChaos      int                         `json:"price_chaos" bson:"price_chaos"`
	PriceSoul       int                         `json:"price_soul" bson:"price_soul"`
	PriceBless      int                         `json:"price_bless" bson:"price_bless"`
	PriceLife       int                         `json:"price_life" bson:"price_life"`

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

func AddItemToMarketplace(fighter *Fighter, inv *Inventory, slot *InventorySlot, priceGold, priceCredits, priceChaos, priceSoul, priceBless, priceLife int) error {
	if inv.GetType() != "vault" {
		log.Printf("[AddItemToMarketplace] Item not from vault %v", slot.ItemHash)
		return fmt.Errorf("[AddItemToMarketplace] Item not from vault")
	}

	// Create a new MarketplaceItem
	newMarketplaceItem := &MarketplaceItem{
		Item:           slot.GetAttributes(),
		ItemHash: 		slot.GetItemHash(),
		Qty:			slot.GetQty(),
		OwnerAccountID: fighter.GetAccountID(),  // Example owner ID (you'll want to fetch real owner details)
		FighterTokenId: inv.OwnerId,              // Assuming owner is the fighter token ID
		PriceGold:      priceGold,
		PriceCredits:   priceCredits,
		PriceChaos:     priceChaos,
		PriceSoul:      priceSoul,
		PriceBless:     priceBless,
		PriceLife:      priceLife,
	}

	// Add the new marketplace item to the database
	err := AddMarketplaceItemToDB(newMarketplaceItem)
	if err != nil {
		log.Printf("[AddItemToMarketplace] Error adding item to marketplace: %v", err)
		// Revert the inTrade status if there's an error
		slot.InTrade = false
		return err
	}

	inv.RemoveItemByHash(slot.GetItemHash())

	log.Printf("[AddItemToMarketplace] Item successfully added to marketplace: %s", slot.ItemHash)
	return nil
}

func RemoveItemFromMarketplace(fighter *Fighter, itemHash string) error {
	itemToRemove := MarketplaceItems.FindItem(itemHash)

	// If no item found, return an error
	if itemToRemove == nil {
		log.Printf("[RemoveItemFromMarketplace] Item not found in marketplace for FighterTokenId: %d, ItemHash: %s", fighterTokenId, itemHash)
		return fmt.Errorf("[RemoveItemFromMarketplace] item not found in marketplace")
	}

	vault := fighter.GetVault()
	_, _, err := vault.AddItem(itemToRemove.GetItem(), itemToRemove.GetQty(), itemToRemove.GetItemHash())
	if err != nil {
		log.Printf("[RemoveItemFromMarketplace] Error moving item to vault: %v", err)
		return err
	}

	// Call the function to remove the item from DB and from the in-memory list
	err = RemoveMarketplaceItemFromDB(itemToRemove)
	if err != nil {
		log.Printf("[RemoveItemFromMarketplace] Error removing item from marketplace: %v", err)
		return err
	}

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

	// Lock the MarketplaceItems to safely update the in-memory list
	MarketplaceItems.Lock()
	MarketplaceItems.Items = append(MarketplaceItems.Items, item)
	MarketplaceItems.Unlock()

	log.Printf("[AddMarketplaceItemToDB] Marketplace item successfully added to in-memory MarketplaceItems")
	return nil
}

func RemoveMarketplaceItemFromDB(item *MarketplaceItem) error {
	collection := db.Client.Database("game").Collection("marketplace")
	ctx := context.Background()

	// Remove the MarketplaceItem from the database
	_, err := collection.DeleteOne(ctx, bson.M{
		"fighter_token_id": item.FighterTokenId,
		"item.item_hash":   item.Item.ItemHash,  // Assuming item has a unique identifier (ItemHash)
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
		if marketplaceItem.FighterTokenId == item.FighterTokenId && marketplaceItem.Item.ItemHash == item.Item.ItemHash {
			// Remove item from the slice
			MarketplaceItems.Items = append(MarketplaceItems.Items[:i], MarketplaceItems.Items[i+1:]...)
			log.Printf("[RemoveMarketplaceItemFromDB] Marketplace item successfully removed from in-memory map")
			break
		}
	}

	return nil
}




