// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./Fighters.sol";
import "./FightersAtts.sol";

contract FightersHelper is FightersAtts {
    Fighters private _fighters;

    constructor(address fightersContract) {
        _fighters = Fighters(fightersContract);
    }

    function createFighter(address owner, string calldata name, FightersAtts.FighterClass fighterClass) external returns (uint256) {
        return _fighters.createFighter(owner, name, fighterClass);
    }

    function getTokenAttributes(uint256 tokenId) public view returns (Attributes memory) {
        return _fighters.getTokenAttributes(tokenId);
    }

    function getLevel(uint256 tokenId) public view returns (uint256) {
        return _fighters.getLevel(tokenId);
    }

    function increaseExperience(uint256 tokenId, uint256 addExp) external {
        _fighters.increaseExperience(tokenId, addExp);
    }

    function getOwner(uint256 tokenId) public view returns (address) {
        return _fighters.getOwner(tokenId);
    }

    function getDropRarityLevel(uint256 tokenId) public view returns (uint256) {
        return _fighters.getDropRarityLevel(tokenId);
    }

    function getFighterStats(uint256 tokenId) public view returns (FighterStats memory) {
        return _fighters.getFighterStats(tokenId);
    }
}