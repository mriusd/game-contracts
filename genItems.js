async function loadItems() {
  const response = await fetch('itemsList.json');
  const items = await response.json();
  return items;
}

(async () => {
  console.log("Connecting to the contract...");
  const metadataJSON = await remix.call('fileManager', 'getFile', 'browser/artifacts/Items.json');
  console.log("Metadata retrieved from file");

  const metadata = JSON.parse(metadataJSON);
  console.log("Metadata parsed into json");

  contract = new web3.eth.Contract(metadata.abi, contractAddress);
  console.log("Connected to the contract: ", contract);
  console.log("Contract ABI: ", JSON.stringify(metadata.abi));
  console.log("Contract Address: ", contractAddress);


  try {
    for (let i = 0; i < 2; i++) {
      await generateRandomItem(contract);
    }
  } catch (error) {
    console.error("Error during sending:", error);
  }
})();
