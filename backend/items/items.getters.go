// items.getters.go

package items

import (

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

func GetDropItemByName(name string) ItemAttributes {
	for _, item := range BaseItemAttributes {
		if 	item.Name == name {
			return item
		}
	}

	return ItemAttributes{}
}