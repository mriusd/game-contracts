// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

import "truffle/Console.sol";
import "./ExcellentItems.sol";
import "./FighterHelper.sol";
import "./MoneyHelper.sol";

contract Items is ERC721Enumerable, ExcellentItems {
    using Counters for Counters.Counter;
    FighterHelper private _fighterHelper;
    MoneyHelper private _moneyHelper;

    address public owner;

    constructor(address fighterHelperContract, address moneyHelperContract) ERC721("MRIUSD", "Item") {
        owner = msg.sender;
        _fighterHelper = FighterHelper(fighterHelperContract);  
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
    mapping (uint256 => uint256[]) public Misc; // mapping of rarityLevel => tokenId[]

    mapping (bytes32 => uint256) public dropHashes;



    struct DropParams {
        uint256 weaponsDropRate;
        uint256 armoursDropRate;
        uint256 jewelsDropRate;
        uint256 miscDropRate;
        uint256 boxDropRate;

        uint256 excDropRate;
        uint256 boxId;

        uint256 minItemLevel;
        uint256 maxItemLevel;
        uint256 maxAddPoints;

        uint256 blockCrated;
    }

    mapping (uint256 => DropParams) public DropParamsList; // mapping of rarityLevels to drop parameters
    mapping (uint256 => DropParams) public BoxDropPramsList; // mapping of rarityLevels to drop parameters

    event ItemGenerated(uint256 itemId, string name);
    event LogError(uint8, uint256);
    event DropParametersChange(uint256 rarityLevel, DropParams params);
    event BoxDropParametersChange(uint256 rarityLevel, DropParams params);
    
    event ItemCrafted(uint256 tokenId, address owner);
    event BoxOppened(uint256 tokenId, uint256 rarityLevel, address owner);
    event ItemBoughtFromShop(uint256 tokenId, uint256 itemId, address owner, string itemName);

    event ItemLevelUpgrade(uint256 tokenId, uint256 newLevel);
    event ItemAddPointsUpdate(uint256 tokenId, uint256 newAddPoints);

    event ItemDropped(bytes32 itemHash, ItemAttributes item, uint256 qty, uint256 tokenId);
    event ItemPicked(uint256 tokenId, uint256 fighterId, uint256 qty);

    function setAdditionalPoints(uint256 tokenId, uint256 points) external {
        require(points <= maxAdditionalPoints, "Max points reached");
        require(_tokenAttributes[tokenId].isWeapon || _tokenAttributes[tokenId].isArmour, "Max points reached");
        if (_tokenAttributes[tokenId].isWeapon) {
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

    function openBox(uint256 tokenId) external returns (uint256) {
        require(_exists(tokenId), "Token doesn't exist");
        ItemAttributes memory box = _tokenAttributes[tokenId];
        DropParams memory params = BoxDropPramsList[box.itemRarityLevel];

        ItemAttributes memory itemAtts = getDropItem(box.itemRarityLevel, params);

        uint256 qty = 1;

        bytes32 itemHash = keccak256(abi.encode(itemAtts, qty, block.number));
        dropHashes[itemHash] = qty;
         emit ItemDropped(itemHash, itemAtts, qty, 0);
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

    function getDropQty(bytes32 dropHash) external returns(uint256) {
        return dropHashes[dropHash];
    }

    function createDropHash(bytes32 itemHash, uint256 qty) external {
        dropHashes[itemHash] = qty;
    }

    function itemExists(uint256 tokenId) public view returns(bool) {
        return _exists(tokenId);
    }

    function setDropParams(uint256 rarityLevel, DropParams memory params) external {
        DropParamsList[rarityLevel] = params;
        emit DropParametersChange(rarityLevel, params);
    }

    function setBoxDropParams(uint256 rarityLevel, DropParams memory params) external {
        BoxDropPramsList[rarityLevel] = params;
        emit BoxDropParametersChange(rarityLevel, params);
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
         } 

        emit ItemGenerated(itemId, atts.name);
        return itemId;
    }

    function getDropItem(uint256 rarityLevel, DropParams memory params) internal returns (ItemAttributes memory) {

        uint256 randomNumber = getRandomNumberMax(0, 100);
        uint256 randomItem;
        
        //emit LogDropValues(randomNumber, params.jewelsDropRate, params.weaponsDropRate, params.armoursDropRate, params.miscDropRate, params.boxDropRate);
        
        uint256 cumulativeRate = params.jewelsDropRate;
        if (randomNumber < cumulativeRate) {
            randomItem = returnRandomItemFromDropList(100, Jewels[rarityLevel]);
        } else {
            cumulativeRate += params.armoursDropRate;
            if (randomNumber < cumulativeRate) {
                randomItem = returnRandomItemFromDropList(101, Armours[rarityLevel]);
            } else {
                cumulativeRate += params.weaponsDropRate;
                if (randomNumber < cumulativeRate) {
                    randomItem = returnRandomItemFromDropList(102, Weapons[rarityLevel]);
                } else {
                    cumulativeRate += params.miscDropRate;
                    if (randomNumber < cumulativeRate) {
                        randomItem = returnRandomItemFromDropList(103, Misc[rarityLevel]);
                    } else {
                        cumulativeRate += params.boxDropRate;
                        if (randomNumber < cumulativeRate) {
                            randomItem = params.boxId;
                        } else {
                            randomItem = 1;
                        }
                    }
                }
            }
        }


        ItemAttributes memory itemAtts = _itemAttributes[randomItem];
        itemAtts.itemAttributesId = randomItem;


        if (!itemAtts.isJewel && !itemAtts.isMisc && !itemAtts.isBox && getRandomNumber(1) <= luckDropRate && itemAtts.itemAttributesId != 1 && itemAtts.itemAttributesId != 2) {
            itemAtts.luck  = true;
        }

        if (itemAtts.isWeapon && getRandomNumber(2) <= skillDropRate) {
            itemAtts.skill  = true;
        }

        if (itemAtts.isBox) {
            itemAtts.itemLevel = rarityLevel;
        } else if (!itemAtts.isMisc && !itemAtts.isJewel && itemAtts.itemAttributesId != 1 && itemAtts.itemAttributesId != 2) {
            itemAtts.itemLevel = params.minItemLevel + getRandomNumber(4) % (params.maxItemLevel-params.minItemLevel+1);

            if (itemAtts.isWeapon) {
                itemAtts.additionalDamage = 4 * (getRandomNumber(5) % (params.maxAddPoints/4+1));
            } else if (itemAtts.isArmour) {
                itemAtts.additionalDefense = 4 * (getRandomNumber(5) % (params.maxAddPoints/4+1));
            }
        }

        if ((itemAtts.isWeapon || itemAtts.isArmour) && getRandomNumberMax(3, 1000) <= params.excDropRate) {
            itemAtts = addExcelentOption(itemAtts);
        }

        return itemAtts;  
    }

    function dropItem(uint256 rarityLevel, uint256 fighterId, uint256 experience) external returns (bytes32) {
        DropParams memory params = DropParamsList[rarityLevel];
        require(params.blockCrated > 0, "No drop parameters for rarityLevel");
        require(experience > 0, "Experience must be non zero");

        ItemAttributes memory itemAtts = getDropItem(rarityLevel, params);        

        uint256 randomNumber = getRandomNumber(200);

        uint256 qty = 1;

        if (randomNumber > generalDropRate) {
            itemAtts = _itemAttributes[1];
            emit ItemDropped(dummyHash, itemAtts, 0, 0);
            return 0;
        }

        if (itemAtts.itemAttributesId == goldItemId) {
            qty = max(1, experience/_moneyHelper.getExperienceDivider());
        }

        bytes32 itemHash = keccak256(abi.encode(itemAtts, qty, block.number));
        dropHashes[itemHash] = qty;
        emit ItemDropped(itemHash, itemAtts, qty, 0);
        return itemHash;
    }


    // This function is for development and testing only
    // Should be removed in Production
    function makeItem(ItemAttributes memory itemAtts) external returns (bytes32) {
        bytes32 itemHash = keccak256(abi.encode(itemAtts, 1, block.number));
        dropHashes[itemHash] = 1;
        emit ItemDropped(itemHash, itemAtts, 1, 0);
        return itemHash;
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

    // function buyItemFromShop(uint256 itemId, uint256 fighterId) external {
    //     require(_itemAttributes[itemId].inShop, "Item not in shop or doesn't exist");

    //     // money logic

    //     _tokenIdCounter.increment();
    //     uint256 tokenId = _tokenIdCounter.current();
    //     _safeMint(msg.sender, tokenId);
    //     _setTokenAttributes(tokenId, _itemAttributes[itemId]);

    //     _tokenAttributes[tokenId].tokenId = tokenId;      
    //     _tokenAttributes[tokenId].itemAttributesId = itemId;      
    //     _tokenAttributes[tokenId].fighterId = fighterId;      
    //     _tokenAttributes[tokenId].lastUpdBlock = block.number;      

    //     emit ItemBoughtFromShop(tokenId, itemId, msg.sender, _tokenAttributes[tokenId].name);
    // }

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

    function returnRandomItemFromDropList(uint256 seed, uint256[] memory items) internal returns (uint256) {
        uint256 randomNumber = getRandomNumber(seed);
        uint256 len = items.length;

        return items[randomNumber % len];

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