// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

contract ItemAttributes {
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
        uint goldAfterMonsterIncrease;
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
        uint isAncient;
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
        uint fighterId;
        uint lastUpdBlock;
        uint itemRarityLevel;

        uint itemAttributesId;

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
}