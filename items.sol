// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract Items is ERC721 {
    using Counters for Counters.Counter;

    address public owner;

    mapping (uint256 => string) public itemName;

    

    struct ItemAttributes {
        uint tokenId;
        string name;
        uint itemLevel;
        uint maxLevel;
        uint durability;
        uint classRequired; // Dark Knight - 1, Dark Wizard - 2, Fairy Elf - 3, Magic Gladiator - 4
        uint strengthRequired;
        uint agilityRequired;
        uint energyRequired;
        uint vitalityRequired;
        uint itemWidth;
        uint itemHeight;
        uint acceptableSlot1;
        uint acceptableSlot2;

        /*
            1. helmet
            2. armour
            3. pants
            4. gloves
            5. boots
            6. left hand
            7. right hand
            8. left ring
            9, right ring
            10, pendant
            11. wings

        */
        uint physicalDamage;
        uint magicDamage;
        uint defense;
        uint attackSpeed;
        uint defenseSuccessRate;
        uint additionalDamage;
        uint additionalDefense;
        uint increasedExperienceGain;
        bool luck;
        bool skill;
        uint damageIncrease;
        uint defenseSuccessRateIncrease;
        uint lifeAfterMonsterIncrease;
        uint manaAfterMonsterIncrease;
        uint zenAfterMonsterIncrease;
        uint doubleDamageProbabilityIncrease;
        uint excellentDamageProbabilityIncrease;
        uint ignoreOpponentsDefenseRateIncrease;
        uint reflectDamage;
        uint maxLifeIncrease;
        uint maxManaIncrease;
        uint excellentDamageRateIncrease;
        uint doubleDamageRateIncrease;
        uint ignoreOpponentsDefenseSuccessRateIncrease;
        uint attackDamageIncrease;
        uint defenseSuccessRateIncreaseAncient;
        uint reflectDamageRateIncrease;
        uint decreaseDamageRateIncrease;
        uint hpRecoveryRateIncrease;
        uint mpRecoveryRateIncrease;
        uint hpIncreaseRing;
        uint mpIncreaseRing;
        uint increaseDefenseRateRing;
        uint increaseStrengthRing;
        uint increaseAgilityRing;
        uint increaseEnergyRing;
        uint increaseVitalityRing;
        uint attackSpeedIncrease;
        bool isBox;
        bool isWeapon;
        bool isArmour;
        bool isJewel;
        bool isMisc;
        bool isConsumable;
        bool inShop;

    }

    mapping (uint256 => ItemAttributes) private _tokenAttributes;
    mapping (uint256 => ItemAttributes) private _itemAttributes;
    Counters.Counter private _tokenIdCounter;

    ItemAttributes boxAttributes;

    constructor() ERC721("Combats", "Item") {
        owner = msg.sender;

        boxAttributes.isBox = true;
        boxAttributes.itemWidth = 1;
        boxAttributes.itemHeight = 1;
    }

    struct DropList {
        uint256 level;
        uint256 weaponsDropRate;
        uint256[] weaponsList;
        uint256 armourDropRate;
        uint256[] armourList;
        uint256 jewelsDropRate;
        uint256[] jewelsList;
        uint256 miscDropRate;
        uint256[] miscList;
    }

    mapping (uint => DropList) public boxDropList;

    event BoxDropUpdated(uint level, uint256 weaponsDropRate, uint256 armourDropRate, uint256 jewelsDropRate, uint256 miscDropRate);

    function setBoxDrop (
        uint256 level, 
        uint256[4] calldata dropRates,
        uint256[] calldata weaponsList,        
        uint256[] calldata armourList, 
        uint256[] calldata jewelsList, 
        uint256[] calldata miscList
    ) external {
        DropList memory newList = DropList({
            level: level,
            weaponsDropRate: dropRates[0],
            armourDropRate: dropRates[1],
            jewelsDropRate: dropRates[2],
            miscDropRate: dropRates[3],
            weaponsList: weaponsList,
            armourList: armourList,
            jewelsList: jewelsList,
            miscList: miscList
        });
        boxDropList[level] = newList;
        emit BoxDropUpdated(level, dropRates[0], dropRates[1], dropRates[2], dropRates[3]);
    }

    event BoxDropped(address owner, uint256 tokenId);
    function dropBox (uint256 tokenId) external returns (uint256) {

        require(_exists(tokenId), "Token does not exist");
        require(_tokenAttributes[tokenId].isBox, "Not a box");
        uint256 boxLevel = _tokenAttributes[tokenId].itemLevel;
        
        DropList memory dropList = boxDropList[boxLevel];

        uint256 totalProbability = safeAdd(dropList.weaponsDropRate, safeAdd(dropList.armourDropRate, safeAdd(dropList.jewelsDropRate, dropList.miscDropRate)));

        uint256 randomNum = getRandomNumber();

        uint256 droppedItemId = 0;

        bool isWeapon;
        bool isArmour; 
        bool isJewel; 
        bool isMisc;

        if (randomNum <= totalProbability) {
            if (randomNum <= dropList.weaponsDropRate)    {
                isWeapon = true;
                droppedItemId = returnRandomItemFromDropList(dropList.weaponsList);
            } else if (randomNum <= safeAdd(dropList.weaponsDropRate, dropList.armourDropRate)) {
                isArmour = true;
                droppedItemId = returnRandomItemFromDropList(dropList.armourList);
            } else if (randomNum <= safeAdd(dropList.jewelsDropRate, safeAdd(dropList.weaponsDropRate, dropList.armourDropRate))) {
                isJewel = true;
                droppedItemId = returnRandomItemFromDropList(dropList.jewelsList);
            } else {
                isMisc = true;
                droppedItemId = returnRandomItemFromDropList(dropList.miscList);
            }
        } else {
            return 0;
        }

        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        _safeMint(msg.sender, tokenId);
        _setTokenAttributes(tokenId, _itemAttributes[droppedItemId]);

        // Determine Excelent option
        //uint excelentProbability = boxDropList[_tokenAttributes[tokenId].itemLevel];

        _tokenAttributes[tokenId].tokenId = tokenId;       

        emit BoxDropped(msg.sender, tokenId);

        return tokenId;
    }

    function mintBox(address winner, uint256 opponentLevel) external {
        uint256 boxLevel = opponentLevel/50;

        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        _safeMint(winner, tokenId);

        ItemAttributes memory newBoxAttributes = boxAttributes;

        boxAttributes.itemLevel = boxLevel; 

        _setTokenAttributes(tokenId, boxAttributes);
    }

    function _setTokenAttributes(uint256 tokenId, ItemAttributes memory attrs) internal {
        require(_exists(tokenId), "Token does not exist");

        _tokenAttributes[tokenId] = attrs;
    }

    function returnRandomItemFromDropList(uint256[] memory items) internal returns (uint256) {
        uint256 randomNumber = getRandomNumber();
        uint256 len = items.length;

        return items[randomNumber % len];

    }

    function getRandomNumber () internal returns (uint256) {
        uint256 randomNumber = block.prevrandao;

        uint256 modulus = randomNumber % 100;

        return randomNumber;
    }


    // Returns the smaller of two values
    function min(uint a, uint b) private pure returns (uint) {
        return a < b ? a : b;
    }

    // Returns the largest of the two values
    function max(uint a, uint b) private pure returns (uint) {
        return a > b ? a : b;
    }

    // Safe Multiply Function - prevents integer overflow 
    function safeMul(uint a, uint b) public returns (uint) {
        uint c = a * b;
        assert(a == 0 || c / a == b);
        return c;
    }

    // Safe Subtraction Function - prevents integer overflow 
    function safeSub(uint a, uint b) public returns (uint) {
        assert(b <= a);
        return a - b;
    }

    // Safe Addition Function - prevents integer overflow 
    function safeAdd(uint a, uint b) public returns (uint) {
        uint c = a + b;
        assert(c>=a && c>=b);
        return c;
    }

}