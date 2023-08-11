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
)


type SolidityItemAtts struct {
	Name                            string   `json:"name"`
	TokenId                         *big.Int `json:"tokenId"`
	ItemLevel                       *big.Int `json:"itemLevel"`
	MaxLevel                        *big.Int `json:"maxLevel"`
	AdditionalDamage                *big.Int `json:"additionalDamage"`
	AdditionalDefense               *big.Int `json:"additionalDefense"`
	FighterId                       *big.Int `json:"fighterId"`
	LastUpdBlock                    *big.Int `json:"lastUpdBlock"`
	ItemRarityLevel                 *big.Int `json:"itemRarityLevel"`
	PackSize                        *big.Int `json:"packSize"`
	Luck                            bool     `json:"luck"`
	Skill                           bool     `json:"skill"`
	IsPackable                      bool     `json:"isPackable"`
	IsBox                           bool     `json:"isBox"`
	IsWeapon                        bool     `json:"isWeapon"`
	IsArmour                        bool     `json:"isArmour"`
	IsJewel                         bool     `json:"isJewel"`
	IsWings                         bool     `json:"isWings"`
	IsMisc                          bool     `json:"isMisc"`
	IsConsumable                    bool     `json:"isConsumable"`
	InShop                          bool     `json:"inShop"`

	IsExcellent                     bool     `json:"isExcellent"`
	IncreaseAttackSpeedPoints       *big.Int `json:"increaseAttackSpeedPoints"`
	ReflectDamagePercent            *big.Int `json:"reflectDamagePercent"`
	RestoreHPChance                 *big.Int `json:"restoreHPChance"`
	RestoreMPChance                 *big.Int `json:"restoreMPChance"`
	DoubleDamageChance              *big.Int `json:"doubleDamageChance"`
	IgnoreOpponentDefenseChance     *big.Int `json:"ignoreOpponentDefenseChance"`
	LifeAfterMonsterIncrease        *big.Int `json:"lifeAfterMonsterIncrease"`
	ManaAfterMonsterIncrease        *big.Int `json:"manaAfterMonsterIncrease"`
	ExcellentDamageProbabilityIncrease *big.Int `json:"excellentDamageProbabilityIncrease"`
	AttackSpeedIncrease             *big.Int `json:"attackSpeedIncrease"`
	AttackLvl20                     *big.Int `json:"attackLvl20"`
	AttackIncreasePercent           *big.Int `json:"attackIncreasePercent"`
	DefenseSuccessRateIncrease      *big.Int `json:"defenseSuccessRateIncrease"`
	GoldAfterMonsterIncrease        *big.Int `json:"goldAfterMonsterIncrease"`
	ReflectDamage                   *big.Int `json:"reflectDamage"`
	MaxLifeIncrease                 *big.Int `json:"maxLifeIncrease"`
	MaxManaIncrease                 *big.Int `json:"maxManaIncrease"`
	HpRecoveryRateIncrease          *big.Int `json:"hpRecoveryRateIncrease"`
	MpRecoveryRateIncrease          *big.Int `json:"mpRecoveryRateIncrease"`
	DecreaseDamageRateIncrease      *big.Int `json:"decreaseDamageRateIncrease"`
}

type TokenAttributes struct {
	Name            	string 		`json:"name"`
	TokenId         	*big.Int 	`json:"tokenId"`
	ItemLevel       	*big.Int 	`json:"itemLevel"`
	AdditionalDamage 	*big.Int 	`json:"additionalDamage"`
	AdditionalDefense 	*big.Int 	`json:"additionalDefense"`
	FighterId       	*big.Int 	`json:"fighterId"`
	LastUpdBlock    	*big.Int 	`json:"lastUpdBlock"`
	PackSize        	*big.Int 	`json:"packSize"`
	Luck            	bool   		`json:"luck"`
	Skill           	bool   		`json:"skill"`

	ItemAttributes  *ItemAttributes `json:"itemAttributes"`
	ItemParameters 		*ItemParameters `json:"itemParameters"`
	ExcellentItemAttributes *ExcellentItemAttributes `json:"excellentItemAttributes"`
}


