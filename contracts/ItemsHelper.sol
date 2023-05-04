// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./ItemAtts.sol";
import "./Items.sol";

contract ItemsHelper is ItemAtts {
    Items private _items;

    constructor(address itemsAddress) {
        _items = Items(itemsAddress);
    }

    // function craftItem(uint256 itemId, address itemOwner, uint256 maxLevel, uint256 maxAddPoints) external returns (uint256) {
    //     return _items.craftItem(itemId, itemOwner, maxLevel, maxAddPoints);
    // }

    function getTokenAttributes(uint256 tokenId) external returns (ItemAttributes memory) {
        return _items.getTokenAttributes(tokenId);
    }

    function burnItem(uint256 tokenId) external {
        return _items.burnItem(tokenId);
    }

    function setItemLevel(uint256 tokenId, uint256 level)  external {
        return _items.setItemLevel(tokenId, level);
    }   

    function setAdditionalPoints(uint256 tokenId, uint256 points) external {
        return _items.setAdditionalPoints(tokenId, points);
    }

    function dropItem(uint256 rarityLevel, uint256 fighterId, uint256 experience) external returns (bytes32) {
        return _items.dropItem(rarityLevel, fighterId, experience);
    }

    function safeMint(address owner) external returns (uint256) {
        return _items.safeMint(owner);
    }

    function setTokenAttributes(uint256 tokenId, ItemAttributes memory atts) external {
        return _items.setTokenAttributes(tokenId, atts);
    }

    function getDropQty(bytes32 itemHash) external returns (uint256)  {
        return _items.getDropQty(itemHash);
    }

    function createDropHash(bytes32 itemHash, uint256 qty) external {
        _items.createDropHash(itemHash, qty);
    }

    function itemExists(uint256 tokenId) external returns (bool) {
        return _items.itemExists(tokenId);
    }
}