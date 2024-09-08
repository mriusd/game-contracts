// items.list.go

package items

import (
	"log"

	"encoding/json"
	"io/ioutil"
)



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
		Binding        	string          `json:"binding"`
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

		switch item.Params.AcceptableSlot1 {
			case 1: 	item.Params.Type = "helm"
			case 2: 	item.Params.Type = "armour"
			case 3: 	item.Params.Type = "pants"
			case 4: 	item.Params.Type = "gloves"
			case 5: 	item.Params.Type = "boots"
			case 6: 	item.Params.Type = "weapon"
			case 7: 	item.Params.Type = "weapon"
			case 8: 	item.Params.Type = "pendant"
			case 9: 	item.Params.Type = "ring"
			case 10: 	item.Params.Type = "ring"
			case 11: 	item.Params.Type = "wings"
		}

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
			Binding:          item.Binding,
		}
	}


	// log.Printf("[loadItems] BaseItemParameters=%v", BaseItemParameters["Magic Box"])
	// log.Printf("[loadItems] BaseItemAttributes=%v", BaseItemAttributes["Magic Box"])


}