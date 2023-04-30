// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./Money.sol";
contract MoneyHelper {
    Money private _money;
    uint256 public experienceDivider = 1;

    constructor(address moneyAddress) {
        _money = Money(moneyAddress);
    }

    function mintGold(address playerAddress, uint256 amount) external {
        _money.mintGold(playerAddress, max(1, amount)*1e18);
    }

    // Returns the largest of the two values
    function max(uint a, uint b) private pure returns (uint) {
        return a > b ? a : b;
    }

    function getExperienceDivider() public view returns(uint256) {
        return experienceDivider;
    }
}