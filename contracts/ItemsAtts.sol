// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;


contract ItemsAtts {

    uint generalDropRate = 100;

    uint maxAdditionalPoints = 28;
    uint jolSuccessRate = 50;
    uint josSuccessRate = 50;
    uint luckDropRate = 30;
    uint skillDropRate = 30;

    uint256 goldTokenId = 2;
    uint256 goldItemId = 1;
    bytes32 dummyHash = bytes32("");

    struct TokenAttributes {
        string name;

        uint256 tokenId;        
        uint256 itemLevel;

        uint256 additionalDamage;
        uint256 additionalDefense;
    
        uint256 fighterId;
        uint256 lastUpdBlock;
        
        uint256 packSize;

        bool luck;
        bool skill;
    }

    struct ItemAttributes {
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
    }
}