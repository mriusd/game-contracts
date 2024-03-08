// fighterModel.go

package main

import (
	"sync"	
    "math/big"
    "github.com/ethereum/go-ethereum/common"
)

type Fighter struct {
	ID    					string  		    `json:"id"`
    MaxHealth     			int64 			    `json:"maxHealth"`
    Name           			string 			    `json:"name"`
    IsNpc         			bool 			    `json:"isNpc"`
    IsDead         			bool 			    `json:"isDead"`
    CanFight 				bool 			    `json:"canFight"`
    LastDmgTimestamp 		int64  			    `json:"lastDmgTimestamp"`
    HealthAfterLastDmg 		int64  			    `json:"healthAfterLastDmg"`

    TokenID                 int64               `json:"tokenId"`
    BirthBlock              int64               `json:"birthBlock"`
    Location                string              `json:"location"`
    
    DamageReceived          []Damage            `json:"damageDealt"`
    OwnerAddress            string              `json:"ownerAddress"`
    Coordinates             Coordinate          `json:"coordinates"`
    MovementSpeed           int64               `json:"movementSpeed"` // squares per minute
    Skill                   int64               `json:"skill"`
    SpawnCoords             Coordinate          `json:"spawnCoords"`
    

    LastMoveTimestamp       int64               `json:"lastMoveTimestamp"` // milliseconds


    // Fighter stats
    Strength                int64               `json:"strength"`
    Agility                 int64               `json:"agility"`
    Energy                  int64               `json:"energy"`
    Vitality                int64               `json:"vitality"`


    // Fighter dynamic paramters
    CurrentHealth           int64               `json:"currentHealth"`
    CurrentMana             int64               `json:"currentMana"`


    // Fighter parameters with equipped items
    Damage                  int64               `json:"damage"`
    Defence                 int64               `json:"defence"`
    AttackSpeed             int64               `json:"attackSpeed"` 
    HpRegenerationRate      float64             `json:"hpRegenerationRate"`
    HpRegenerationBonus     float64             `json:"hpRegenerationBonus"`

    // Damage type rates
    CriticalDmgRate         int64               `json:"criticalDmgRate"`
    ExcellentDmgRate        int64               `json:"excellentDmgRate"`
    DoubleDmgRate           int64               `json:"doubleDmgRate"`
    IgnoreDefRate           int64               `json:"ignoreDefRate"`


    Level                   int64               `json:"level"`
    Experience              int64               `json:"experience"`

    Direction               Direction           `json:"direction"`

    Skills                  map[int64]*Skill    `json:"skills"`
    Backpack                *Inventory           `json:"-"`
    Vault                   *Inventory           `json:"-"`
    Equipment               map[int64]*InventorySlot `json:"equipment"`

    LastChatMsg             string              `json:"lastChatMessage"`
    LastChatMsgTimestamp    int64               `json:"lastChatMsgTimestamp"`

    Credits                 int64               `json:"credits"`

    sync.RWMutex
}

func (i *Fighter) gName() string {
    i.RLock()
    i.RUnlock()

    return i.Name
}

func (i *Fighter) gID() string {
    i.RLock()
    i.RUnlock()

    return i.ID
}

func (i *Fighter) gCoordinates() Coordinate {
    i.RLock()
    i.RUnlock()

    return i.Coordinates
}

func (i *Fighter) gDamageReceived() []Damage {
    i.RLock()
    i.RUnlock()

    return i.DamageReceived
}

func (i *Fighter) gHealthAfterLastDmg() int64 {
    i.RLock()
    i.RUnlock()

    return i.HealthAfterLastDmg
}

func (i *Fighter) gHpRegenerationRate() float64 {
    i.RLock()
    i.RUnlock()

    return i.HpRegenerationRate
}

func (i *Fighter) gMovementSpeed() int64 {
    i.RLock()
    i.RUnlock()

    return i.MovementSpeed
}

func (i *Fighter) gLastDmgTimestamp() int64 {
    i.RLock()
    i.RUnlock()

    return i.LastDmgTimestamp
}

func (i *Fighter) gMaxHealth() int64 {
    i.RLock()
    i.RUnlock()

    return i.MaxHealth
}

func (i *Fighter) gDefence() int64 {
    i.RLock()
    i.RUnlock()

    return i.Defence
}

func (i *Fighter) gDamage() int64 {
    i.RLock()
    i.RUnlock()

    return i.Damage
}

func (i *Fighter) gIgnoreDefRate() int64 {
    i.RLock()
    i.RUnlock()

    return i.IgnoreDefRate
}

func (i *Fighter) gExcellentDmgRate() int64 {
    i.RLock()
    i.RUnlock()

    return i.ExcellentDmgRate
}

func (i *Fighter) gCriticalDmgRate() int64 {
    i.RLock()
    i.RUnlock()

    return i.CriticalDmgRate
}

func (i *Fighter) gDoubleDmgRate() int64 {
    i.RLock()
    i.RUnlock()

    return i.DoubleDmgRate
}

func (i *Fighter) gSkill() int64 {
    i.RLock()
    i.RUnlock()

    return i.Skill
}

