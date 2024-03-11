// events.go

package main

import (
	"math/big"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/mriusd/game-contracts/maps" 
	"github.com/mriusd/game-contracts/blockchain" 
	"github.com/mriusd/game-contracts/items" 
	"github.com/mriusd/game-contracts/fighters" 

)

func handleItemDroppedEvent(logEntry *types.Log, blockNumber *big.Int, coords maps.Coordinate, killer *big.Int) {
	// Parse the contract ABI
	parsedABI := blockchain.LoadABI("Items")

	// Iterate through logs and unpack the event data

	event := items.ItemDroppedEventSolidity{}

	err := parsedABI.UnpackIntoInterface(&event, "ItemDropped", logEntry.Data)
	if err != nil {
		log.Printf("[handleItemDroppedEvent] Failed to unpack log data: %v", err)
		return
	}

	log.Printf("[handleItemDroppedEvent] ItemHash: %v", event.ItemHash)

	event.BlockNumber = blockNumber
    event.Coords = coords
	event.OwnerId = killer

	// Add a self-destruct timer to remove the item from the map after 30 seconds
	time.AfterFunc(30*time.Second, func() {
		// DroppedItemsMutex.Lock() // Use a mutex if needed to protect access to the map
		// delete(DroppedItems, event.ItemHash)
		// DroppedItemsMutex.Unlock()

		items.DroppedItems.Remove(event.ItemHash)
		log.Printf("[handleItemDroppedEvent] Item with hash %v has been removed after 30 seconds", event.ItemHash)
		broadcastDropMessage()
	})

    // DroppedItemsMutex.Lock()
	// DroppedItems[event.ItemHash] = &event
    // DroppedItemsMutex.Unlock()
    items.DroppedItems.Add(event.ItemHash, &event)

	broadcastDropMessage()
}


func HandleItemPickedEvent(itemHash common.Hash, logEntry *types.Log, fighter *fighters.Fighter) {

	// Parse the contract ABI
	parsedABI := blockchain.LoadABI("Bqckpack")

	// Iterate through logs and unpack the event data

	event := items.ItemPickedEvent{}

	log.Printf("[handleItemPickedEvent] logEntry: %v", logEntry)

	err := parsedABI.UnpackIntoInterface(&event, "ItemPicked", logEntry.Data)
	if err != nil {
		log.Printf("[handleItemPickedEvent] Failed to unpack log data: %v", err)
		return
	}

	// DroppedItemsMutex.Lock()
	// dropEvent := DroppedItems[itemHash]
    // DroppedItemsMutex.Unlock()
    dropEvent := items.DroppedItems.Find(itemHash)
	
	item := dropEvent.Item
	item.TokenId = event.TokenId
	tokenAtts := items.ConvertSolidityItemToGoItem(item);
	items.SaveItemAttributesToDB(tokenAtts)

	// DroppedItemsMutex.Lock()
	// delete(DroppedItems, itemHash)
    // DroppedItemsMutex.Unlock()
    items.DroppedItems.Remove(itemHash)

	if item.Name != "Gold" {
        _, _, err := fighter.Backpack.AddItem(tokenAtts, dropEvent.Qty.Int64(), itemHash)
        saveBackpackToDB(fighter)
        if err != nil {
            log.Printf("[handleItemPickedEvent] Inventory full: %v", itemHash)
            sendErrorMessage(fighter, "Inventory full")
        }
    }

	log.Printf("[handleItemPickedEvent] event: %+v\n", event)   
	broadcastPickupMessage(fighter, tokenAtts, event.Qty)
}