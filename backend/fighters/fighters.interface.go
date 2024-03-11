// fighterModel.go

package fighters

import (
	"sync"	
    "math/big"
    "github.com/ethereum/go-ethereum/common"

    "github.com/mriusd/game-contracts/battle"
    "github.com/mriusd/game-contracts/maps"
    "github.com/mriusd/game-contracts/inventory"
    "github.com/mriusd/game-contracts/skill"
)

type Fighter struct {
	ID    					string  		    `json:"id" bson:"-"`
    Class                   string              `json:"class" bson:"class"`
    MaxHealth     			int64 			    `json:"maxHealth" bson:"-"`
    Name           			string 			    `json:"name" bson:"name"`
    IsNpc         			bool 			    `json:"isNpc" bson:"-"`
    IsDead         			bool 			    `json:"isDead" bson:"-"`
    CanFight 				bool 			    `json:"canFight" bson:"-"`
    LastDmgTimestamp 		int64  			    `json:"lastDmgTimestamp" bson:"lastDmgTimestamp"`
    HealthAfterLastDmg 		int64  			    `json:"healthAfterLastDmg" bson:"healthAfterLastDmg"`

    TokenID                 int64               `json:"tokenId" bson:"tokenId"`
    Location                string              `json:"location" bson:"location"`
    
    DamageReceived          []battle.Damage     `json:"damageDealt" bson:"-"`
    OwnerAddress            string              `json:"ownerAddress" bson:"owner_address"`
    Coordinates             maps.Coordinate     `json:"coordinates" bson:"coordinates"`
    MovementSpeed           int64               `json:"movementSpeed" bson:"-"` // squares per minute
    Skill                   int64               `json:"skill" bson:"skill"`
    SpawnCoords             maps.Coordinate     `json:"spawnCoords" bson:"-"`
    

    LastMoveTimestamp       int64               `json:"lastMoveTimestamp" bson:"-"` // milliseconds


    // Fighter stats
    Strength                int64               `json:"strength" bson:"strength"`
    Agility                 int64               `json:"agility" bson:"agility"`
    Energy                  int64               `json:"energy" bson:"energy"`
    Vitality                int64               `json:"vitality" bson:"vitality"`


    // Fighter dynamic paramters
    CurrentHealth           int64               `json:"currentHealth" bson:"-"`
    CurrentMana             int64               `json:"currentMana" bson:"-"`


    // Fighter parameters with equipped items
    Damage                  int64               `json:"damage" bson:"-"`
    Defence                 int64               `json:"defence" bson:"-"`
    AttackSpeed             int64               `json:"attackSpeed" bson:"-"` 
    HpRegenerationRate      float64             `json:"hpRegenerationRate" bson:"-"`
    HpRegenerationBonus     float64             `json:"hpRegenerationBonus" bson:"-"`

    // Damage type rates
    CriticalDmgRate         int64               `json:"criticalDmgRate" bson:"-"`
    ExcellentDmgRate        int64               `json:"excellentDmgRate" bson:"-"`
    DoubleDmgRate           int64               `json:"doubleDmgRate" bson:"-"`
    IgnoreDefRate           int64               `json:"ignoreDefRate" bson:"-"`


    Level                   int64               `json:"level" bson:"-"`
    Experience              int64               `json:"experience" bson:"experience"`

    Direction               maps.Direction           `json:"direction" bson:"-"`

    Skills                  map[int64]skill.Skill    `json:"skills" bson:"skills"`
    Backpack                *inventory.Inventory           `json:"-" bson:"-"`
    Vault                   *inventory.Inventory           `json:"-" bson:"-"`
    Equipment               *inventory.Equipment          `json:"equipment" bson:"-"`

    LastChatMsg             string              `json:"lastChatMessage" bson:"-"`
    LastChatMsgTimestamp    int64               `json:"lastChatMsgTimestamp" bson:"-"`

    Credits                 int64               `json:"credits" bson:"-"`

    sync.RWMutex                                `json:"-" bson:"-"`
}

func (i *Fighter) GetVault() *inventory.Inventory {
    i.RLock()
    defer i.RUnlock()

    return i.Vault
}

func (i *Fighter) GetBackpack() *inventory.Inventory {
    i.RLock()
    defer i.RUnlock()

    return i.Backpack
}


func (i *Fighter) GetName() string {
    i.RLock()
    defer i.RUnlock()

    return i.Name
}


func (i *Fighter) GetOwnerAddress() string {
    i.RLock()
    defer i.RUnlock()

    return i.OwnerAddress
}

func (i *Fighter) GetID() string {
    i.RLock()
    defer i.RUnlock()

    return i.ID
}

