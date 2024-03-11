// items.go

package items

import (
	"context"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"encoding/json"
	"io/ioutil"

	"fmt"
	"crypto/sha256"
    "encoding/hex"
    "strings"



    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"

    "github.com/mriusd/game-contracts/db"
    "github.com/mriusd/game-contracts/maps"
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

	ItemAttributes  	*ItemAttributes `json:"itemAttributes"`
	ItemParameters 		*ItemParameters `json:"itemParameters"`
	ExcellentItemAttributes *ExcellentItemAttributes `json:"excellentItemAttributes"`

	sync.RWMutex
}

func (i *TokenAttributes) GetTokenId() *big.Int {
	i.RLock()
	defer i.RUnlock()

	return big.NewInt(0).Set(i.TokenId)
}

func (i *TokenAttributes) GetName() string {
	i.RLock()
	defer i.RUnlock()

	return i.Name
}

func (i *TokenAttributes) GetItemParameters() *ItemParameters {
	i.RLock()
	defer i.RUnlock()

	return i.ItemParameters
}

type SafeItemAttributesCache struct {
	Map map[int64]*TokenAttributes
	sync.RWMutex
}

var ItemAttributesCache = &SafeItemAttributesCache{Map: make(map[int64]*TokenAttributes)}

func (i *SafeItemAttributesCache) Find(index int64) *TokenAttributes {
	i.RLock()
	defer i.RUnlock()

	item, ok := i.Map[index]
	if !ok {
		return nil
	}

	return item
}

func (i *SafeItemAttributesCache) Add(index int64, atts *TokenAttributes) {
	i.Lock()
	defer i.Unlock()

	i.Map[index] = atts
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

	sync.RWMutex
}

func (i *ItemParameters) GetItemHeight() int64 {
	i.RLock()
	defer i.RUnlock()

	return i.ItemHeight
}

func (i *ItemParameters) GetItemWidth() int64 {
	i.RLock()
	defer i.RUnlock()

	return i.ItemWidth
}


type SafeItemParametersMap struct {
	Map map[string]*ItemParameters
	sync.RWMutex
}

var BaseItemParameters = &SafeItemParametersMap{Map: make(map[string]*ItemParameters)}

func (i *SafeItemParametersMap) Add(k string, v *ItemParameters) {
	i.Lock()
	defer i.Unlock()

	i.Map[strings.ToLower(k)] = v
}

func (i *SafeItemParametersMap) Find(k string) *ItemParameters {
	i.RLock()
	defer i.RUnlock()

	v, exists := i.Map[strings.ToLower(k)]
	if !exists {
		return nil
	}

	return v
}


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


type SafeItemAttributesMap struct {
	Map map[string]*ItemAttributes
	sync.RWMutex
}

var BaseItemAttributes = &SafeItemAttributesMap{Map: make(map[string]*ItemAttributes)}

func (i *SafeItemAttributesMap) Add(k string, v *ItemAttributes) {
	i.Lock()
	defer i.Unlock()

	i.Map[strings.ToLower(k)] = v
}

func (i *SafeItemAttributesMap) Find(k string) *ItemAttributes {
	i.RLock()
	defer i.RUnlock()

	log.Printf("[SafeItemAttributesMap.Find] k=%v i.Map[k]=%v", k, i.Map[k])

	v, exists := i.Map[strings.ToLower(k)]
	if !exists {
		return nil
	}

	return v
}

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
	Item        *TokenAttributes `json:"item"`
	Qty         *big.Int       `json:"qty"`
	BlockNumber *big.Int       `json:"blockNumber"`
	Coords      maps.Coordinate     `json:"coords"`
    OwnerId     *big.Int       `json:"ownerId"`
    TokenId     *big.Int       `json:"tokenId"`
}

