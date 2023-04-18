const Items = artifacts.require("Items");
const UpgradeItem = artifacts.require("UpgradeItem");
const ChaosMachine = artifacts.require("ChaosMachine");

module.exports = async function (deployer, network, accounts) {
  // Deploy the Items contract
  await deployer.deploy(Items);
  const itemsInstance = await Items.deployed();
  const itemsContractAddress = itemsInstance.address; // Use the address of the deployed Items contract
  
  // Deploy the UpgradeItem contract
  await deployer.deploy(UpgradeItem, itemsContractAddress);

  // Deploy the ChaosMachine contract
  await deployer.deploy(ChaosMachine, itemsContractAddress);
};