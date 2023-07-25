## TRUFFLE CONSOLE
var contr = await FighterAttributes.deployed()
var ev = await contr.getPastEvents('FighterCreated')
console.log(ev)

