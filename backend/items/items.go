// items.go

package items

import (
	"sync"
	"fmt"
	"encoding/json"

	"crypto/sha256"
    "encoding/hex"
)

var MAX_ITEM_LEVEL = 15

type ItemParameters struct {
	Durability         		int     	`json:"durability"`
	ClassRequired      		string  	`json:"classRequired"`
	StrengthRequired   		int     	`json:"strengthRequired"`
	AgilityRequired    		int     	`json:"agilityRequired"`
	EnergyRequired     		int     	`json:"energyRequired"`
	VitalityRequired   		int     	`json:"vitalityRequired"`
	ItemWidth          		int     	`json:"itemWidth"`
	ItemHeight         		int     	`json:"itemHeight"`
	AcceptableSlot1    		int     	`json:"acceptableSlot1"`
	AcceptableSlot2    		int     	`json:"acceptableSlot2"`
	MinPhysicalDamage 		int  		`json:"minPhysicalDamage"`
	MaxPhysicalDamage 		int  		`json:"maxPhysicalDamage"`
	MinMagicDamage  		int     	`json:"minMagicDamage"`
	MaxMagicDamage  		int     	`json:"maxMagicDamage"`
	Defense        			int     	`json:"defense"`
	AttackSpeed        		int     	`json:"attackSpeed"`
	Type 					string 		`json:"type"`
}

var BaseItemParameters = make(map[string]ItemParameters)

type ItemAttributes struct {
	Name                    string   	`json:"name" bson:"name"`
	MaxLevel                int 		`json:"maxLevel" bson:"maxLevel"`

	ItemRarityLevel         int 		`json:"itemRarityLevel" bson:"itemRarityLevel"`

	IsPackable              bool     	`json:"isPackable" bson:"isPackable"`

	IsBox                   bool     	`json:"isBox" bson:"isBox"`
	IsWeapon                bool     	`json:"isWeapon" bson:"isWeapon"`
	IsArmour                bool     	`json:"isArmour" bson:"isArmour"`
	IsJewel                 bool     	`json:"isJewel" bson:"isJewel"`
	IsWings                 bool     	`json:"isWings" bson:"isWings"`
	IsMisc                  bool     	`json:"isMisc" bson:"isMisc"`
	IsConsumable            bool     	`json:"isConsumable" bson:"isConsumable"`
	InShop                  bool     	`json:"inShop" bson:"inShop"`
	Binding                 string     	`json:"binding" bson:"binding"`

	Price  					int 		`json:"price" bson:"price"`
}

var BaseItemAttributes = make(map[string]ItemAttributes)

type TokenAttributes struct {
	Name            		string 					`json:"name" bson:"name"`
	TokenId         		int 					`json:"tokenId" bson:"token_id"`
	ItemLevel       		int 					`json:"itemLevel" bson:"item_level"`
	AdditionalDamage 		int 					`json:"additionalDamage" bson:"additional_damage"`
	AdditionalDefense 		int 					`json:"additionalDefense" bson:"additional_defence"`
	FighterId       		int 					`json:"fighterId" bson:"fighter_id"`
	PackSize        		int 					`json:"packSize" bson:"pack_size"`
	Luck            		bool   					`json:"luck" bson:"luck"`
	Skill           		bool   					`json:"skill" bson:"skill"`

	CreatedAt				int 					`json:"createdAt" bson:"created_at"`

	ItemAttributes  		ItemAttributes 			`json:"itemAttributes" bson:"-"`
	ItemParameters 			ItemParameters 			`json:"itemParameters" bson:"-"`
	ExcellentItemAttributes ExcellentItemAttributes `json:"excellentItemAttributes" bson:"excellent_item_attributes"`

	sync.RWMutex									`json:"-" bson:"-"`
}

func (i *TokenAttributes) GetLuck() bool {
	i.RLock()
	defer i.RUnlock()

	return i.Luck
}


func (i *TokenAttributes) GetItemLevel() int {
	i.RLock()
	defer i.RUnlock()

	return i.ItemLevel
}

func (i *TokenAttributes) IncreaseItemLevel() error {
	i.Lock()
	defer i.Unlock()

	if i.ItemLevel == MAX_ITEM_LEVEL {
		return fmt.Errorf("[IncreaseItemLevel] Item at max level")
	}

	i.ItemLevel++
	return nil
}

