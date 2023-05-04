const fs = require('fs');

const FighterAttributes = artifacts.require("FighterAttributes");
const FighterHelper = artifacts.require("FighterHelper");
const Money = artifacts.require("Money");
const MoneyHelper = artifacts.require("MoneyHelper");
const Battle = artifacts.require("Battle");
const Items = artifacts.require("Items");
const ItemsHelper = artifacts.require("ItemsHelper");
const UpgradeItem = artifacts.require("UpgradeItem");
const ChaosMachine = artifacts.require("ChaosMachine");
const Backpack = artifacts.require("Backpack");


module.exports = async function (deployer, network, accounts) {
  // Deploy the FighterAttributes contract
  await deployer.deploy(FighterAttributes);
  const fighterAttributesInstance = await FighterAttributes.deployed();
  const fighterAttributesContractAddress = fighterAttributesInstance.address; 

  // Deploy the FighterHelper contract
  await deployer.deploy(FighterHelper, fighterAttributesContractAddress);
  const fighterHelperInstance = await FighterHelper.deployed();
  const fighterHelperContractAddress = fighterHelperInstance.address;

  // Deploy the Money contract
  await deployer.deploy(Money);
  const moneyInstance = await Money.deployed();
  const moneyContractAddress = moneyInstance.address; 
    // Deploy the ItemsHelper contract
  await deployer.deploy(MoneyHelper, moneyContractAddress);
  const moneyHelperInstance = await MoneyHelper.deployed();
  const moneyHelperContractAddress = moneyHelperInstance.address;



  // Deploy the Items contract
  await deployer.deploy(Items, fighterHelperContractAddress, moneyHelperContractAddress);
  const itemsInstance = await Items.deployed();
  const itemsContractAddress = itemsInstance.address; // Use the address of the deployed Items contract
  
  // Deploy the ItemsHelper contract
  await deployer.deploy(ItemsHelper, itemsContractAddress);
  const itemsHelperInstance = await ItemsHelper.deployed();
  const itemsHelperContractAddress = itemsHelperInstance.address;

  // Deploy the UpgradeItem contract
  await deployer.deploy(UpgradeItem, itemsContractAddress, itemsHelperContractAddress);
  const upgradeItemInstance = await UpgradeItem.deployed();
  const upgradeItemContractAddress = upgradeItemInstance.address;

  // Deploy the ChaosMachine contract
  await deployer.deploy(ChaosMachine, itemsContractAddress, itemsHelperContractAddress);
  const chaosMachineInstance = await ChaosMachine.deployed();
  const chaosMachineContractAddress = chaosMachineInstance.address;


  // Deploy the Battle contract
  await deployer.deploy(Battle, fighterHelperContractAddress, itemsHelperContractAddress);
  const battleInstance = await Battle.deployed();
  const battleContractAddress = battleInstance.address; 


  // Deploy the Backpack contract
  await deployer.deploy(Backpack, fighterHelperContractAddress, itemsHelperContractAddress, moneyHelperContractAddress);
  const backpackInstance = await Backpack.deployed();
  const backpackContractAddress = backpackInstance.address; 


  // Perform initial Transactions
  // Create Items
  // Read the item list from itemList.json
  const itemList = JSON.parse(fs.readFileSync('./itemsList.json'));

  //console.log(itemList);

  // Loop through each item in the list and call the createItem function
  for (let i = 0; i < itemList.length; i++) {

    // Call the createItem function and verify that it creates a new item
    const result = await itemsInstance.createItem(itemList[i], { from: accounts[0] });

    const iteatts = await itemsInstance.getItemAttributes.call(result.logs[0].args.itemId);
    
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
  var result = await itemsInstance.setDropParams(0, {
      weaponsDropRate:  5,
      armoursDropRate: 10,
      jewelsDropRate:   1,
      miscDropRate:     0,
      boxDropRate:      3,

      excDropRate:      20,
      boxId:             2,

      minItemLevel:      0,
      maxItemLevel:      3,
      maxAddPoints:      8,

      blockCrated:       1

  }, { from: accounts[0] });
  console.log("Created Drop Paramete");


  var result = await itemsInstance.setBoxDropParams(0, {
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
  await fighterAttributesInstance.createFighter("StonedApe", 1, { gas: 3000000, from: accounts[0] });
  //fighterAtts = await fighterAttributesInstance.getTokenAttributes.call(1);
  console.log("StonedApe Created")   

  // Create NPCs
  const npcList = JSON.parse(fs.readFileSync('./npcList.json'));

  for (let i = 0; i < npcList.length; i++) {   
    var npc = npcList[i];
    //console.log("Npc", npc);
    var result = await fighterAttributesInstance.createNPC(npc.name, npc.strength, npc.agility, npc.energy, npc.vitality, npc.attackSpeed, npc.level, npc.dropRarityLevel, { gas: 3000000, from: accounts[0] });
    console.log("NPC Created: ", result.logs[0].args.tokenId.toString())     

  }

  // Load the current environment variables from the .env file
  const currentEnvContent = fs.readFileSync('.env', 'utf-8');
  const envLines = currentEnvContent.split('\n');
  const envVars = {};

  envLines.forEach((line) => {
    const [key, value] = line.split('=');
    envVars[key] = value;
  });

  // Update the contract addresses
  envVars.FIGHTER_ATTRIBUTES_CONTRACT = fighterAttributesContractAddress;
  envVars.ITEMS_CONTRACT = itemsContractAddress;
  envVars.BATTLE_CONTRACT = battleContractAddress;
  envVars.MONEY_CONTRACT = moneyContractAddress;
  envVars.UPGRADE_ITEM_CONTRACT = upgradeItemContractAddress;
  envVars.CHAOS_MACHINE_CONTRACT = chaosMachineContractAddress;
  envVars.BACKPACK_CONTRACT = backpackContractAddress;

  console.log("FighterAttributes: ", fighterAttributesInstance.address);  
  console.log("Items:             ", itemsContractAddress);
  console.log("Battle:            ", battleContractAddress);
  console.log("Money:             ", moneyContractAddress);
  console.log("Upgrade Item:      ", upgradeItemContractAddress);
  console.log("Chaos Machine:     ", chaosMachineContractAddress);
  console.log("Backpack:     ", backpackContractAddress);

  // Convert the updated environment variables back to the file content
  const updatedEnvContent = Object.entries(envVars)
    .map(([key, value]) => `${key}=${value}`)
    .join('\n');

  fs.writeFileSync('.env', updatedEnvContent);
};
