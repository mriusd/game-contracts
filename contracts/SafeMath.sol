pragma solidity ^0.8.18;


contract SafeMath {
    // Returns the smaller of two values
    function min(uint a, uint b) public  view returns (uint) {
        return a < b ? a : b;
    }

    // Returns the largest of the two values
    function max(uint a, uint b) public  view returns (uint) {
        return a > b ? a : b;
    }

    // Safe Multiply Function - prevents integer overflow 
    function safeMul(uint a, uint b) public view returns (uint) {
        uint c = a * b;
        assert(a == 0 || c / a == b);
        return c;
    }

    // Safe Subtraction Function - prevents integer overflow 
    function safeSub(uint a, uint b) public view returns (uint) {
        assert(b <= a);
        return a - b;
    }

    // Safe Addition Function - prevents integer overflow 
    function safeAdd(uint a, uint b) public view returns (uint) {
        uint c = a + b;
        assert(c>=a && c>=b);
        return c;
    }

    // Cube Root function
    function sqrt(uint y) internal  view returns (uint z) {
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

    function safePow(uint256 base, uint256 exponent) internal  view returns (uint256) {
        if (exponent == 0) {
            return 1;
        } else if (exponent == 1) {
            return base;
        }

        uint256 result = base;
        for (uint256 i = 1; i < exponent; i++) {
            result = safeMul(result, base); // result.mul(base);
        }
        return result;
    }

    function getRandomNumber (uint256 seed) public view returns (uint256) {
        uint256 randomNumber = uint256(keccak256(abi.encodePacked(block.prevrandao, seed, block.number, block.timestamp, msg.sender)));

        return randomNumber % 100;
    }

    function getRandomNumberMax (uint256 seed, uint256 maxNum)  public view returns (uint256) {
        uint256 randomNumber = uint256(keccak256(abi.encodePacked(block.prevrandao, seed, block.number, block.timestamp, msg.sender)));

        return randomNumber % maxNum;
    }

    function stringsEqual(string memory a, string memory b) public view returns (bool) {
        if (keccak256(abi.encodePacked(a)) == keccak256(abi.encodePacked(b))) {
            return true;
        } else {
            return false;
        }
    }
}