// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract Items is ERC721 {
    using Counters for Counters.Counter;

    address public owner;

    mapping (uint256 => string) public itemName;

    

    struct ItemAttributes {
        string name;

        uint tokenId;        
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
        uint hpIncrease;
        uint mpIncrease;
        uint increaseDefenseRate;
        uint increaseStrength;
        uint increaseAgility;
        uint increaseEnergy;
        uint increaseVitality;
        uint attackSpeedIncrease; 

        bool luck;
        bool skill;
        bool isBox;
        bool isWeapon;
        bool isArmour;
        bool isJewel;
        bool isMisc;
        bool isConsumable;
        bool inShop;
        

    }

    function createItem(
        string calldata name,
        uint[49] calldata uintValues,
        bool[9] calldata boolValues
        
    ) external returns (uint256 tokenId)
    {
        if (uintValues[0] > 0) uintValues[0] == 0;
        ItemAttributes memory atts;
        
        atts.name = name;

        atts.tokenId = uintValues[0];
        atts.itemLevel = uintValues[1];
        atts.maxLevel = uintValues[2];
        atts.durability = uintValues[3];
        atts.classRequired = uintValues[4];
        atts.strengthRequired = uintValues[5];
        atts.agilityRequired = uintValues[6];
        atts.energyRequired = uintValues[7];
        atts.vitalityRequired = uintValues[8];
        atts.itemWidth = uintValues[9];
        atts.itemHeight = uintValues[10];
        atts.acceptableSlot1 = uintValues[11];
        atts.acceptableSlot2 = uintValues[12];

        atts.physicalDamage = uintValues[13];
        atts.magicDamage = uintValues[14];
        atts.defense = uintValues[15];
        atts.attackSpeed = uintValues[16];
        atts.defenseSuccessRate = uintValues[17];
        atts.additionalDamage = uintValues[18];
        atts.additionalDefense = uintValues[19];
        atts.increasedExperienceGain = uintValues[20];

        atts.damageIncrease = uintValues[21];
        atts.defenseSuccessRateIncrease = uintValues[22];
        atts.lifeAfterMonsterIncrease = uintValues[23];
        atts.manaAfterMonsterIncrease = uintValues[24];
        atts.zenAfterMonsterIncrease = uintValues[25];
        atts.doubleDamageProbabilityIncrease = uintValues[26];
        atts.excellentDamageProbabilityIncrease = uintValues[27];
        atts.ignoreOpponentsDefenseRateIncrease = uintValues[28];
        atts.reflectDamage = uintValues[29];
        atts.maxLifeIncrease = uintValues[30];
        atts.maxManaIncrease = uintValues[31];
        atts.excellentDamageRateIncrease = uintValues[32];
        atts.doubleDamageRateIncrease = uintValues[33];
        atts.ignoreOpponentsDefenseSuccessRateIncrease = uintValues[34];
        atts.attackDamageIncrease = uintValues[35];
        atts.defenseSuccessRateIncreaseAncient = uintValues[36];
        atts.reflectDamageRateIncrease = uintValues[37];
        atts.decreaseDamageRateIncrease = uintValues[38];
        atts.hpRecoveryRateIncrease = uintValues[39];
        atts.mpRecoveryRateIncrease = uintValues[40];
        atts.hpIncrease = uintValues[41];
        atts.mpIncrease = uintValues[42];
        atts.increaseDefenseRate = uintValues[43];
        atts.increaseStrength = uintValues[44];
        atts.increaseAgility = uintValues[45];
        atts.increaseEnergy = uintValues[46];
        atts.increaseVitality = uintValues[47];
        atts.attackSpeedIncrease = uintValues[48];

        atts.luck = boolValues[0];
        atts.skill = boolValues[1];
        atts.isBox = boolValues[2];
        atts.isWeapon = boolValues[3];
        atts.isArmour = boolValues[4];
        atts.isJewel = boolValues[5];
        atts.isMisc = boolValues[6];
        atts.isConsumable = boolValues[7];
        atts.inShop = boolValues[8];


        uint256 itemId = _itemAttributes.length;
        _itemAttributes[itemId] = atts;      


        return itemId;
    }

    event ItemBoughtFromShop(uint256 tokenId, uint256 itemId, address owner, string itemName);

    function buyItemFromShop(uint256 itemId, uint256 fighterId) external 
    {
        require(_itemAttributes[itemId].durability > 0, "Item doesn't exist");

        // money logic

        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        _safeMint(msg.sender, tokenId);
        _setTokenAttributes(tokenId, _itemAttributes[itemId]);

        _tokenAttributes[tokenId].tokenId = tokenId;      

        emit ItemBoughtFromShop(tokenId, itemId, msg.sender, _tokenAttributes[tokenId].name);
    }

    mapping (uint256 => ItemAttributes) private _tokenAttributes;
    ItemAttributes[] private _itemAttributes;
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