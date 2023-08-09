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

var RpcClinetAddress            = "http://127.0.0.1:7545"
var RPCNetworkID                = big.NewInt(1337)

var HealthRegenerationDivider = 8;
var ManaRegenerationDivider = 8;
var AgilityPerDefence = 4;
var StrengthPerDamage = 8;
var EnergyPerDamage = 8;
var MaxExperience = 291342500;
var ExperienceDivider = 5;


var GoldTokenId int64 = 2;
var GoldItemId int64 = 1;

var Environment = "demo"


func main() {
	loadEnv()   	
    lastBlockNumber()  
    loadItems()
    loadMaps() 
    loadNPCs()

    http.HandleFunc("/ws", handleWebSocket)


    if Environment == "demo" {
        // Start the server
        log.Fatal(http.ListenAndServe(":8080", nil))   
    } else {
        certPath := "/etc/letsencrypt/live/mriusd.com/fullchain.pem"
        keyPath := "/etc/letsencrypt/live/mriusd.com/privkey.pem"

        // Serve over HTTPS
        log.Fatal(http.ListenAndServeTLS(":443", certPath, keyPath, nil))        
    }
    
    defer client.Disconnect(context.TODO())
    
}


