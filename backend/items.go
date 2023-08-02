package main

import (
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"encoding/json"
	"io/ioutil"

	"fmt"
	"crypto/sha256"
    "encoding/hex"

    "strings"
)



// type ItemAttributes struct {
// 	Name                                      string   `json:"name" bson:"name"`
// 	TokenId                                   *big.Int `json:"tokenId" bson:"tokenId"`
// 	ItemLevel                                 *big.Int `json:"itemLevel" bson:"itemLevel"`
// 	MaxLevel                                  *big.Int `json:"maxLevel" bson:"maxLevel"`
// 	Durability                                *big.Int `json:"durability" bson:"durability"`
// 	ClassRequired                             *big.Int `json:"classRequired" bson:"classRequired"`
// 	StrengthRequired                          *big.Int `json:"strengthRequired" bson:"strengthRequired"`
// 	AgilityRequired                           *big.Int `json:"agilityRequired" bson:"agilityRequired"`
// 	EnergyRequired                            *big.Int `json:"energyRequired" bson:"energyRequired"`
// 	VitalityRequired                          *big.Int `json:"vitalityRequired" bson:"vitalityRequired"`
// 	ItemWidth                                 *big.Int `json:"itemWidth" bson:"itemWidth"`
// 	ItemHeight                                *big.Int `json:"itemHeight" bson:"itemHeight"`
// 	AcceptableSlot1                           *big.Int `json:"acceptableSlot1" bson:"acceptableSlot1"`
// 	AcceptableSlot2                           *big.Int `json:"acceptableSlot2" bson:"acceptableSlot2"`

