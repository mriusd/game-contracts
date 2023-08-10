// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./SafeMath.sol";

contract Money is IERC20, SafeMath {
    string public constant name = "MRIUSD Gold";
    string public constant symbol = "MRIUSG";
    uint8 public constant decimals = 18;
    uint256 private _totalSupply;
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;

    event Mint(address receiver, uint256 amount);
    event Burn(address receiver, uint256 amount);
    
    
    constructor() {
        _totalSupply = 0;
    }    

    function mintGold(address playerAddress, uint256 amount) external {
        _totalSupply = safeAdd(_totalSupply, amount);
        _balances[playerAddress] = safeAdd(_balances[playerAddress], amount);
        emit Mint(playerAddress, amount);
        emit Transfer(address(0), playerAddress, _balances[playerAddress]);
    }

    function burnGold(address playerAddress, uint256 amount) external {
        _totalSupply = _totalSupply - amount;
        _balances[playerAddress] = _balances[playerAddress] - amount;
        emit Burn(playerAddress, amount);
        emit Transfer(playerAddress, address(0), _balances[playerAddress]);
    }

    function totalSupply() public view override returns (uint256) {
        return _totalSupply;
    }
    
    function balanceOf(address account) public view override returns (uint256) {
        return _balances[account];
    }
    
    function transfer(address recipient, uint256 amount) public override returns (bool) {
        _transfer(msg.sender, recipient, amount);
        return true;
    }

    function transferThroughTrade(address sender, address recipient, uint256 amount) public returns (bool) {
        _transfer(sender, recipient, amount);
        return true;
    }
    
    function allowance(address owner, address spender) public view override returns (uint256) {
        return _allowances[owner][spender];
    }
    
    function approve(address spender, uint256 amount) public override returns (bool) {
        _approve(msg.sender, spender, amount);
        return true;
    }
    
    function transferFrom(address sender, address recipient, uint256 amount) public override returns (bool) {
        _transfer(sender, recipient, amount);
        _approve(sender, msg.sender, _allowances[sender][msg.sender] - amount);
        return true;
    }
    
    function _transfer(address sender, address recipient, uint256 amount) internal {
        require(sender != address(0), "ERC20: transfer from the zero address");
        require(recipient != address(0), "ERC20: transfer to the zero address");
        
        _balances[sender] -= amount;
        _balances[recipient] += amount;
        emit Transfer(sender, recipient, amount);
    }
    
    function _approve(address owner, address spender, uint256 amount) internal {
        require(owner != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");
        
        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }
}
