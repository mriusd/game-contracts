// population.go

package main

import (
    "sync"
    "log"

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
    log.Printf("[PopulationMap] Adding %v to %v", fighter.Name, town)
	i.Lock()
	defer i.Unlock()


    log.Printf("[PopulationMap] 1 Adding %v to %v", fighter.Name, town)
	if i.Map[town] == nil {
        i.Map[town] = make([]*fighters.Fighter, 0)
    }

    log.Printf("[PopulationMap] 2 Adding %v to %v", fighter.Name, town)

    i.Map[town] = append(i.Map[town], fighter)
    log.Printf("[PopulationMap] Fighter %v added to %v", fighter.Name, town)
}

func (i *SafePopultationMap) Remove(fighter *fighters.Fighter) {
    name := fighter.GetName()
    log.Printf("[PopulationMap] Removing %v ", name)
    town := fighter.GetLocation()

	i.Lock()
	defer i.Unlock()

    

    pop, exists := i.Map[town]
    if !exists {
        return // Town does not exist, so fighter cannot be removed
    }

	for index, f := range pop {
        if f == fighter {
            // Remove the fighter from the slice.
            i.Map[town] = append(i.Map[town][:index], i.Map[town][index+1:]...)
            return
        }
    }

    log.Printf("[PopulationMap] Removed %v ", name)
}

func (i *SafePopultationMap) Find(town string, v string) *fighters.Fighter {
    for _, fighter := range i.GetTownMap(town) {
        if fighter.GetID() == v {
            return fighter
        }
    }
    return nil
}

var PopulationMap = &SafePopultationMap{Map: make(map[string][]*fighters.Fighter)}