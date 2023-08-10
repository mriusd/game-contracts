// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "./ItemsAtts.sol";
contract ItemsExcellentAtts is ItemsAtts {

    struct ExcellentItemAtts {
        string name;

        uint256 tokenId;        
        uint256 itemLevel;
        uint256 maxLevel;

        uint256 additionalDamage;
        uint256 additionalDefense;
    
        uint256 fighterId;
        uint256 lastUpdBlock;
        uint256 itemRarityLevel;

        uint256 packSize;

        bool luck;
        bool skill;
        bool isPackable;
        
        bool isBox;
        bool isWeapon;
        bool isArmour;
        bool isJewel;
        bool isWings;
        bool isMisc;
        bool isConsumable;
        bool inShop;


        // Excellent
        bool isExcellent;
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