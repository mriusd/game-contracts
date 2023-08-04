## TRUFFLE CONSOLE
var contr = await DropHelper.deployed()
var ev = await contr.getPastEvents()
console.log(ev)

