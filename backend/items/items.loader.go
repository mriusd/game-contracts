// items.list.go

package items

import (
	"log"

	"encoding/json"
	"io/ioutil"
)

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
}

var BaseItemAttributes = make(map[string]ItemAttributes)

func LoadItems() {
	log.Printf("[loadItems]")
	file, err := ioutil.ReadFile("./items/items.json")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}



	var items []struct {
		Name          	string          `json:"name"`
		MaxLevel      	int        	`json:"maxLevel"`
		ItemRarityLevel int        	`json:"itemRarityLevel"`
		IsPackable    	bool            `json:"isPackable"`
		IsBox         	bool            `json:"isBox"`
		IsWeapon      	bool            `json:"isWeapon"`
		IsArmour      	bool            `json:"isArmour"`
		IsJewel       	bool            `json:"isJewel"`
		IsWings       	bool            `json:"isWings"`
		IsMisc        	bool            `json:"isMisc"`
		IsConsumable  	bool            `json:"isConsumable"`
		InShop        	bool            `json:"inShop"`
		Params        	ItemParameters  `json:"params"`
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
		BaseItemParameters[item.Name] = item.Params

		// Populate BaseItemAttributes
		BaseItemAttributes[item.Name] = ItemAttributes{
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