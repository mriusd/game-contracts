// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";


import "./SafeMath.sol";
import "./FightersClasses.sol";
import "./FightersAtts.sol";

contract Fighters is ERC721Enumerable, FightersAtts, FightersClasses, SafeMath {
    using Counters for Counters.Counter;

    // Initial attributes by class
    mapping (uint256 => Attributes) private _initialAttributes;
    mapping (uint256 => Attributes) private _tokenAttributes;
    mapping (string => bool) public names;

    // Events
    event FighterCreated(uint256 tokenId, address owner, string fighterClass, string name);
    event NPCCreated(address indexed creator, uint256 tokenId, string className, string npcName);
    event StatsUpdated(uint256 tokenId, uint256 strength, uint256 agility, uint256 energy, uint256 vitality);
    event ItemEquiped(uint256 tokenId, uint256 itemId, uint256 slot);

    // Counter for token IDs
    Counters.Counter private _tokenIdCounter;

    constructor() ERC721("MRIUS", "Fighter") {
        // Set initial attributes for each class
        // _initialAttributes[uint256(FighterClass.DarkKnight)]        = Attributes("", 0, 42, 21,  5, 20, 0, 1,     5,  5, 5, 1, 5, 27, 7, 0, 0);
        // _initialAttributes[uint256(FighterClass.DarkWizard)]        = Attributes("", 0, 15, 20, 50, 20, 0, 2,     3, 10, 3, 5, 5, 16, 5, 0, 0);
        // _initialAttributes[uint256(FighterClass.FairyElf)]          = Attributes("", 0, 20, 25, 15, 20, 0, 3,     3,  5, 3, 4, 5, 12, 5, 0, 0);
        // _initialAttributes[uint256(FighterClass.MagicGladiator)]    = Attributes("", 0, 28, 14, 20, 20, 0, 4,     4,  7, 6, 2, 7, 23, 7, 0, 0);


        // DK
        updateFighterClass("Warrior",   FightersClassAttributes(42, 21,  5, 20,  5,  5, 5, 1, 5, 27, 7,  0, 0));
        updateFighterClass("Wizard",    FightersClassAttributes(15, 20, 50, 20,  3, 10, 3, 5, 5, 16, 5,  0, 0));
    }

    function updateFighterStats(uint256 tokenId, uint256 strength, uint256 agility, uint256 energy, uint256 vitality)  external  {
        // Check msg.sender is the NFT owner
        require(ownerOf(tokenId) == msg.sender, "Access Denied");
        
        require(strength >= _tokenAttributes[tokenId].strength, "Strength can't be decreased");
        require(agility >= _tokenAttributes[tokenId].agility, "Agility can't be decreased");
        require(energy >= _tokenAttributes[tokenId].energy, "Energy can't be decreased");
        require(vitality >= _tokenAttributes[tokenId].vitality, "Vitality can't be decreased");


        uint256 maxPoints = getMaxStatPoints(tokenId);
        uint256 usedStatPoints = getTotalUsedStatPoints(tokenId);

        uint256 newTotalStats = safeAdd(strength, safeAdd(agility, safeAdd(energy, vitality)));

        require(newTotalStats <= maxPoints, "Trying too add too many points");

        _tokenAttributes[tokenId].strength = strength;
        _tokenAttributes[tokenId].agility = agility;
        _tokenAttributes[tokenId].energy = energy;
        _tokenAttributes[tokenId].vitality = vitality;

        emit StatsUpdated(tokenId, strength, agility, energy, vitality);
    }

    function getMaxStatPoints(uint256 tokenId) public view returns (uint256) {
        uint256 levelPoints = safeMul(getLevel(tokenId), FighterClasses[_tokenAttributes[tokenId].class].statPointsPerLevel);

        return levelPoints;
    }

    function getTotalUsedStatPoints(uint256 tokenId) public view returns (uint256) {
        return safeAdd(_tokenAttributes[tokenId].strength,
                safeAdd(_tokenAttributes[tokenId].agility, 
                safeAdd(_tokenAttributes[tokenId].energy,
                _tokenAttributes[tokenId].vitality)));
    }

    function getUserFighters(address userAddress) external view returns (uint256[] memory) {        
        uint256 numTokens = balanceOf(userAddress);
        uint256[] memory tokenIds = new uint256[](numTokens);

        for (uint256 i = 0; i < numTokens; i++) {
            tokenIds[i] = tokenOfOwnerByIndex(userAddress, i);
        }

        return tokenIds;
    }

    // Create a new fighter NFT with initial attributes
    function createFighter(address owner, string calldata name, string calldata fighterClass) external returns (uint256) {
        // Make sure the class is valid
        require(fighterClassExists(fighterClass), "Invalid fighter class");

        validateFighterName(name);
        

        // Mint the NFT with the initial attributes
        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();       

        emit FighterCreated(tokenId, owner, fighterClass, name);

        

        _safeMint(owner, tokenId);
        

        _tokenAttributes[tokenId] = Attributes({
            name: name,
            class: fighterClass,
            tokenId: tokenId,
            strength: 0,
            agility: 0,
            energy: 0,
            vitality: 0,
            experience: 0,
            hpPerVitalityPoint: 0,
            manaPerEnergyPoint: 0,
            hpIncreasePerLevel: 0,
            manaIncreasePerLevel: 0,
            statPointsPerLevel: 0,
            attackSpeed: 0,
            agilityPointsPerSpeed: 0,
            isNpc: 0,
            dropRarityLevel: 0
        });  

        names[name] = true; 

        return tokenId;
    }

    function validateFighterName(string calldata name) public view {
        require(!names[name], "Name taken");   
        require(bytes(name).length <= 13, "Name too long");

        // Check if name contains only A-Z, a-z, 0-9
        for (uint i = 0; i < bytes(name).length; i++) {
            bytes1 char = bytes(name)[i];
            require((char >= '0' && char <= '9') || 
                    (char >= 'A' && char <= 'Z') || 
                    (char >= 'a' && char <= 'z'), "Name contains invalid characters");
        }
    }

    function createNPC(string calldata npcName, string calldata className, address ownerAddress) external returns (uint256) {
        validateFighterName(npcName);
        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        _safeMint(ownerAddress, tokenId);
        _tokenAttributes[tokenId].tokenId = tokenId;
        _tokenAttributes[tokenId].name = npcName;
        _tokenAttributes[tokenId].class = className;
        _tokenAttributes[tokenId].experience = getExpFromLevel(FighterClasses[className].dropRarityLevel);      

        names[npcName] = true;

        emit NPCCreated(ownerAddress, tokenId, npcName, className);

        return tokenId;
    }

    // Get the attributes for a fighter NFT
    function getTokenAttributes(uint256 tokenId) public view returns (Attributes memory) {
        require(_exists(tokenId), "Token does not exist");

        Attributes memory atts = _tokenAttributes[tokenId];

        atts.strength   = safeAdd(atts.strength,    FighterClasses[atts.class].baseStrength);
        atts.agility    = safeAdd(atts.agility,     FighterClasses[atts.class].baseAgility);
        atts.energy     = safeAdd(atts.energy,      FighterClasses[atts.class].baseEnergy);
        atts.vitality   = safeAdd(atts.vitality,    FighterClasses[atts.class].baseVitality);

        atts.hpPerVitalityPoint     = FighterClasses[atts.class].hpPerVitalityPoint;
        atts.manaPerEnergyPoint     = FighterClasses[atts.class].manaPerEnergyPoint;
        atts.hpIncreasePerLevel     = FighterClasses[atts.class].hpIncreasePerLevel;
        atts.statPointsPerLevel     = FighterClasses[atts.class].statPointsPerLevel;
        atts.attackSpeed            = FighterClasses[atts.class].attackSpeed;
        atts.agilityPointsPerSpeed  = FighterClasses[atts.class].agilityPointsPerSpeed;
        atts.isNpc                  = FighterClasses[atts.class].isNpc;
        atts.dropRarityLevel        = FighterClasses[atts.class].dropRarityLevel; 

        return atts;
    }

    function getOwner(uint256 tokenId) public view returns (address) {
        return ownerOf(tokenId);
    }

    function getDropRarityLevel(uint256 tokenId) public view returns (uint256) {
        return _tokenAttributes[tokenId].dropRarityLevel;
    }

    // Set the attributes for a fighter NFT
    function _setTokenAttributes(uint256 tokenId, Attributes memory atts) internal {
        require(_exists(tokenId), "Token does not exist");

        _tokenAttributes[tokenId].strength  = atts.strength;
        _tokenAttributes[tokenId].agility   = atts.agility;
        _tokenAttributes[tokenId].energy    = atts.energy;
        _tokenAttributes[tokenId].vitality  = atts.vitality;
    }

    // Increase fighter experience
    function increaseExperience(uint256 tokenId, uint256 addExp) external  {
        require(_exists(tokenId), "Token does not exist");

        _tokenAttributes[tokenId].experience = min(maxExperience, safeAdd(_tokenAttributes[tokenId].experience, addExp));
    }


    function getMaxHealth(uint256 tokenId) public view returns (uint256) 
    {
        require(_exists(tokenId), "Token does not exist");
        uint256 increasePerLevel = 0;
        uint256 hpPerVitalityPoint = 0;
        uint256 vitalityPoints;

        increasePerLevel = FighterClasses[_tokenAttributes[tokenId].class].hpIncreasePerLevel;
        hpPerVitalityPoint = FighterClasses[_tokenAttributes[tokenId].class].hpPerVitalityPoint;
        

        return safeAdd(safeMul(increasePerLevel, getLevel(tokenId)), safeMul(hpPerVitalityPoint, _tokenAttributes[tokenId].vitality));
    }


    function getMaxMana(uint256 tokenId) public view returns (uint256) 
    {
        require(_exists(tokenId), "Token does not exist");
        uint256 increasePerLevel = 0;
        uint256 mpPerEnergyPoint = 0;
        uint256 energyPoints;

        
        increasePerLevel = FighterClasses[_tokenAttributes[tokenId].class].manaIncreasePerLevel;
        mpPerEnergyPoint = FighterClasses[_tokenAttributes[tokenId].class].manaPerEnergyPoint;
        

        return safeAdd(safeMul(mpPerEnergyPoint, _tokenAttributes[tokenId].energy), safeMul(increasePerLevel, getLevel(tokenId)));
    }

    

    function getFighterStats(uint256 tokenId) public view returns (FighterStats memory)
    {
        require(_exists(tokenId), "Token does not exist");
        FighterStats memory stats = FighterStats({
            tokenId: tokenId,
            maxHealth: getMaxHealth(tokenId),
            maxMana: getMaxMana(tokenId),
            level: getLevel(tokenId),
            exp: getExperience(tokenId),
            totalStatPoints: getTotalUsedStatPoints(tokenId),
            maxStatPoints: getMaxStatPoints(tokenId)
        });

        return stats;
    }

    function getExperience(uint256 tokenId) public view returns (uint256)
    {
        require(_exists(tokenId), "Token does not exist");

        // Health = Health_After_Last_Damage + (Current Block Number - Last Damage Block Number) * Health Regeneration Rate
        return _tokenAttributes[tokenId].experience;
    }

    // Get player level
    function getLevel(uint256 tokenId) public view returns (uint256) 
    {
        require(_exists(tokenId), "Token does not exist");

        uint256 exp = _tokenAttributes[tokenId].experience;

        return safeSub(sqrt(safeAdd(safeMul(experienceDivider, exp), 125)), 5) / 10;
    }

    // Reverse function for getLevel
    function getExpFromLevel(uint256 level) public view returns (uint256) {
        uint256 exp = safeAdd(safeMul(level, 10), 5); // level.mul(10).add(5) 
        exp = safePow(exp, 2);
        exp = safeSub(exp, 125)/experienceDivider; // exp.sub(125).div(experienceDivider);
        return exp;
    }

    uint256 healthRegenerationDivider = 8;
    uint256 manaRegenerationDivider = 8;
    uint256 agilityPerDefence = 4;
    uint256 strengthPerDamage = 8;
    uint256 energyPerDamage = 8;
    uint256 maxExperience = 291342500;
    uint256 experienceDivider = 5;


    event updateHealthRegenerationDivider(uint256);
    event updateManaRegenerationDivider(uint256);
    event updateAgilityPerDefence(uint256);
    event updateStrengthPerDamage(uint256);
    event updateEnergyPerDamage(uint256);
    event updateMaxExperience(uint256);
    event updateExperienceDivider(uint256);
    

    function setHealthRegenerationDivider(uint256 val) external {
        healthRegenerationDivider = val;

        emit updateHealthRegenerationDivider(val);
    }

   function setManaRegenerationDivider(uint256 val) external {
        manaRegenerationDivider = val;

        emit updateManaRegenerationDivider(val);
    }

   function setAgilityPerDefence(uint256 val) external {
        agilityPerDefence = val;

        emit updateAgilityPerDefence(val);
    }

   function setStrengthPerDamage(uint256 val) external {
        strengthPerDamage = val;

        emit updateStrengthPerDamage(val);
    }

   function setEnergyPerDamage(uint256 val) external  {
        energyPerDamage = val;

        emit updateEnergyPerDamage(val);
    }

    function setMaxExperience(uint256 val) external  {
        maxExperience = val;

        emit updateMaxExperience(val);
    }

    function setExperienceDivider(uint256 val) external  {
        experienceDivider = val;

        emit updateExperienceDivider(val);
    }
}
