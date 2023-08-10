const fs = require('fs');



const Money = artifacts.require("Money");
const MoneyHelper = artifacts.require("MoneyHelper");



const UpgradeItem = artifacts.require("UpgradeItem");
const ChaosMachine = artifacts.require("ChaosMachine");

const Trade = artifacts.require("Trade");
const TradeHelper = artifacts.require("TradeHelper");
const Drop = artifacts.require("Drop");
const DropHelper = artifacts.require("DropHelper");


module.exports = async function (deployer, network, accounts) {
  // Load the current environment variables from the .env file
  const currentEnvContent = fs.readFileSync('.env', 'utf-8');
  const envLines = currentEnvContent.split('\n');
  const envVars = {};

  envLines.forEach((line) => {
    const [key, value] = line.split('=');
    envVars[key] = value;
  });

  // Deploy the Fighters contract
  const Fighters = artifacts.require("Fighters");
  await deployer.deploy(Fighters);
  const FightersInstance = await Fighters.deployed();
  const FightersContractAddress = FightersInstance.address; 

  envVars.FIGHTERS_CONTRACT = FightersContractAddress;
  
  console.log("Fighters: ", FightersInstance.address); 
  

  // Deploy the FightersHelper contract
  const FightersHelper = artifacts.require("FightersHelper");
  await deployer.deploy(FightersHelper, FightersContractAddress);
  const FightersHelperInstance = await FightersHelper.deployed();
  const FightersHelperContractAddress = FightersHelperInstance.address;

  envVars.FIGHTERS_HELPER_CONTRACT = FightersHelperContractAddress;
  console.log("FightersHelper:        ", FightersHelperContractAddress); 






  // Deploy the Money contract
  await deployer.deploy(Money);
  const moneyInstance = await Money.deployed();
  const moneyContractAddress = moneyInstance.address; 
  // Deploy the MoneyHelper contract
  await deployer.deploy(MoneyHelper, moneyContractAddress);
  const moneyHelperInstance = await MoneyHelper.deployed();
  const moneyHelperContractAddress = moneyHelperInstance.address;



  // Deploy the Items contract
  const ItemsBase = artifacts.require("ItemsBase");
  await deployer.deploy(ItemsBase);
  const ItemsBaseInstance = await ItemsBase.deployed();
  const ItemsBaseContractAddress = ItemsBaseInstance.address; // Use the address of the deployed Items contract

  envVars.ITEMS_BASE_CONTRACT = ItemsBaseContractAddress;
  console.log("ItemsBase:             ", ItemsBaseContractAddress);

  const Items = artifacts.require("Items");
  await deployer.deploy(Items, ItemsBaseContractAddress, FightersHelperContractAddress);
  const itemsInstance = await Items.deployed();
  const itemsContractAddress = itemsInstance.address; // Use the address of the deployed Items contract

  envVars.ITEMS_CONTRACT = itemsContractAddress;
  console.log("Items:             ", itemsContractAddress);

  const ItemsExcellent = artifacts.require("ItemsExcellent");
  await deployer.deploy(ItemsExcellent, itemsContractAddress);
  const ItemsExcellentInstance = await ItemsExcellent.deployed();
  const ItemsExcellentContractAddress = ItemsExcellentInstance.address; // Use the address of the deployed Items contract

  envVars.ITEMS_EXCELLENT_CONTRACT = ItemsExcellentContractAddress;
  console.log("ItemsExcellent:             ", ItemsExcellentContractAddress);
  
  // Deploy the ItemsHelper contract
  const ItemsHelper = artifacts.require("ItemsHelper");
  await deployer.deploy(ItemsHelper, ItemsBaseContractAddress, itemsContractAddress, ItemsExcellentContractAddress);
  const itemsHelperInstance = await ItemsHelper.deployed();
  const itemsHelperContractAddress = itemsHelperInstance.address;

  envVars.ITEMS_HELPER_CONTRACT = itemsHelperContractAddress;  
  console.log("ItemsHelper:        ", itemsHelperContractAddress);






  // Deploy the Drop contract
  await deployer.deploy(Drop, itemsHelperContractAddress, moneyHelperContractAddress);
  const dropInstance = await Drop.deployed();
  const dropContractAddress = dropInstance.address; 

  // Deploy the DropHelper contract
  await deployer.deploy(DropHelper, dropContractAddress);
  const DropHelperInstance = await DropHelper.deployed();
  const DropHelperContractAddress = DropHelperInstance.address; 

  // Deploy the UpgradeItem contract
  await deployer.deploy(UpgradeItem, itemsContractAddress, itemsHelperContractAddress);
  const upgradeItemInstance = await UpgradeItem.deployed();
  const upgradeItemContractAddress = upgradeItemInstance.address;

  // Deploy the ChaosMachine contract
  await deployer.deploy(ChaosMachine, itemsContractAddress, itemsHelperContractAddress);
  const chaosMachineInstance = await ChaosMachine.deployed();
  const chaosMachineContractAddress = chaosMachineInstance.address;





  // Deploy the Battle contract
  const Battle = artifacts.require("Battle");
  await deployer.deploy(Battle, FightersHelperContractAddress, DropHelperContractAddress);
  const battleInstance = await Battle.deployed();
  const battleContractAddress = battleInstance.address; 

  envVars.BATTLE_CONTRACT = battleContractAddress;
  console.log("Battle:            ", battleContractAddress);

  const BattleHelper = artifacts.require("BattleHelper");
  await deployer.deploy(BattleHelper, battleContractAddress);
  const battleHelperInstance = await BattleHelper.deployed();
  const battleHelperContractAddress = battleHelperInstance.address; 

  envVars.BATTLE_HELPER_CONTRACT = battleHelperContractAddress;
  console.log("BattleHelper:  ", battleHelperContractAddress);










  // Deploy the Backpack contract
  const Backpack = artifacts.require("Backpack");
  await deployer.deploy(Backpack, FightersHelperContractAddress, itemsHelperContractAddress, moneyHelperContractAddress, DropHelperContractAddress);
  const backpackInstance = await Backpack.deployed();
  const backpackContractAddress = backpackInstance.address; 

  envVars.BACKPACK_CONTRACT = backpackContractAddress;
  console.log("Backpack:          ", backpackContractAddress);

  const BackpackHelper = artifacts.require("BackpackHelper");
  await deployer.deploy(BackpackHelper, backpackContractAddress);
  const BackpackHelperInstance = await BackpackHelper.deployed();
  const BackpackHelperContractAddress = BackpackHelperInstance.address; 

  envVars.BACKPACK_HELPER_CONTRACT = BackpackHelperContractAddress;
  console.log("BackpackHelper:          ", BackpackHelperContractAddress);



  // Deploy the Trade contract
  await deployer.deploy(Trade, itemsHelperContractAddress, moneyHelperContractAddress);
  const tradeInstance = await Trade.deployed();
  const tradeContractAddress = tradeInstance.address; 

  // Deploy the TradeHelper contract
  await deployer.deploy(TradeHelper, tradeContractAddress);
  const tradeHelperInstance = await TradeHelper.deployed();
  const tradeHelperContractAddress = tradeHelperInstance.address; 



  // Deploy the Backpack contract
  const Credits = artifacts.require("Credits");
  await deployer.deploy(Credits);
  const CreditsInstance = await Credits.deployed();
  const CreditsContractAddress = CreditsInstance.address; 

  envVars.CREDITS_CONTRACT = CreditsContractAddress;
  console.log("Credits:          ", CreditsContractAddress);

  const CreditsHelper = artifacts.require("CreditsHelper");
  await deployer.deploy(CreditsHelper, CreditsContractAddress);
  const CreditsHelperInstance = await CreditsHelper.deployed();
  const CreditsHelperContractAddress = CreditsHelperInstance.address; 

  envVars.CREDITS_HELPER_CONTRACT = CreditsHelperContractAddress;
  console.log("CreditsHelper:          ", CreditsHelperContractAddress);



  // Deploy Shop contract
  // Set Shop Price List
  var PriceList = {
    weaponBasePrice: 100,
    armourBasePrice: 50, 
    wingBasePrice:  3000000,

    jocPrice:  100000,
    josPrice: 1000000,
    jobPrice: 2000000,
    jolPrice: 4000000,

    rarityMultiplierPct: 5,
    levelMultiplierPct: 4,
    addPointsMultiplierPct: 3,
    luckMultiplierPct: 10,
    exceMultiplierPct: 200,

    buySellMultiplier: 5
  }

  envVars.PL_WEAPON_BASE_PRICE = PriceList.weaponBasePrice;
  envVars.PL_ARMOUR_BASE_PRICE = PriceList.armourBasePrice;
  envVars.PL_WING_BASE_PRICE = PriceList.wingBasePrice;
  envVars.PL_JOC_PRICE = PriceList.jocPrice;
  envVars.PL_JOS_PRICE = PriceList.josPrice;
  envVars.PL_JOB_PRICE = PriceList.jobPrice;
  envVars.PL_JOL_PRICE = PriceList.jolPrice;
  envVars.PL_RARITY_MULTIPLIER_PCT = PriceList.rarityMultiplierPct;
  envVars.PL_LEVEL_MULTIPLIER_PCT = PriceList.levelMultiplierPct;
  envVars.PL_ADDPOINTS_MULTIPLIER_PCT = PriceList.addPointsMultiplierPct;
  envVars.PL_LUCK_MULTIPLIER_PCT = PriceList.luckMultiplierPct;
  envVars.PL_EXCE_MULTIPLIER_PCT = PriceList.exceMultiplierPct;
  envVars.PL_BUYSELL_MULTIPLIER_PCT = PriceList.buySellMultiplier;

  const Shop = artifacts.require("Shop");
  await deployer.deploy(Shop, itemsHelperContractAddress, moneyHelperContractAddress, FightersHelperContractAddress, PriceList);
  const ShopInstance = await Shop.deployed();
  const ShopContractAddress = ShopInstance.address;




  // Perform initial Transactions
  // Create Items
  // Read the item list from itemList.json
  const itemList = JSON.parse(fs.readFileSync('./game_items.json'));

  //console.log(itemList);

  // Loop through each item in the list and call the createItem function
  for (let i = 0; i < itemList.length; i++) {

    var item = itemList[i];

    delete(item.params);

    // Call the createItem function and verify that it creates a new item
    const result = await itemsInstance.createItem(itemList[i], { from: accounts[0] });

    //const iteatts = await itemsInstance.getItemAttributes.call(result.logs[0].args.name);
    
    console.log("Added item "+(i+1)+"/"+itemList.length)
    
  }

  // const weapons = await itemsInstance.getWeaponsLength.call(0);
  // const armours= await itemsInstance.getArmoursLength.call(0);
  // const jewels = await itemsInstance.getJewelsLength.call(0);
  // const miscs = await itemsInstance.getMiscsLength.call(0);
  // console.log("Weapons ", weapons.toString());
  // console.log("Armours ", armours.toString());
  // console.log("Jewels ", jewels.toString());
  // console.log("Miscs ", miscs.toString());



  // Create drop parameters
  var result = await DropHelperInstance.setDropParams(0, {
      weaponsDropRate:  5,
      armoursDropRate: 10,
      jewelsDropRate:   1,
      miscDropRate:     0,
      boxDropRate:      50,

      excDropRate:      20,
      boxId:             2,

      minItemLevel:      0,
      maxItemLevel:      3,
      maxAddPoints:      8,

      blockCrated:       1

  }, { from: accounts[0] });
  console.log("Created Drop Paramete");


  var result = await DropHelperInstance.setBoxDropParams(0, {
    weaponsDropRate:    30,
    armoursDropRate:    60,
    jewelsDropRate:     10,
    miscDropRate:         0,
    boxDropRate:          0,

    luckDropRate:       50,
    skillDropRate:      50,
    excDropRate:        15,
    boxId:                0,

    minItemLevel:         4,
    maxItemLevel:         7,
    maxAddPoints:        12,

    blockCrated:          1

  }, { from: accounts[0] });

  console.log("Box Drop params");



          


  // Create StonedApe
  await FightersHelperInstance.createFighter(accounts[0], "StonedApe", "Warrior", { gas: 3000000, from: accounts[0] });
  fighterAtts = await FightersHelperInstance.getTokenAttributes.call(1);
  console.log("StonedApe Created = ", fighterAtts);   

  // Create NPCs
  const npcList = JSON.parse(fs.readFileSync('./npcList.json'));

  for (let i = 0; i < npcList.length; i++) {   
    var npc = npcList[i];
    //console.log("Npc", npc);

    // Create NPC class

    /*
        uint256 hpPerVitalityPoint;
        uint256 manaPerEnergyPoint;
        uint256 hpIncreasePerLevel;
        uint256 manaIncreasePerLevel;
        uint256 statPointsPerLevel;
        uint256 attackSpeed;
        uint256 agilityPointsPerSpeed;
        uint256 isNpc;
        uint256 dropRarityLevel; // for npcs        
    */
    var result = await FightersHelperInstance.updateFighterClass(npc.name, [npc.strength, npc.agility, npc.energy, npc.vitality, 0,0,0,0,0, npc.attackSpeed, 0,1, npc.dropRarityLevel] );

    console.log("NPC class created");
    var result = await FightersHelperInstance.createNPC(npc.name, npc.name, accounts[0], { gas: 3000000, from: accounts[0] });
    console.log("NPC Created: ")     

  }



  // Update the contract addresses



 





  envVars.MONEY_CONTRACT = moneyContractAddress;
  envVars.MONEY_HELPER_CONTRACT = moneyHelperContractAddress;
  console.log("Money:             ", moneyContractAddress);
  console.log("MoneyHelper:       ", moneyHelperContractAddress);


  envVars.UPGRADE_ITEM_CONTRACT = upgradeItemContractAddress;
  console.log("Upgrade Item:      ", upgradeItemContractAddress);


  envVars.CHAOS_MACHINE_CONTRACT = chaosMachineContractAddress;
  console.log("Chaos Machine:     ", chaosMachineContractAddress);





  envVars.TRADE_CONTRACT = tradeContractAddress;
  envVars.TRADE_HELPER_CONTRACT = tradeHelperContractAddress;
  console.log("Trade:       ", tradeContractAddress);
  console.log("TradeHelper:       ", tradeHelperContractAddress);


  envVars.DROP_CONTRACT = dropContractAddress;
  envVars.DROP_HELPER_CONTRACT = DropHelperContractAddress;
  console.log("Drop:        ", dropContractAddress);
  console.log("DropHelper:        ", DropHelperContractAddress);
  

  // Convert the updated environment variables back to the file content
  const updatedEnvContent = Object.entries(envVars)
    .map(([key, value]) => `${key}=${value}`)
    .join('\n');

  fs.writeFileSync('.env', updatedEnvContent);
};
