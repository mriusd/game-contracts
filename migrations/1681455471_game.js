const FighterAttributes = artifacts.require("FighterAttributes");
const Money = artifacts.require("Money");
const Battle = artifacts.require("Battle");
const Items = artifacts.require("Items");
const ItemsHelper = artifacts.require("ItemsHelper");
const UpgradeItem = artifacts.require("UpgradeItem");
const ChaosMachine = artifacts.require("ChaosMachine");

module.exports = async function (deployer, network, accounts) {
  // Deploy the FighterAttributes contract
  await deployer.deploy(FighterAttributes);
  const fighterAttributesInstance = await FighterAttributes.deployed();
  const fighterAttributesContractAddress = fighterAttributesInstance.address; 

  // Deploy the Money contract
  await deployer.deploy(Money);
  const moneyInstance = await Money.deployed();
  const moneyContractAddress = moneyInstance.address; 



  // Deploy the Items contract
  await deployer.deploy(Items);
  const itemsInstance = await Items.deployed();
  const itemsContractAddress = itemsInstance.address; // Use the address of the deployed Items contract
  
  // Deploy the ItemsHelper contract
  await deployer.deploy(ItemsHelper, itemsContractAddress);
  const itemsHelperInstance = await ItemsHelper.deployed();
  const itemsHelperContractAddress = itemsHelperInstance.address;

  // Deploy the UpgradeItem contract
  await deployer.deploy(UpgradeItem, itemsContractAddress, itemsHelperContractAddress);

  // Deploy the ChaosMachine contract
  await deployer.deploy(ChaosMachine, itemsContractAddress, itemsHelperContractAddress);



  // Deploy the Battle contract
  await deployer.deploy(Battle, fighterAttributesContractAddress, moneyContractAddress, itemsHelperContractAddress);
  const battleInstance = await Money.deployed();
  const battleContractAddress = battleInstance.address; 
};
