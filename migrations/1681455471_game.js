const Items = artifacts.require("Items");

module.exports = function(deployer) {
  deployer.deploy(Items);
};
