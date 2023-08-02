// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./Backpack.sol";
import "./ItemsExcellentAtts.sol";
contract BackpackHelper is ItemsExcellentAtts {
    Backpack private _backpack;

    constructor(address backpackAddress) {
        _backpack = Backpack(backpackAddress);
    }    

    function dropBackpackItem(uint256 tokenId, uint256 qty) external {
        return _backpack.dropBackpackItem(tokenId, qty);
    }

    function pickupItem(bytes32 itemHash, ExcellentItemAtts memory itemAtts, uint256 dropBlock, uint256 fighterId) external {
        return _backpack.pickupItem(itemHash, itemAtts, dropBlock, fighterId);
    }
}