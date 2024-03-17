// items.getters.go

package items

import (
	"fmt"
)

func GetDropItems(rarityLevel int,  category string) []ItemAttributes {
	selectedItems := make([]ItemAttributes, 0)
	for _, item := range BaseItemAttributes {
		if 	item.ItemRarityLevel == rarityLevel && (
			(category == "jewel" && item.IsJewel) ||
			(category == "armour" && item.IsArmour) ||
			(category == "weapon" && item.IsWeapon) ||
			(category == "misc" && item.IsMisc) ||
			(category == "box" && item.IsBox)) {

			selectedItems = append(selectedItems, item)
		}
	}

	return selectedItems
}

func GetItemAttributesByName(name string) (ItemAttributes, error) {
	for _, item := range BaseItemAttributes {
		if 	item.Name == name {
			return item, nil
		}
	}

	return ItemAttributes{}, fmt.Errorf("[GetItemAttributesByName] Item not found: %v", name)
}

func GetItemParametersByName(name string) (ItemParameters, error) {
	for itemName, item := range BaseItemParameters {
		if 	itemName == name {
			return item, nil
		}
	}

	return ItemParameters{}, fmt.Errorf("[GetItemParametersByName] Item not found: %v", name)
}