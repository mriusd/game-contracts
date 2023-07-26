// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;
import "./ItemsHelper.sol";
import "./MoneyHelper.sol";

contract Trade {
    ItemsHelper private _itemsHelper;
    MoneyHelper private _moneyHelper;

    mapping(bytes32 => bool) public tradeHashes;

    event LogError(uint8, uint256);
    event BattleRecorded(uint256 killedFighter, bytes32 battleHash, uint256 battleNonce);


    constructor(address itemsHelperContract, address moneyHelperContract) {     
        _itemsHelper = ItemsHelper(itemsHelperContract);
        _moneyHelper = MoneyHelper(moneyHelperContract);
    }

    
    function trade (address from, address to, uint256[] calldata fromItems, uint256[] calldata toItems, uint256 fromMoney, uint256 toMoney) external {
        for (uint i = 0; i < fromItems.length; i++) {
            // use fromItems[i] here to access each item
            _itemsHelper.transferItem(fromItems[i], to);
        }
        for (uint i = 0; i < toItems.length; i++) {
            // use fromItems[i] here to access each item
            _itemsHelper.transferItem(toItems[i], from);
        }

        _moneyHelper.transferThroughTrade(from, to, fromMoney);
        _moneyHelper.transferThroughTrade(to, from, toMoney);
    }
}