// 	PhysicalDamage                            *big.Int `json:"physicalDamage" bson:"physicalDamage"`
// 	MagicDamage                               *big.Int `json:"magicDamage" bson:"magicDamage"`
// 	Defense                                   *big.Int `json:"defense" bson:"defense"`
// 	AttackSpeed                               *big.Int `json:"attackSpeed" bson:"attackSpeed"`
// 	DefenseSuccessRate                        *big.Int `json:"defenseSuccessRate" bson:"defenseSuccessRate"`
// 	AdditionalDamage                          *big.Int `json:"additionalDamage" bson:"additionalDamage"`
// 	AdditionalDefense                         *big.Int `json:"additionalDefense" bson:"additionalDefense"`
// 	IncreasedExperienceGain                   *big.Int `json:"increasedExperienceGain" bson:"increasedExperienceGain"`
// 	DamageIncrease                            *big.Int `json:"damageIncrease" bson:"damageIncrease"`
// 	DefenseSuccessRateIncrease                *big.Int `json:"defenseSuccessRateIncrease" bson:"defenseSuccessRateIncrease"`
// 	LifeAfterMonsterIncrease                  *big.Int `json:"lifeAfterMonsterIncrease" bson:"lifeAfterMonsterIncrease"`
// 	ManaAfterMonsterIncrease                  *big.Int `json:"manaAfterMonsterIncrease" bson:"manaAfterMonsterIncrease"`
// 	GoldAfterMonsterIncrease                  *big.Int `json:"goldAfterMonsterIncrease" bson:"goldAfterMonsterIncrease"`
// 	DoubleDamageProbabilityIncrease           *big.Int `json:"doubleDamageProbabilityIncrease" bson:"doubleDamageProbabilityIncrease"`
// 	ExcellentDamageProbabilityIncrease        *big.Int `json:"excellentDamageProbabilityIncrease" bson:"excellentDamageProbabilityIncrease"`
// 	IgnoreOpponentsDefenseRateIncrease        *big.Int `json:"ignoreOpponentsDefenseRateIncrease" bson:"ignoreOpponentsDefenseRateIncrease"`
// 	ReflectDamage                             *big.Int `json:"reflectDamage" bson:"reflectDamage"`
// 	MaxLifeIncrease                           *big.Int `json:"maxLifeIncrease" bson:"maxLifeIncrease"`
// 	MaxManaIncrease                           *big.Int `json:"maxManaIncrease" bson:"maxManaIncrease"`
// 	ExcellentDamageRateIncrease               *big.Int `json:"excellentDamageRateIncrease" bson:"excellentDamageRateIncrease"`
// 	DoubleDamageRateIncrease                  *big.Int `json:"doubleDamageRateIncrease" bson:"doubleDamageRateIncrease"`
// 	IgnoreOpponentsDefenseSuccessRateIncrease *big.Int `json:"ignoreOpponentsDefenseSuccessRateIncrease" bson:"ignoreOpponentsDefenseSuccessRateIncrease"`
// 	AttackDamageIncrease                      *big.Int `json:"attackDamageIncrease" bson:"attackDamageIncrease"`
// 	IsAncient                                 *big.Int `json:"isAncient" bson:"isAncient"`
// 	ReflectDamageRateIncrease                 *big.Int `json:"reflectDamageRateIncrease" bson:"reflectDamageRateIncrease"`
// 	DecreaseDamageRateIncrease                *big.Int `json:"decreaseDamageRateIncrease" bson:"decreaseDamageRateIncrease"`
// 	HpRecoveryRateIncrease                    *big.Int `json:"hpRecoveryRateIncrease" bson:"hpRecoveryRateIncrease"`
// 	MpRecoveryRateIncrease                    *big.Int `json:"mpRecoveryRateIncrease" bson:"mpRecoveryRateIncrease"`
// 	DefenceIncreasePerLevel                   *big.Int `json:"defenceIncreasePerLevel" bson:"defenceIncreasePerLevel"`
// 	DamageIncreasePerLevel                    *big.Int `json:"damageIncreasePerLevel" bson:"damageIncreasePerLevel"`
// 	IncreaseDefenseRate                       *big.Int `json:"increaseDefenseRate" bson:"increaseDefenseRate"`
// 	StrengthReqIncreasePerLevel               *big.Int `json:"strengthReqIncreasePerLevel" bson:"strengthReqIncreasePerLevel"`
// 	AgilityReqIncreasePerLevel                *big.Int `json:"agilityReqIncreasePerLevel" bson:"agilityReqIncreasePerLevel"`
// 	EnergyReqIncreasePerLevel                 *big.Int `json:"energyReqIncreasePerLevel" bson:"energyReqIncreasePerLevel"`
// 	VitalityReqIncreasePerLevel               *big.Int `json:"vitalityReqIncreasePerLevel" bson:"vitalityReqIncreasePerLevel"`
// 	AttackSpeedIncrease                       *big.Int `json:"attackSpeedIncrease" bson:"attackSpeedIncrease"`

// 	FighterId       *big.Int `json:"fighterId" bson:"fighterId"`
// 	LastUpdBlock    *big.Int `json:"lastUpdBlock" bson:"lastUpdBlock"`
// 	ItemRarityLevel *big.Int `json:"itemRarityLevel" bson:"itemRarityLevel"`

// 	ItemAttributesId *big.Int `json:"itemAttributesId" bson:"itemAttributesId"`

// 	Luck         bool `json:"luck" bson:"luck"`
// 	Skill        bool `json:"skill" bson:"skill"`
// 	IsBox        bool `json:"isBox" bson:"isBox"`
// 	IsWeapon     bool `json:"isWeapon" bson:"isWeapon"`
// 	IsArmour     bool `json:"isArmour" bson:"isArmour"`
// 	IsJewel      bool `json:"isJewel" bson:"isJewel"`
// 	IsMisc       bool `json:"isMisc" bson:"isMisc"`
// 	IsConsumable bool `json:"isConsumable" bson:"isConsumable"`
// 	InShop       bool `json:"inShop" bson:"inShop"`
// }