func (i *Fighter) gTokenID() int64 {
    i.RLock()
    i.RUnlock()

    return i.TokenID
}

func (i *Fighter) gLocation() string {
    i.RLock()
    i.RUnlock()

    return i.Location
}

func (i *Fighter) gIsDead() bool {
    i.RLock()
    i.RUnlock()

    return i.IsDead
}

func (i *Fighter) gIsNpc() bool {
    i.RLock()
    i.RUnlock()

    return i.IsNpc
}


func (i *Fighter) gSpawnCoords() Coordinate {
    i.RLock()
    i.RUnlock()

    return i.SpawnCoords
}

func (i *Fighter) sDirection(v Direction) {
    i.Lock()
    defer i.Unlock()

    i.Direction = v
}

func (i *Fighter) sCoordinates(v Coordinate) {
    i.Lock()
    defer i.Unlock()

    i.Coordinates = v
}

func (i *Fighter) sCurrentHealth(v int64) {
    i.Lock()
    defer i.Unlock()

    i.CurrentHealth = v
}

func (i *Fighter) sIsDead(v bool) {
    i.Lock()
    defer i.Unlock()

    i.IsDead = v
}

func (i *Fighter) sLastDmgTimestamp(v int64) {
    i.Lock()
    defer i.Unlock()

    i.LastDmgTimestamp = v
}

func (i *Fighter) sHealthAfterLastDmg(v int64) {
    i.Lock()
    defer i.Unlock()

    i.HealthAfterLastDmg = v
}


type SafeFightersMap struct {
    Map map[string]*Fighter
    sync.RWMutex
}

var FightersMap = &SafeFightersMap{Map: make(map[string]*Fighter)}


func (i *SafeFightersMap) gMap() map[string]*Fighter {
    i.RLock()
    defer i.RUnlock()

    copy := make(map[string]*Fighter, len(i.Map))
    for key, val := range i.Map {
        copy[key] = val
    }
    return copy
}


func (i *SafeFightersMap) Find(id string) *Fighter {
    i.RLock()
    defer i.RUnlock()

    fighter, exists := FightersMap.Map[id]
    if exists {
        return fighter
    }

    return nil
}

func (i *SafeFightersMap) Add(id string, f *Fighter) {
    i.Lock()
    defer i.Unlock()

    FightersMap.Map[id] = f
}




type FighterAttributes struct {
    Name                    string `json:"Name"`
    Class                   string `json:"Class"`
    TokenID                 *big.Int `json:"TokenID"`
    BirthBlock              *big.Int `json:"BirthBlock"`
    Strength                *big.Int `json:"Strength"`
    Agility                 *big.Int `json:"Agility"`
    Energy                  *big.Int `json:"Energy"`
    Vitality                *big.Int `json:"Vitality"`
    Experience              *big.Int `json:"Experience"`
    HpPerVitalityPoint      *big.Int `json:"HpPerVitalityPoint"`
    ManaPerEnergyPoint      *big.Int `json:"ManaPerEnergyPoint"`
    HpIncreasePerLevel      *big.Int `json:"HpIncreasePerLevel"`
    ManaIncreasePerLevel    *big.Int `json:"manaIncreasePerLevel"`
    StatPointsPerLevel      *big.Int `json:"statPointsPerLevel"`
    AttackSpeed             *big.Int `json:"attackSpeed"`
    AgilityPointsPerSpeed   *big.Int `json:"agilityPointsPerSpeed"`
    IsNpc                   *big.Int `json:"isNpc"`
    DropRarityLevel         *big.Int `json:"dropRarityLevel"`

    sync.RWMutex
}

func (i *FighterAttributes) gVitality() *big.Int {
    i.RLock()
    defer i.RUnlock()

    return new(big.Int).Set(i.Vitality)
}

type SafeFightersAttributesCacheMap struct {
    Map map[int64]*FighterAttributes
    sync.RWMutex

}

func (i *SafeFightersAttributesCacheMap) Find(k int64) *FighterAttributes {
    i.RLock()
    defer i.RUnlock()

    v, exists := i.Map[k]
    if !exists {
        return nil
    }

    return v
}

func (i *SafeFightersAttributesCacheMap) Add(k int64, v *FighterAttributes) {
    i.Lock()
    defer i.Unlock()

    i.Map[k] = v
}

var FighterAttributesCache = &SafeFightersAttributesCacheMap{Map: make(map[int64]*FighterAttributes)}

type FighterStats struct {
    TokenID                 *big.Int `json:"TokenID"`
    MaxHealth               *big.Int `json:"maxHealth"`
    MaxMana                 *big.Int `json:"maxMana"`
    Level                   *big.Int `json:"level"`
    Exp                     *big.Int `json:"exp"`
    TotalStatPoints         *big.Int `json:"totalStatPoints"`
    MaxStatPoints           *big.Int `json:"maxStatPoints"`
}

type FighterCreatedEvent struct {
    TokenID         *big.Int            `json:"tokenId"`
    Owner           common.Address      `json:"owner"`
    
    FighterClass    string            `json:"fighterClass"`
    Name            string              `json:"name"`
}




















