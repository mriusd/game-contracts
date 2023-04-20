const fs = require('fs');
const FighterAttributes = artifacts.require('FighterAttributes'); // Import the ABI and bytecode for your contract
const Battle = artifacts.require('Battle'); // Import the ABI and bytecode for your contract
const Items = artifacts.require('Items'); // Import the ABI and bytecode for your contract


contract('FighterAttributes', (accounts) => {
  it.skip('create first fighter', async () => {
    const fighterAttributesInstance = await FighterAttributes.deployed();

    try {
      var result = await fighterAttributesInstance.createFighter("StonedApe", 1, { gas: 3000000, from: accounts[0] });
      console.log("Fighter Created: ", result.logs[0].args.tokenId.toString())
    } catch (error) {
      console.log("Error: ", error);
    }
  });

  it('print initial npcs', async () => {
    const fighterAttributesInstance = await FighterAttributes.deployed();
    // Read the npc list from npcList.json
    const npcList = JSON.parse(fs.readFileSync('./npcList.json'));

    for (let i = 0; i < npcList.length; i++) {   
      var npc = npcList[i];
      //console.log("Npc", npc);
      //var result = await fighterAttributesInstance.createNPC(npc.name, npc.strength, npc.agility, npc.energy, npc.vitality, npc.attackSpeed, npc.level, npc.dropRarityLevel, { gas: 3000000, from: accounts[0] });
      //console.log("NPC Created: ", result.logs[0].args.tokenId.toString())     

      var npcAtts = await fighterAttributesInstance.getTokenAttributes.call(i+2);
      console.log("NPC: ", npcAtts)   

    }
  });  
  it.skip('record npc kill', async () => {
    const fighterAttributesInstance = await FighterAttributes.deployed();
    const battleInstance = await Battle.deployed();

    var fighter = await fighterAttributesInstance.getTokenAttributes.call(1);
    console.log("Fighter Before: ", fighter)

    var result = await battleInstance.recordKill(2, [[1, 100, 1]], 1, { gas: 3000000, from: accounts[0] });

    var fighter = await fighterAttributesInstance.getTokenAttributes.call(1);
    console.log("Fighter After: ", fighter)
  });
});