type ItemAttributes struct {
	Name                                      string   `json:"name" bson:"name"`
	TokenId                                   *big.Int `json:"tokenId" bson:"tokenId"`
	ItemLevel                                 *big.Int `json:"itemLevel" bson:"itemLevel"`
	MaxLevel                                  *big.Int `json:"maxLevel" bson:"maxLevel"`
	Durability                                *big.Int `json:"durability" bson:"durability"`
	ClassRequired                             *big.Int `json:"classRequired" bson:"classRequired"`
	StrengthRequired                          *big.Int `json:"strengthRequired" bson:"strengthRequired"`
	AgilityRequired                           *big.Int `json:"agilityRequired" bson:"agilityRequired"`
	EnergyRequired                            *big.Int `json:"energyRequired" bson:"energyRequired"`
	VitalityRequired                          *big.Int `json:"vitalityRequired" bson:"vitalityRequired"`
	ItemWidth                                 *big.Int `json:"itemWidth" bson:"itemWidth"`
	ItemHeight                                *big.Int `json:"itemHeight" bson:"itemHeight"`
	AcceptableSlot1                           *big.Int `json:"acceptableSlot1" bson:"acceptableSlot1"`
	AcceptableSlot2                           *big.Int `json:"acceptableSlot2" bson:"acceptableSlot2"`
	
	BaseMinPhysicalDamage                     *big.Int `json:"baseMinPhysicalDamage" bson:"baseMinPhysicalDamage"`
	BaseMaxPhysicalDamage                     *big.Int `json:"baseMaxPhysicalDamage" bson:"baseMaxPhysicalDamage"`
	BaseMinMagicDamage                        *big.Int `json:"baseMinMagicDamage" bson:"baseMinMagicDamage"`
	BaseMaxMagicDamage                        *big.Int `json:"baseMaxMagicDamage" bson:"baseMaxMagicDamage"`
	BaseDefense                               *big.Int `json:"baseDefense" bson:"baseDefense"`
	AttackSpeed                               *big.Int `json:"attackSpeed" bson:"attackSpeed"`
	AdditionalDamage                          *big.Int `json:"additionalDamage" bson:"additionalDamage"`
	AdditionalDefense                         *big.Int `json:"additionalDefense" bson:"additionalDefense"`
	
	FighterId                                 *big.Int `json:"fighterId" bson:"fighterId"`
	LastUpdBlock                              *big.Int `json:"lastUpdBlock" bson:"lastUpdBlock"`
	ItemRarityLevel                           *big.Int `json:"itemRarityLevel" bson:"itemRarityLevel"`
	ItemAttributesId                          *big.Int `json:"itemAttributesId" bson:"itemAttributesId"`
	PackSize                          		  *big.Int `json:"packSize" bson:"packSize"`
	
	Luck                                      bool     `json:"luck" bson:"luck"`
	Skill                                     bool     `json:"skill" bson:"skill"`
	IsPackable                                bool     `json:"isPackable" bson:"isPackable"`

	IsBox                                     bool     `json:"isBox" bson:"isBox"`
	IsWeapon                                  bool     `json:"isWeapon" bson:"isWeapon"`
	IsArmour                                  bool     `json:"isArmour" bson:"isArmour"`
	IsJewel                                   bool     `json:"isJewel" bson:"isJewel"`
	IsWings                                   bool     `json:"isWings" bson:"isWings"`
	IsMisc                                    bool     `json:"isMisc" bson:"isMisc"`
	IsConsumable                              bool     `json:"isConsumable" bson:"isConsumable"`
	InShop                                    bool     `json:"inShop" bson:"inShop"`

	//// Excellent
	// Wings
	IncreaseAttackSpeedPoints                 *big.Int `json:"increaseAttackSpeedPoints" bson:"increaseAttackSpeedPoints"`
	ReflectDamagePercent                      *big.Int `json:"reflectDamagePercent" bson:"reflectDamagePercent"`
	RestoreHPChance                           *big.Int `json:"restoreHPChance" bson:"restoreHPChance"`
	RestoreMPChance                           *big.Int `json:"restoreMPChance" bson:"restoreMPChance"`
	DoubleDamageChance                        *big.Int `json:"doubleDamageChance" bson:"doubleDamageChance"`
	IgnoreOpponentDefenseChance               *big.Int `json:"ignoreOpponentDefenseChance" bson:"ignoreOpponentDefenseChance"`
	
	// Weapons
	LifeAfterMonsterIncrease                  *big.Int `json:"lifeAfterMonsterIncrease" bson:"lifeAfterMonsterIncrease"`
	ManaAfterMonsterIncrease                  *big.Int `json:"manaAfterMonsterIncrease" bson:"manaAfterMonsterIncrease"`
	ExcellentDamageProbabilityIncrease        *big.Int `json:"excellentDamageProbabilityIncrease" bson:"excellentDamageProbabilityIncrease"`
	AttackSpeedIncrease                       *big.Int `json:"attackSpeedIncrease" bson:"attackSpeedIncrease"`
	AttackLvl20                               *big.Int `json:"attackLvl20" bson:"attackLvl20"`
	AttackIncreasePercent                     *big.Int `json:"attackIncreasePercent" bson:"attackIncreasePercent"`
	
	// Armours
	DefenseSuccessRateIncrease                *big.Int `json:"defenseSuccessRateIncrease" bson:"defenseSuccessRateIncrease"`
	GoldAfterMonsterIncrease                  *big.Int `json:"goldAfterMonsterIncrease" bson:"goldAfterMonsterIncrease"`
	ReflectDamage                             *big.Int `json:"reflectDamage" bson:"reflectDamage"`
	MaxLifeIncrease                           *big.Int `json:"maxLifeIncrease" bson:"maxLifeIncrease"`
	MaxManaIncrease                           *big.Int `json:"maxManaIncrease" bson:"maxManaIncrease"`
	HpRecoveryRateIncrease                    *big.Int `json:"hpRecoveryRateIncrease" bson:"hpRecoveryRateIncrease"`
	MpRecoveryRateIncrease                    *big.Int `json:"mpRecoveryRateIncrease" bson:"mpRecoveryRateIncrease"`
	DecreaseDamageRateIncrease                *big.Int `json:"decreaseDamageRateIncrease" bson:"decreaseDamageRateIncrease"`
}
 


