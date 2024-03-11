// equipment.go

package inventory

import (
	"sync"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type Equipment struct {
	Map map[int64]*InventorySlot
	sync.RWMutex
}

func (i *Equipment) GetMap() map[int64]*InventorySlot {
	i.RLock()
	i.RUnlock()

	return i.Map
}

func (i *Equipment) SetMap(v map[int64]*InventorySlot)  {
	i.Lock()
	defer i.Unlock()

	i.Map = v
}

func (i *Equipment) Find (k int64) *InventorySlot {
	i.RLock()
	defer i.RUnlock()

	return i.Map[k]
}

func (i *Equipment) Dress (slotId int64, item *InventorySlot) {
	i.Lock()
	defer i.Unlock()

	i.Map[slotId] = item
}

func (i *Equipment) RemoveByHash (hash common.Hash) {
	i.Lock()
	defer i.Unlock()

	for slotId, slot := range i.Map {
		if slot.GetItemHash() == hash {
			delete(i.Map, slotId)
			return
		}
	}
}

func (i *Equipment) FindByHash (hash common.Hash) *InventorySlot {
	i.RLock()
	defer i.RUnlock()

	for _, slot := range i.Map {
		if slot.GetItemHash() == hash {
			return slot
		}
	}

	return nil
}

func (e *Equipment) MarshalJSON() ([]byte, error) {
    e.RLock()
    defer e.RUnlock()

    // Marshal the Map field directly
    return json.Marshal(e.Map)
}

func NewEquipment() *Equipment {
	return &Equipment{Map: make(map[int64]*InventorySlot)}
}







