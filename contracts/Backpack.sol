// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "./FightersHelper.sol";
import "./MoneyHelper.sol";
import "./ItemsHelper.sol";
import "./DropHelper.sol";

contract Backpack is ItemsExcellentAtts { 
    ItemsHelper private _itemsHelper;
    DropHelper private _dropHelper;
    MoneyHelper private _moneyHelper;
    FightersHelper private _fightersHelper;


    constructor ( address fightersHelperContract, address itemsHelperContract, address moneyHelperContract, address dropHelperContract ) {
        _itemsHelper = ItemsHelper(itemsHelperContract);
        _moneyHelper = MoneyHelper(moneyHelperContract);
        _fightersHelper = FightersHelper(fightersHelperContract);
        _dropHelper = DropHelper(dropHelperContract);
    }

    event BackpackItemDropped(bytes32 itemHash, ExcellentItemAtts item, uint256 qty, uint256 tokenId);
    event ItemPicked(uint256 tokenId, uint256 fighterId, uint256 qty);

    function dropBackpackItem(uint256 tokenId) external {
        require(_itemsHelper.itemExists(tokenId), "Token not found");
        uint256 qty = 1;

        ExcellentItemAtts memory tokenAttributes = _itemsHelper.getTokenAttributes(tokenId);

        tokenAttributes.fighterId = 0;
        tokenAttributes.tokenId = 0;

        bytes32 itemHash = keccak256(abi.encode(tokenAttributes, 1, block.number));
        _dropHelper.createDropHash(itemHash, 1);

        _itemsHelper.burnItem(tokenId);
        emit BackpackItemDropped(itemHash, tokenAttributes, 1, tokenId);
    }

    function pickupItem(bytes32 itemHash, ExcellentItemAtts memory itemAtts, uint256 dropBlock, uint256 fighterId) external {
        uint256 dropQty = _dropHelper.getDropQty(itemHash);
        require(dropQty > 0, "Item hash not found");

        bytes32 genHash = keccak256(abi.encode(itemAtts, dropQty, dropBlock));
        require(genHash == itemHash, "Item hash doest match");

        address fighterOwner = _fightersHelper.getOwner(fighterId);
        

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