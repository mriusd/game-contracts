// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

contract FightersAtts {

    // Struct for fighter attributes
    struct Attributes {
        string name;
        string class;
        uint256 tokenId;
        uint256 strength;
        uint256 agility;
        uint256 energy;
        uint256 vitality;
        uint256 experience;  

        
        uint256 hpPerVitalityPoint;
        uint256 manaPerEnergyPoint;
        uint256 hpIncreasePerLevel;
        uint256 manaIncreasePerLevel;
        uint256 statPointsPerLevel;
        uint256 attackSpeed;
        uint256 agilityPointsPerSpeed;
        uint256 isNpc;
        uint256 dropRarityLevel; // for npcs                  
    }

    struct FighterStats {
        uint256 tokenId;
        uint256 maxHealth;
        uint256 maxMana;
        uint256 level;
        uint256 exp;
        uint256 totalStatPoints;
        uint256 maxStatPoints;
    }

}