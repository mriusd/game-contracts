// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

import "truffle/Console.sol";
import "./ExcellentItems.sol";
contract Items is ERC721Enumerable, ExcellentItems {
    using Counters for Counters.Counter;

    address public owner;

    uint maxAdditionalPoints = 28;
    uint jolSuccessRate = 50;
    uint josSuccessRate = 50;

    mapping (uint256 => string) public itemName;

    mapping (uint256 => uint256[]) public Weapons; // mapping of rarityLevel => tokenId[]
    mapping (uint256 => uint256[]) public Armours; // mapping of rarityLevel => tokenId[]
    mapping (uint256 => uint256[]) public Jewels; // mapping of rarityLevel => tokenId[]
    mapping (uint256 => uint256[]) public Misc; // mapping of rarityLevel => tokenId[]

    struct DropParams {
        uint weaponsDropRate;
        uint armoursDropRate;
        uint jewelsDropRate;
        uint miscDropRate;
        uint boxDropRate;

        uint luckDropRate;
        uint skillDropRate;
        uint excDropRate;
        uint boxId;

        uint minItemLevel;
        uint maxItemLevel;
        uint maxAddPoints;
    }

    mapping (uint256 => DropParams) public DropParamsList; // mapping of rarityLevels to drop parameters
    mapping (uint256 => DropParams) public BoxDropPramsList; // mapping of rarityLevels to drop parameters

    event ItemGenerated(uint256 itemId, string name);
    event LogError(uint8, uint256);
    event DropParametersChange(uint256 rarityLevel, DropParams params);
    event BoxDropParametersChange(uint256 rarityLevel, DropParams params);
    event ItemDropped(uint256 tokenId, uint256 rarityLevel, address owner);
    event BoxOppened(uint256 tokenId, uint256 rarityLevel, address owner);
    event ItemBoughtFromShop(uint256 tokenId, uint256 itemId, address owner, string itemName);

    event ItemLevelUpgrade(uint256 tokenId, uint256 newLevel);
    event ItemAddPointsUpdate(uint256 tokenId, uint256 newAddPoints);

    // function upgradeItemLevel(uint256 itemTokenId, uint256 jewelTokenId) external {
    //     ItemAttributes memory item = _tokenAttributes[itemTokenId];
    //     ItemAttributes memory jewel = _tokenAttributes[jewelTokenId];
    //     uint256 luckPoints = 0;
    //     bool success = false;

    //     require(item.isWeapon || item.isArmour, "Invalid item");
    //     require(item.itemLevel < 9, "Item can be upgraded only up to level 9, use Chaos Machine further");
    //     require(jewel.itemAttributesId == 2 ||  jewel.itemAttributesId == 3, "Invalid jewel");

    //     if (item.luck) {
    //         luckPoints = 25;
    //     }

    //     if (item.itemLevel < 6) { // bless
    //         require(jewel.itemAttributesId == 2, "Required jewel of bless");
    //         _tokenAttributes[itemTokenId].itemLevel = _tokenAttributes[itemTokenId].itemLevel + 1;
    //     }

    //     if (item.itemLevel >= 6) { // soul
    //         require(jewel.itemAttributesId == 3, "Required jewel of soul");


    //         if (getRandomNumber(20) <= josSuccessRate+luckPoints) {
    //             _tokenAttributes[itemTokenId].itemLevel = _tokenAttributes[itemTokenId].itemLevel + 1;

    //         } else {
    //             if (item.itemLevel == 6) {
    //                 _tokenAttributes[itemTokenId].itemLevel = _tokenAttributes[itemTokenId].itemLevel - 1;
    //             } else {
    //                 _tokenAttributes[itemTokenId].itemLevel = 0;
    //             }
                
    //         }
    //     }
    //     _burn(jewelTokenId);
    //     emit ItemLevelUpgrade(itemTokenId, _tokenAttributes[itemTokenId].itemLevel);        
    // }

    function updateItemAdditionalPoints(uint256 itemTokenId, uint256 jolTokenId) external returns (bool) {
        ItemAttributes memory item = _tokenAttributes[itemTokenId];
        ItemAttributes memory jewel = _tokenAttributes[jolTokenId];
        uint256 luckPoints = 0;
        bool success = false;

        if (item.luck) {
            luckPoints = 25;
        }

        require(item.isWeapon || item.isArmour, "Invalid item");
        require(item.additionalDamage + item.additionalDefense < maxAdditionalPoints, "Max additional points reached");
        require(jewel.itemAttributesId == 4, "Invalid jewel");

        
        if (item.isWeapon) {
            if (getRandomNumber(30) <= jolSuccessRate+luckPoints) {
                _tokenAttributes[itemTokenId].additionalDamage = _tokenAttributes[itemTokenId].additionalDamage+4;
            } else {
                _tokenAttributes[itemTokenId].additionalDamage = 0;
            }
            
        }

        if (item.isArmour) {
            if (getRandomNumber(30) <= jolSuccessRate+luckPoints) {
                _tokenAttributes[itemTokenId].additionalDefense = _tokenAttributes[itemTokenId].additionalDefense+4;
            } else {
                _tokenAttributes[itemTokenId].additionalDefense = 0;
            }
        }

        _burn(jolTokenId);
        emit ItemAddPointsUpdate(itemTokenId, _tokenAttributes[itemTokenId].additionalDefense + _tokenAttributes[itemTokenId].additionalDamage);
    }


    function openBox(uint256 tokenId) external returns (uint256) {
        ItemAttributes memory box = _tokenAttributes[tokenId];
        DropParams memory params = BoxDropPramsList[box.itemRarityLevel];

        ItemAttributes memory dropItem = getDropItem(box.itemRarityLevel, params);

        _tokenIdCounter.increment();
        uint256 newTokenId = _tokenIdCounter.current();
        _safeMint(msg.sender, newTokenId);
        _setTokenAttributes(newTokenId, dropItem);

        _tokenAttributes[newTokenId].tokenId = newTokenId;      
        _tokenAttributes[newTokenId].fighterId = 0;      
        _tokenAttributes[newTokenId].lastUpdBlock = block.number; 

        emit BoxOppened(newTokenId, dropItem.itemRarityLevel, msg.sender);

        return newTokenId;
    }

    

    function setDropParams(uint256 rarityLevel, DropParams memory params) external {
        DropParamsList[rarityLevel] = params;
        emit DropParametersChange(rarityLevel, params);
    }

    function setBoxDropParams(uint256 rarityLevel, DropParams memory params) external {
        BoxDropPramsList[rarityLevel] = params;
        emit BoxDropParametersChange(rarityLevel, params);
    }

    function createItem(ItemAttributes memory atts) external returns (uint256 tokenId)
    {                
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

        uint256 randomNumber = getRandomNumber(0);
        uint256 randomItem;
        

        if (randomNumber <= params.weaponsDropRate) {
            randomItem =  returnRandomItemFromDropList(Weapons[rarityLevel]);
        } else if (randomNumber <= safeAdd(params.weaponsDropRate, params.armoursDropRate)) {
            randomItem =  returnRandomItemFromDropList(Armours[rarityLevel]);
        } else if (randomNumber <= safeAdd(params.weaponsDropRate, safeAdd(params.armoursDropRate, params.jewelsDropRate))) {
            randomItem =  returnRandomItemFromDropList(Jewels[rarityLevel]);
        } else if (randomNumber <= safeAdd(params.weaponsDropRate, safeAdd(params.armoursDropRate, safeAdd(params.jewelsDropRate, params.miscDropRate)))) {
            randomItem =  returnRandomItemFromDropList(Misc[rarityLevel]);
        } else if (randomNumber <= safeAdd(params.weaponsDropRate, safeAdd(params.armoursDropRate, safeAdd(params.jewelsDropRate, safeAdd(params.miscDropRate, params.boxDropRate))))) {
            randomItem =  params.boxId;
        } else {
            randomItem = 0;
        }


        ItemAttributes memory itemAtts = _itemAttributes[randomItem];


        if (!itemAtts.isJewel && !itemAtts.isMisc && !itemAtts.isBox && getRandomNumber(1) <= params.luckDropRate) {
            itemAtts.luck  = true;
        }

        if (itemAtts.isWeapon && getRandomNumber(2) <= params.skillDropRate) {
            itemAtts.skill  = true;
        }

        if (itemAtts.isBox) {
            itemAtts.itemLevel = rarityLevel;
        } else if (!itemAtts.isMisc && !itemAtts.isJewel) {
            itemAtts.itemLevel = params.minItemLevel + getRandomNumber(4) % (params.maxItemLevel-params.minItemLevel+1);

            if (itemAtts.isWeapon) {
                itemAtts.additionalDamage = 4 * (getRandomNumber(5) % (params.maxAddPoints/4+1));
            } else if (itemAtts.isArmour) {
                itemAtts.additionalDefense = 4 * (getRandomNumber(5) % (params.maxAddPoints/4+1));
            }
        }

        if ((itemAtts.isWeapon || itemAtts.isArmour) && getRandomNumber(3) <= params.excDropRate) {
            itemAtts = addExcelentOption(itemAtts);
        }

        return itemAtts;  
    }

    function dropItem(uint256 rarityLevel) external returns (uint256) {
        DropParams memory params = DropParamsList[rarityLevel];

        ItemAttributes memory itemAtts = getDropItem(rarityLevel, params);

        if (itemAtts.tokenId == 1) {
            emit ItemDropped(0, rarityLevel, msg.sender);
            return 0;
        }

        _tokenIdCounter.increment();
        uint256 newTokenId = _tokenIdCounter.current();
        _safeMint(msg.sender, newTokenId);
        _setTokenAttributes(newTokenId, itemAtts);

        _tokenAttributes[newTokenId].tokenId = newTokenId;      
        _tokenAttributes[newTokenId].fighterId = 0;      
        _tokenAttributes[newTokenId].lastUpdBlock = block.number; 

        emit ItemDropped(newTokenId, rarityLevel, msg.sender);

        return newTokenId;
    }    

    function getWeaponsLength(uint256 rarityLevel) public view returns (uint256) {
        return Weapons[rarityLevel].length;
    }

    function getArmoursLength(uint256 rarityLevel) public view returns (uint256) {
        return Armours[rarityLevel].length;
    }

    function getJewelsLength(uint256 rarityLevel) public view returns (uint256) {
        return Jewels[rarityLevel].length;
    }

    function getMiscsLength(uint256 rarityLevel) public view returns (uint256) {
        return Misc[rarityLevel].length;
    }    

    function buyItemFromShop(uint256 itemId, uint256 fighterId) external 
    {
        require(_itemAttributes[itemId].inShop, "Item not in shop or doesn't exist");

        // money logic

        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        _safeMint(msg.sender, tokenId);
        _setTokenAttributes(tokenId, _itemAttributes[itemId]);

        _tokenAttributes[tokenId].tokenId = tokenId;      
        _tokenAttributes[tokenId].itemAttributesId = itemId;      
        _tokenAttributes[tokenId].fighterId = fighterId;      
        _tokenAttributes[tokenId].lastUpdBlock = block.number;      

        emit ItemBoughtFromShop(tokenId, itemId, msg.sender, _tokenAttributes[tokenId].name);
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

    mapping (uint256 => ItemAttributes) private _tokenAttributes;
    ItemAttributes[] private _itemAttributes;
    Counters.Counter private _tokenIdCounter;

    ItemAttributes boxAttributes;

    constructor() ERC721("Combats", "Item") {
        owner = msg.sender;
    }

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

    function returnRandomItemFromDropList(uint256[] memory items) internal returns (uint256) {
        uint256 randomNumber = getRandomNumber(block.number);
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