package main

import (
	"sync"
	"errors"
	"log"
	//"fmt"
	//"time"

	"github.com/mriusd/game-contracts/inventory"
	"github.com/mriusd/game-contracts/fighters"
)

type Trade struct {	
	Fighter1  	*fighters.Fighter
	Fighter2  	*fighters.Fighter
	TradeGrid1 	*inventory.Inventory
	TradeGrid2  *inventory.Inventory
	sync.RWMutex
}

type SafeTradesMap struct {
	Trades []*Trade
	sync.RWMutex
}

var TradesMap = &SafeTradesMap{Trades: make([]*Trade, 0)}


func (i *Trade) gFighter1() *fighters.Fighter {
	i.RLock()
	i.RUnlock()

	return i.Fighter1
}

func (i *Trade) gFighter2() *fighters.Fighter {
	i.RLock()
	i.RUnlock()

	return i.Fighter2
}

func (i *SafeTradesMap) getFighterTrade(f *fighters.Fighter) *Trade {
	i.RLock()
	defer i.RUnlock()

	for _, trade := range i.Trades {
		if trade.gFighter1() == f || trade.gFighter2() == f {
			return trade
		}
	}

	return nil
}

func (i *SafeTradesMap) Delete(trade *Trade) {
	i.RLock()
	i.RUnlock()

	for index, t := range i.Trades {
		if t == trade {
			i.Trades = append(i.Trades[:index], i.Trades[index+1:]...)
			log.Printf("[SafeTradeMap.Delete] trade deleted index=%v", index)
			return
		}
	}
}


func StartTrade(fighter1 *fighters.Fighter, fighter2 *fighters.Fighter) (*Trade, error) {
	// Check if the players are not already trading
	if TradesMap.getFighterTrade(fighter1) != nil {
		return nil, errors.New("fighter1 is already trading")
	}

	if TradesMap.getFighterTrade(fighter2) != nil {
		return nil, errors.New("fighter2 is already trading")
	}

	// Initialize the trade
	trade := &Trade{
		Fighter1: fighter1,
		Fighter2: fighter2,
		TradeGrid1: inventory.NewInventory(8, 4),
		TradeGrid2: inventory.NewInventory(8, 4),
	}

	// Start the trade
	TradesMap.Lock()
	TradesMap.Trades = append(TradesMap.Trades, trade)
	TradesMap.Unlock()

	return trade, nil
}

func CancelTrade(fighter *fighters.Fighter) error {
	trade := TradesMap.getFighterTrade(fighter)

	if trade == nil {
		return errors.New("[CancelTrade] Fighter has no open trades")
	}

	// move items from trade to backpack

	// delete trade
	TradesMap.Delete(trade)
	return nil
}




// // AddItem allows a player to add an item to their trade grid.
// func (t *Trade) AddItem(player *Fighter, item *InventorySlot, x, y int) error {
// 	t.Lock()
// 	defer t.Unlock()

// 	// Check if the position and item size are in range.
// 	width := int(item.Attributes.ItemParameters.ItemWidth)
// 	height := int(item.Attributes.ItemParameters.ItemHeight)
// 	if x < 0 || x+width > gridWidth || y < 0 || y+height > gridHeight {
// 		return errors.New("position or item size out of range")
// 	}

// 	// Decide which grid to put the item in based on the player.
// 	var grid [][]bool
// 	var items map[string]*InventorySlot
// 	if player == t.Fighter1 {
// 		grid = t.ItemGrids[t.Fighter1.ID]
// 		items = t.Items[t.Fighter1.ID]
// 	} else if player == t.Fighter2 {
// 		grid = t.ItemGrids[t.Fighter2.ID]
// 		items = t.Items[t.Fighter2.ID]
// 	} else {
// 		return errors.New("player not in this trade")
// 	}

// 	// Check if the positions required by the item are already occupied.
// 	for i := 0; i < height; i++ {
// 		for j := 0; j < width; j++ {
// 			if grid[y+i][x+j] == true {
// 				return errors.New("some positions required by the item are already occupied")
// 			}
// 		}
// 	}

// 	// Add the item to the grid and items map.
// 	for i := 0; i < height; i++ {
// 		for j := 0; j < width; j++ {
// 			grid[y+i][x+j] = true
// 			items[fmt.Sprintf("%d,%d", x+j, y+i)] = item
// 		}
// 	}
// 	return nil
// }




// // RemoveItem allows a player to remove an item from their trade grid.
// func (t *Trade) RemoveItem(player *Fighter, x, y int) error {
// 	t.Lock()
// 	defer t.Unlock()

// 	// Decide which grid to remove the item from based on the player.
// 	var grid [][]bool
// 	var items map[string]*InventorySlot
// 	if player == t.Fighter1 {
// 		grid = t.ItemGrids[t.Fighter1.ID]
// 		items = t.Items[t.Fighter1.ID]
// 	} else if player == t.Fighter2 {
// 		grid = t.ItemGrids[t.Fighter2.ID]
// 		items = t.Items[t.Fighter2.ID]
// 	} else {
// 		return errors.New("player not in this trade")
// 	}

