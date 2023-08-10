// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;
import "./ItemsHelper.sol";
import "./ItemsExcellentAtts.sol";
import "./MoneyHelper.sol";
import "./FightersHelper.sol";

contract Shop is ItemsExcellentAtts, SafeMath {
    ItemsHelper private _itemsHelper;
    MoneyHelper private _moneyHelper;
    FightersHelper private _fightersHelper;

    uint256 WEI_MULTIPLE = 10**18;

    struct PriceList {
        uint256 weaponBasePrice;
        uint256 armourBasePrice; 
        uint256 wingBasePrice;


        uint256 jocPrice;
        
        uint256 josPrice;
        uint256 jobPrice;
        uint256 jolPrice;

        

        uint256 rarityMultiplierPct;
        uint256 levelMultiplierPct;
        uint256 addPointsMultiplierPct;
        uint256 luckMultiplierPct;
        uint256 exceMultiplierPct;

        uint256 buySellMultiplier;
    }
    

    PriceList public pl;

    event ItemBought(uint256 tokenId, string itemName, address ownerAddress);
    event ItemSold(uint256 tokenId, string itemName, uint256 itemPrice, address ownerAddress);

    constructor(address itemsHelperContract, address moneyHelperContract, address fightersHelperContract, PriceList memory _pl) {     
        _itemsHelper = ItemsHelper(itemsHelperContract);
        _moneyHelper = MoneyHelper(moneyHelperContract);
        _fightersHelper = FightersHelper(fightersHelperContract);

        pl = _pl;
    }
    
    function buyItemFromShop(string calldata itemName, uint256 fighterId) external {
        require(_itemsHelper.isItemInShop(itemName), "Item not in shop or doesn't exist");

        ExcellentItemAtts memory atts = _itemsHelper.getItemAttributes(itemName);   

        uint256 itemPrice = pl.buySellMultiplier * calculateItemSellingPrice(atts);
        address ownerAddress = _fightersHelper.getOwner(atts.fighterId);
        _moneyHelper.burnGold(ownerAddress, itemPrice);

        // money logic

        uint256 tokenId = _itemsHelper.craftItemForShop(itemName, fighterId);   

        emit ItemBought(tokenId, itemName, _fightersHelper.getOwner(fighterId));
    }

    function sellItemToShop(uint256 tokenId) external {
        ExcellentItemAtts memory atts = _itemsHelper.getTokenAttributes(tokenId);

        uint256 itemPrice = calculateItemSellingPrice(atts);

        address ownerAddress = _fightersHelper.getOwner(atts.fighterId);

        _itemsHelper.burnItem(tokenId);
        _moneyHelper.mintGold(ownerAddress, itemPrice);

        emit ItemSold(tokenId, atts.name, itemPrice, ownerAddress);
    }

    function calculateItemSellingPrice(ExcellentItemAtts memory atts)  public view returns (uint256) {
        bool isExc = _itemsHelper.isItemExcellent(atts.tokenId);
        uint256 basePrice;

        if (atts.isWeapon) {
            basePrice = pl.weaponBasePrice;
        } else if (atts.isArmour) {
            basePrice = pl.armourBasePrice;
        } else if (atts.isJewel) {
            if (stringsEqual(atts.name, "Jewel of Chaos")) {
                basePrice = pl.jocPrice;
            } else if (stringsEqual(atts.name, "Jewel of Soul")) {
                basePrice = pl.josPrice;
            } else if (stringsEqual(atts.name, "Jewel of Bless")) {
                basePrice = pl.jobPrice;
            } else if (stringsEqual(atts.name, "Jewel of Life")) {
                basePrice = pl.jolPrice;
            }
        } else if (atts.isWings) {
            basePrice = pl.wingBasePrice;
        }


        uint256 rarityMultiplier = 1 + ((pl.rarityMultiplierPct * atts.itemRarityLevel * WEI_MULTIPLE) / 100);
        uint256 levelMultiplier = 1 + ((pl.levelMultiplierPct * atts.itemLevel * WEI_MULTIPLE) / 100);
        uint256 optionMultiplier = 1 + ((pl.addPointsMultiplierPct * (atts.additionalDamage + atts.additionalDefense)/4 * WEI_MULTIPLE) / 100);
        uint256 luckMultiplier = atts.luck ? 1 +((pl.luckMultiplierPct * WEI_MULTIPLE) / 100) : WEI_MULTIPLE; // 1.1 in wei format
        uint256 excellentMultiplier = isExc ? 1 +((pl.exceMultiplierPct * WEI_MULTIPLE) / 100) : WEI_MULTIPLE; // 1.15 in wei format

        uint256 totalPrice = (basePrice * rarityMultiplier * levelMultiplier * optionMultiplier * luckMultiplier * excellentMultiplier) / (WEI_MULTIPLE**5);


        return totalPrice;
    }
}
