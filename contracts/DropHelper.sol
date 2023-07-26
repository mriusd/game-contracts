// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./Drop.sol";
import "./ItemsExcellentAtts.sol";
contract DropHelper  is ItemsExcellentAtts {
    Drop private _drop;

    constructor(address dropAddress) {
        _drop = Drop(dropAddress);
    }    

    function makeItem(ExcellentItemAtts memory itemAtts) external returns (bytes32) {
        return _drop.makeItem(itemAtts);
    }

    function setDropParams(uint256 rarityLevel, Drop.DropParams memory params) external {
        return _drop.setDropParams(rarityLevel, params);
    }

    function setBoxDropParams(uint256 rarityLevel, Drop.DropParams memory params) external {
        return _drop.setBoxDropParams(rarityLevel, params);
    }

    function dropItem(uint256 rarityLevel, uint256 fighterId, uint256 experience) external returns (bytes32) {
        return _drop.dropItem(rarityLevel, fighterId, experience);
    }

    function openBox(uint256 tokenId) external returns (uint256) {
        return _drop.openBox(tokenId);
    }

    function getDropQty(bytes32 itemHash) external returns (uint256)  {
        return _drop.getDropQty(itemHash);
    }

    function createDropHash(bytes32 itemHash, uint256 qty) external {
        _drop.createDropHash(itemHash, qty);
    }
}