// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "./ItemsExcellentAtts.sol";
import "./Items.sol";
contract ItemsExcellent is ItemsExcellentAtts {  
    Items private _items;

    constructor(address itemsAddress) {
        _items = Items(itemsAddress);
    }

    struct WingsExcOptions {
        uint256 increaseAttackSpeedPoints;
        uint256 reflectDamagePercent;
        uint256 restoreHPChance;
        uint256 restoreMPChance;
        uint256 doubleDamageChance;
        uint256 ignoreOpponentDefenseChance;

        uint256 doubleOptionChance;
    }

    WingsExcOptions WingsExcOptionsDropRates = WingsExcOptions({
        increaseAttackSpeedPoints: 15,
        reflectDamagePercent: 15,
        restoreHPChance: 30,
        restoreMPChance: 30,
        doubleDamageChance: 5,
        ignoreOpponentDefenseChance: 3,

        doubleOptionChance: 15
    });


    struct WeaponsExcOptions {
        uint256 lifeAfterMonsterIncrease;
        uint256 manaAfterMonsterIncrease;
        uint256 excellentDamageProbabilityIncrease;
        uint256 attackSpeedIncrease;
        uint256 attackLvl20; 
        uint256 attackIncreasePercent; 

        uint256 doubleOptionChance;
    }

    WeaponsExcOptions WeaponsExcOptionsDropRates = WeaponsExcOptions({
        lifeAfterMonsterIncrease: 35,
        manaAfterMonsterIncrease: 35,
        excellentDamageProbabilityIncrease: 10,
        attackSpeedIncrease: 10,
        attackLvl20: 10,
        attackIncreasePercent: 10,

        doubleOptionChance: 15
    });

    struct ArmoursExcOptions {
        uint256 defenseSuccessRateIncrease;
        uint256 goldAfterMonsterIncrease;
        uint256 reflectDamage;
        uint256 maxLifeIncrease;
        uint256 maxManaIncrease;
        uint256 hpRecoveryRateIncrease;
        uint256 mpRecoveryRateIncrease;
        uint256 decreaseDamageRateIncrease;

        uint256 doubleOptionChance;
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


    mapping (uint256 => WingsExcOptions) public WingsExcOpts;
    mapping (uint256 => WeaponsExcOptions) public WeaponsExcOpts;
    mapping (uint256 => ArmoursExcOptions) public ArmoursExcOpts;

    function convertToExcellent(ItemAttributes memory item) public view returns (ExcellentItemAtts memory) {
        ExcellentItemAtts memory excellentItem;

        // Copy original properties
        excellentItem.name = item.name;
        excellentItem.tokenId = item.tokenId;        
        excellentItem.itemLevel = item.itemLevel;
        excellentItem.maxLevel = item.maxLevel;
        excellentItem.durability = item.durability;
        excellentItem.classRequired = item.classRequired;
        excellentItem.strengthRequired = item.strengthRequired;
        excellentItem.agilityRequired = item.agilityRequired;
        excellentItem.energyRequired = item.energyRequired;
        excellentItem.vitalityRequired = item.vitalityRequired;
        excellentItem.itemWidth = item.itemWidth;
        excellentItem.itemHeight = item.itemHeight;
        excellentItem.acceptableSlot1 = item.acceptableSlot1;
        excellentItem.acceptableSlot2 = item.acceptableSlot2;
        excellentItem.baseMinPhysicalDamage = item.baseMinPhysicalDamage;
        excellentItem.baseMaxPhysicalDamage = item.baseMaxPhysicalDamage;
        excellentItem.baseMinMagicDamage = item.baseMinMagicDamage;
        excellentItem.baseMaxMagicDamage = item.baseMaxMagicDamage;
        excellentItem.baseDefense = item.baseDefense;
        excellentItem.attackSpeed = item.attackSpeed;
        excellentItem.additionalDamage = item.additionalDamage;
        excellentItem.additionalDefense = item.additionalDefense;
        excellentItem.fighterId = item.fighterId;
        excellentItem.lastUpdBlock = item.lastUpdBlock;
        excellentItem.itemRarityLevel = item.itemRarityLevel;
        excellentItem.itemAttributesId = item.itemAttributesId;
        excellentItem.luck = item.luck;
        excellentItem.skill = item.skill;
        excellentItem.isBox = item.isBox;
        excellentItem.isWeapon = item.isWeapon;
        excellentItem.isArmour = item.isArmour;
        excellentItem.isJewel = item.isJewel;
        excellentItem.isWings = item.isWings;
        excellentItem.isMisc = item.isMisc;
        excellentItem.inShop = item.inShop;

        // Wings
        excellentItem.increaseAttackSpeedPoints = WingsExcOpts[item.tokenId].increaseAttackSpeedPoints;
        excellentItem.reflectDamagePercent = WingsExcOpts[item.tokenId].reflectDamagePercent;
        excellentItem.restoreHPChance = WingsExcOpts[item.tokenId].restoreHPChance;
        excellentItem.restoreMPChance = WingsExcOpts[item.tokenId].restoreMPChance;
        excellentItem.doubleDamageChance = WingsExcOpts[item.tokenId].doubleDamageChance;
        excellentItem.ignoreOpponentDefenseChance = WingsExcOpts[item.tokenId].ignoreOpponentDefenseChance;

        // Weapons
        excellentItem.lifeAfterMonsterIncrease = WeaponsExcOpts[item.tokenId].lifeAfterMonsterIncrease;
        excellentItem.manaAfterMonsterIncrease = WeaponsExcOpts[item.tokenId].manaAfterMonsterIncrease;
        excellentItem.excellentDamageProbabilityIncrease = WeaponsExcOpts[item.tokenId].excellentDamageProbabilityIncrease;
        excellentItem.attackSpeedIncrease = WeaponsExcOpts[item.tokenId].attackSpeedIncrease;
        excellentItem.attackLvl20 = WeaponsExcOpts[item.tokenId].attackLvl20;
        excellentItem.attackIncreasePercent = WeaponsExcOpts[item.tokenId].attackIncreasePercent;

        // Armour
        excellentItem.defenseSuccessRateIncrease = ArmoursExcOpts[item.tokenId].defenseSuccessRateIncrease;
        excellentItem.goldAfterMonsterIncrease = ArmoursExcOpts[item.tokenId].goldAfterMonsterIncrease;
        excellentItem.reflectDamage = ArmoursExcOpts[item.tokenId].reflectDamage;
        excellentItem.maxLifeIncrease = ArmoursExcOpts[item.tokenId].maxLifeIncrease;
        excellentItem.maxManaIncrease = ArmoursExcOpts[item.tokenId].maxManaIncrease;
        excellentItem.hpRecoveryRateIncrease = ArmoursExcOpts[item.tokenId].hpRecoveryRateIncrease;
        excellentItem.mpRecoveryRateIncrease = ArmoursExcOpts[item.tokenId].mpRecoveryRateIncrease;
        excellentItem.decreaseDamageRateIncrease = ArmoursExcOpts[item.tokenId].decreaseDamageRateIncrease;

        return excellentItem;
    }

    function assignExcellentAtts(ExcellentItemAtts memory a, uint256 tokenId) internal returns (ExcellentItemAtts memory) {
        a.increaseAttackSpeedPoints = WingsExcOpts[tokenId].increaseAttackSpeedPoints;
        a.reflectDamagePercent = WingsExcOpts[tokenId].reflectDamagePercent;
        a.restoreHPChance = WingsExcOpts[tokenId].restoreHPChance;
        a.restoreMPChance = WingsExcOpts[tokenId].restoreMPChance;
        a.doubleDamageChance = WingsExcOpts[tokenId].doubleDamageChance;
        a.ignoreOpponentDefenseChance = WingsExcOpts[tokenId].ignoreOpponentDefenseChance;


        a.lifeAfterMonsterIncrease = WeaponsExcOpts[tokenId].lifeAfterMonsterIncrease;
        a.manaAfterMonsterIncrease = WeaponsExcOpts[tokenId].manaAfterMonsterIncrease;
        a.excellentDamageProbabilityIncrease = WeaponsExcOpts[tokenId].excellentDamageProbabilityIncrease;
        a.attackSpeedIncrease = WeaponsExcOpts[tokenId].attackSpeedIncrease;
        a.attackLvl20 = WeaponsExcOpts[tokenId].attackLvl20;
        a.attackIncreasePercent = WeaponsExcOpts[tokenId].attackIncreasePercent;


        a.defenseSuccessRateIncrease = ArmoursExcOpts[tokenId].defenseSuccessRateIncrease;
        a.goldAfterMonsterIncrease = ArmoursExcOpts[tokenId].goldAfterMonsterIncrease;
        a.reflectDamage = ArmoursExcOpts[tokenId].reflectDamage;
        a.maxLifeIncrease = ArmoursExcOpts[tokenId].maxLifeIncrease;
        a.maxManaIncrease = ArmoursExcOpts[tokenId].maxManaIncrease;
        a.hpRecoveryRateIncrease = ArmoursExcOpts[tokenId].hpRecoveryRateIncrease;
        a.mpRecoveryRateIncrease = ArmoursExcOpts[tokenId].mpRecoveryRateIncrease;
        a.decreaseDamageRateIncrease = ArmoursExcOpts[tokenId].decreaseDamageRateIncrease;

        return a;
    }

    function convertToItemAttributes(ExcellentItemAtts memory eia) internal pure returns (ItemAttributes memory) {
        ItemAttributes memory ia;

        ia.name = eia.name;
        ia.tokenId = eia.tokenId;
        ia.itemLevel = eia.itemLevel;
        ia.maxLevel = eia.maxLevel;
        ia.durability = eia.durability;
        ia.classRequired = eia.classRequired;
        ia.strengthRequired = eia.strengthRequired;
        ia.agilityRequired = eia.agilityRequired;
        ia.energyRequired = eia.energyRequired;
        ia.vitalityRequired = eia.vitalityRequired;
        ia.itemWidth = eia.itemWidth;
        ia.itemHeight = eia.itemHeight;
        ia.acceptableSlot1 = eia.acceptableSlot1;
        ia.acceptableSlot2 = eia.acceptableSlot2;

        ia.baseMinPhysicalDamage = eia.baseMinPhysicalDamage;
        ia.baseMaxPhysicalDamage = eia.baseMaxPhysicalDamage;
        ia.baseMinMagicDamage = eia.baseMinMagicDamage;
        ia.baseMaxMagicDamage = eia.baseMaxMagicDamage;
        ia.baseDefense = eia.baseDefense;
        ia.attackSpeed = eia.attackSpeed;

        ia.additionalDamage = eia.additionalDamage;
        ia.additionalDefense = eia.additionalDefense;

        ia.fighterId = eia.fighterId;
        ia.lastUpdBlock = eia.lastUpdBlock;
        ia.itemRarityLevel = eia.itemRarityLevel;

        ia.itemAttributesId = eia.itemAttributesId;

        ia.luck = eia.luck;
        ia.skill = eia.skill;
        ia.isBox = eia.isBox;
        ia.isWeapon = eia.isWeapon;
        ia.isArmour = eia.isArmour;
        ia.isJewel = eia.isJewel;
        ia.isWings = eia.isWings;
        ia.isMisc = eia.isMisc;
        ia.inShop = eia.inShop;

        return ia;
    }


    function reverseAssignExcellentAtts(uint256 tokenId, ExcellentItemAtts memory a) internal {
        WingsExcOpts[tokenId].increaseAttackSpeedPoints = a.increaseAttackSpeedPoints;
        WingsExcOpts[tokenId].reflectDamagePercent = a.reflectDamagePercent;
        WingsExcOpts[tokenId].restoreHPChance = a.restoreHPChance;
        WingsExcOpts[tokenId].restoreMPChance = a.restoreMPChance;
        WingsExcOpts[tokenId].doubleDamageChance = a.doubleDamageChance;
        WingsExcOpts[tokenId].ignoreOpponentDefenseChance = a.ignoreOpponentDefenseChance;

        WeaponsExcOpts[tokenId].lifeAfterMonsterIncrease = a.lifeAfterMonsterIncrease;
        WeaponsExcOpts[tokenId].manaAfterMonsterIncrease = a.manaAfterMonsterIncrease;
        WeaponsExcOpts[tokenId].excellentDamageProbabilityIncrease = a.excellentDamageProbabilityIncrease;
        WeaponsExcOpts[tokenId].attackSpeedIncrease = a.attackSpeedIncrease;
        WeaponsExcOpts[tokenId].attackLvl20 = a.attackLvl20;
        WeaponsExcOpts[tokenId].attackIncreasePercent = a.attackIncreasePercent;

        ArmoursExcOpts[tokenId].defenseSuccessRateIncrease = a.defenseSuccessRateIncrease;
        ArmoursExcOpts[tokenId].goldAfterMonsterIncrease = a.goldAfterMonsterIncrease;
        ArmoursExcOpts[tokenId].reflectDamage = a.reflectDamage;
        ArmoursExcOpts[tokenId].maxLifeIncrease = a.maxLifeIncrease;
        ArmoursExcOpts[tokenId].maxManaIncrease = a.maxManaIncrease;
        ArmoursExcOpts[tokenId].hpRecoveryRateIncrease = a.hpRecoveryRateIncrease;
        ArmoursExcOpts[tokenId].mpRecoveryRateIncrease = a.mpRecoveryRateIncrease;
        ArmoursExcOpts[tokenId].decreaseDamageRateIncrease = a.decreaseDamageRateIncrease;
    }




    function getItemAttributes(uint256 itemId) external returns (ExcellentItemAtts memory) {
        ItemAttributes memory itemAtts = _items.getItemAttributes(itemId);
        return convertToExcellent(itemAtts);
    }

    function getTokenAttributes(uint256 tokenId) external returns (ExcellentItemAtts memory) {
        ItemAttributes memory tokenAtts = _items.getTokenAttributes(tokenId);
        ExcellentItemAtts memory excAtts = convertToExcellent(tokenAtts);

        return assignExcellentAtts(excAtts, tokenId);
    }

    function setTokenAttributes(uint256 tokenId, ExcellentItemAtts memory atts) external {
        require(_items.itemExists(tokenId), "Token does not exist");

        ItemAttributes memory baseAtts = convertToItemAttributes(atts);
        _items.setTokenAttributes(tokenId, baseAtts);
        reverseAssignExcellentAtts(tokenId, atts);
    }


    function getRandomArmourExcOption(ExcellentItemAtts memory item,  uint256 seed) internal returns (ExcellentItemAtts memory) {
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



    function addExcellentOption(ExcellentItemAtts memory excItem) external returns (ExcellentItemAtts memory newItem) {


        excItem.itemLevel = 0;
        if (excItem.isWeapon) {
            excItem.skill = true;
            // choose from weapon exc options            
            excItem = getRandomWeaponExcOption(excItem, 11);
            if (getRandomNumber(12) <= WeaponsExcOptionsDropRates.doubleOptionChance) {
                excItem = getRandomWeaponExcOption(excItem, 13);
            }
        } else if (excItem.isArmour) {
            // choose from armout exc options
            excItem = getRandomArmourExcOption(excItem, 14);
            if (getRandomNumber(15) <= ArmourExcOptionsDropRates.doubleOptionChance) {
                excItem = getRandomArmourExcOption(excItem, 16);
            }
        }

        return excItem;
    }

    function getRandomWeaponExcOption(ExcellentItemAtts memory item,  uint256 seed) internal returns (ExcellentItemAtts memory) {
        uint256 randomNumber = getRandomNumber(seed);
        uint256 totalProbability = WeaponsExcOptionsDropRates.lifeAfterMonsterIncrease +
            WeaponsExcOptionsDropRates.manaAfterMonsterIncrease +
            WeaponsExcOptionsDropRates.excellentDamageProbabilityIncrease +
            WeaponsExcOptionsDropRates.attackSpeedIncrease +
            WeaponsExcOptionsDropRates.attackLvl20 +
            WeaponsExcOptionsDropRates.attackIncreasePercent;
        

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
        } else if (randomIndex <= WeaponsExcOptionsDropRates.lifeAfterMonsterIncrease 
            + WeaponsExcOptionsDropRates.manaAfterMonsterIncrease
            + WeaponsExcOptionsDropRates.excellentDamageProbabilityIncrease
            + WeaponsExcOptionsDropRates.attackSpeedIncrease
            + WeaponsExcOptionsDropRates.attackLvl20 ) {
            item.attackLvl20 = 1;
        } else {
            item.attackIncreasePercent = 1;
        }

        return item;
    }

}