func (i *TokenAttributes) DecreaseItemLevel() {
	i.Lock()
	defer i.Unlock()

	if i.ItemLevel == 0 {
		return
	}

	if i.ItemLevel < 7 {
		i.ItemLevel--
	} else {
		i.ItemLevel = 0
	}
}

func (i *TokenAttributes) GetTokenId() int {
	i.RLock()
	defer i.RUnlock()

	return i.TokenId
}

func (i *TokenAttributes) GetName() string {
	i.RLock()
	defer i.RUnlock()

	return i.Name
}

func (i *TokenAttributes) GetItemParameters() ItemParameters {
	i.RLock()
	defer i.RUnlock()

	return i.ItemParameters
}

func (i *TokenAttributes) GetItemAttributes() ItemAttributes {
	i.RLock()
	defer i.RUnlock()

	return i.ItemAttributes
}

type SafeItemAttributesCache struct {
	Map map[int]*TokenAttributes
	sync.RWMutex
}

var ItemAttributesCache = &SafeItemAttributesCache{Map: make(map[int]*TokenAttributes)}

func (i *SafeItemAttributesCache) Find(index int) *TokenAttributes {
	i.RLock()
	defer i.RUnlock()

	item, ok := i.Map[index]
	if !ok {
		return nil
	}

	return item
}

func (i *SafeItemAttributesCache) Add(index int, atts *TokenAttributes) {
	i.Lock()
	defer i.Unlock()

	i.Map[index] = atts
}


type ExcellentItemAttributes struct {
	IsExcellent                     		  bool     `json:"IsExcellent"						bson:"is_excellent"`

	// Wings
	IncreaseAttackSpeedPoints                 int `json:"increaseAttackSpeedPoints" 			bson:"increase_attack_speed_points"`
	ReflectDamagePercent                      int `json:"reflectDamagePercent" 				bson:"reflect_damage_percent"`
	RestoreHPChance                           int `json:"restoreHPChance" 					bson:"restore_hp_chance"`
	RestoreMPChance                           int `json:"restoreMPChance" 					bson:"restore_mp_cance"`
	DoubleDamageChance                        int `json:"doubleDamageChance" 					bson:"double_damage_chance"`
	IgnoreOpponentDefenseChance               int `json:"ignoreOpponentDefenseChance" 		bson:"ignore_opponent_defense_chance"`
	
	// Weapons
	LifeAfterMonsterIncrease                  int `json:"lifeAfterMonsterIncrease" 			bson:"life_after_monster_increase"`
	ManaAfterMonsterIncrease                  int `json:"manaAfterMonsterIncrease" 			bson:"mana_after_monster_increase"`
	ExcellentDamageProbabilityIncrease        int `json:"excellentDamageProbabilityIncrease" 	bson:"excellent_damage_probability_increase"`
	AttackSpeedIncrease                       int `json:"attackSpeedIncrease" 				bson:"attack_speed_increase"`
	AttackLvl20                               int `json:"attackLvl20" 						bson:"attack_lvl_20"`
	AttackIncreasePercent                     int `json:"attackIncreasePercent" 				bson:"attack_increase_percent"`
	
	// Armours
	DefenseSuccessRateIncrease                int `json:"defenseSuccessRateIncrease" 			bson:"defense_success_rate_increase"`
	GoldAfterMonsterIncrease                  int `json:"goldAfterMonsterIncrease" 			bson:"gold_after_monster_increase"`
	ReflectDamage                             int `json:"reflectDamage" 						bson:"reflect_damage"`
	MaxLifeIncrease                           int `json:"maxLifeIncrease" 					bson:"max_life_increase"`
	MaxManaIncrease                           int `json:"maxManaIncrease" 					bson:"max_mana_increase"`
	HpRecoveryRateIncrease                    int `json:"hpRecoveryRateIncrease" 				bson:"hp_recovery_rate_increase"`
	MpRecoveryRateIncrease                    int `json:"mpRecoveryRateIncrease" 				bson:"mp_recovery_rate_increase"`
	DecreaseDamageRateIncrease                int `json:"decreaseDamageRateIncrease" 			bson:"decrease_damage_rate_increase"`
}

func HashItemAttributes(attributes *TokenAttributes) (string, error) {
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
 



