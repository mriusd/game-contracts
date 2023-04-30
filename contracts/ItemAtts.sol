// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;


contract ItemAtts {

    uint generalDropRate = 100;

    uint maxAdditionalPoints = 28;
    uint jolSuccessRate = 50;
    uint josSuccessRate = 50;
    uint luckDropRate = 30;
    uint skillDropRate = 30;

    struct ItemAttributes {
        string name;

        uint256 tokenId;        
        uint256 itemLevel;
        uint256 maxLevel;
        uint256 durability;
        uint256 classRequired; // Dark Knight - 1, Dark Wizard - 2, Fairy Elf - 3, Magic Gladiator - 4
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
        uint256 physicalDamage;
        uint256 magicDamage;
        uint256 defense;
        uint256 attackSpeed;
        uint256 defenseSuccessRate;
        uint256 additionalDamage;
        uint256 additionalDefense;
        uint256 increasedExperienceGain;
        
        uint256 damageIncrease;
        uint256 defenseSuccessRateIncrease;
        uint256 lifeAfterMonsterIncrease;
        uint256 manaAfterMonsterIncrease;
        uint256 goldAfterMonsterIncrease;
        uint256 doubleDamageProbabilityIncrease;
        uint256 excellentDamageProbabilityIncrease;
        uint256 ignoreOpponentsDefenseRateIncrease;
        uint256 reflectDamage;
        uint256 maxLifeIncrease;
        uint256 maxManaIncrease;
        uint256 excellentDamageRateIncrease;
        uint256 doubleDamageRateIncrease;
        uint256 ignoreOpponentsDefenseSuccessRateIncrease;
        uint256 attackDamageIncrease;
        uint256 isAncient;
        uint256 reflectDamageRateIncrease;
        uint256 decreaseDamageRateIncrease;
        uint256 hpRecoveryRateIncrease;
        uint256 mpRecoveryRateIncrease;
        uint256 defenceIncreasePerLevel;
        uint256 damageIncreasePerLevel;
        uint256 increaseDefenseRate;
        uint256 strengthReqIncreasePerLevel;
        uint256 agilityReqIncreasePerLevel;
        uint256 energyReqIncreasePerLevel;
        uint256 vitalityReqIncreasePerLevel;
        uint256 attackSpeedIncrease; 
        
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
        bool isMisc;
        bool isConsumable;
        bool inShop;
    }


    function getRandomNumber (uint256 seed) internal returns (uint256) {
        uint256 randomNumber = uint256(keccak256(abi.encodePacked(block.prevrandao, seed, block.number, block.timestamp, msg.sender)));

        return randomNumber % 100;
    }

    function getRandomNumberMax (uint256 seed, uint256 max) internal returns (uint256) {
        uint256 randomNumber = uint256(keccak256(abi.encodePacked(block.prevrandao, seed, block.number, block.timestamp, msg.sender)));

        return randomNumber % max;
    }
}