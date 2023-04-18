// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "./ItemAtts.sol";

abstract contract MainItems {
    function setAdditionalPoints(uint256 tokenId, uint256 points) external virtual;
    function setItemLevel(uint256 tokenId, uint256 level)  external virtual;
    function burnItem(uint256 itemId)  external virtual;
    function getTokenAttributes(uint256 tokenId) external virtual view returns (ItemAtts.ItemAttributes memory);
    function craftItem(uint256 itemId, address itemOwner) external virtual returns (uint256);
}

contract ChaosMachine is ItemAtts {
    address ItemsContractAddress;   

    constructor (address itemsContract) {
        ItemsContractAddress = itemsContract;
    }

    struct CombinationRecipe {
        uint256[] itemIds;
        uint256 successRate;
        uint256[] successItemIds;
    }

    CombinationRecipe[] Recipes;

    event NewRecipeCreated(uint256 RecipeId, CombinationRecipe);
    event ItemCrafted(uint256 tokenId, address itemOwner);
    event LogUint(uint256 logValue);

    function createRecipe(CombinationRecipe memory newRecipe) external returns (uint256) {
        Recipes.push(newRecipe);
        emit NewRecipeCreated(Recipes.length -1, newRecipe);
        return Recipes.length -1;
    }

    function combineItems(uint256[] memory tokenIds, uint256 recipeId, address itemOwner) external {
        require(Recipes.length > recipeId, "Recipe not found");

        CombinationRecipe memory recipe = Recipes[recipeId];
        ItemAttributes memory item;

        
        for (uint256 i = 0; i <recipe.itemIds.length; i++) {
            item =  MainItems(ItemsContractAddress).getTokenAttributes(tokenIds[i]);

            require(item.itemAttributesId == recipe.itemIds[i], "Invalid item");

            MainItems(ItemsContractAddress).burnItem(item.tokenId);            
        }

        // if (getRandomNumber(42) < recipe.successRate) {
            // get random item from successItems
            uint256 randomItem = recipe.successItemIds[getRandomNumberMax(41, recipe.successItemIds.length)];
            uint256 newTokenId = MainItems(ItemsContractAddress).craftItem(randomItem, msg.sender);
            emit ItemCrafted(newTokenId, msg.sender);
        // } else {
        //     emit ItemCrafted(0, msg.sender);
        // }
    }
}