type ItemDroppedEvent struct {
	ItemHash    common.Hash    `json:"itemHash"`
	Item        ItemAttributes `json:"item"`
	Qty         *big.Int       `json:"qty"`
	BlockNumber *big.Int       `json:"blockNumber"`
	Coords      Coordinate     `json:"coords"`
    OwnerId     *big.Int       `json:"ownerId"`
    TokenId     *big.Int       `json:"tokenId"`
}

type ItemPickedEvent struct {
	TokenId   *big.Int `json:"tokenId"`
	FighterId *big.Int `json:"fighterId"`
	Qty       *big.Int `json:"qty"`
}

type ItemListEntry struct {
    Name           string
    ItemsAttributes ItemAttributes
}

var ItemAttributesCache = make(map[int64]ItemAttributes)
var DroppedItems = make(map[common.Hash]*ItemDroppedEvent)
var DroppedItemsMutex sync.RWMutex


func handleItemDroppedEvent(logEntry *types.Log, blockNumber *big.Int, coords Coordinate, killer *big.Int) {
	// Parse the contract ABI
	parsedABI := loadABI("Items")

	// Iterate through logs and unpack the event data

	event := ItemDroppedEvent{}

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
		DroppedItemsMutex.Lock() // Use a mutex if needed to protect access to the map
		delete(DroppedItems, event.ItemHash)
		DroppedItemsMutex.Unlock()
		log.Printf("[handleItemDroppedEvent] Item with hash %v has been removed after 30 seconds", event.ItemHash)
		broadcastDropMessage()
	})

    DroppedItemsMutex.Lock()
	DroppedItems[event.ItemHash] = &event
    DroppedItemsMutex.Unlock()

	broadcastDropMessage()
}

