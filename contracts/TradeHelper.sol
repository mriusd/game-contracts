// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./Trade.sol";
contract TradeHelper {
    Trade private _trade;

    constructor(address tradeAddress) {
        _trade = Trade(tradeAddress);
    }    

    function trade (address from, address to, uint256[] calldata fromItems, uint256[] calldata toItems, uint256 fromMoney, uint256 toMoney) external {
        _trade.trade(from, to, fromItems, toItems, fromMoney, toMoney);
    }
}