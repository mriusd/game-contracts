// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./ItemsHelper.sol";
import "./ItemsExcellentAtts.sol";

contract UpgradeItem is ItemsExcellentAtts {
	ItemsHelper private _itemsHelper;

	constructor (address itemsContract, address itemsHelperContract) {
		_itemsHelper = ItemsHelper(itemsHelperContract);
	}

	function upgradeItemLevel(uint256 itemTokenId, uint256 jewelTokenId) external {
        ExcellentItemAtts memory item = _itemsHelper.getTokenAttributes(itemTokenId);
        ExcellentItemAtts memory jewel = _itemsHelper.getTokenAttributes(jewelTokenId);
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
            _itemsHelper.setItemLevel(itemTokenId, item.itemLevel + 1);
        }

        if (item.itemLevel >= 6) { // soul
            require(jewel.itemAttributesId == 3, "Required jewel of soul");


            if (getRandomNumber(20) <= josSuccessRate+luckPoints) {
                _itemsHelper.setItemLevel(itemTokenId, item.itemLevel + 1);

            } else {
                if (item.itemLevel == 6) {
                    _itemsHelper.setItemLevel(itemTokenId, item.itemLevel - 1);
                } else {
                    _itemsHelper.setItemLevel(itemTokenId, 0);
                }
                
            }
        }

        _itemsHelper.burnItem(jewelTokenId);     
    }

    function updateItemAdditionalPoints(uint256 itemTokenId, uint256 jolTokenId) external returns (bool) {
        ExcellentItemAtts memory item = _itemsHelper.getTokenAttributes(itemTokenId);
        ExcellentItemAtts memory jewel = _itemsHelper.getTokenAttributes(jolTokenId);
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
                _itemsHelper.setAdditionalPoints(itemTokenId, item.additionalDamage+4);
            } else {
                _itemsHelper.setAdditionalPoints(itemTokenId, 0);
            }
            
        }

        if (item.isArmour) {
            if (getRandomNumber(30) <= jolSuccessRate+luckPoints) {
                _itemsHelper.setAdditionalPoints(itemTokenId, item.additionalDefense+4);
            } else {
                _itemsHelper.setAdditionalPoints(itemTokenId, 0);
            }
        }

        _itemsHelper.burnItem(jolTokenId);
    }
}