func getDroppedItemsSafely(fighter *Fighter) map[common.Hash]*ItemDroppedEvent {
    DroppedItemsMutex.RLock()
    defer DroppedItemsMutex.RUnlock()

    items := DroppedItems

    return items
}

func hashItemAttributes(attributes *ItemAttributes) (string, error) {
    // Marshal attributes into a JSON byte slice
    attributesJSON, err := json.Marshal(attributes)
    if err != nil {
        return "", fmt.Errorf("Error marshaling ItemAttributes: %v", err)
    }

    // Generate a SHA-256 hash
    hash := sha256.Sum256(attributesJSON)

    // Convert the hash into a string
    hashString := hex.EncodeToString(hash[:])

    return hashString, nil
}

func generateItem(fighter *Fighter, itemName string, level, additionalPoints int64, luck, excellent bool) {
    // Load items list from JSON file
    itemsListJSON, err := ioutil.ReadFile("../itemsList.json")
    if err != nil {
        log.Printf("[generateItem] Error reading itemsList.json: %v", err)
        sendErrorMessage(fighter, "Error reading items list")
        return
    }

    // Unmarshal JSON into a slice of []interface{}
    var itemsList [][]interface{}
    err = json.Unmarshal(itemsListJSON, &itemsList)
    if err != nil {
        log.Printf("[generateItem] Error unmarshalling itemsList.json: %v", err)
        sendErrorMessage(fighter, "Error processing items list")
        return
    }

    // Find the item by name
	var item ItemAttributes
	for key, entry := range itemsList {
	    itemNameInEntry, ok := entry[0].(string)
	    if !ok {
	        log.Printf("[generateItem] Error asserting entry[0] as string: %v", entry[0])
	        continue
	    }
	    if strings.ToLower(itemNameInEntry) == strings.ToLower(itemName) {
	        item = getItemAttributes(int64(key+2))
	        break
	    }
	}

    if item.ItemAttributesId == nil {
        log.Printf("[generateItem] Item not found: %s", itemName)
        sendErrorMessage(fighter, "Item not found")
        return
    }

    // Update item attributes based on the drop command
    item.ItemLevel = big.NewInt(level)

    if item.IsWeapon {
    	item.AdditionalDamage = big.NewInt(additionalPoints)
    	item.Skill = true
    } 

    if item.IsArmour {
    	item.AdditionalDefense = big.NewInt(additionalPoints)
    } 
    
    item.Luck = luck

    
	if excellent {
    	item.IncreaseAttackSpeedPoints = big.NewInt(1)
    	item.ManaAfterMonsterIncrease = big.NewInt(1)
    	item.LifeAfterMonsterIncrease = big.NewInt(1)
    	item.GoldAfterMonsterIncrease = big.NewInt(1)
    	item.ReflectDamagePercent = big.NewInt(1)
    	item.RestoreHPChance = big.NewInt(1)
    	item.RestoreMPChance = big.NewInt(1)
    	item.DoubleDamageChance = big.NewInt(1)
    	item.IgnoreOpponentDefenseChance = big.NewInt(1)
    	item.ExcellentDamageProbabilityIncrease = big.NewInt(1)
    	item.AttackSpeedIncrease = big.NewInt(1)
    	item.AttackLvl20 = big.NewInt(1)
    	item.AttackIncreasePercent = big.NewInt(1)
    	item.DefenseSuccessRateIncrease = big.NewInt(1)
    	item.ReflectDamage = big.NewInt(1)
    	item.MaxLifeIncrease = big.NewInt(1)
    	item.MaxManaIncrease = big.NewInt(1)
    	item.DecreaseDamageRateIncrease = big.NewInt(1)
    	item.HpRecoveryRateIncrease = big.NewInt(1)
    	item.MpRecoveryRateIncrease = big.NewInt(1)
    }    


    MakeItem(fighter, &item)   
}
