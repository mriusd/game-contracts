// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./ItemsExcellentAtts.sol";
import "./ItemsExcellent.sol";

contract ItemsHelper is ItemsExcellentAtts {
    Items private _items;
    ItemsExcellent private _itemsExcellent;



    constructor(address itemsAddress, address itemsExcellentAddress) {
        _items = Items(itemsAddress);
        _itemsExcellent = ItemsExcellent(itemsExcellentAddress);
    }


    // Contract Calls (These are called only from allowed contracts)
    function burnItem(uint256 tokenId) external {
        return _items.burnItem(tokenId);
    }

    function setItemLevel(uint256 tokenId, uint256 level)  external {
        return _items.setItemLevel(tokenId, level);
    }   

    function setAdditionalPoints(uint256 tokenId, uint256 points) external {
        return _items.setAdditionalPoints(tokenId, points);
    }

    function safeMint(address owner) external returns (uint256) {
        return _items.safeMint(owner);
    }

    function setTokenAttributes(uint256 tokenId, ExcellentItemAtts memory atts) external {
        return _itemsExcellent.setTokenAttributes(tokenId, atts);
    }

    function itemExists(uint256 tokenId) external returns (bool) {
        return _items.itemExists(tokenId);
    }

    function transferItem(uint256 tokenId, address to) external {
        _items.transferItem(tokenId, to);
    }

    function addExcellentOption(ExcellentItemAtts memory item) external returns (ExcellentItemAtts memory) {
        return _itemsExcellent.addExcellentOption(item);
    }

    function convertToExcellent(ItemAttributes memory item) public view returns (ExcellentItemAtts memory) {
        return _itemsExcellent.convertToExcellent(item);
    }

    // function craftItem(uint256 itemId, address itemOwner, uint256 maxLevel, uint256 maxAddPoints) external returns (uint256) {
    //     return _items.craftItem(itemId, itemOwner, maxLevel, maxAddPoints);
    // }



    // RPC Calls (These cane be called by the backend)
    function burnConsumable(uint256 tokenId) extarnal {
        return _items.burnConsumable(tokenId);
    }

    function buyItemFromShop(uint256 itemId, uint256 fighterId) external { 
        return _items.buyItemFromShop(itemId, fighterId);
    }


    function packItems(uint256 itemId, uint256[] memory tokenIds, uint256 packSize, uint256 fighterId) external returns (uint256 newTokenId) {
        return _items.packItems(itemId, tokenIds, packSize, fighterId);
    }

    function unpackItems(uint256 tokenId) external returns(uint256[] memory tokenIds) {
        return _items.unpackItems(tokenId);
    }


    // Getters (These can be called by anyone)
    function getWeapons(uint256 rarityLevel) public view returns(uint256[] memory) {
        return _items.getWeapons(rarityLevel);
    }

    function getArmours(uint256 rarityLevel) public view returns(uint256[] memory) {
        return _items.getArmours(rarityLevel);
    }

    function getJewels(uint256 rarityLevel) public view returns(uint256[] memory) {
        return _items.getJewels(rarityLevel);
    }

    function getMisc(uint256 rarityLevel) public view returns(uint256[] memory) {
        return _items.getMisc(rarityLevel);
    }

    function getItemAttributes(uint256 itemId) external returns (ExcellentItemAtts memory) {
        return _itemsExcellent.getItemAttributes(itemId);
    }

    function getTokenAttributes(uint256 tokenId) external returns (ExcellentItemAtts memory) {
        return _itemsExcellent.getTokenAttributes(tokenId);
    }

    function getFighterItems(address userAddress, uint256 fighterId) external view returns (uint256[2][] memory) { 
        return _items.getFighterItems(userAddress, fighterId);
    }
}