const fs = require('fs');
const Items = artifacts.require('Items'); // Import the ABI and bytecode for your contract




contract('Items', (accounts) => {

  it('should create all items in the list', async () => {
    // const itemsContractInstance = await Items.deployed();
    const networkId = await web3.eth.net.getId();
    const contractAddress = Items.networks[networkId].address;
    const itemsContractInstance = await Items.at(contractAddress);
    console.log("itemsContractInstance", networkId, itemsContractInstance.address)

    // Read the item list from itemList.json
    const itemList = JSON.parse(fs.readFileSync('./itemsList.json'));

    //console.log(itemList);

    // Loop through each item in the list and call the createItem function
    for (let i = 0; i < itemList.length; i++) {
      const { name, uintValues, boolValues } = itemList[i];

      // Call the createItem function and verify that it creates a new item
      const result = await itemsContractInstance.createItem(name, uintValues, boolValues, { from: accounts[0] });

      const iteatts = await itemsContractInstance.getItemAttributes.call(result.logs[0].args.itemId);
      console.log("iteatts ", iteatts)
    }
  });
});