type ItemParameters struct {
	Durability         		int64     	`json:"durability"`
	ClassRequired      		string  	`json:"classRequired"`
	StrengthRequired   		int64     	`json:"strengthRequired"`
	AgilityRequired    		int64     	`json:"agilityRequired"`
	EnergyRequired     		int64     	`json:"energyRequired"`
	VitalityRequired   		int64     	`json:"vitalityRequired"`
	ItemWidth          		int64     	`json:"itemWidth"`
	ItemHeight         		int64     	`json:"itemHeight"`
	AcceptableSlot1    		int64     	`json:"acceptableSlot1"`
	AcceptableSlot2    		int64     	`json:"acceptableSlot2"`
	MinPhysicalDamage 		int64  		`json:"minPhysicalDamage"`
	MaxPhysicalDamage 		int64  		`json:"maxPhysicalDamage"`
	MinMagicDamage  		int64     	`json:"minMagicDamage"`
	MaxMagicDamage  		int64     	`json:"maxMagicDamage"`
	Defense        			int64     	`json:"defense"`
	AttackSpeed        		int64     	`json:"attackSpeed"`
}

var BaseItemParameters = make(map[string]*ItemParameters)


type ItemAttributes struct {
	Name                    string   `json:"name" bson:"name"`
	MaxLevel                *big.Int `json:"maxLevel" bson:"maxLevel"`

	ItemRarityLevel         *big.Int `json:"itemRarityLevel" bson:"itemRarityLevel"`

	IsPackable              bool     `json:"isPackable" bson:"isPackable"`

	IsBox                   bool     `json:"isBox" bson:"isBox"`
	IsWeapon                bool     `json:"isWeapon" bson:"isWeapon"`
	IsArmour                bool     `json:"isArmour" bson:"isArmour"`
	IsJewel                 bool     `json:"isJewel" bson:"isJewel"`
	IsWings                 bool     `json:"isWings" bson:"isWings"`
	IsMisc                  bool     `json:"isMisc" bson:"isMisc"`
	IsConsumable            bool     `json:"isConsumable" bson:"isConsumable"`
	InShop                  bool     `json:"inShop" bson:"inShop"`
}

var BaseItemAttributes = make(map[string]*ItemAttributes)


