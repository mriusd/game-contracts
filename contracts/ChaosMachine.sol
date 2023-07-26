// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "./ItemsHelper.sol";

contract ChaosMachine is ItemsAtts { 
    ItemsHelper private _itemsHelper;

    constructor (address itemsContract, address itemsHelperContract) {
        _itemsHelper = ItemsHelper(itemsHelperContract);
    }

    struct RecipeItem {
        uint256 itemId;
        uint256 minLevel;
        uint256 maxLevel;
        uint256 minAddPoints;
        uint256 maxAddPoints;
    }


    struct CombinationRecipe {
        RecipeItem[] inputItems;
        RecipeItem[] outputItems;
        uint256 successRate;
    }

    CombinationRecipe[] public Recipes;

    event NewRecipeCreated(uint256 RecipeId, CombinationRecipe);
    event ItemCrafted(uint256 tokenId, address itemOwner);
    event LogUint(uint256 logValue);

    // function createRecipe(CombinationRecipe memory newRecipe) external returns (uint256) {
    //     RecipeItem[] memory inputItems = newRecipe.inputItems;
    //     RecipeItem[] memory outputItems = newRecipe.outputItems;
        
    //     CombinationRecipe storage recipe = Recipes.push();
    //     recipe.successRate = newRecipe.successRate;
        
    //     for (uint256 i = 0; i < inputItems.length; i++) {
    //         recipe.inputItems.push(inputItems[i]);
    //     }

    //     for (uint256 i = 0; i < outputItems.length; i++) {
    //         recipe.outputItems.push(outputItems[i]);
    //     }

    //     emit NewRecipeCreated(Recipes.length - 1, newRecipe);
    //     return Recipes.length - 1;
    // }

    // function combineItems(uint256[] memory tokenIds, uint256 recipeId, address itemOwner) external {
    //     require(Recipes.length > recipeId, "Recipe not found");

    //     CombinationRecipe memory recipe = Recipes[recipeId];
    //     ItemAttributes memory item;
        
    //     for (uint256 i = 0; i <recipe.inputItems.length; i++) {
    //         item =  _itemsHelper.getTokenAttributes(tokenIds[i]);

    //         require(item.itemAttributesId == recipe.inputItems[i].itemId, "Invalid item");
    //         require(item.itemLevel >= recipe.inputItems[i].minLevel, "Item level low");
    //         require(item.additionalDamage + item.additionalDefense >= recipe.inputItems[i].minAddPoints, "Item add points low");

    //         _itemsHelper.burnItem(item.tokenId);            
    //     }

    //     if (getRandomNumber(42) < recipe.successRate) {
    //         // get random item from successItems
    //         RecipeItem memory outputItem = recipe.outputItems[getRandomNumberMax(41, recipe.outputItems.length)];
    //         uint256 randomItem = outputItem.itemId;
    //         uint256 newTokenId = _itemsHelper.craftItem(randomItem, msg.sender, outputItem.maxLevel,  outputItem.maxAddPoints);
    //         emit ItemCrafted(newTokenId, msg.sender);
    //     } else {
    //         emit ItemCrafted(0, msg.sender);
    //     }
    // }
}