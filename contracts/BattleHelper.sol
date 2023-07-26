// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./Battle.sol";
contract BattleHelper {
    Battle private _battle;

    constructor(address battleAddress) {
        _battle = Battle(battleAddress);
    }    

    function recordKill(uint256 killedFighter, Battle.Damage[] memory damageDealt, uint256 battleNonce) external  {
        return _battle.recordKill(killedFighter, damageDealt, battleNonce);
    }
}