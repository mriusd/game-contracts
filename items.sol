// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract Items is ERC721 {
    using Counters for Counters.Counter;

    address public owner;

    struct ItemAttributes {
        uint itemLevel;
        uint durability;
        uint classRequired;
        uint strengthRequired;
        uint agilityRequired;
        uint energyRequired;
        uint vitalityRequired;
        string slot;
        uint itemWidth;
        uint itemHeight;
        string acceptableSlot;
        uint physicalDamage;
        uint magicDamage;
        uint defense;
        uint attackSpeed;
        uint defenseSuccessRate;
        uint additionalDamage;
        uint additionalDefense;
        bool increasedExperienceGain;
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
    }

    constructor() ERC721("Combats", "Item") {
        owner = msg.sender;
    }
}