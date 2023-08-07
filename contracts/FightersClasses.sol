// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;


contract FightersClasses {

    // Struct for fighter class attributes
    struct FightersClassAttributes {
        uint256 baseStrength;
        uint256 baseAgility;
        uint256 baseEnergy;
        uint256 baseVitality;

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

    mapping (string => FightersClassAttributes) public FighterClasses;

    function updateFighterClass(string memory className, FightersClassAttributes memory atts) public {
        FighterClasses[className] = atts;
    }

    function fighterClassExists(string memory className) public view returns (bool) {
        if (FighterClasses[className].baseStrength == 0) {
            return false;
        }

        return true;
    }
}