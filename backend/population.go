// population.go

package main

import (
    "sync"

    "github.com/mriusd/game-contracts/fighters"
)


type SafePopultationMap struct {
	Map map[string][]*fighters.Fighter
	sync.RWMutex
}

func (i *SafePopultationMap) GetTownMap(town string) []*fighters.Fighter {
    i.RLock()
    defer i.RUnlock()

    if pop, exists := i.Map[town]; exists {
        copy := make([]*fighters.Fighter, 0, len(pop))
        for _, fighter := range pop {
            copy = append(copy, fighter)
        }
        return copy
    }
    return nil
}

func (i *SafePopultationMap) Add(town string, fighter *fighters.Fighter) {
	i.Lock()
	defer i.Unlock()

	if i.Map[town] == nil {
        i.Map[town] = make([]*fighters.Fighter, 0)
    }

    i.Map[town] = append(i.Map[town], fighter)
}

func (i *SafePopultationMap) Remove(fighter *fighters.Fighter) {
	i.Lock()
	defer i.Unlock()

	for town, fighters := range i.Map {
        for index, f := range fighters {
            if f == fighter {
                // Remove the fighter from the slice.
                i.Map[town] = append(fighters[:index], fighters[index+1:]...)
                return
            }
        }
    }
}

func (i *SafePopultationMap) Find(town string, v int64) *fighters.Fighter {
    for id, fighter := range i.GetTownMap(town) {
        if int64(id) == v {
            return fighter
        }
    }
    return nil
}

var PopulationMap = &SafePopultationMap{Map: make(map[string][]*fighters.Fighter)}