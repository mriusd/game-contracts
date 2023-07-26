// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "./ItemsAtts.sol";
contract ItemsExcellentAtts is ItemsAtts {

    struct ExcellentItemAtts {
        string name;

        uint256 tokenId;        
        uint256 itemLevel;
        uint256 maxLevel;
        uint256 durability;
        uint256 classRequired; 
        uint256 strengthRequired;
        uint256 agilityRequired;
        uint256 energyRequired;
        uint256 vitalityRequired;
        uint256 itemWidth;
        uint256 itemHeight;
        uint256 acceptableSlot1;
        uint256 acceptableSlot2;

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

        uint256 baseMinPhysicalDamage;
        uint256 baseMaxPhysicalDamage;
        uint256 baseMinMagicDamage;
        uint256 baseMaxMagicDamage;
        uint256 baseDefense;
        uint256 attackSpeed;

        uint256 additionalDamage;
        uint256 additionalDefense;
    
        uint256 fighterId;
        uint256 lastUpdBlock;
        uint256 itemRarityLevel;

        uint256 itemAttributesId;

        bool luck;
        bool skill;
        bool isBox;
        bool isWeapon;
        bool isArmour;
        bool isJewel;
        bool isWings;
        bool isMisc;
        bool inShop;


        // Excellent
        uint256 increaseAttackSpeedPoints;
        uint256 reflectDamagePercent;
        uint256 restoreHPChance;
        uint256 restoreMPChance;
        uint256 doubleDamageChance;
        uint256 ignoreOpponentDefenseChance;

        uint256 lifeAfterMonsterIncrease;
        uint256 manaAfterMonsterIncrease;
        uint256 excellentDamageProbabilityIncrease;
        uint256 attackSpeedIncrease;
        uint256 attackLvl20; 
        uint256 attackIncreasePercent; 

        uint256 defenseSuccessRateIncrease;
        uint256 goldAfterMonsterIncrease;
        uint256 reflectDamage;
        uint256 maxLifeIncrease;
        uint256 maxManaIncrease;
        uint256 hpRecoveryRateIncrease;
        uint256 mpRecoveryRateIncrease;
        uint256 decreaseDamageRateIncrease;
    }
}