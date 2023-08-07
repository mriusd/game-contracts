// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;


contract ItemsBase {



    mapping (string => BaseItemAtts) public baseItemAttributes;

    mapping (uint256 => string[]) public Weapons;
    mapping (uint256 => string[]) public Armours;
    mapping (uint256 => string[]) public Jewels;
    mapping (uint256 => string[]) public Wings;
    mapping (uint256 => string[]) public Boxes; 
    mapping (uint256 => string[]) public Consumables; // potions
    mapping (uint256 => string[]) public Misc; // scrolls

    struct BaseItemAtts {
        string name;

        
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
        
        uint256 itemRarityLevel;

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

    event ItemGenerated(string name, ItemsBase.BaseItemAtts atts);

    function getWeapons(uint256 rarityLevel) public view returns(string[] memory) {
        return Weapons[rarityLevel];
    }
    function getArmours(uint256 rarityLevel) public view returns(string[] memory) {
        return Armours[rarityLevel];
    }
    function getJewels(uint256 rarityLevel) public view returns(string[] memory) {
        return Jewels[rarityLevel];
    }
    function getWings(uint256 rarityLevel) public view returns(string[] memory) {
        return Wings[rarityLevel];
    }
    function getBoxes(uint256 rarityLevel) public view returns(string[] memory) {
        return Boxes[rarityLevel];
    }
    function getConsumables(uint256 rarityLevel) public view returns(string[] memory) {
        return Consumables[rarityLevel];
    }
    function getMisc(uint256 rarityLevel) public view returns(string[] memory) {
        return Misc[rarityLevel];
    }


    function addWeapons(uint256 rarityLevel, string memory name) public returns(string[] memory) {
        Weapons[rarityLevel].push(name);
    }
    function addArmours(uint256 rarityLevel, string memory name) public returns(string[] memory) {
        Armours[rarityLevel].push(name);
    }
    function addJewels(uint256 rarityLevel, string memory name) public returns(string[] memory) {
        Jewels[rarityLevel].push(name);
    }
    function addWings(uint256 rarityLevel, string memory name) public returns(string[] memory) {
        Wings[rarityLevel].push(name);
    }
    function addBoxes(uint256 rarityLevel, string memory name) public returns(string[] memory) {
        Boxes[rarityLevel].push(name);
    }
    function addConsumables(uint256 rarityLevel, string memory name) public returns(string[] memory) {
        Consumables[rarityLevel].push(name);
    }
    function addMisc(uint256 rarityLevel, string memory name) public returns(string[] memory) {
        Misc[rarityLevel].push(name);
    }




    function createBaseItem(BaseItemAtts memory atts) public returns (string memory name) {  
        baseItemAttributes[atts.name] = atts;

        if (atts.isWeapon) {
            Weapons[atts.itemRarityLevel].push(atts.name);
        } else if (atts.isArmour) {
            Armours[atts.itemRarityLevel].push(atts.name);
        } else if (atts.isJewel) {
            Jewels[atts.itemRarityLevel].push(atts.name);
        } else if (atts.isMisc) {
            Misc[atts.itemRarityLevel].push(atts.name);
        } else if (atts.isWings) {
            Wings[atts.itemRarityLevel].push(atts.name);
        } else if (atts.isConsumable) {
            Consumables[atts.itemRarityLevel].push(atts.name);
        } else if (atts.isBox) {
            Boxes[atts.itemRarityLevel].push(atts.name);
        } 

        emit ItemGenerated(atts.name, atts);

        return atts.name;
    }

    function gatBaseItemAtts(string memory name) public view returns (BaseItemAtts memory) {  
        return baseItemAttributes[name];
    }

}