type ItemDroppedEventSolidity struct {
	ItemHash    common.Hash    		`json:"itemHash"`
	Item        SolidityItemAtts 	`json:"item"`
	Qty         *big.Int       		`json:"qty"`
	BlockNumber *big.Int       		`json:"blockNumber"`
	Coords      maps.Coordinate     		`json:"coords"`
    OwnerId     *big.Int       		`json:"ownerId"`
    TokenId     *big.Int       		`json:"tokenId"`
    Map 		string 				`json:"map"`
    CX 			*big.Int   			`json:"cX"`
    CY 			*big.Int  			`json:"cY"`

    sync.RWMutex
}


func (i *ItemDroppedEventSolidity) GetItem() SolidityItemAtts {
	i.RLock()
	defer i.RUnlock()

	return i.Item
}

func (i *ItemDroppedEventSolidity) GetBlockNumber() *big.Int {
	i.RLock()
	defer i.RUnlock()

	return new(big.Int).Set(i.BlockNumber)
}

func (i *ItemDroppedEventSolidity) SetTokenId(v *big.Int) {
	i.Lock()
	defer i.Unlock()

	i.TokenId = v
}


type SafeDroppedItemsMap struct {
	Map map[common.Hash]*ItemDroppedEventSolidity
	sync.RWMutex
}

var DroppedItems = &SafeDroppedItemsMap{Map: make(map[common.Hash]*ItemDroppedEventSolidity)}



func (i *SafeDroppedItemsMap) GetMap() map[common.Hash]*ItemDroppedEventSolidity {
	i.RLock()
	defer i.RUnlock()

	copy := make(map[common.Hash]*ItemDroppedEventSolidity, len(i.Map))
    for key, val := range i.Map {
        copy[key] = val
    }
    return copy
}

func (i *SafeDroppedItemsMap) Remove(hash common.Hash) {
	i.Lock()
	defer i.Unlock()

	delete(DroppedItems.Map, hash)
}

func (i *SafeDroppedItemsMap) Add(hash common.Hash, item *ItemDroppedEventSolidity) {
	i.Lock()
	defer i.Unlock()

	DroppedItems.Map[hash] = item
}