type ExcellentItemAttributes struct {
	IsExcellent                     		  bool     `json:"IsExcellent"`

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
 


type ItemDroppedEventGo struct {
	ItemHash    common.Hash    `json:"itemHash"`
	Item        TokenAttributes `json:"item"`
	Qty         *big.Int       `json:"qty"`
	BlockNumber *big.Int       `json:"blockNumber"`
	Coords      Coordinate     `json:"coords"`
    OwnerId     *big.Int       `json:"ownerId"`
    TokenId     *big.Int       `json:"tokenId"`
}

type ItemDroppedEventSolidity struct {
	ItemHash    common.Hash    `json:"itemHash"`
	Item        SolidityItemAtts `json:"item"`
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

var ItemAttributesCache = make(map[int64]TokenAttributes)
var DroppedItems = make(map[common.Hash]*ItemDroppedEventSolidity)
var DroppedItemsMutex sync.RWMutex


func getDroppedItemsInGo() map[common.Hash]*ItemDroppedEventGo {
    DroppedItemsMutex.RLock()  // Acquire read lock
    defer DroppedItemsMutex.RUnlock() // Ensure the lock is released after function execution

    // Clear the DroppedItemsGo map first (in case there are stale entries)
    DroppedItemsGo := make(map[common.Hash]*ItemDroppedEventGo)

    // Iterate over DroppedItems and convert them
    for hash, solItem := range DroppedItems {
        DroppedItemsGo[hash] = &ItemDroppedEventGo{
            ItemHash:    solItem.ItemHash,
            Item:        convertSolidityItemToGoItem(solItem.Item),
            Qty:         solItem.Qty,
            BlockNumber: solItem.BlockNumber,
            Coords:      solItem.Coords,
            OwnerId:     solItem.OwnerId,
            TokenId:     solItem.TokenId,
        }
    }

    return DroppedItemsGo;
}

func convertSolidityDroppedEventToGo(sol ItemDroppedEventSolidity) ItemDroppedEventGo {
	return ItemDroppedEventGo{
		ItemHash: sol.ItemHash,
		Item: convertSolidityItemToGoItem(sol.Item),
		Qty: sol.Qty,
		BlockNumber: sol.BlockNumber,
		Coords: sol.Coords,
		OwnerId: sol.OwnerId,
		TokenId: sol.TokenId,		
	}
}

func generateSolidityItem(itemName string) SolidityItemAtts {
	// Fetch data from the base maps
	itemAttrs, ok := BaseItemAttributes[itemName]

	if !ok {
		// Handle error: No such item found in base maps
		// You can return an empty SolidityItemAtts or handle it differently
		return SolidityItemAtts{}
	}

	// Create the SolidityItemAtts object
	return SolidityItemAtts{
		Name:                itemName,
		MaxLevel:            itemAttrs.MaxLevel,
		IsPackable:          itemAttrs.IsPackable,
		IsBox:               itemAttrs.IsBox,
		IsWeapon:            itemAttrs.IsWeapon,
		IsArmour:            itemAttrs.IsArmour,
		IsJewel:             itemAttrs.IsJewel,
		IsWings:             itemAttrs.IsWings,
		IsMisc:              itemAttrs.IsMisc,
		IsConsumable:        itemAttrs.IsConsumable,
		InShop:              itemAttrs.InShop,

		// Set all other fields to their zero values (including all Excellent fields)
		TokenId:                              big.NewInt(0),
		ItemLevel:                            big.NewInt(0),
		AdditionalDamage:                     big.NewInt(0),
		AdditionalDefense:                    big.NewInt(0),
		FighterId:                            big.NewInt(0),
		LastUpdBlock:                         big.NewInt(0),
		PackSize:                             big.NewInt(0),
		Luck:                                 false,
		Skill:                                false,
		IncreaseAttackSpeedPoints:            big.NewInt(0),
		ReflectDamagePercent:                 big.NewInt(0),
		RestoreHPChance:                      big.NewInt(0),
		RestoreMPChance:                      big.NewInt(0),
		DoubleDamageChance:                   big.NewInt(0),
		IgnoreOpponentDefenseChance:          big.NewInt(0),
		LifeAfterMonsterIncrease:             big.NewInt(0),
		ManaAfterMonsterIncrease:             big.NewInt(0),
		ExcellentDamageProbabilityIncrease:   big.NewInt(0),
		AttackSpeedIncrease:                  big.NewInt(0),
		AttackLvl20:                          big.NewInt(0),
		AttackIncreasePercent:                big.NewInt(0),
		DefenseSuccessRateIncrease:           big.NewInt(0),
		GoldAfterMonsterIncrease:             big.NewInt(0),
		ReflectDamage:                        big.NewInt(0),
		MaxLifeIncrease:                      big.NewInt(0),
		MaxManaIncrease:                      big.NewInt(0),
		HpRecoveryRateIncrease:               big.NewInt(0),
		MpRecoveryRateIncrease:               big.NewInt(0),
		DecreaseDamageRateIncrease:           big.NewInt(0),
	}
}

func convertSolidityItemToGoItem(solidityItem SolidityItemAtts) TokenAttributes {
	log.Printf("[convertSolidityItemToGoItem] solidityItem=%v solidityItem.Name=%v", solidityItem, solidityItem.Name)
	itemParams := getItemParameters(solidityItem.Name) 

	itemAttributes := &ItemAttributes{
		Name:        solidityItem.Name,
		MaxLevel:    solidityItem.MaxLevel,
		IsPackable:  solidityItem.IsPackable,
		IsBox:       solidityItem.IsBox,
		IsWeapon:    solidityItem.IsWeapon,
		IsArmour:    solidityItem.IsArmour,
		IsJewel:     solidityItem.IsJewel,
		IsWings:     solidityItem.IsWings,
		IsMisc:      solidityItem.IsMisc,
		IsConsumable: solidityItem.IsConsumable,
		InShop:      solidityItem.InShop,
	}

	excellentItemAttributes := &ExcellentItemAttributes{
		IncreaseAttackSpeedPoints:       solidityItem.IncreaseAttackSpeedPoints,
		ReflectDamagePercent:            solidityItem.ReflectDamagePercent,
		RestoreHPChance:                 solidityItem.RestoreHPChance,
		RestoreMPChance:                 solidityItem.RestoreMPChance,
		DoubleDamageChance:              solidityItem.DoubleDamageChance,
		IgnoreOpponentDefenseChance:     solidityItem.IgnoreOpponentDefenseChance,
		LifeAfterMonsterIncrease:        solidityItem.LifeAfterMonsterIncrease,
		ManaAfterMonsterIncrease:        solidityItem.ManaAfterMonsterIncrease,
		ExcellentDamageProbabilityIncrease: solidityItem.ExcellentDamageProbabilityIncrease,
		AttackSpeedIncrease:             solidityItem.AttackSpeedIncrease,
		AttackLvl20:                     solidityItem.AttackLvl20,
		AttackIncreasePercent:           solidityItem.AttackIncreasePercent,
		DefenseSuccessRateIncrease:      solidityItem.DefenseSuccessRateIncrease,
		GoldAfterMonsterIncrease:        solidityItem.GoldAfterMonsterIncrease,
		ReflectDamage:                   solidityItem.ReflectDamage,
		MaxLifeIncrease:                 solidityItem.MaxLifeIncrease,
		MaxManaIncrease:                 solidityItem.MaxManaIncrease,
		HpRecoveryRateIncrease:          solidityItem.HpRecoveryRateIncrease,
		MpRecoveryRateIncrease:          solidityItem.MpRecoveryRateIncrease,
		DecreaseDamageRateIncrease:      solidityItem.DecreaseDamageRateIncrease,
	}

	return TokenAttributes{
		Name:                  solidityItem.Name,
		TokenId:               solidityItem.TokenId,
		ItemLevel:             solidityItem.ItemLevel,
		AdditionalDamage:      solidityItem.AdditionalDamage,
		AdditionalDefense:     solidityItem.AdditionalDefense,
		FighterId:             solidityItem.FighterId,
		LastUpdBlock:          solidityItem.LastUpdBlock,
		PackSize:              solidityItem.PackSize,
		Luck:                  solidityItem.Luck,
		Skill:                 solidityItem.Skill,
		ItemAttributes:        itemAttributes,
		ItemParameters: 		itemParams,
		ExcellentItemAttributes: excellentItemAttributes,
	}
}

func getItemParameters(itemName string) *ItemParameters {
	log.Printf("[getItemParameters] itemName=%v params=%v", itemName, BaseItemParameters[itemName])
	return BaseItemParameters[itemName]
}


func loadItems() {
	log.Printf("[loadItems]")
	file, err := ioutil.ReadFile("../game_items.json")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}



	var items []struct {
		Name          string          `json:"name"`
		MaxLevel      *big.Int        `json:"maxLevel"`
		ItemRarityLevel      *big.Int        `json:"itemRarityLevel"`
		IsPackable    bool            `json:"isPackable"`
		IsBox         bool            `json:"isBox"`
		IsWeapon      bool            `json:"isWeapon"`
		IsArmour      bool            `json:"isArmour"`
		IsJewel       bool            `json:"isJewel"`
		IsWings       bool            `json:"isWings"`
		IsMisc        bool            `json:"isMisc"`
		IsConsumable  bool            `json:"isConsumable"`
		InShop        bool            `json:"inShop"`
		Params        ItemParameters  `json:"params"`
	}

	err = json.Unmarshal(file, &items)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// log.Printf("[loadItems] file=%v ", file)
	// log.Printf("[loadItems] items=%v ", items)

	for _, item := range items {
		//log.Printf("[loadItems] item.Name=%v, item.Params=%v, item=%v ", item.Name, item.Params, item)

		// Create a new variable for this iteration.
		currentParams := item.Params

		// Populate BaseItemParameters
		BaseItemParameters[item.Name] = &currentParams

		// Populate BaseItemAttributes
		BaseItemAttributes[item.Name] = &ItemAttributes{
			Name:             item.Name,
			MaxLevel:         item.MaxLevel,
			ItemRarityLevel:  item.ItemRarityLevel,
			IsPackable:       item.IsPackable,
			IsBox:            item.IsBox,
			IsWeapon:         item.IsWeapon,
			IsArmour:         item.IsArmour,
			IsJewel:          item.IsJewel,
			IsWings:          item.IsWings,
			IsMisc:           item.IsMisc,
			IsConsumable:     item.IsConsumable,
			InShop:           item.InShop,
		}
	}


	// log.Printf("[loadItems] BaseItemParameters=%v", BaseItemParameters["Magic Box"])
	// log.Printf("[loadItems] BaseItemAttributes=%v", BaseItemAttributes["Magic Box"])


}

func handleItemDroppedEvent(logEntry *types.Log, blockNumber *big.Int, coords Coordinate, killer *big.Int) {
	// Parse the contract ABI
	parsedABI := loadABI("Items")

	// Iterate through logs and unpack the event data

	event := ItemDroppedEventSolidity{}

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

func getDroppedItemsSafely(fighter *Fighter) map[common.Hash]*ItemDroppedEventSolidity {
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
    // Find the item by name
	item := generateSolidityItem(itemName)
	
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
