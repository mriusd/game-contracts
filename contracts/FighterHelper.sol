// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./FighterAttributes.sol";

contract FighterHelper {
    FighterAttributes private _fighter;

    constructor(address fighterAttributesContract) {
        _fighter = FighterAttributes(fighterAttributesContract);
    }

    function getLevel(uint256 tokenId) public view returns (uint256) {
        return _fighter.getLevel(tokenId);
    }

    function increaseExperience(uint256 tokenId, uint256 addExp) external {
        _fighter.increaseExperience(tokenId, addExp);
    }
}