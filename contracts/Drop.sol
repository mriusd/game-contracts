// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "./SafeMath.sol";
import "./ItemsHelper.sol";
import "./MoneyHelper.sol";

contract Drop is ItemsExcellentAtts, SafeMath {
    ItemsHelper private _itemsHelper;
    MoneyHelper private _moneyHelper;

    constructor(address itemsHelperContract, address moneyHelperContract) {     
        _itemsHelper = ItemsHelper(itemsHelperContract);
        _moneyHelper = MoneyHelper(moneyHelperContract);
    }

    mapping (bytes32 => uint256) public dropHashes;

    event ItemDropped(bytes32 itemHash, ExcellentItemAtts item, uint256 qty, uint256 tokenId);
    event ItemPicked(uint256 tokenId, uint256 fighterId, uint256 qty);
    event BoxOppened(uint256 tokenId, uint256 rarityLevel, address owner);
    event DropParametersChange(uint256 rarityLevel, DropParams params);
    event BoxDropParametersChange(uint256 rarityLevel, DropParams params);


    struct DropParams {
        uint256 weaponsDropRate;
        uint256 armoursDropRate;
        uint256 jewelsDropRate;
        uint256 miscDropRate;
        uint256 boxDropRate;

        uint256 excDropRate;

        uint256 minItemLevel;
        uint256 maxItemLevel;
        uint256 maxAddPoints;

        uint256 blockCrated;
    }

    mapping (uint256 => DropParams) public DropParamsList; // mapping of rarityLevels to drop parameters
    mapping (uint256 => DropParams) public BoxDropPramsList; // mapping of rarityLevels to drop parameters
    
    // This function is for development and testing only
    // Should be disabled in Production
    function makeItem(ExcellentItemAtts memory itemAtts) external returns (bytes32) {
        bytes32 itemHash = keccak256(abi.encode(itemAtts, 1, block.number));
        dropHashes[itemHash] = 1;
        emit ItemDropped(itemHash, itemAtts, 1, 0);
        return itemHash;
    }

    function setDropParams(uint256 rarityLevel, DropParams memory params) external {
        DropParamsList[rarityLevel] = params;
        emit DropParametersChange(rarityLevel, params);
    }

    function setBoxDropParams(uint256 rarityLevel, DropParams memory params) external {
        BoxDropPramsList[rarityLevel] = params;
        emit BoxDropParametersChange(rarityLevel, params);
    }

    function openBox(uint256 tokenId) external returns (uint256) {
        require(_itemsHelper.itemExists(tokenId), "Token doesn't exist");
        ExcellentItemAtts memory box = _itemsHelper.getTokenAttributes(tokenId);
        DropParams memory params = BoxDropPramsList[box.itemRarityLevel];

        ExcellentItemAtts memory itemAtts = getDropItem(box.itemRarityLevel, params);

        uint256 qty = 1;

        bytes32 itemHash = keccak256(abi.encode(itemAtts, qty, block.number));
        dropHashes[itemHash] = qty;
         emit ItemDropped(itemHash, itemAtts, qty, 0);
        _itemsHelper.burnItem(tokenId);
       
    }  

    function dropItem(uint256 rarityLevel, uint256 fighterId, uint256 experience) external returns (bytes32) {
        DropParams memory params = DropParamsList[rarityLevel];
        require(params.blockCrated > 0, "No drop parameters for rarityLevel");
        require(experience > 0, "Experience must be non zero");

        ExcellentItemAtts memory itemAtts = getDropItem(rarityLevel, params);        

        uint256 randomNumber = getRandomNumber(200);

        uint256 qty = 1;

        if (randomNumber > generalDropRate) {
            itemAtts = _itemsHelper.getItemAttributes("Empty item");
            emit ItemDropped(dummyHash, itemAtts, 0, 0);
            return 0;
        }

        if (stringsEqual(itemAtts.name, "Gold")) {
            qty = max(1, experience/_moneyHelper.getExperienceDivider());
        }

        bytes32 itemHash = keccak256(abi.encode(itemAtts, qty, block.number));
        dropHashes[itemHash] = qty;
        emit ItemDropped(itemHash, itemAtts, qty, 0);
        return itemHash;
    }



    function getDropItem(uint256 rarityLevel, DropParams memory params) internal returns (ExcellentItemAtts memory) {

        uint256 randomNumber = getRandomNumberMax(0, 100);
        string memory randomItem;
        
        //emit LogDropValues(randomNumber, params.jewelsDropRate, params.weaponsDropRate, params.armoursDropRate, params.miscDropRate, params.boxDropRate);
        
        uint256 cumulativeRate = params.jewelsDropRate;
        if (randomNumber < cumulativeRate) {
            randomItem = returnRandomItemFromDropList(100, _itemsHelper.getJewels(rarityLevel));
        } else {
            cumulativeRate += params.armoursDropRate;
            if (randomNumber < cumulativeRate) {
                randomItem = returnRandomItemFromDropList(101, _itemsHelper.getArmours(rarityLevel));
            } else {
                cumulativeRate += params.weaponsDropRate;
                if (randomNumber < cumulativeRate) {
                    randomItem = returnRandomItemFromDropList(102, _itemsHelper.getWeapons(rarityLevel));
                } else {
                    cumulativeRate += params.miscDropRate;
                    if (randomNumber < cumulativeRate) {
                        randomItem = returnRandomItemFromDropList(103, _itemsHelper.getMisc(rarityLevel));
                    } else {
                        cumulativeRate += params.boxDropRate;
                        if (randomNumber < cumulativeRate) {
                            randomItem = returnRandomItemFromDropList(104, _itemsHelper.getBoxes(rarityLevel));
                        } else {
                            randomItem = "Gold";
                        }
                    }
                }
            }
        }


        ExcellentItemAtts memory itemAtts = _itemsHelper.getItemAttributes(randomItem);


        if (!itemAtts.isJewel && !itemAtts.isMisc && !itemAtts.isBox && getRandomNumber(1) <= luckDropRate && !stringsEqual(itemAtts.name, "Empty item") && !stringsEqual(itemAtts.name, "Gold")) {
            itemAtts.luck  = true;
        }

        if (itemAtts.isWeapon && getRandomNumber(2) <= skillDropRate) {
            itemAtts.skill  = true;
        }

        if (itemAtts.isBox) {
            itemAtts.itemLevel = rarityLevel;
        } else if (!itemAtts.isMisc && !itemAtts.isJewel && !stringsEqual(itemAtts.name, "Empty item") && !stringsEqual(itemAtts.name, "Gold")) {
            itemAtts.itemLevel = params.minItemLevel + getRandomNumber(4) % (params.maxItemLevel-params.minItemLevel+1);

            if (itemAtts.isWeapon) {
                itemAtts.additionalDamage = 4 * (getRandomNumber(5) % (params.maxAddPoints/4+1));
            } else if (itemAtts.isArmour) {
                itemAtts.additionalDefense = 4 * (getRandomNumber(5) % (params.maxAddPoints/4+1));
            }
        }

        if ((itemAtts.isWeapon || itemAtts.isArmour) && getRandomNumberMax(3, 1000) <= params.excDropRate) {
            itemAtts = _itemsHelper.addExcellentOption(itemAtts);
        }


        
        return itemAtts;  
    }


    

    function getDropQty(bytes32 dropHash) external returns (uint256) {
        return dropHashes[dropHash];
    }

    function createDropHash(bytes32 itemHash, uint256 qty) external {
        dropHashes[itemHash] = qty;
    }

    function returnRandomItemFromDropList(uint256 seed, string[] memory items) internal returns (string memory) {
        uint256 randomNumber = getRandomNumber(seed);
        uint256 len = items.length;

        return items[randomNumber % len];

    }

    
}
