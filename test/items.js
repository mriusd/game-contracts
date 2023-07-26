const fs = require('fs');
const FighterAttributes = artifacts.require('FighterAttributes'); // Import the ABI and bytecode for your contract
const Items = artifacts.require('Items'); // Import the ABI and bytecode for your contract
const UpgradeItem = artifacts.require('UpgradeItem'); // Import the ABI and bytecode for your contract
const ChaosMachine = artifacts.require('ChaosMachine'); // Import the ABI and bytecode for your contract
const Backpack = artifacts.require('Backpack'); // Import the ABI and bytecode for your contract
const TradeHelper = artifacts.require('TradeHelper'); // Import the ABI and bytecode for your contract

function isExcellent(item) {
  if (item.lifeAfterMonsterIncrease == 1 || 
      item.manaAfterMonsterIncrease == 1 || 
      item.excellentDamageProbabilityIncrease == 1 || 
      item.attackSpeedIncrease == 1 ||
      item.damageIncrease == 1 ||

      item.defenseSuccessRateIncrease == 1 ||
      item.goldAfterMonsterIncrease == 1 ||
      item.reflectDamage == 1 ||
      item.maxLifeIncrease == 1 ||
      item.maxManaIncrease == 1 ||
      item.hpRecoveryRateIncrease == 1 ||
      item.mpRecoveryRateIncrease == 1 ||
      item.decreaseDamageRateIncrease == 1
  ) {
    return true;
  }

  return false;
}

function generateItemName(item) {
  var itemName = item.name;
  
  if (item.itemLevel > 0) {
    itemName += " +"+item.itemLevel;
  }

  if (item.luck) {
    itemName += " +Luck";
  }

  if (item.skill) {
    itemName += " +Skill";
  }

  if (isExcellent(item)) {
    itemName = "Exc "+itemName;
  }

  if (item.additionalDamage > 0) {
    itemName += " +"+item.additionalDamage;
  } 

  if (item.additionalDefense > 0) {
    itemName += " +"+item.additionalDefense;
  } 

  return itemName;
}


async function updateItemLevel(item, level, itemsContractInstance, upgradeItemsInstance, accounts) {

  var totalBlesses = 0;
  var totalSouls = 0;
  var blesses = 0;
  var souls = 0;
  var counter = 0;

  while (item.itemLevel < level) {
    if (item.itemLevel < 6) {
      // buy bless
      result = await itemsContractInstance.buyItemFromShop(2, 1, { from: accounts[1] });
      var jewelId = result.logs[0].args.tokenId;
      //const jewelId = await itemsContractInstance.getTokenAttributes.call(result.logs[0].args.tokenId);


      result = await upgradeItemsInstance.upgradeItemLevel(item.tokenId, jewelId, { from: accounts[1] });
      item = await itemsContractInstance.getTokenAttributes.call(item.tokenId);

      console.log("Bless throun (B= , S= ", generateItemName(item));
      blesses++;
      totalBlesses++;
    } else {
      result = await itemsContractInstance.buyItemFromShop(3, 1, { from: accounts[1] });
      var jewelId = result.logs[0].args.tokenId;

      result = await upgradeItemsInstance.upgradeItemLevel(item.tokenId, jewelId, { from: accounts[1] });
      item = await itemsContractInstance.getTokenAttributes.call(item.tokenId);
      console.log("Soul throun ", generateItemName(item));
      souls++;
      totalSouls++;
    }
    counter++;
  }

  console.log("Item upgraded to +9 (Luck: "+item.luck+", B: "+totalBlesses+", S: "+totalSouls+")");
}

async function updateAddPoints(item, points, itemsContractInstance, upgradeItemsInstance, accounts) {
  var totalJols = 0;
  var jols = 0;

  while (parseInt(item.additionalDamage)+parseInt(item.additionalDefense) < points) {
    
      result = await itemsContractInstance.buyItemFromShop(4, 1, { from: accounts[1] });
      var jewelId = result.logs[0].args.tokenId;

      result = await upgradeItemsInstance.updateItemAdditionalPoints(item.tokenId, jewelId, { from: accounts[1] });
      item = await itemsContractInstance.getTokenAttributes.call(item.tokenId);
      
      jols++;
      totalJols++;
      console.log("JOLs throun ("+totalJols+") ", generateItemName(item));
  }
} 



