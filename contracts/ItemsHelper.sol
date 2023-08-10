// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./ItemsExcellentAtts.sol";
import "./ItemsBase.sol";
import "./ItemsExcellent.sol";

contract ItemsHelper is ItemsExcellentAtts {
    ItemsBase private _base;
    Items private _items;
    ItemsExcellent private _itemsExcellent;



    constructor(address itemsBaseAddress, address itemsAddress, address itemsExcellentAddress) {
        _base = ItemsBase(itemsBaseAddress);
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
        //return _items.setAdditionalPoints(tokenId, points);
    }

    function safeMint(address owner) external returns (uint256) {
        return _items.safeMint(owner);
    }

    function setTokenAttributes(uint256 tokenId, ExcellentItemAtts memory atts) external {
        return _itemsExcellent.setTokenAttributes(tokenId, atts);
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

    function craftItemForShop(string calldata itemName, uint256 fighterId) external returns (uint256) {
        return _items.craftItemForShop(itemName, fighterId);
    }

    // function craftItem(uint256 itemId, address itemOwner, uint256 maxLevel, uint256 maxAddPoints) external returns (uint256) {
    //     return _items.craftItem(itemId, itemOwner, maxLevel, maxAddPoints);
    // }





    // RPC Calls (These cane be called by the backend)
    function burnConsumable(uint256 tokenId) external {
        return _items.burnConsumable(tokenId);
    }

    function packItems(string memory itemName, uint256[] memory tokenIds, uint256 packSize, uint256 fighterId) external returns (uint256 newTokenId) {
        //return _items.packItems(itemName, tokenIds, packSize, fighterId);
    }

    function unpackItems(uint256 tokenId) external returns(uint256[] memory tokenIds) {
        //return _items.unpackItems(tokenId);
    }







    // Getters (These can be called by anyone)
    function itemExists(uint256 tokenId) public view returns (bool) {
        return _items.itemExists(tokenId);
    }

    function isItemInShop(string memory name) public view returns (bool) {
        return _items.isItemInShop(name);
    }

    function isItemExcellent(uint256 tokenId) public view returns (bool) {
        return _itemsExcellent.isItemExcellent(tokenId);
    }




    function getWeapons(uint256 rarityLevel) public view returns(string[] memory) {
        return _base.getWeapons(rarityLevel);
    }

    function getArmours(uint256 rarityLevel) public view returns(string[] memory) {
        return _base.getArmours(rarityLevel);
    }

    function getJewels(uint256 rarityLevel) public view returns(string[] memory) {
        return _base.getJewels(rarityLevel);
    }

    function getWings(uint256 rarityLevel) public view returns(string[] memory) {
        return _base.getWings(rarityLevel);
    }
    function getBoxes(uint256 rarityLevel) public view returns(string[] memory) {
        return _base.getBoxes(rarityLevel);
    }
    function getConsumables(uint256 rarityLevel) public view returns(string[] memory) {
        return _base.getConsumables(rarityLevel);
    }

    function getMisc(uint256 rarityLevel) public view returns(string[] memory) {
        return _base.getMisc(rarityLevel);
    }

    function getItemAttributes(string memory itemName) external returns (ExcellentItemAtts memory) {
        return _itemsExcellent.getItemAttributes(itemName);
    }

    function getTokenAttributes(uint256 tokenId) external returns (ExcellentItemAtts memory) {
        return _itemsExcellent.getTokenAttributes(tokenId);
    }

    function getFighterItems(address userAddress, uint256 fighterId) external view returns (uint256[2][] memory) { 
        return _items.getFighterItems(userAddress, fighterId);
    }
}