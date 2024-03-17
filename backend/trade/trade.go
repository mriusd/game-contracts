package trade

import (
	"sync"
	"errors"

	"github.com/mriusd/game-contracts/inventory"
	"github.com/mriusd/game-contracts/fighters"
)

type Trade struct {	
	Fighter1 *fighters.Fighter
	Fighter2 *fighters.Fighter
	Inventory1 	*inventory.Inventory `json:"inventory_1"`
	Inventory2  *inventory.Inventory `json:"inventory_2"`
	Approve1 	bool `json:"approve_1"`
	Approve2 	bool `json:"approve_2"`

	sync.RWMutex `json:"-"`
}


type SafeTradesMap struct {
	Map []*Trade

	sync.RWMutex
}

var TradesMap = &SafeTradesMap{Map: make([]*Trade, 0)}

func (i *SafeTradesMap) Add(v *Trade) {
	i.Lock()
	defer i.Unlock()

	i.Map = append(i.Map, v)
}

func (i *SafeTradesMap) RemoveItem(v *Trade) {
    i.Lock()
    defer i.Unlock()

    for index, item := range i.Map {
        if item == v {
            i.Map = append(i.Map[:index], i.Map[index+1:]...)
            break
        }
    }
}

func (i *SafeTradesMap) FindByFighter (v *fighters.Fighter) *Trade {
	i.RLock()
	defer i.RUnlock()

	for _, trade := range i.Map {
		if trade.GetInventory1().GetOwnerId() == v.GetTokenID() || trade.GetInventory2().GetOwnerId() == v.GetTokenID() {
			return trade
		}
	}

	return nil
}


func (i *Trade) GetFighter1() *fighters.Fighter {
	i.RLock()
	defer i.RUnlock()

	return i.Fighter1
}

func (i *Trade) GetFighter2() *fighters.Fighter {
	i.RLock()
	defer i.RUnlock()

	return i.Fighter2
}



func (i *Trade) GetInventory1() *inventory.Inventory {
	i.RLock()
	defer i.RUnlock()

	return i.Inventory1
}

func (i *Trade) GetInventory2() *inventory.Inventory {
	i.RLock()
	defer i.RUnlock()

	return i.Inventory2
}


func (i *Trade) GetApprove1() bool {
	i.RLock()
	defer i.RUnlock()

	return i.Approve1
}

func (i *Trade) GetApprove2() bool {
	i.RLock()
	defer i.RUnlock()

	return i.Approve2
}

func (i *Trade) SetApprove1(v bool) {
	i.Lock()
	defer i.Unlock()

	i.Approve1 = v
}

func (i *Trade) SetApprove2(v bool) {
	i.Lock()
	defer i.Unlock()

	i.Approve2 = v
}


func Initiate(fighter1, fighter2 *fighters.Fighter) (*Trade, error) {
	if TradesMap.FindByFighter(fighter1) != nil {
		return nil, errors.New("Player already in trade")
	}

	if TradesMap.FindByFighter(fighter2) != nil {
		return nil, errors.New("Player already in trade")
	}

	trade := &Trade{
		Fighter1: fighter1,
		Fighter2: fighter2,
		Inventory1: inventory.NewInventory(fighter1.GetTokenID(), "trade"),
		Inventory2: inventory.NewInventory(fighter2.GetTokenID(), "trade"),
	}

	TradesMap.Add(trade)

	return trade, nil
}

func SetGold(fighter *fighters.Fighter, amount int) error {
	trade := TradesMap.FindByFighter(fighter)

	if trade == nil {
		return errors.New("No open trades")
	}

	backpack := fighter.GetBackpack()
	if backpack == nil {
		return errors.New("Backpack not found")
	}

	if backpack.GetGold() < amount {
		return errors.New("Not enough Gold")
	}


	if trade.GetFighter1() == fighter {
		trade.GetInventory1().SetGold(amount)
	}

	if trade.GetFighter2() == fighter {
		trade.GetInventory2().SetGold(amount)		
	}

	return nil
}

func AddItem(fighter *fighters.Fighter, itemHash string) error {
	trade := TradesMap.FindByFighter(fighter)
	if trade == nil {
		return errors.New("No open trades")
	}

	backpack := fighter.GetBackpack()
	equipment := fighter.GetEquipment()
	if backpack == nil || equipment == nil {
		return errors.New("Backpack/Equipment not found")
	}

	item := backpack.FindByHash(itemHash)
	if item == nil {

		item = equipment.FindByHash(itemHash)

		if item == nil {
			return errors.New("Item not found on player")
		}		
	}

	item.SetInTrade(true)

	if trade.GetFighter1() == fighter {
		_, _, err := trade.GetInventory1().AddItem(item.GetAttributes(), item.GetQty(), itemHash)
		return err
	}

	if trade.GetFighter2() == fighter {
		_, _, err := trade.GetInventory2().AddItem(item.GetAttributes(), item.GetQty(), itemHash)	
		return err	
	}

	return nil
}


