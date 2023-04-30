package main

import (
    "context"
    "log"
    "math/big"
    "net/http"
)

var BlockedDamageReduction = 0.1
var HeadAtackDamageIncrease = 1.3
var BodyAtackDamageIncrease = 1.0
var LegsAtackDamageIncrease = 0.7

var FighterAttributesContract string
var BattleContract string
var ItemsContract string
var MoneyContract string
var PrivateKey string


var RpcClinetAddress            = "http://127.0.0.1:7545"
var RPCNetworkID                = big.NewInt(1337)

var HealthRegenerationDivider = 8;
var ManaRegenerationDivider = 8;
var AgilityPerDefence = 4;
var StrengthPerDamage = 8;
var EnergyPerDamage = 8;
var MaxExperience = 291342500;
var ExperienceDivider = 5;


func main() {
	loadEnv()   	
    lastBlockNumber()   
    loadNPCs()   
    http.HandleFunc("/ws", handleWebSocket)

    // Start the server
    log.Fatal(http.ListenAndServe(":8080", nil))    

    defer client.Disconnect(context.TODO())
}


