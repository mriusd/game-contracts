// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract Items is ERC721 {
    using Counters for Counters.Counter;

    address public owner;

    mapping (uint256 => string) public itemName;

    struct ItemAttributes {
        uint itemLevel;
        uint maxLevel;
        uint durability;
        uint classRequired;
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
        bool isBox;
    }

    constructor() ERC721("Combats", "Item") {
        owner = msg.sender;
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

    mapping (uint256 => DropList) public boxDropList;

    event BoxDropUpdated(uint256 level, uint256 weaponsDropRate, uint256 armourDropRate, uint256 jewelsDropRate, uint256 miscDropRate);

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

    

}