func RemoveItem(fighter *fighters.Fighter, itemHash string) error {
	trade := TradesMap.FindByFighter(fighter)
	if trade == nil {
		return errors.New("No open trades")
	}

	backpack := fighter.GetBackpack()
	equipment := fighter.GetEquipment()
	if backpack == nil || equipment == nil {
		return errors.New("Backpack/Equipment not found")
	}

	item := backpack.FindByHash(itemHash)
	if item == nil {
		item = equipment.FindByHash(itemHash)
		if item == nil {
			return errors.New("Item not found on player")
		}		
	}

	item.SetInTrade(false)

	if trade.GetFighter1() == fighter {
		trade.SetApprove2(false)
		ok := trade.GetInventory1().RemoveItemByHash(itemHash)
		if !ok {
			return errors.New("Failed to remove item")
		}		
	}

	if trade.GetFighter2() == fighter {
		trade.SetApprove1(false)
		ok := trade.GetInventory2().RemoveItemByHash(itemHash)
		if !ok {
			return errors.New("Failed to remove item")
		}
	}

	return nil
}


func Approve(fighter *fighters.Fighter) error {
	trade := TradesMap.FindByFighter(fighter)
	if trade == nil {
		return errors.New("No open trades")
	}

	backpack := fighter.GetBackpack()
	equipment := fighter.GetEquipment()
	if backpack == nil || equipment == nil {
		return errors.New("Backpack/Equipment not found")
	}

	if trade.GetFighter1() == fighter {
		trade.SetApprove1(true)
	}

	if trade.GetFighter2() == fighter {
		trade.SetApprove2(true)
	}

	if trade.GetApprove1() && trade.GetApprove2() {
		return Execute(trade)
	}

	return nil
}

func (i *Trade) Cancel() {
	tradeGrid1 := i.GetInventory1()
	tradeGrid2 := i.GetInventory2()

	for _, inventorySlot := range tradeGrid1.GetItems() {
		inventorySlot.SetInTrade(false)
	}

	for _, inventorySlot := range tradeGrid2.GetItems() {
		inventorySlot.SetInTrade(false)
	}

	TradesMap.RemoveItem(i)
}

func Execute(trade *Trade) error {
	if !trade.GetApprove1() || !trade.GetApprove2() {
		return errors.New("Trade not confirmed")
	}

	fighter1 := trade.GetFighter1()
	fighter2 := trade.GetFighter2()

	inventory1 := trade.GetInventory1()
	inventory2 := trade.GetInventory2()

	// check wnough space in player1 backpack
	backpack1 := fighter1.GetBackpack()
	if !backpack1.IsEnoughSpaceForMultipleItems(inventory2.GetItems()) {
		trade.Cancel()
		return errors.New("Not enough space")
	}

	backpack2 := fighter2.GetBackpack()
	if !backpack2.IsEnoughSpaceForMultipleItems(inventory1.GetItems()) {
		trade.Cancel()
		return errors.New("Not enough space")
	}

	equipment1 := fighter1.GetEquipment()
	equipment2 := fighter2.GetEquipment()

	

	// add items to fighter1
	for itemHash, itemSlot := range inventory2.GetItems() {
		item := backpack2.FindByHash(itemHash)
		if item != nil {
			backpack2.RemoveItemByHash(itemHash)
			
		} else {
			item = equipment2.FindByHash(itemHash)

			if item == nil {
				return errors.New("Item not found on player")
			}

			equipment2.RemoveByHash(itemHash)
		}

		backpack1.AddItem(itemSlot.GetAttributes(), itemSlot.GetQty(), itemHash)
		itemSlot.SetInTrade(false)
	}

	backpack1.SetGold(backpack1.GetGold() + inventory2.GetGold())
	backpack2.SetGold(backpack2.GetGold() - inventory2.GetGold())


	// add items to fighter2
	for itemHash, itemSlot := range inventory1.GetItems() {
		item := backpack1.FindByHash(itemHash)
		if item != nil {
			backpack1.RemoveItemByHash(itemHash)
			
		} else {
			item = equipment1.FindByHash(itemHash)

			if item == nil {
				return errors.New("Item not found on player")
			}

			equipment1.RemoveByHash(itemHash)
		}

		backpack2.AddItem(itemSlot.GetAttributes(), itemSlot.GetQty(), itemHash)
		itemSlot.SetInTrade(false)
	}

	backpack2.SetGold(backpack2.GetGold() + inventory1.GetGold())
	backpack1.SetGold(backpack1.GetGold() - inventory1.GetGold())


	TradesMap.RemoveItem(trade)

	return nil
}







