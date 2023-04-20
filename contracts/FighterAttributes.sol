// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/utils/math/SafeMath.sol";

contract FighterAttributes is ERC721 {
    using Counters for Counters.Counter;
    using SafeMath for uint256;

    // Class enumeration
    enum FighterClass {
        None,
        DarkKnight,
        DarkWizard,
        FairyElf,
        MagicGladiator
    }

    // Struct for fighter attributes
    struct Attributes {
        string name;
        uint256 tokenId;
        uint256 strength;
        uint256 agility;
        uint256 energy;
        uint256 vitality;
        uint256 experience;
        uint256 class;


        uint256 hpPerVitalityPoint;
        uint256 manaPerEnergyPoint;
        uint256 hpIncreasePerLevel;
        uint256 manaIncreasePerLevel;
        uint256 statPointsPerLevel;
        uint256 attackSpeed;
        uint256 agilityPointsPerSpeed;
        uint256 isNpc;
        uint256 dropRarityLevel;
        

        // item slots
        /*
            1. helmet
            2. armour
            3. pants
            4. gloves
            5. boots
            6. left hand
            7. right hand
            8. left ring
            9, right ring
            10, pendant
            11. wings

        */
        uint256 helmSlot;
        uint256 armourSlot;
        uint256 pantsSlot;
        uint256 glovesSlot;
        uint256 bootsSlot;
        uint256 leftHandSlot;
        uint256 rightHandSlot;
        uint256 leftRingSlot;
        uint256 rightRingSlot;
        uint256 pendSlot;
        uint256 wingsSlot;
        
        
    }

    // Initial attributes by class
    mapping (uint256 => Attributes) private _initialAttributes;
    mapping (uint256 => Attributes) private _tokenAttributes;
    mapping (string => bool) public names;

    // Events
    event FighterCreated(address indexed owner, uint256 tokenId, uint8 class);
    event NPCCreated(address indexed creator, uint256 tokenId, uint8 class);
    event StatsUpdated(uint256 tokenId, uint256 strength, uint256 agility, uint256 energy, uint256 vitality);
    event ItemEquiped(uint256 tokenId, uint256 itemId, uint256 slot);

    // Counter for token IDs
    Counters.Counter private _tokenIdCounter;

    constructor() ERC721("Combats", "Fighter") {
        // Set initial attributes for each class
        _initialAttributes[uint256(FighterClass.DarkKnight)]        = Attributes("", 0, 42, 21,  5, 20, 0, 1,     5,  5, 5, 1, 5, 27, 7, 0, 0,    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);
        _initialAttributes[uint256(FighterClass.DarkWizard)]        = Attributes("", 0, 15, 20, 50, 20, 0, 2,     3, 10, 3, 5, 5, 16, 5, 0, 0,    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);
        _initialAttributes[uint256(FighterClass.FairyElf)]          = Attributes("", 0, 20, 25, 15, 20, 0, 3,     3,  5, 3, 4, 5, 12, 5, 0, 0,    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);
        _initialAttributes[uint256(FighterClass.MagicGladiator)]    = Attributes("", 0, 28, 14, 20, 20, 0, 4,     4,  7, 6, 2, 7, 23, 7, 0, 0,    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);

        owner = msg.sender;
    }

    function equipItem(uint256 tokenId, uint256 itemId, uint256 slot) external {
        require(_exists(tokenId), "Char does not exist");
        require(ownerOf(tokenId) == msg.sender, "Access Denied");
        require(slot > 0, "Invalid slot");

        if (slot == 1) _tokenAttributes[tokenId].helmSlot = itemId;
        if (slot == 2) _tokenAttributes[tokenId].armourSlot = itemId;
        if (slot == 3) _tokenAttributes[tokenId].pantsSlot = itemId;
        if (slot == 4) _tokenAttributes[tokenId].glovesSlot = itemId;
        if (slot == 5) _tokenAttributes[tokenId].bootsSlot = itemId;
        if (slot == 6) _tokenAttributes[tokenId].leftHandSlot = itemId;
        if (slot == 7) _tokenAttributes[tokenId].rightHandSlot = itemId;
        if (slot == 8) _tokenAttributes[tokenId].leftRingSlot = itemId;
        if (slot == 9) _tokenAttributes[tokenId].rightRingSlot = itemId;
        if (slot == 10) _tokenAttributes[tokenId].pendSlot = itemId;
        if (slot == 11) _tokenAttributes[tokenId].wingsSlot = itemId;

        emit ItemEquiped(tokenId, itemId, slot);
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
        uint256 initialPoints = safeAdd(_initialAttributes[_tokenAttributes[tokenId].class].strength,
                safeAdd(_initialAttributes[_tokenAttributes[tokenId].class].agility, 
                safeAdd(_initialAttributes[_tokenAttributes[tokenId].class].energy,
                _initialAttributes[_tokenAttributes[tokenId].class].vitality)));

        uint256 levelPoints = safeMul(getLevel(tokenId), _tokenAttributes[tokenId].statPointsPerLevel);

        return safeAdd(initialPoints, levelPoints);
    }

    function getTotalUsedStatPoints(uint256 tokenId) public view returns (uint256) {
        return safeAdd(_tokenAttributes[tokenId].strength,
                safeAdd(_tokenAttributes[tokenId].agility, 
                safeAdd(_tokenAttributes[tokenId].energy,
                _tokenAttributes[tokenId].vitality)));
    }

    // Create a new fighter NFT with initial attributes
    function createFighter(string calldata name, FighterClass fighterClass) external returns (uint256) {
        // Make sure the class is valid
        require(fighterClass != FighterClass.None, "Invalid fighter class");
        require(!names[name], "Name taken");

        // Get the initial attributes for the class
        Attributes memory initialAttrs = _initialAttributes[uint256(fighterClass)];
        initialAttrs.name = name;
       

        // Mint the NFT with the initial attributes
        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        _safeMint(msg.sender, tokenId);
        _setTokenAttributes(tokenId, initialAttrs);

        _tokenAttributes[tokenId].tokenId = tokenId;
        

        emit FighterCreated(msg.sender, tokenId, uint8(fighterClass));

        return tokenId;
    }

    function createNPC(string calldata name, uint256 strength, uint256 agility, uint256 energy, uint256 vitality, uint256 attackSpeed, uint256 level, uint256 dropRarityLevel) external returns (uint256) {
        Attributes memory atts = Attributes(name, 0, strength, agility, energy,vitality, getExpFromLevel(level), 0,        0, 0, 0, 0, 0, attackSpeed,      0, 1, dropRarityLevel,     0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);
        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        _safeMint(msg.sender, tokenId);
        _setTokenAttributes(tokenId, atts);
        _tokenAttributes[tokenId].tokenId = tokenId;

        emit NPCCreated(msg.sender, tokenId, 0);

        return tokenId;
    }

    // Get the attributes for a fighter NFT
    function getTokenAttributes(uint256 tokenId) public view returns (Attributes memory) {
        require(_exists(tokenId), "Token does not exist");

        return _tokenAttributes[tokenId];
    }

    // Set the attributes for a fighter NFT
    function _setTokenAttributes(uint256 tokenId, Attributes memory attrs) internal {
        require(_exists(tokenId), "Token does not exist");

        _tokenAttributes[tokenId] = attrs;
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

        increasePerLevel = _tokenAttributes[tokenId].hpIncreasePerLevel;
        hpPerVitalityPoint = _tokenAttributes[tokenId].hpPerVitalityPoint;
        

        return safeAdd(safeMul(increasePerLevel, getLevel(tokenId)), safeMul(hpPerVitalityPoint, _tokenAttributes[tokenId].vitality));
    }


    function getMaxMana(uint256 tokenId) public view returns (uint256) 
    {
        require(_exists(tokenId), "Token does not exist");
        uint256 increasePerLevel = 0;
        uint256 mpPerEnergyPoint = 0;
        uint256 energyPoints;

        
        increasePerLevel = _tokenAttributes[tokenId].manaIncreasePerLevel;
        mpPerEnergyPoint = _tokenAttributes[tokenId].manaPerEnergyPoint;
        

        return safeAdd(safeMul(mpPerEnergyPoint, _tokenAttributes[tokenId].energy), safeMul(increasePerLevel, getLevel(tokenId)));
    }

    struct FighterStats {
        uint256 tokenId;
        uint256 maxHealth;
        uint256 maxMana;
        uint256 level;
        uint256 exp;
        uint256 totalStatPoints;
        uint256 maxStatPoints;
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
        uint256 exp = level.mul(10).add(5);
        exp = safePow(exp, 2);
        exp = exp.sub(125).div(experienceDivider);
        return exp;
    }

    address public owner; // holds the address of the contract owner

    // Event fired when the owner of the contract is changed
    event SetOwner(address indexed previousOwner, address indexed newOwner);
    event SetBattleContract(address battleContract, bool isActive);

    // Allows only the owner of the contract to execute the function
    modifier onlyOwner {
        require(msg.sender == owner, "Access Denied");
        _;
    }

    // Allows only the owner of the contract to execute the function
    modifier onlyGM {
        require(isGM[msg.sender] || msg.sender == owner, "Access Denied");
        _;
    }

     // Allows only the owner of the contract to execute the function
    modifier onlyBattleContract {
        require(isBattleContract[msg.sender], "Access Denied");
        _;
    }

 

    // Changes the owner of the contract
    function setOwner(address newOwner) public onlyOwner {
        emit SetOwner(owner, newOwner);
        owner = newOwner;
    }

    mapping (address => bool) public isGM;
    event SetGM(address, bool);
    // Changes the owner of the contract
    function setGM(address gmAddress, bool val) public onlyOwner {
        isGM[gmAddress] = val;
        emit SetGM(gmAddress, val);
    }

    mapping (address => bool) private isBattleContract;

    // Set battle contract
    function setBattleContract(address contr, bool isActive) external onlyOwner
    {
        isBattleContract[contr] = isActive;
        emit SetBattleContract(contr, isActive);
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
    

    function setHealthRegenerationDivider(uint256 val) external onlyOwner {
        healthRegenerationDivider = val;

        emit updateHealthRegenerationDivider(val);
    }

   function setManaRegenerationDivider(uint256 val) external onlyOwner {
        manaRegenerationDivider = val;

        emit updateManaRegenerationDivider(val);
    }

   function setAgilityPerDefence(uint256 val) external onlyOwner {
        agilityPerDefence = val;

        emit updateAgilityPerDefence(val);
    }

   function setStrengthPerDamage(uint256 val) external onlyOwner {
        strengthPerDamage = val;

        emit updateStrengthPerDamage(val);
    }

   function setEnergyPerDamage(uint256 val) external onlyOwner {
        energyPerDamage = val;

        emit updateEnergyPerDamage(val);
    }

    function setMaxExperience(uint256 val) external onlyOwner {
        maxExperience = val;

        emit updateMaxExperience(val);
    }

    function setExperienceDivider(uint256 val) external onlyOwner {
        experienceDivider = val;

        emit updateExperienceDivider(val);
    }



    // Cube Root function
    function sqrt(uint y) internal pure returns (uint z) {
        if (y > 3) {
            z = y;
            uint x = y / 2 + 1;
            while (x < z) {
                z = x;
                x = (y / x + x) / 2;
            }
        } else if (y != 0) {
            z = 1;
        }

        return z;
    }

    // Returns the smaller of two values
    function min(uint a, uint b) private pure returns (uint) {
        return a < b ? a : b;
    }

    // Returns the largest of the two values
    function max(uint a, uint b) private pure returns (uint) {
        return a > b ? a : b;
    }

    // Safe Multiply Function - prevents integer overflow 
    function safeMul(uint a, uint b) public view returns (uint) {
        uint c = a * b;
        require(a == 0 || c / a == b, "safeMul Failed");
        return c;
    }

    // Safe Subtraction Function - prevents integer overflow 
    function safeSub(uint a, uint b) public view returns (uint) {
        require(b <= a, "safeSub Failed");
        return a - b;
    }

    // Safe Addition Function - prevents integer overflow 
    function safeAdd(uint a, uint b) public view returns (uint) {
        uint c = a + b;
        require(c>=a && c>=b, "safeAdd Failed");
        return c;
    }

    function safePow(uint256 base, uint256 exponent) internal pure returns (uint256) {
        if (exponent == 0) {
            return 1;
        } else if (exponent == 1) {
            return base;
        }

        uint256 result = base;
        for (uint256 i = 1; i < exponent; i++) {
            result = result.mul(base);
        }
        return result;
    }
}
