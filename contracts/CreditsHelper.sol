// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "./Credits.sol";

contract CreditsHelper is SafeMath {
    Credits private _credits;

    uint256 maxFaucetAmount = 1000;
    uint256 dailyFaucetLimit = 150;

    constructor(address baseContract) {
        _credits = Credits(baseContract);
    }

    event Faucet(address owner, uint256 amount);   

    function balanceOf(address account) public view returns (uint256) {
        return _credits.balanceOf(account);
    }

    function mintCredits(address playerAddress, uint256 amount) external returns (bool)  {
        return _credits.mintCredits(playerAddress, amount);
    }

    function faucetCredits(address playerAddress, uint256 amount) external  {
        require(_credits.balanceOf(playerAddress) < maxFaucetAmount, "Max credits reached for faucet");

        amount = min(safeSub(maxFaucetAmount, amount), dailyFaucetLimit);
        _credits.mintCredits(playerAddress, amount);
        emit Faucet(playerAddress, amount);
    }

    function burnCredits(address playerAddress, uint256 amount) external returns (bool) {
        return _credits.burnCredits(playerAddress, amount);
    }

    function transferThroughTrade(address sender, address recipient, uint256 amount) public  returns (bool) {
        return _credits.transferThroughTrade(sender, recipient, amount);
    }
}