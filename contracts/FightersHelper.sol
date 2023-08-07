// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./Fighters.sol";
import "./FightersAtts.sol";
import "./FightersClasses.sol";

contract FightersHelper is FightersAtts {
    Fighters private _fighters;

    constructor(address fightersContract) {
        _fighters = Fighters(fightersContract);
    }

    function updateFighterClass(string memory className, FightersClasses.FightersClassAttributes memory atts) public {
        _fighters.updateFighterClass(className, atts);
    }

    function createNPC(string calldata npcName, string calldata className, address ownerAddress) external returns (uint256) {
        return _fighters.createNPC(npcName, className, ownerAddress);
    }

    function createFighter(address owner, string calldata name, string calldata fighterClass) external returns (uint256) {
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

    function getUserFighters(address userAddress) external view returns (uint256[] memory) {
        return _fighters.getUserFighters(userAddress);
    }
}