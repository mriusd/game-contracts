// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./Items.sol";

abstract contract FighterContract {
    function increaseExperience(uint256 tokenId, uint256 addExp) external virtual;
    function decreaseHealth(uint256 tokenId, uint256 damage) external virtual;
    function getHealth(uint256 tokenId) public virtual returns (uint256);
    function getMana(uint256 tokenId) public virtual returns (uint256);
    function getLevel(uint256 tokenId) public virtual returns (uint256);
    function getDefence (uint256 tokenId) public virtual returns (uint256);
    function getDamage (uint256 tokenId) public virtual returns (uint256);
    function decreaseMana(uint256 tokenId, uint256 manaUsed) external virtual;
    function decreaseStamina(uint256 tokenId, uint256 staminaUsed) external virtual;
}

abstract contract  FighterMoneyContract {
    function payoutBattlePurse(address playerAddress, uint256 amount) external virtual;
}

contract Battle {
    ItemsHelper private _itemsHelper;

    struct Battle {
        uint256 blockNumber;
        uint256 opponent1;
        uint256 opponent2;
        uint256 winner; 
    }

   
    mapping(bytes32 => Battle) public battles;
    
    address fighterContractAddress; 
    address fighterMoneyContractAddress; 


    struct FightVars {
        uint256 level1;
        uint256 level2;
        address owner1;
        address owner2;

    }

    event LogError(uint8, uint256);
    event BattleRecorded(uint256 player1, uint256 player2, uint256 winner, uint256 purse);


    constructor(address fighterContractAddress_, address fighterMoneyContractAddress_, address itemsHelperContract) {
        fighterContractAddress = fighterContractAddress_;
        fighterMoneyContractAddress = fighterMoneyContractAddress_;
        owner = msg.sender;
        isAdmin[owner] = true;
        _itemsHelper = ItemsHelper(itemsHelperContract);
    }

    function recordFight(
        uint256[7] calldata vals1,
        uint256[7] calldata vals2,
        // uint8[2] calldata v,
        // bytes32[4] calldata rs,
        uint256 battleNonce
    ) external onlyAdmin
    {
        /*
            Opponent Vals
            [0] opponentId
            [1] health
            [2] mana
            [3] damageDealt
            [4] damageReceived
            [5] manaConsumed
            [6] win (0 - loss, 1 - win, 2 - draw)

        */

        bytes32 battleHash = generateBattleHash(vals1[0], vals2[0], battleNonce);
        require(battles[battleHash].blockNumber == 0, "Battle already recorded");

        

        FightVars memory t  = FightVars({
            level1: FighterContract(fighterContractAddress).getLevel(vals1[0]),
            level2: FighterContract(fighterContractAddress).getLevel(vals2[0]),
            owner1: ERC721(fighterContractAddress).ownerOf(vals1[0]),
            owner2: ERC721(fighterContractAddress).ownerOf(vals2[0])
        });



        /*
            vals1: [1, 100, 50, 10, 5, 7, 0]
            vals2: [2, 100, 50, 5, 10, 14, 1]


        */
        

        // check owner signatures

        // Generate owner hash
        // hash: this, battleNonce, opponentId, health, mana
        // bytes32 owner1hash = keccak256(abi.encode(this, battleNonce, vals2[0], vals1[1], vals1[2]));
        // require(ecrecover(keccak256(abi.encode("\x19Ethereum Signed Message:\n32", owner1hash)), v[0], rs[0], rs[1]) == t.owner1, "Wrong Opponent1 signature");

        // bytes32 owner2hash = keccak256(abi.encode(this, battleNonce, vals1[0], vals2[1], vals2[2]));
        // require(ecrecover(keccak256(abi.encode("\x19Ethereum Signed Message:\n32", owner1hash)), v[0], rs[0], rs[1]) == t.owner2, "Wrong Opponent2 signature");

   

        // update health of fighters
        FighterContract(fighterContractAddress).decreaseHealth(vals1[0], vals1[4]);
        FighterContract(fighterContractAddress).decreaseHealth(vals2[0], vals2[4]);

             

        // update fighters mana
        FighterContract(fighterContractAddress).decreaseMana(vals1[0], vals1[5]);
        FighterContract(fighterContractAddress).decreaseMana(vals2[0], vals2[5]);


        
 

        // get fighter levels
        uint256 exp1 = calculateExpirience(t.level1, t.level2, vals1[3]);

        emit LogError(21,exp1);
       // return; 
        uint256 exp2 = calculateExpirience(t.level2, t.level1, vals2[3]);



        emit LogError(3,exp1);
        emit LogError(4,exp2);
        
        //return;


        // increase fighter experience
        FighterContract(fighterContractAddress).increaseExperience(vals1[0], exp1);
        FighterContract(fighterContractAddress).increaseExperience(vals2[0], exp2);

        emit LogError(5,0);
        
        uint256 winner = 0;
        if (vals1[6] == 1) 
        {
            winner = vals1[0];
        }
        
        if (vals2[6] == 1)
        {
            winner = vals2[0];
        }
        

        emit LogError(6,0);
        //return;

        uint256 purse = 0;

        if (winner == vals1[0]) 
        {
            purse = safeMul(exp1, 1e18) / moneyExpDivider;
            FighterMoneyContract(fighterMoneyContractAddress).payoutBattlePurse(t.owner1, purse);
        } else if (winner == vals2[0])
        {
            purse = safeMul(exp2, 1e18) / moneyExpDivider;
            FighterMoneyContract(fighterMoneyContractAddress).payoutBattlePurse(t.owner2, purse);
        }

        
        


        battles[battleHash] = Battle({
            blockNumber: block.number,
            opponent1: vals1[0],
            opponent2: vals2[0],
            winner: winner
        });

        emit BattleRecorded(vals1[0], vals2[0], winner, purse);
    }


    function calculateExpirience(uint256 fighterLevel, uint256 opponentLevel, uint256 damageDealt) public returns (uint256) {
        uint256 levelDifference;
        uint256 exp;
        if (fighterLevel > opponentLevel)
        {
            levelDifference = fighterLevel - opponentLevel;

            exp = safeMul(damageDealt, safeSub(1e18, safeMul(levelDifference, 1e18) / fighterLevel)) / 1e18;
        }
        else if (fighterLevel < opponentLevel)
        {
            levelDifference = opponentLevel - fighterLevel;

            exp = safeMul(damageDealt, safeAdd(safeMul(levelDifference, 1e18) / fighterLevel, 1e18)) / 1e18;
        } 
        else
        {
            exp = damageDealt;
        }

        return exp;
    }

    function generateBattleHash (
        uint256 opponent1id,
        uint256 opponent2id,
        uint256 nonce
    ) internal returns (bytes32)
    {

        return keccak256(abi.encode(this, opponent1id, opponent2id, nonce));
    }


    
    address public owner; // holds the address of the contract owner

    // Event fired when the owner of the contract is changed
    event SetOwner(address previousOwner, address newOwner);
    event SetAdmin(address  previousOwner, bool isActive);

    // Allows only the owner of the contract to execute the function
    modifier onlyOwner {
        assert(msg.sender == owner);
        _;
    }


    // Allows only the owner of the contract to execute the function
    modifier onlyAdmin {
        assert(isAdmin[msg.sender]);
        _;
    }


    // Changes the owner of the contract
    function setOwner(address newOwner) public onlyOwner {
        emit SetOwner(owner, newOwner);
        owner = newOwner;
    }


    mapping (address => bool) private isAdmin;

    // Set battle contract
    function setAdmin(address addy, bool isActive) external onlyOwner
    {
        isAdmin[addy] = isActive;
        emit SetAdmin(addy, isActive);
    }








    
    uint256 moneyExpDivider = 100;


    event updatedMoneyExpDivider(uint256);

    function setMoneyExpDivider(uint256 val) external onlyOwner {
        moneyExpDivider = val;

        emit updatedMoneyExpDivider(val);
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
    function safeMul(uint a, uint b) public returns (uint) {
        uint c = a * b;
        assert(a == 0 || c / a == b);
        return c;
    }

    // Safe Subtraction Function - prevents integer overflow 
    function safeSub(uint a, uint b) public returns (uint) {
        assert(b <= a);
        return a - b;
    }

    // Safe Addition Function - prevents integer overflow 
    function safeAdd(uint a, uint b) public returns (uint) {
        uint c = a + b;
        assert(c>=a && c>=b);
        return c;
    }

    address public itemsContract;



}