func (i *SafeDroppedItemsMap) Find(hash common.Hash) *ItemDroppedEventSolidity {
	i.RLock()
	defer i.RUnlock()

	item, ok := i.Map[hash]
    if !ok {
        return nil
    }

    return item
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








func GetDroppedItemsInGo() map[common.Hash]*ItemDroppedEventGo {
    // Clear the DroppedItemsGo map first (in case there are stale entries)
    DroppedItemsGo := make(map[common.Hash]*ItemDroppedEventGo)

    // Iterate over DroppedItems and convert them
    for hash, solItem := range DroppedItems.GetMap() {
        DroppedItemsGo[hash] = &ItemDroppedEventGo{
            ItemHash:    solItem.ItemHash,
            Item:        ConvertSolidityItemToGoItem(solItem.Item),
            Qty:         solItem.Qty,
            BlockNumber: solItem.BlockNumber,
            Coords:      solItem.Coords,
            OwnerId:     solItem.OwnerId,
            TokenId:     solItem.TokenId,
        }
    }

    return DroppedItemsGo;
}

func ConvertSolidityDroppedEventToGo(sol *ItemDroppedEventSolidity) ItemDroppedEventGo {
	return ItemDroppedEventGo{
		ItemHash: sol.ItemHash,
		Item: ConvertSolidityItemToGoItem(sol.Item),
		Qty: sol.Qty,
		BlockNumber: sol.BlockNumber,
		Coords: sol.Coords,
		OwnerId: sol.OwnerId,
		TokenId: sol.TokenId,		
	}
}

func GenerateSolidityItem(itemName string) (SolidityItemAtts, error) {
	// Fetch data from the base maps
	itemAttrs := BaseItemAttributes.Find(itemName)

	if itemAttrs == nil {
		// Handle error: No such item found in base maps
		// You can return an empty SolidityItemAtts or handle it differently
		return SolidityItemAtts{}, fmt.Errorf("item with name %s not found in base attributes", itemName)
	}

	// Create the SolidityItemAtts object
	return SolidityItemAtts{
		Name:                itemAttrs.Name,
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
		ItemRarityLevel:                      big.NewInt(0),
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
	}, nil
}

func ConvertSolidityItemToGoItem(solidityItem SolidityItemAtts) *TokenAttributes {
	itemParams := BaseItemParameters.Find(solidityItem.Name) 
	log.Printf("[convertSolidityItemToGoItem] solidityItem=%v solidityItem.Name=%v itemParams=%v", solidityItem, solidityItem.Name, itemParams)
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

	return &TokenAttributes{
		Name:                  		solidityItem.Name,
		TokenId:               		solidityItem.TokenId,
		ItemLevel:             		solidityItem.ItemLevel,
		AdditionalDamage:      		solidityItem.AdditionalDamage,
		AdditionalDefense:     		solidityItem.AdditionalDefense,
		FighterId:             		solidityItem.FighterId,
		LastUpdBlock:          		solidityItem.LastUpdBlock,
		PackSize:             		solidityItem.PackSize,
		Luck:                  		solidityItem.Luck,
		Skill:                 		solidityItem.Skill,
		ItemAttributes:        		itemAttributes,
		ItemParameters: 			itemParams,
		ExcellentItemAttributes: 	excellentItemAttributes,
	}
}

// func getItemParameters(itemName string) *ItemParameters {
// 	log.Printf("[getItemParameters] itemName=%v params=%v", itemName, BaseItemParameters[strings.ToLower(itemName)])
// 	return BaseItemParameters[strings.ToLower(itemName)]
// }


func LoadItems() {
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
		Params        *ItemParameters  `json:"params"`
	}

	err = json.Unmarshal(file, &items)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// log.Printf("[loadItems] file=%v ", file)
	// log.Printf("[loadItems] items=%v ", items)

	for _, item := range items {
		log.Printf("[loadItems] item.Name=%v, item.Params=%v, item=%v ", item.Name, item.Params, item)

		// Populate BaseItemParameters
		BaseItemParameters.Add(item.Name,  item.Params)

		// Populate BaseItemAttributes
		BaseItemAttributes.Add(item.Name, &ItemAttributes{
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
		})
	}


	// log.Printf("[loadItems] BaseItemParameters=%v", BaseItemParameters["Magic Box"])
	// log.Printf("[loadItems] BaseItemAttributes=%v", BaseItemAttributes["Magic Box"])


}


// func getDroppedItemsSafely(fighter *Fighter) map[common.Hash]*ItemDroppedEventSolidity {
//     DroppedItemsMutex.RLock()
//     defer DroppedItemsMutex.RUnlock()

//     items := DroppedItems

//     return items
// }

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



func SaveItemAttributesToDB(item *TokenAttributes) {
    log.Printf("[saveItemAttributesToDB] item=%v", item)
    collection := db.Client.Database("game").Collection("items")

    item.RLock()
    jsonData, _ := json.Marshal(item)
    item.RUnlock()

    filter := bson.M{"tokenId": item.TokenId.Int64()}
    update := bson.M{"$set": bson.M{"attributes": string(jsonData)}}
    opts := options.Update().SetUpsert(true)
    _, err := collection.UpdateOne(context.Background(), filter, update, opts)

    if err != nil {
        log.Fatal("[saveItemAttributesToDB] ", err)
    }
}

func GetItemAttributesFromDB(itemId int64) (*TokenAttributes, bool) {
    collection := db.Client.Database("game").Collection("items")

    var itemWithAttributes struct {
        Attributes string `bson:"attributes"`
    }

    filter := bson.M{"tokenId": itemId}
    err := collection.FindOne(context.Background(), filter).Decode(&itemWithAttributes)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return &TokenAttributes{}, false
        }
        log.Fatal("[getItemAttributesFromDB] ", err)
    }

    var item TokenAttributes
    err = json.Unmarshal([]byte(itemWithAttributes.Attributes), &item)
    if err != nil {
        log.Fatal("[getItemAttributesFromDB] JSON unmarshal error: ", err)
    }

    return &item, true
}
