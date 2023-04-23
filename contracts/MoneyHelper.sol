// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./Money.sol";
contract MoneyHelper {
    Money private _money;
    uint256 experienceDivider = 100;

    constructor(address moneyAddress) {
        _money = Money(moneyAddress);
    }

    function mintGold(address playerAddress, uint256 experience) external {
        _money.mintGold(playerAddress, max(1, experience/experienceDivider));
    }

    // Returns the largest of the two values
    function max(uint a, uint b) private pure returns (uint) {
        return a > b ? a : b;
    }
}