func (i *Fighter) GetCoordinates() maps.Coordinate {
    i.RLock()
    defer i.RUnlock()

    return i.Coordinates
}

func (i *Fighter) GetDamageReceived() []battle.Damage {
    i.RLock()
    defer i.RUnlock()

    return i.DamageReceived
}

func (i *Fighter) GetHealthAfterLastDmg() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.HealthAfterLastDmg
}

func (i *Fighter) GetLastChatMsgTimestamp() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.LastChatMsgTimestamp
}

func (i *Fighter) GetHpRegenerationRate() float64 {
    i.RLock()
    defer i.RUnlock()

    return i.HpRegenerationRate
}

func (i *Fighter) GetMovementSpeed() int64 {
    i.RLock()
    i.RUnlock()

    return i.MovementSpeed
}

func (i *Fighter) GetLastDmgTimestamp() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.LastDmgTimestamp
}

func (i *Fighter) GetMaxHealth() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.MaxHealth
}

func (i *Fighter) GetDefence() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.Defence
}

func (i *Fighter) GetDamage() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.Damage
}

func (i *Fighter) GetStrength() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.Strength
}

func (i *Fighter) GetAgility() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.Agility
}

func (i *Fighter) GetEnergy() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.Energy
}

func (i *Fighter) GetVitality() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.Vitality
}

func (i *Fighter) GetIgnoreDefRate() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.IgnoreDefRate
}

func (i *Fighter) GetExcellentDmgRate() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.ExcellentDmgRate
}

func (i *Fighter) GetCriticalDmgRate() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.CriticalDmgRate
}

func (i *Fighter) GetDoubleDmgRate() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.DoubleDmgRate
}

func (i *Fighter) GetSkill() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.Skill
}

func (i *Fighter) GetTokenID() int64 {
    i.RLock()
    defer i.RUnlock()

    return i.TokenID
}

func (i *Fighter) GetLocation() string {
    i.RLock()
    defer i.RUnlock()

    return i.Location
}

func (i *Fighter) GetIsDead() bool {
    i.RLock()
    defer i.RUnlock()

    return i.IsDead
}

func (i *Fighter) GetIsNpc() bool {
    i.RLock()
    defer i.RUnlock()

    return i.IsNpc
}


func (i *Fighter) GetSpawnCoords() maps.Coordinate {
    i.RLock()
    defer i.RUnlock()

    return i.SpawnCoords
}

func (i *Fighter) GetEquipment() *inventory.Equipment {
    i.RLock()
    defer i.RUnlock()

    return i.Equipment
}

func (i *Fighter) SetDirection(v maps.Direction) {
    i.Lock()
    defer i.Unlock()

    i.Direction = v
}

func (i *Fighter) SetCoordinates(v maps.Coordinate) {
    i.Lock()
    defer i.Unlock()

    i.Coordinates = v
}

func (i *Fighter) SetCurrentHealth(v int64) {
    i.Lock()
    defer i.Unlock()

    i.CurrentHealth = v
}

func (i *Fighter) SetIsDead(v bool) {
    i.Lock()
    defer i.Unlock()

    i.IsDead = v
}

func (i *Fighter) SetLastDmgTimestamp(v int64) {
    i.Lock()
    defer i.Unlock()

    i.LastDmgTimestamp = v
}

func (i *Fighter) SetHealthAfterLastDmg(v int64) {
    i.Lock()
    defer i.Unlock()

    i.HealthAfterLastDmg = v
}

func (i *Fighter) SetLastChatMsg(v string) {
    i.Lock()
    defer i.Unlock()

    i.LastChatMsg = v
}

func (i *Fighter) SetEquipment(v *inventory.Equipment) {
    i.Lock()
    defer i.Unlock()

    i.Equipment = v
}


// type SafeFightersMap struct {
//     Map map[string]*Fighter
//     sync.RWMutex
// }

// var FightersMap = &SafeFightersMap{Map: make(map[string]*Fighter)}


// func (i *SafeFightersMap) GetMap() map[string]*Fighter {
//     i.RLock()
//     defer i.RUnlock()

//     copy := make(map[string]*Fighter, len(i.Map))
//     for key, val := range i.Map {
//         copy[key] = val
//     }
//     return copy
// }


// func (i *SafeFightersMap) Find(id string) *Fighter {
//     i.RLock()
//     defer i.RUnlock()

//     fighter, exists := FightersMap.Map[id]
//     if exists {
//         return fighter
//     }

//     return nil
// }

// func (i *SafeFightersMap) Add(id string, f *Fighter) {
//     i.Lock()
//     defer i.Unlock()

//     FightersMap.Map[id] = f
// }




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

func (i *FighterAttributes) GetVitality() *big.Int {
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




















