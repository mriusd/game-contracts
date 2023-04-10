const generateRandomItem = async (contractInstance) => {
  const accounts = await web3.eth.getAccounts();
  const sender = accounts[0];

  const name = `Item_${Math.random().toString(36).substring(7)}`;
  const uintValues = Array.from({ length: 52 }, () => Math.floor(Math.random() * 100));
  const boolValues = Array.from({ length: 9 }, () => Math.random() < 0.5);

  try {
    const result = await contractInstance.methods.createItem(name, uintValues, boolValues).send({ from: sender });
    console.log(`Item created successfully: `, result);
  } catch (error) {
    console.error(`Item creation failed: `, error);
  }
};

(async () => {
    console.log("Deploying...")
    const metadata = JSON.parse(await remix.call('fileManager', 'getFile', 'browser/artifacts/Items.json'))
    let contract = new web3.eth.Contract(metadata.abi)
    console.log("Deploying 3...")
    contract = contract.deploy({
      data: metadata.data.bytecode.object,
      arguments: null
    })
    console.log("Deploying 2...")
    newContractInstance = await contract.send({
      from: accounts[0],
      gas: 1500000,
      gasPrice: '30000000000'
    })
    console.log("newContractInstance", newContractInstance)

    for (let i = 0; i < 2; i++) {
        await generateRandomItem(newContractInstance);
    }
})();
