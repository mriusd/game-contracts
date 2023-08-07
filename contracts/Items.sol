// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

import "./SafeMath.sol";
import "./ItemsBase.sol";
import "./ItemsAtts.sol";

contract Items is ERC721Enumerable, ItemsAtts, SafeMath {
    ItemsBase private _base;
    using Counters for Counters.Counter;

    address public owner;

    constructor(address ItemsBaseContract) ERC721("MRIUSD", "Item") {
        _base = ItemsBase(ItemsBaseContract);
        owner = msg.sender;

        // create empty item
        ItemsBase.BaseItemAtts memory emptyBaseAtts;
        emptyBaseAtts.name = "Empty item";
        createItem(emptyBaseAtts);
        generateInitialTokens(emptyBaseAtts);

        // create item for gold
        emptyBaseAtts.name = "Gold";
        createItem(emptyBaseAtts);
        generateInitialTokens(emptyBaseAtts);
    }







   
    // event LogError(uint8, uint256);
    
    // event ItemCrafted(uint256 tokenId, address owner);
    // event ItemBoughtFromShop(uint256 tokenId, string itemName, address owner);

    event ItemLevelUpgrade(uint256 tokenId, uint256 newLevel);
    // event ItemAddPointsUpdate(uint256 tokenId, uint256 newAddPoints);

    // function setAdditionalPoints(uint256 tokenId, uint256 points) external {
    //     require(points <= maxAdditionalPoints, "Max points reached");

    //     ItemAttributes memory item = getTokenAttributes(tokenId);

    //     require(item.isWeapon || item.isArmour || item.isWings , "Max points reached");
    //     if (item.isWeapon || item.isWings) {
    //         item.additionalDamage = points;
    //     } else if (item.isArmour) { 
    //         item.additionalDefense = points;
    //     } 

    //     emit ItemAddPointsUpdate(tokenId, item.additionalDefense + item.additionalDamage);
    // }

    function generateInitialTokens(ItemsBase.BaseItemAtts memory itemAtts) internal {
        _tokenIdCounter.increment();
        uint256 newTokenId = _tokenIdCounter.current();

        _safeMint(msg.sender, newTokenId);

        _tokenAttributes[newTokenId].name = itemAtts.name;      
        _tokenAttributes[newTokenId].tokenId = newTokenId;      
        _tokenAttributes[newTokenId].fighterId = 0;      
        _tokenAttributes[newTokenId].lastUpdBlock = block.number; 
    }

    function setItemLevel(uint256 tokenId, uint256 level)  external {
        ItemAttributes memory item = getTokenAttributes(tokenId);

        require(level <= item.maxLevel, "Max item level reached");
        require(level <= item.itemLevel+1, "Item can be upgrade one level at a time only");

        _tokenAttributes[tokenId].itemLevel = level;
        emit ItemLevelUpgrade(tokenId, _tokenAttributes[tokenId].itemLevel);  
    }

    function burnItem(uint256 tokenId) external {
        require(tokenId > 1, "Token cannot be burnt");
        require(_exists(tokenId), "Token doesn't exist");
        _burn(tokenId);
    }  

    function safeMint(address ownerAddres) external returns (uint256) {
        _tokenIdCounter.increment();
        uint256 newTokenId = _tokenIdCounter.current();
        _safeMint(ownerAddres, newTokenId);

        return newTokenId;
    }

    function setTokenAttributes(uint256 tokenId, ItemAttributes memory atts) external {
        _setTokenAttributes(tokenId, atts);
    }

    function itemExists(uint256 tokenId) public view returns(bool) {
        return _exists(tokenId);
    }

    function createItem(ItemsBase.BaseItemAtts memory atts) public returns (string memory name) {                
        _base.createBaseItem(atts);        
        return atts.name;
    }

    function transferItem(uint256 tokenId, address to) external {
        require(_exists(tokenId), "Token does not exist");
        address from = _ownerOf(tokenId);
        _transfer(from, to, tokenId);
    }   

    // event PackedItems(string itemName, uint256 packSize, uint256 fighterId, uint256 newTokenId);
    // event UnpackedItems(uint256 tokenId);

    // function packItems(string calldata itemName, uint256[] memory tokenIds, uint256 packSize, uint256 fighterId) external returns (uint256 newTokenId) {
    //     ItemAttributes memory atts = getItemAttributes(itemName);
    //     require(atts.isPackable, "Item not packable");

    //     // Ensure that the number of tokens is equal to the packSize
    //     require(tokenIds.length == packSize, "Invalid pack size");
        
    //     for (uint i = 0; i < tokenIds.length; i++) {
    //         require(keccak256(abi.encodePacked(_tokenAttributes[tokenIds[i]].name)) == keccak256(abi.encodePacked(itemName)), "Invalid item");
    //         require(_tokenAttributes[tokenIds[i]].fighterId == fighterId, "Item does not belong to fighter");

    //         _burn(tokenIds[i]);
    //     }

    //     _tokenIdCounter.increment();
    //     newTokenId = _tokenIdCounter.current();
    //     _safeMint(_fightersHelper.getOwner(fighterId), newTokenId);

    //     _setTokenAttributes(newTokenId, atts);

    //     _tokenAttributes[newTokenId].name = itemName;      
    //     _tokenAttributes[newTokenId].tokenId = newTokenId;      
    //     _tokenAttributes[newTokenId].fighterId = fighterId;      
    //     _tokenAttributes[newTokenId].lastUpdBlock = block.number; 
    //     _tokenAttributes[newTokenId].packSize = packSize; 

    //     emit PackedItems(itemName, packSize, fighterId, newTokenId);

    //     // Return the ItemAttributes object
    //     return newTokenId;
    // }

    // function unpackItems(uint256 tokenId) external returns(uint256[] memory tokenIds) {
    //     ItemAttributes memory atts = getTokenAttributes(tokenId);
    //     require(atts.packSize > 1, "Item not packed");

    //     // Initialize tokenIds with size atts.packSize
    //     tokenIds = new uint256[](atts.packSize);

    //     for (uint i = 0; i < atts.packSize; i++) {
    //         _tokenIdCounter.increment();
    //         uint256 newTokenId = _tokenIdCounter.current();
    //         _safeMint(_fightersHelper.getOwner(atts.fighterId), newTokenId);

    //         _setTokenAttributes(newTokenId, atts);

    //         _tokenAttributes[newTokenId].tokenId = newTokenId;      
    //         _tokenAttributes[newTokenId].fighterId = atts.fighterId;      
    //         _tokenAttributes[newTokenId].lastUpdBlock = block.number; 
    //         _tokenAttributes[newTokenId].packSize = 1; 

    //         tokenIds[i] = newTokenId;
    //     }

    //     emit UnpackedItems(tokenId);

    //     return tokenIds;

    // }

    // function craftItem(uint256 itemId, address itemOwner, uint256 maxLevel, uint256 maxAddPoints) public returns (uint256) {
    //     require(_itemAttributes.length > itemId, "Item attributes not found");
    //     ItemAttributes memory itemAtts = _itemAttributes[itemId];

    //     if ((itemAtts.isWeapon || itemAtts.isArmour) && getRandomNumber(50) < luckDropRate)
    //     {
    //         itemAtts.luck = true;

    //     }

    //     if (itemAtts.isWeapon && getRandomNumber(51) < skillDropRate)
    //     {
    //         itemAtts.skill = true;
    //     }

    //     if (itemAtts.isWeapon || itemAtts.isArmour)
    //     {
    //         if (maxLevel > 0)
    //         {
    //             itemAtts.itemLevel = getRandomNumber(53) % (maxLevel + 1);
    //         }

    //         if (maxAddPoints > 0)
    //         {
    //             if (itemAtts.isWeapon) {
    //                 itemAtts.additionalDamage = (getRandomNumber(54) % (maxAddPoints/4 + 1)) * 4;
    //             } else if (itemAtts.isArmour) {
    //                 itemAtts.additionalDefense = (getRandomNumber(54) % (maxAddPoints/4 + 1)) * 4;
    //             }
    //         }                
    //     }        

    //     _tokenIdCounter.increment();
    //     uint256 newTokenId = _tokenIdCounter.current();
        
    //     _safeMint(itemOwner, newTokenId);
 
    //     _setTokenAttributes(newTokenId, itemAtts);

    //     _tokenAttributes[newTokenId].tokenId = newTokenId;      
    //     _tokenAttributes[newTokenId].itemAttributesId = itemId;      
    //     _tokenAttributes[newTokenId].fighterId = 0;      
    //     _tokenAttributes[newTokenId].lastUpdBlock = block.number; 

    //     emit ItemCrafted(newTokenId, itemOwner);

    //     return newTokenId;
    // }     

    // function getWeaponsLength(uint256 rarityLevel) public view returns (uint256) {
    //     return Weapons[rarityLevel].length;
    // }

    // function getArmoursLength(uint256 rarityLevel) public view returns (uint256) {
    //     return Armours[rarityLevel].length;
    // }

    // function getJewelsLength(uint256 rarityLevel) public view returns (uint256) {
    //     return Jewels[rarityLevel].length;
    // }

    // function getMiscsLength(uint256 rarityLevel) public view returns (uint256) {
    //     return Misc[rarityLevel].length;
    // }    

    // function buyItemFromShop(string calldata itemName, uint256 fighterId) external {
    //     require(baseItemAttributes[itemName].inShop, "Item not in shop or doesn't exist");

    //     // money logic

    //     _tokenIdCounter.increment();
    //     uint256 tokenId = _tokenIdCounter.current();
    //     _safeMint(_fightersHelper.getOwner(fighterId), tokenId);

    //     _tokenAttributes[tokenId].tokenId = tokenId;      
    //     _tokenAttributes[tokenId].name = itemName;      
    //     _tokenAttributes[tokenId].fighterId = fighterId;      
    //     _tokenAttributes[tokenId].lastUpdBlock = block.number;      

    //     emit ItemBoughtFromShop(tokenId, itemName, _fightersHelper.getOwner(fighterId));
    // }

    function getTokenAttributes(uint256 tokenId) public view returns (ItemAttributes memory) {
        require(_exists(tokenId), "Token does not exist");

        TokenAttributes memory tokenAtts = _tokenAttributes[tokenId];
        ItemsBase.BaseItemAtts memory baseAtts = _base.gatBaseItemAtts(tokenAtts.name);
        
        return ItemAttributes({
            name: baseAtts.name,
            tokenId: tokenAtts.tokenId,
            itemLevel: tokenAtts.itemLevel,
            maxLevel: baseAtts.maxLevel,
            durability: baseAtts.durability,
            classRequired: baseAtts.classRequired,
            strengthRequired: baseAtts.strengthRequired,
            agilityRequired: baseAtts.agilityRequired,
            energyRequired: baseAtts.energyRequired,
            vitalityRequired: baseAtts.vitalityRequired,
            itemWidth: baseAtts.itemWidth,
            itemHeight: baseAtts.itemHeight,
            acceptableSlot1: baseAtts.acceptableSlot1,
            acceptableSlot2: baseAtts.acceptableSlot2,
            baseMinPhysicalDamage: baseAtts.baseMinPhysicalDamage,
            baseMaxPhysicalDamage: baseAtts.baseMaxPhysicalDamage,
            baseMinMagicDamage: baseAtts.baseMinMagicDamage,
            baseMaxMagicDamage: baseAtts.baseMaxMagicDamage,
            baseDefense: baseAtts.baseDefense,
            attackSpeed: baseAtts.attackSpeed,
            additionalDamage: tokenAtts.additionalDamage,
            additionalDefense: tokenAtts.additionalDefense,
            fighterId: tokenAtts.fighterId,
            lastUpdBlock: tokenAtts.lastUpdBlock,
            itemRarityLevel: baseAtts.itemRarityLevel,
            packSize: tokenAtts.packSize,
            luck: tokenAtts.luck,
            skill: tokenAtts.skill,
            isPackable: baseAtts.isPackable,
            isBox: baseAtts.isBox,
            isWeapon: baseAtts.isWeapon,
            isArmour: baseAtts.isArmour,
            isJewel: baseAtts.isJewel,
            isWings: baseAtts.isWings,
            isMisc: baseAtts.isMisc,
            isConsumable: baseAtts.isConsumable,
            inShop: baseAtts.inShop
        });
    }


    // Get the attributes for a fighter NFT
    function getItemAttributes(string memory name) public view returns (ItemAttributes memory) {
        ItemsBase.BaseItemAtts memory baseAtts = _base.gatBaseItemAtts(name);
        require(baseAtts.itemWidth > 0, "Item does not exist");

        ItemAttributes memory itemAtts;

        return ItemAttributes({
            name: baseAtts.name,
            tokenId: 0,
            itemLevel: 0,
            maxLevel: baseAtts.maxLevel,
            durability: baseAtts.durability,
            classRequired: baseAtts.classRequired,
            strengthRequired: baseAtts.strengthRequired,
            agilityRequired: baseAtts.agilityRequired,
            energyRequired: baseAtts.energyRequired,
            vitalityRequired: baseAtts.vitalityRequired,
            itemWidth: baseAtts.itemWidth,
            itemHeight: baseAtts.itemHeight,
            acceptableSlot1: baseAtts.acceptableSlot1,
            acceptableSlot2: baseAtts.acceptableSlot2,
            baseMinPhysicalDamage: baseAtts.baseMinPhysicalDamage,
            baseMaxPhysicalDamage: baseAtts.baseMaxPhysicalDamage,
            baseMinMagicDamage: baseAtts.baseMinMagicDamage,
            baseMaxMagicDamage: baseAtts.baseMaxMagicDamage,
            baseDefense: baseAtts.baseDefense,
            attackSpeed: baseAtts.attackSpeed,
            additionalDamage: 0,
            additionalDefense: 0,
            fighterId: 0,
            lastUpdBlock: 0,
            itemRarityLevel: baseAtts.itemRarityLevel,
            packSize: 1,
            luck: false,
            skill: false,
            isPackable: baseAtts.isPackable,
            isBox: baseAtts.isBox,
            isWeapon: baseAtts.isWeapon,
            isArmour: baseAtts.isArmour,
            isJewel: baseAtts.isJewel,
            isWings: baseAtts.isWings,
            isMisc: baseAtts.isMisc,
            isConsumable: baseAtts.isConsumable,
            inShop: baseAtts.inShop
        });  
    }

    function burnConsumable(uint256 tokenId) external {
        ItemAttributes memory atts = getTokenAttributes(tokenId);

        require(atts.isConsumable, "Item not a consumable");
        _burn(tokenId);
    }

    mapping (uint256 => TokenAttributes) private _tokenAttributes;
    ItemAttributes[] private _itemAttributes;
    Counters.Counter private _tokenIdCounter;

    ItemAttributes boxAttributes;



    function getUserItems(address userAddress) external view returns (uint256[] memory) {
        
        uint256 numTokens = balanceOf(userAddress);
        uint256[] memory tokenIds = new uint256[](numTokens);

        for (uint256 i = 0; i < numTokens; i++) {
            tokenIds[i] = tokenOfOwnerByIndex(userAddress, i);
        }

        return tokenIds;
    }


    function getFighterItems(address userAddress, uint256 fighterId) external view returns (uint256[2][] memory) {
        
        uint256 numTokens = balanceOf(userAddress);
        uint256[2][] memory tokenIds = new uint256[2][](numTokens);

        uint256 id;
        uint counter = 0;

        for (uint256 i = 0; i < numTokens; i++) {
            id = tokenOfOwnerByIndex(userAddress, i);
            if (_tokenAttributes[id].fighterId == fighterId)
            {
                tokenIds[counter] = [id, _tokenAttributes[id].lastUpdBlock];
            }
            
            counter++;
        }

        return tokenIds;
    }


    function _setTokenAttributes(uint256 tokenId, ItemAttributes memory atts) internal {
        require(_exists(tokenId), "Token does not exist");

        _tokenAttributes[tokenId].itemLevel = atts.itemLevel;
        _tokenAttributes[tokenId].additionalDamage = atts.additionalDamage;
        _tokenAttributes[tokenId].additionalDefense = atts.additionalDefense;
        _tokenAttributes[tokenId].fighterId = atts.fighterId;
        _tokenAttributes[tokenId].lastUpdBlock = atts.lastUpdBlock;
        _tokenAttributes[tokenId].packSize = atts.packSize;
        _tokenAttributes[tokenId].luck = atts.luck;
        _tokenAttributes[tokenId].skill = atts.skill;
    }

}