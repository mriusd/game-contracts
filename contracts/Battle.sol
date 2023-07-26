// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;
import "./DropHelper.sol";
import "./FightersHelper.sol";

contract Battle {
    DropHelper private _dropHelper;
    FightersHelper private _fightersHelper;


    mapping(bytes32 => uint256) public battles;

    event LogError(uint8, uint256);
    event BattleRecorded(uint256 killedFighter, bytes32 battleHash, uint256 battleNonce);


    constructor(address fightersHelperContract, address dropHelperContract) {
        _fightersHelper = FightersHelper(fightersHelperContract);        
        _dropHelper = DropHelper(dropHelperContract);
    }

    struct Damage {
        uint256 fighterId;
        uint256 damage;
    }

    function recordKill(
        uint256 killedFighter,
        Damage[] memory damageDealt,
        uint256 battleNonce
    ) external 
    {
        bytes32 battleHash = generateBattleHash(killedFighter, battleNonce);
        require(battles[battleHash] == 0, "Battle already recorded");        
        require(damageDealt.length > 0, "Damages empty");
        
        uint256 killLevel = _fightersHelper.getLevel(killedFighter);

        uint256 killerLevel;
        uint256 killerExp;
        for (uint256 i =0; i<damageDealt.length; i++) {
            killerLevel = _fightersHelper.getLevel(damageDealt[i].fighterId);
            killerExp = calculateExperience(killerLevel, killLevel, damageDealt[i].damage);
            _fightersHelper.increaseExperience(damageDealt[i].fighterId, killerExp);

            if (i == 0) {
                _dropHelper.dropItem(_fightersHelper.getDropRarityLevel(killedFighter), damageDealt[i].fighterId, killerExp);
            }
        }       

        battles[battleHash] = block.number;
        
        emit BattleRecorded(killedFighter, battleHash, battleNonce);
    }


    function calculateExperience(uint256 fighterLevel, uint256 opponentLevel, uint256 damageDealt) public returns (uint256) {
        uint256 levelDifference;
        uint256 exp;

        fighterLevel = max(1, fighterLevel);
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
        uint256 killedFighter,
        uint256 nonce
    ) internal returns (bytes32)
    {
        return keccak256(abi.encode(this, killedFighter, nonce));
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

}