contract('Items', (accounts) => {
  it.skip('create all items in the itemList', async () => {
    const itemsContractInstance = await Items.deployed();


    // Read the item list from itemList.json
    const itemList = JSON.parse(fs.readFileSync('./itemsList.json'));

    //console.log(itemList);

    // Loop through each item in the list and call the createItem function
    for (let i = 0; i < itemList.length; i++) {

      // Call the createItem function and verify that it creates a new item
      const result = await itemsContractInstance.createItem(itemList[i], { from: accounts[1] });

      const iteatts = await itemsContractInstance.getItemAttributes.call(result.logs[0].args.itemId);
      
      console.log("Added item "+(i+1)+"/"+itemList.length)
      
    }

    const weapons = await itemsContractInstance.getWeaponsLength.call(0);
    const armours= await itemsContractInstance.getArmoursLength.call(0);
    const jewels = await itemsContractInstance.getJewelsLength.call(0);
    const miscs = await itemsContractInstance.getMiscsLength.call(0);
    console.log("Weapons ", weapons.toString());
    console.log("Armours ", armours.toString());
    console.log("Jewels ", jewels.toString());
    console.log("Miscs ", miscs.toString());
  });

  it.skip('should craft 10 items', async () => {
      const itemsContractInstance = await Items.deployed();

      for (var i=0; i<10;i++) {
        var result = await itemsContractInstance.craftItem(6, accounts[0], { from: accounts[1] });
        var item = await itemsContractInstance.getTokenAttributes.call(result.logs[0].args.tokenId);
        console.log("Item creafted: ", generateItemName(item));
      }
      
  });

  it.skip('buy two items in the shop and combine them in the chaos machine', async () => {
    const itemsContractInstance = await Items.deployed();
    const upgradeItemsInstance = await UpgradeItem.deployed();

    console.log("itemsContractInstance.address", itemsContractInstance.address);
    const chaosMachineInstance = await ChaosMachine.deployed();

    var newRecipe = {
      inputItems: [{
          itemId: 6,
          minLevel: 0,
          maxLevel: 15,
          minAddPoints: 12,
          maxAddPoints: 28
        }, 
        {
          itemId: 2,
          minLevel: 0,
          maxLevel: 15,
          minAddPoints: 0,
          maxAddPoints: 28
        }],
      successRate: 100,
      outputItems: [
        {
          itemId: 11,
          minLevel: 0,
          maxLevel: 4,
          minAddPoints: 0,
          maxAddPoints: 12
        },
        {
          itemId: 12,
          minLevel: 0,
          maxLevel: 4,
          minAddPoints: 0,
          maxAddPoints: 12
        },
        {
          itemId: 13,
          minLevel: 0,
          maxLevel: 4,
          minAddPoints: 0,
          maxAddPoints: 12
        }
      ]
    }

    var result = await chaosMachineInstance.createRecipe(newRecipe, { from: accounts[1] });


    for (var i=0; i<100;i++) {
      // Call the createItem function and verify that it creates a new item
      result = await itemsContractInstance.buyItemFromShop(6, 1, { from: accounts[1] });

      var item = await itemsContractInstance.getTokenAttributes.call(result.logs[0].args.tokenId);

      await updateItemLevel(item, 3, itemsContractInstance, upgradeItemsInstance, accounts);
      await updateAddPoints(item, 12, itemsContractInstance, upgradeItemsInstance, accounts);

      result = await itemsContractInstance.buyItemFromShop(2, 1, { from: accounts[1] });
      var jewelId = result.logs[0].args.tokenId;
      console.log("Combining "+item.tokenId+" with "+jewelId); 
     
     result = await chaosMachineInstance.combineItems([item.tokenId, jewelId], 0, accounts[1], { gas: 3000000, from: accounts[1] });

     var newItemId = result.logs[0].args.tokenId;

      if (newItemId == 0) {
        console.log("Combinnation failed");
      } else {
        //console.log("Getting item attributes: ", result.logs);
        var newItem = await itemsContractInstance.getTokenAttributes.call(newItemId);

        console.log("Combinnation ended in: ", generateItemName(newItem));
      }


          
        
      
    }
    
  });

  it.skip('should buy an item from the shop and upgrade it to +9', async () => {
    const itemsContractInstance = await Items.deployed();
    const upgradeItemsInstance = await UpgradeItem.deployed();

    var totalBlesses = 0;
    var totalSouls = 0;
    for (var i=0; i<1;i++) {
      // Call the createItem function and verify that it creates a new item
      var result = await itemsContractInstance.buyItemFromShop(6, 1, { from: accounts[1] });

      var item = await itemsContractInstance.getTokenAttributes.call(result.logs[0].args.tokenId);
        
      await updateItemLevel(item, 9, itemsContractInstance, upgradeItemsInstance, accounts);

      var item = await itemsContractInstance.getTokenAttributes.call(result.logs[0].args.tokenId);
    }
    
  });

  it.skip('should buy an item from the shop and upgrade it to +28 add points', async () => {
    const itemsContractInstance = await Items.deployed();
    const upgradeItemsInstance = await UpgradeItem.deployed();

    var totalJols = 0;
    for (var i=0; i<1;i++) {
      // Call the createItem function and verify that it creates a new item
      var result = await itemsContractInstance.buyItemFromShop(6, 1, { from: accounts[1] });

      var item = await itemsContractInstance.getTokenAttributes.call(result.logs[0].args.tokenId);
        
      await updateAddPoints(item, 28, itemsContractInstance, upgradeItemsInstance, accounts);

      var item = await itemsContractInstance.getTokenAttributes.call(result.logs[0].args.tokenId);

      console.log("Item Upgraded to +28 add points (Luck: "+item.luck+")", generateItemName(item));
    }    
  });

  it.skip('should create dropPramas for rarity 0', async () => {
    const itemsContractInstance = await Items.deployed();
    const result = await itemsContractInstance.setDropParams(0, {
        weaponsDropRate:  50,
        armoursDropRate: 100,
        jewelsDropRate:    5,
        miscDropRate:      0,
        boxDropRate:      25,

        excDropRate:      20,
        boxId:             1,

        minItemLevel:      0,
        maxItemLevel:      3,
        maxAddPoints:      8,

        blockCrated:       1

    }, { from: accounts[1] });

    console.log("Drop params", result.logs[0].args.params);
  });

  it.skip('should create boxDropPramas for rarity 0', async () => {
    const itemsContractInstance = await Items.deployed();
    const result = await itemsContractInstance.setBoxDropParams(0, {
      weaponsDropRate:    300,
      armoursDropRate:    600,
      jewelsDropRate:     100,
      miscDropRate:         0,
      boxDropRate:          0,

      luckDropRate:       500,
      skillDropRate:      500,
      excDropRate:        150,
      boxId:                0,

      minItemLevel:         4,
      maxItemLevel:         7,
      maxAddPoints:        12,

      blockCrated:          1

    }, { from: accounts[1] });

    console.log("Box Drop params", result.logs[0].args.params);
  });

  var boxes = [];
  it('should drop a random item', async () => {
    const itemsContractInstance = await Items.deployed();
    const backpackContractInstance = await Backpack.deployed();
    var drops = {};
    var itemsDropped = 0;
    var luck = 0;
    var categoryDrops = {
      weapons: 0,
      armours: 0,
      jewels: 0,
      miscs: 0,
      boxes: 0,
      gold: 0,
      nothing: 0,

      luck: 0,
      skill: 0,
      exc: 0,
      addPoints: 0
    };
    for (var i = 0; i<1000; i++) {
      
      var result = await itemsContractInstance.dropItem(0, 1, 2000, { gas: 3000000, from: accounts[1] });
      //console.log("Drop Logs: ", result.logs[0]);
      const droppedItemHash = result.logs[0].args.itemHash;
      var item = result.logs[0].args.item;
      var qty = result.logs[0].args.qty.toString();
      const blockNumber = result.receipt.blockNumber;
      // console.log("Drop droppedItemHash: ", droppedItemHash);
      //console.log("Drop blockNumber: ", blockNumber);
      
      
      itemsDropped++;
      // const item = await itemsContractInstance.getTokenAttributes.call(droppedItemId);

      var itemName = item.name;
    
      if (typeof(drops[item.name]) == 'undefined') {
        drops[item.name]  = 0;
      }

      if (item.luck) {
        categoryDrops.luck++;
      }

      if (item.skill) {
        categoryDrops.skill++;
      }

      if (isExcellent(item)) {
        categoryDrops.exc++;
      } 

      if (item.additionalDamage > 0) {
        categoryDrops.addPoints++;
      } 

      if (item.additionalDefense > 0) {
        categoryDrops.addPoints++;
      } 
      

      drops[item.name]++;

      if (item.isWeapon) { categoryDrops.weapons++; }
      else if (item.isArmour) { categoryDrops.armours++; }
      else if (item.isJewel) { categoryDrops.jewels++; }
      else if (item.isMisc) { categoryDrops.miscs++; }
      else if (item.isBox) { 
        boxes.push(item);
        categoryDrops.boxes++; 
      }

      console.log('['+i+'] Dropped ', qty, ' ', generateItemName(item));    

      // pickup item 
      result = await backpackContractInstance.pickupItem(droppedItemHash, item, blockNumber, 1, { gas: 3000000, from: accounts[1] });
      const droppedItemId = result.logs[0].args.tokenId;
      var pickedItem = await itemsContractInstance.getTokenAttributes.call(droppedItemId);
      //console.log('['+i+'] Picked Up Args  ', result);    
      console.log('['+i+'] Picked Up  ', result.logs[0].args.qty.toString(), ' ', generateItemName(pickedItem));    
    }

    console.log("Total Items: ", itemsDropped);
    console.log("Item Drops: ", drops);
    console.log("Category Drops: ", categoryDrops);
  });

  it.skip('should open all the boxes', async () => {
    const itemsContractInstance = await Items.deployed();
    for (var i =0; i<boxes.length; i++) {
      var box = boxes[i];
      const result = await itemsContractInstance.openBox(box.tokenId, { gas: 3000000, from: accounts[1] });
      const droppedItemId = result.logs[0].args.tokenId;
      const item = await itemsContractInstance.getTokenAttributes.call(droppedItemId);
      console.log("Oppened box ["+i+"]", generateItemName(item));
    }
  });

  it.skip('give a uniform random number distribution', async () => {
    const itemsContractInstance = await Items.deployed();
    var numbers = {};
    var tries = 5000;
    for (var i =0; i<100; i++) {
      numbers[i] = 0;
    }

    for (var i =0; i<tries; i++) {
      result = await itemsContractInstance.genRandomNumber(0, { from: accounts[0] });
      randomNumber = parseInt(result.logs[0].args.randomNumber.toString());
      numbers[randomNumber]++;
      console.log("Random number: ", randomNumber);
    }

    console.log("Distribution: ", numbers);
    
    
  })

  it('should buy an item from the shop and trade it with another char', async () => {
    const itemsContractInstance = await Items.deployed();
    const tradeHelperInstance = await TradeHelper.deployed();

    var totalBlesses = 0;
    var totalSouls = 0;
    for (var i=0; i<1;i++) {
      // Call the createItem function and verify that it creates a new item
      var result = await itemsContractInstance.buyItemFromShop(6, 1, { from: accounts[1] });

      var item = await itemsContractInstance.getTokenAttributes.call(result.logs[0].args.tokenId);
        
      await updateItemLevel(item, 9, itemsContractInstance, upgradeItemsInstance, accounts);

      var item = await itemsContractInstance.getTokenAttributes.call(result.logs[0].args.tokenId);
    }
    
  });
});
