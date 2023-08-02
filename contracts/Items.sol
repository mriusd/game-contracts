// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

import "truffle/Console.sol";
import "./ItemsAtts.sol";
import "./FightersHelper.sol";
import "./MoneyHelper.sol";

contract Items is ERC721Enumerable, ItemsAtts {
    using Counters for Counters.Counter;
    FightersHelper private _fightersHelper;
    MoneyHelper private _moneyHelper;

    address public owner;

    constructor(address fightersHelperContract, address moneyHelperContract) ERC721("MRIUSD", "Item") {
        owner = msg.sender;
        _fightersHelper = FightersHelper(fightersHelperContract);  
        _moneyHelper = MoneyHelper(moneyHelperContract);

        // create empty item
        ItemAttributes memory emptyAtts;
        emptyAtts.name = "Empty item";
        createItem(emptyAtts);
        generateInitialTokens(emptyAtts);

        // create item for gold
        emptyAtts.name = "Gold";
        createItem(emptyAtts);
        generateInitialTokens(emptyAtts);
    }


    mapping (uint256 => uint256[]) public Weapons; // mapping of rarityLevel => tokenId[]
    mapping (uint256 => uint256[]) public Armours; // mapping of rarityLevel => tokenId[]
    mapping (uint256 => uint256[]) public Jewels; // mapping of rarityLevel => tokenId[]
    mapping (uint256 => uint256[]) public Wings; // mapping of rarityLevel => tokenId[]
    mapping (uint256 => uint256[]) public Misc; // mapping of rarityLevel => tokenId[]

    // mapping (bytes32 => uint256) public dropHashes;

    function getWeapons(uint256 rarityLevel) public view returns(uint256[] memory) {
        return Weapons[rarityLevel];
    }
    function getArmours(uint256 rarityLevel) public view returns(uint256[] memory) {
        return Armours[rarityLevel];
    }
    function getJewels(uint256 rarityLevel) public view returns(uint256[] memory) {
        return Jewels[rarityLevel];
    }
    function getMisc(uint256 rarityLevel) public view returns(uint256[] memory) {
        return Misc[rarityLevel];
    }


    event ItemGenerated(uint256 itemId, string name);
    event LogError(uint8, uint256);
    
    event ItemCrafted(uint256 tokenId, address owner);
    event ItemBoughtFromShop(uint256 tokenId, uint256 itemId, address owner, string itemName);

    event ItemLevelUpgrade(uint256 tokenId, uint256 newLevel);
    event ItemAddPointsUpdate(uint256 tokenId, uint256 newAddPoints);

    function setAdditionalPoints(uint256 tokenId, uint256 points) external {
        require(points <= maxAdditionalPoints, "Max points reached");
        require(_tokenAttributes[tokenId].isWeapon || _tokenAttributes[tokenId].isArmour || _tokenAttributes[tokenId].isWings , "Max points reached");
        if (_tokenAttributes[tokenId].isWeapon || _tokenAttributes[tokenId].isWings) {
            _tokenAttributes[tokenId].additionalDamage = points;
        } else if (_tokenAttributes[tokenId].isArmour) { 
            _tokenAttributes[tokenId].additionalDefense = points;
        } 

        emit ItemAddPointsUpdate(tokenId, _tokenAttributes[tokenId].additionalDefense + _tokenAttributes[tokenId].additionalDamage);
    }

    function generateInitialTokens(ItemAttributes memory itemAtts) internal {
        _tokenIdCounter.increment();
        uint256 newTokenId = _tokenIdCounter.current();

        _safeMint(msg.sender, newTokenId);
        _setTokenAttributes(newTokenId, itemAtts);

        _tokenAttributes[newTokenId].tokenId = newTokenId;      
        _tokenAttributes[newTokenId].fighterId = 0;      
        _tokenAttributes[newTokenId].lastUpdBlock = block.number; 
    }

    function setItemLevel(uint256 tokenId, uint256 level)  external {
        require(level <= _tokenAttributes[tokenId].maxLevel, "Max item level reached");
        require(level <= _tokenAttributes[tokenId].itemLevel+1, "Item can be upgrade one level at a time only");

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
        require(_exists(tokenId), "Token does not exist");
        _tokenAttributes[tokenId] = atts;
    }

    function itemExists(uint256 tokenId) public view returns(bool) {
        return _exists(tokenId);
    }

    function createItem(ItemAttributes memory atts) public returns (uint256 tokenId) {                
        _itemAttributes.push(atts);   
         uint256 itemId = _itemAttributes.length - 1;   
         _itemAttributes[itemId].itemAttributesId = itemId;

         if (atts.isWeapon) {
            Weapons[atts.itemRarityLevel].push(itemId);
         } else if (atts.isArmour) {
            Armours[atts.itemRarityLevel].push(itemId);
         } else if (atts.isJewel) {
            Jewels[atts.itemRarityLevel].push(itemId);
         } else if (atts.isMisc) {
            Misc[atts.itemRarityLevel].push(itemId);
         } else if (atts.isWings) {
            Wings[atts.itemRarityLevel].push(itemId);
         }

        emit ItemGenerated(itemId, atts.name);
        return itemId;
    }

    function transferItem(uint256 tokenId, address to) external {
        require(_exists(tokenId), "Token does not exist");
        address from = _ownerOf(tokenId);
        _transfer(from, to, tokenId);
    }   

    event PackedItems(uint256 itemId, uint256 packSize, uint256 fighterId, uint256 newTokenId);
    event UnpackedItems(uint256 tokenId);

    function packItems(uint256 itemId, uint256[] memory tokenIds, uint256 packSize, uint256 fighterId) external returns (uint256 newTokenId) {
        ItemAttributes memory atts = getItemAttributes(itemId);
        require(atts.isPackable, "Item not packable");

        // Ensure that the number of tokens is equal to the packSize
        require(tokenIds.length == packSize, "Invalid pack size");
        
        for (uint i = 0; i < tokenIds.length; i++) {
            // Check that each tokenId has an associated itemAttributeId equal to itemId
            require(_tokenAttributes[tokenIds[i]].itemAttributesId == itemId, "Invalid item");
            require(_tokenAttributes[tokenIds[i]].fighterId == fighterId, "Item does not belong to fighter");

            _burn(tokenIds[i]);
        }

        _tokenIdCounter.increment();
        newTokenId = _tokenIdCounter.current();
        _safeMint(_fightersHelper.getOwner(fighterId), newTokenId);

        _setTokenAttributes(newTokenId, atts);

        _tokenAttributes[newTokenId].tokenId = newTokenId;      
        _tokenAttributes[newTokenId].fighterId = fighterId;      
        _tokenAttributes[newTokenId].lastUpdBlock = block.number; 
        _tokenAttributes[newTokenId].packSize = packSize; 

        emit PackedItems(itemId, packSize, fighterId, newTokenId);

        // Return the ItemAttributes object
        return newTokenId;
    }

    function unpackItems(uint256 tokenId) external returns(uint256[] memory tokenIds) {
        ItemAttributes memory atts = getTokenAttributes(tokenId);
        require(atts.packSize > 1, "Item not packed");

        // Initialize tokenIds with size atts.packSize
        tokenIds = new uint256[](atts.packSize);

        for (uint i = 0; i < atts.packSize; i++) {
            _tokenIdCounter.increment();
            uint256 newTokenId = _tokenIdCounter.current();
            _safeMint(_fightersHelper.getOwner(atts.fighterId), newTokenId);

            _setTokenAttributes(newTokenId, atts);

            _tokenAttributes[newTokenId].tokenId = newTokenId;      
            _tokenAttributes[newTokenId].fighterId = atts.fighterId;      
            _tokenAttributes[newTokenId].lastUpdBlock = block.number; 
            _tokenAttributes[newTokenId].packSize = 1; 

            tokenIds[i] = newTokenId;
        }

        emit UnpackedItems(tokenId);

        return tokenIds;

    }

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

    function buyItemFromShop(uint256 itemId, uint256 fighterId) external {
        require(_itemAttributes[itemId].inShop, "Item not in shop or doesn't exist");

        // money logic

        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        _safeMint(_fightersHelper.getOwner(fighterId), tokenId);
        _setTokenAttributes(tokenId, _itemAttributes[itemId]);

        _tokenAttributes[tokenId].tokenId = tokenId;      
        _tokenAttributes[tokenId].itemAttributesId = itemId;      
        _tokenAttributes[tokenId].fighterId = fighterId;      
        _tokenAttributes[tokenId].lastUpdBlock = block.number;      

        emit ItemBoughtFromShop(tokenId, itemId, _fightersHelper.getOwner(fighterId), _tokenAttributes[tokenId].name);
    }

    // Get the attributes for a fighter NFT
    function getTokenAttributes(uint256 tokenId) public view returns (ItemAttributes memory) {
        require(_exists(tokenId), "Token does not exist");

        return _tokenAttributes[tokenId];
    }

    // Get the attributes for a fighter NFT
    function getItemAttributes(uint256 itemId) public view returns (ItemAttributes memory) {
        require(_itemAttributes.length > itemId, "Item does not exist");

        return _itemAttributes[itemId];
    }

    function burnConsumable(uint256 tokenId) external {
        ItemAttributes memory atts = _tokenAttributes[tokenId];

        require(atts.isConsumable, "Item not a consumable");
        _burn(tokenId);
    }

    mapping (uint256 => ItemAttributes) private _tokenAttributes;
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


    function _setTokenAttributes(uint256 tokenId, ItemAttributes memory attrs) internal {
        require(_exists(tokenId), "Token does not exist");

        _tokenAttributes[tokenId] = attrs;
    }


    // Returns the smaller of two values
    function min(uint a, uint b) private pure returns (uint) {
        return a < b ? a : b;
    }

    // Returns the largest of the two values
    function max(uint a, uint b) private pure returns (uint) {
        return a > b ? a : b;
    }

    // Safe Multiply Function - prevents integer overflow 
    function safeMul(uint a, uint b) public returns (uint) {
        uint c = a * b;
        assert(a == 0 || c / a == b);
        return c;
    }

    // Safe Subtraction Function - prevents integer overflow 
    function safeSub(uint a, uint b) public returns (uint) {
        assert(b <= a);
        return a - b;
    }

    // Safe Addition Function - prevents integer overflow 
    function safeAdd(uint a, uint b) public returns (uint) {
        uint c = a + b;
        assert(c>=a && c>=b);
        return c;
    }

}