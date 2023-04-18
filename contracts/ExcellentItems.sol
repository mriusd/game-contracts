// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "./ItemAtts.sol";
contract ExcellentItems is ItemAtts {  
    struct WeaponsExcOptions {
        uint lifeAfterMonsterIncrease;
        uint manaAfterMonsterIncrease;
        uint excellentDamageProbabilityIncrease;
        uint attackSpeedIncrease;
        uint damageIncrease; // 2%

        uint doubleOptionChance;
    }

    WeaponsExcOptions WeaponsExcOptionsDropRates = WeaponsExcOptions({
        lifeAfterMonsterIncrease: 35,
        manaAfterMonsterIncrease: 35,
        excellentDamageProbabilityIncrease: 10,
        attackSpeedIncrease: 10,
        damageIncrease: 10,

        doubleOptionChance: 15
    });

    struct ArmoursExcOptions {
        uint defenseSuccessRateIncrease;
        uint goldAfterMonsterIncrease;
        uint reflectDamage;
        uint maxLifeIncrease;
        uint maxManaIncrease;
        uint hpRecoveryRateIncrease;
        uint mpRecoveryRateIncrease;
        uint decreaseDamageRateIncrease;

        uint doubleOptionChance;
    }

    ArmoursExcOptions ArmourExcOptionsDropRates = ArmoursExcOptions({
        defenseSuccessRateIncrease: 10,
        goldAfterMonsterIncrease: 20,
        reflectDamage: 10,
        maxLifeIncrease: 10,
        maxManaIncrease: 10,
        hpRecoveryRateIncrease: 10,
        mpRecoveryRateIncrease: 10,
        decreaseDamageRateIncrease: 10,

        doubleOptionChance: 30
    });

    function getRandomArmourExcOption(ItemAttributes memory item,  uint256 seed) internal returns (ItemAttributes memory) {
        uint256 randomNumber = getRandomNumber(seed);
        uint256 totalProbability = ArmourExcOptionsDropRates.defenseSuccessRateIncrease +
            ArmourExcOptionsDropRates.goldAfterMonsterIncrease +
            ArmourExcOptionsDropRates.reflectDamage +
            ArmourExcOptionsDropRates.maxLifeIncrease +
            ArmourExcOptionsDropRates.maxManaIncrease +
            ArmourExcOptionsDropRates.hpRecoveryRateIncrease +
            ArmourExcOptionsDropRates.mpRecoveryRateIncrease +
            ArmourExcOptionsDropRates.decreaseDamageRateIncrease;
        

        uint256 randomIndex = randomNumber % totalProbability;

        if (randomIndex <= ArmourExcOptionsDropRates.defenseSuccessRateIncrease) {
            item.defenseSuccessRateIncrease = 1;
        } else if (randomIndex <= ArmourExcOptionsDropRates.defenseSuccessRateIncrease 
            + ArmourExcOptionsDropRates.goldAfterMonsterIncrease) {
            item.goldAfterMonsterIncrease = 1;
        } else if (randomIndex <= ArmourExcOptionsDropRates.defenseSuccessRateIncrease 
            + ArmourExcOptionsDropRates.goldAfterMonsterIncrease
            + ArmourExcOptionsDropRates.reflectDamage) {
            item.reflectDamage = 1;
        } else if (randomIndex <= ArmourExcOptionsDropRates.defenseSuccessRateIncrease 
            + ArmourExcOptionsDropRates.goldAfterMonsterIncrease
            + ArmourExcOptionsDropRates.reflectDamage
            + ArmourExcOptionsDropRates.maxLifeIncrease) {
            item.maxLifeIncrease = 1;
        } else if (randomIndex <= ArmourExcOptionsDropRates.defenseSuccessRateIncrease 
            + ArmourExcOptionsDropRates.goldAfterMonsterIncrease
            + ArmourExcOptionsDropRates.reflectDamage
            + ArmourExcOptionsDropRates.maxLifeIncrease
            + ArmourExcOptionsDropRates.maxManaIncrease) {
            item.maxManaIncrease = 1;
        } else if (randomIndex <= ArmourExcOptionsDropRates.defenseSuccessRateIncrease 
            + ArmourExcOptionsDropRates.goldAfterMonsterIncrease
            + ArmourExcOptionsDropRates.reflectDamage
            + ArmourExcOptionsDropRates.maxLifeIncrease
            + ArmourExcOptionsDropRates.maxManaIncrease
            + ArmourExcOptionsDropRates.hpRecoveryRateIncrease) {
            item.hpRecoveryRateIncrease = 1;
        } else if (randomIndex <= ArmourExcOptionsDropRates.defenseSuccessRateIncrease 
            + ArmourExcOptionsDropRates.goldAfterMonsterIncrease
            + ArmourExcOptionsDropRates.reflectDamage
            + ArmourExcOptionsDropRates.maxLifeIncrease
            + ArmourExcOptionsDropRates.maxManaIncrease
            + ArmourExcOptionsDropRates.hpRecoveryRateIncrease
            + ArmourExcOptionsDropRates.mpRecoveryRateIncrease) {
            item.mpRecoveryRateIncrease = 1;
        } else {
            item.decreaseDamageRateIncrease = 1;
        }

        return item;
    }



    function addExcelentOption(ItemAttributes memory item) internal returns (ItemAttributes memory) {
        item.itemLevel = 0;
        if (item.isWeapon) {
            item.skill = true;
            // choose from weapon exc options            
            item = getRandomWeaponExcOption(item, 11);
            if (getRandomNumber(12) <= WeaponsExcOptionsDropRates.doubleOptionChance) {
                item = getRandomWeaponExcOption(item, 13);
            }
        } else if (item.isArmour) {
            // choose from armout exc options
            item = getRandomArmourExcOption(item, 14);
            if (getRandomNumber(15) <= ArmourExcOptionsDropRates.doubleOptionChance) {
                item = getRandomArmourExcOption(item, 16);
            }
        }

        return item;
    }

    function getRandomWeaponExcOption(ItemAttributes memory item,  uint256 seed) internal returns (ItemAttributes memory) {
        uint256 randomNumber = getRandomNumber(seed);
        uint256 totalProbability = WeaponsExcOptionsDropRates.lifeAfterMonsterIncrease +
            WeaponsExcOptionsDropRates.manaAfterMonsterIncrease +
            WeaponsExcOptionsDropRates.excellentDamageProbabilityIncrease +
            WeaponsExcOptionsDropRates.attackSpeedIncrease +
            WeaponsExcOptionsDropRates.damageIncrease;
        

        uint256 randomIndex = randomNumber % totalProbability;

        if (randomIndex <= WeaponsExcOptionsDropRates.lifeAfterMonsterIncrease) {
            item.lifeAfterMonsterIncrease = 1;
        } else if (randomIndex <= WeaponsExcOptionsDropRates.lifeAfterMonsterIncrease 
            + WeaponsExcOptionsDropRates.manaAfterMonsterIncrease) {
            item.manaAfterMonsterIncrease = 1;
        } else if (randomIndex <= WeaponsExcOptionsDropRates.lifeAfterMonsterIncrease 
            + WeaponsExcOptionsDropRates.manaAfterMonsterIncrease
            + WeaponsExcOptionsDropRates.excellentDamageProbabilityIncrease) {
            item.excellentDamageProbabilityIncrease = 1;
        } else if (randomIndex <= WeaponsExcOptionsDropRates.lifeAfterMonsterIncrease 
            + WeaponsExcOptionsDropRates.manaAfterMonsterIncrease
            + WeaponsExcOptionsDropRates.excellentDamageProbabilityIncrease
            + WeaponsExcOptionsDropRates.attackSpeedIncrease) {
            item.attackSpeedIncrease = 1;
        } else {
            item.damageIncrease = 1;
        }

        return item;
    }

}