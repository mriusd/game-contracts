// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./ItemAtts.sol";

abstract contract MainItems {
    function setAdditionalPoints(uint256 tokenId, uint256 points) external virtual;
    function setItemLevel(uint256 tokenId, uint256 level)  external virtual;
    function burnItem(uint256 itemId)  external virtual;
    function getTokenAttributes(uint256 tokenId) external virtual view returns (ItemAtts.ItemAttributes memory);
    function craftItem(uint256 itemId, address itemOwner) external virtual returns (uint256);
}

contract UpgradeItem is ItemAtts {
	address ItemsContractAddress;	

	constructor (address itemsContract) {
		ItemsContractAddress = itemsContract;
	}

	function upgradeItemLevel(uint256 itemTokenId, uint256 jewelTokenId) external {
        ItemAttributes memory item = MainItems(ItemsContractAddress).getTokenAttributes(itemTokenId);
        ItemAttributes memory jewel = MainItems(ItemsContractAddress).getTokenAttributes(jewelTokenId);
        uint256 luckPoints = 0;
        bool success = false;

        require(item.isWeapon || item.isArmour, "Invalid item");
        require(item.itemLevel < 9, "Item can be upgraded only up to level 9, use Chaos Machine further");
        require(jewel.itemAttributesId == 2 ||  jewel.itemAttributesId == 3, "Invalid jewel");

        if (item.luck) {
            luckPoints = 25;
        }

        if (item.itemLevel < 6) { // bless
            require(jewel.itemAttributesId == 2, "Required jewel of bless");
            MainItems(ItemsContractAddress).setItemLevel(itemTokenId, item.itemLevel + 1);
        }

        if (item.itemLevel >= 6) { // soul
            require(jewel.itemAttributesId == 3, "Required jewel of soul");


            if (getRandomNumber(20) <= josSuccessRate+luckPoints) {
                MainItems(ItemsContractAddress).setItemLevel(itemTokenId, item.itemLevel + 1);

            } else {
                if (item.itemLevel == 6) {
                    MainItems(ItemsContractAddress).setItemLevel(itemTokenId, item.itemLevel - 1);
                } else {
                    MainItems(ItemsContractAddress).setItemLevel(itemTokenId, 0);
                }
                
            }
        }

        MainItems(ItemsContractAddress).burnItem(jewelTokenId);     
    }

    function updateItemAdditionalPoints(uint256 itemTokenId, uint256 jolTokenId) external returns (bool) {
        ItemAttributes memory item = MainItems(ItemsContractAddress).getTokenAttributes(itemTokenId);
        ItemAttributes memory jewel = MainItems(ItemsContractAddress).getTokenAttributes(jolTokenId);
        uint256 luckPoints = 0;
        bool success = false;

        if (item.luck) {
            luckPoints = 25;
        }

        require(item.isWeapon || item.isArmour, "Invalid item");
        require(item.additionalDamage + item.additionalDefense < maxAdditionalPoints, "Max additional points reached");
        require(jewel.itemAttributesId == 4, "Invalid jewel");

        
        if (item.isWeapon) {
            if (getRandomNumber(30) <= jolSuccessRate+luckPoints) {
                MainItems(ItemsContractAddress).setAdditionalPoints(itemTokenId, item.additionalDamage+4);
            } else {
                MainItems(ItemsContractAddress).setAdditionalPoints(itemTokenId, 0);
            }
            
        }

        if (item.isArmour) {
            if (getRandomNumber(30) <= jolSuccessRate+luckPoints) {
                MainItems(ItemsContractAddress).setAdditionalPoints(itemTokenId, item.additionalDefense+4);
            } else {
                MainItems(ItemsContractAddress).setAdditionalPoints(itemTokenId, 0);
            }
        }

        MainItems(ItemsContractAddress).burnItem(jolTokenId);
    }
}