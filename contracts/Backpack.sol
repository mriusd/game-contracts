// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "./FighterHelper.sol";
import "./MoneyHelper.sol";
import "./ItemsHelper.sol";

contract Backpack is ItemAtts { 
    ItemsHelper private _itemsHelper;
    MoneyHelper private _moneyHelper;
    FighterHelper private _fighterHelper;


    constructor ( address fighterHelperContract, address itemsHelperContract, address moneyHelperContract ) {
        _itemsHelper = ItemsHelper(itemsHelperContract);
        _moneyHelper = MoneyHelper(moneyHelperContract);
        _fighterHelper = FighterHelper(fighterHelperContract);
    }

    event BackpackItemDropped(bytes32 itemHash, ItemAttributes item, uint256 qty, uint256 tokenId);
    event ItemPicked(uint256 tokenId, uint256 fighterId, uint256 qty);

    function dropBackpackItem(uint256 tokenId) external {
        require(_itemsHelper.itemExists(tokenId), "Token not found");
        uint256 qty = 1;

        ItemAttributes memory tokenAttributes = _itemsHelper.getTokenAttributes(tokenId);

        tokenAttributes.fighterId = 0;
        tokenAttributes.tokenId = 0;

        bytes32 itemHash = keccak256(abi.encode(tokenAttributes, 1, block.number));
        _itemsHelper.createDropHash(itemHash, 1);

        _itemsHelper.burnItem(tokenId);
        emit BackpackItemDropped(itemHash, tokenAttributes, 1, tokenId);
    }

    function pickupItem(bytes32 itemHash, ItemAttributes memory itemAtts, uint256 dropBlock, uint256 fighterId) external {
        uint256 dropQty = _itemsHelper.getDropQty(itemHash);
        require(dropQty > 0, "Item hash not found");

        bytes32 genHash = keccak256(abi.encode(itemAtts, dropQty, dropBlock));
        require(genHash == itemHash, "Item hash doest match");

        address fighterOwner = _fighterHelper.getOwner(fighterId);
        

        if (itemAtts.itemAttributesId == goldItemId) {
            
            _moneyHelper.mintGold(fighterOwner, dropQty);
            emit ItemPicked(goldTokenId, fighterId, dropQty);
            
        } else {
            uint256 newTokenId = _itemsHelper.safeMint(fighterOwner);

            

            itemAtts.tokenId = newTokenId;      
            itemAtts.fighterId = fighterId;      
            itemAtts.lastUpdBlock = block.number; 

            _itemsHelper.setTokenAttributes(newTokenId, itemAtts);
            emit ItemPicked(newTokenId, fighterId, dropQty);
            
        } 
    }
}