// 	// Check if there is an item at the given position.
// 	item, ok := items[fmt.Sprintf("%d,%d", x, y)]
// 	if !ok {
// 		return errors.New("no item in the specified position")
// 	}

// 	// Get item size.
// 	width := int(item.Attributes.ItemParameters.ItemWidth)
// 	height := int(item.Attributes.ItemParameters.ItemHeight)

// 	// Check if the item is in range.
// 	if x < 0 || x+width > gridWidth || y < 0 || y+height > gridHeight {
// 		return errors.New("item out of range")
// 	}

// 	// Remove the item from the grid and items map.
// 	for i := 0; i < height; i++ {
// 		for j := 0; j < width; j++ {
// 			grid[y+i][x+j] = false
// 			delete(items, fmt.Sprintf("%d,%d", x+j, y+i))
// 		}
// 	}

// 	return nil
// }

// // MoveItem allows a player to move an item inside the grid from one position to another.
// func (t *Trade) MoveItem(player *Fighter, oldX, oldY, newX, newY int) error {
// 	t.Lock()
// 	defer t.Unlock()

// 	// Decide which grid and items map to use based on the player.
// 	var grid [][]bool
// 	var items map[string]*InventorySlot
// 	if player == t.Fighter1 {
// 		grid = t.ItemGrids[t.Fighter1.ID]
// 		items = t.Items[t.Fighter1.ID]
// 	} else if player == t.Fighter2 {
// 		grid = t.ItemGrids[t.Fighter2.ID]
// 		items = t.Items[t.Fighter2.ID]
// 	} else {
// 		return errors.New("player not in this trade")
// 	}

// 	// Check if there is an item at the old position.
// 	item, ok := items[fmt.Sprintf("%d,%d", oldX, oldY)]
// 	if !ok {
// 		return errors.New("no item in the specified old position")
// 	}

// 	// Get item size.
// 	width := int(item.Attributes.ItemParameters.ItemWidth)
// 	height := int(item.Attributes.ItemParameters.ItemHeight)

// 	// Check if the old and new positions and item size are in range.
// 	if oldX < 0 || oldX+width > gridWidth || oldY < 0 || oldY+height > gridHeight ||
// 		newX < 0 || newX+width > gridWidth || newY < 0 || newY+height > gridHeight {
// 		return errors.New("old position, new position, or item size out of range")
// 	}

// 	// Check if the positions required by the item at the new location are unoccupied.
// 	for i := 0; i < height; i++ {
// 		for j := 0; j < width; j++ {
// 			// Skip checking if old and new positions overlap.
// 			if oldX <= newX+j && newX+j < oldX+width && oldY <= newY+i && newY+i < oldY+height {
// 				continue
// 			}
// 			if grid[newY+i][newX+j] == true {
// 				return errors.New("some positions required by the item at the new location are already occupied")
// 			}
// 		}
// 	}

// 	// Remove the item from the old position and add it at the new position.
// 	for i := 0; i < height; i++ {
// 		for j := 0; j < width; j++ {
// 			grid[oldY+i][oldX+j] = false
// 			delete(items, fmt.Sprintf("%d,%d", oldX+j, oldY+i))
// 			grid[newY+i][newX+j] = true
// 			items[fmt.Sprintf("%d,%d", newX+j, newY+i)] = item
// 		}
// 	}

// 	return nil
// }

// func (t *Trade) AcceptTrade() error {
// 	t.Lock()
// 	defer t.Unlock()

// 	// Check if both players are in the trade
// 	if _, ok := Trades[t.Fighter1.ID]; !ok {
// 		return errors.New("fighter1 is not in a trade")
// 	}
// 	if _, ok := Trades[t.Fighter2.ID]; !ok {
// 		return errors.New("fighter2 is not in a trade")
// 	}

// 	// Exchange items between the two players
// 	exchangeItems(t)
// 	exchangeItems(t)

// 	// Clean up the trade
// 	delete(Trades, t.Fighter1.ID)
// 	delete(Trades, t.Fighter2.ID)

// 	return nil
// }

// func exchangeItems(t *Trade) {
// 	for _, item := range t.Items[t.Fighter1.ID] {
// 		// Remove the item from the from's Backpack
// 		t.Fighter1.Backpack.removeItemByHash(t.Fighter1, item.ItemHash)

// 		// Add the item to the to's Backpack
// 		t.Fighter2.Backpack.AddItem(item.Attributes, item.Qty, item.ItemHash)
// 	}

// 	for _, item := range t.Items[t.Fighter2.ID] {
// 		// Remove the item from the from's Backpack
// 		t.Fighter2.Backpack.removeItemByHash(t.Fighter2, item.ItemHash)

// 		// Add the item to the to's Backpack
// 		t.Fighter1.Backpack.AddItem(item.Attributes, item.Qty, item.ItemHash)
// 	}
// }



// TODO: Add more methods for trading (e.g., AcceptTrade, CancelTrade, etc.).
