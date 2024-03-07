package main

import (
    "context"
    "log"
    "math/big"
    "net/http"
    "os"
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

var GAS_PRICE int64


func main() {
    // Create a log file
    logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
    if err != nil {
        log.Fatalf("[main] error opening log file: %v", err)
    }
    defer logFile.Close()

    // Set the output of the standard logger to the file
    log.SetOutput(logFile)

	loadEnv()   	
    lastBlockNumber()  
    loadItems()
    loadShopPriceList()
    loadMaps() 
    loadNPCs()

    http.HandleFunc("/ws", handleWebSocket)


    if Environment == "demo" {
        // Start the server
        log.Fatal(http.ListenAndServe(":8070", nil))   
    } else {
        certPath := "/etc/letsencrypt/live/mriusd.com/fullchain.pem"
        keyPath := "/etc/letsencrypt/live/mriusd.com/privkey.pem"

        // Serve over HTTPS
        log.Fatal(http.ListenAndServeTLS(":443", certPath, keyPath, nil))        
    }
    
    defer client.Disconnect(context.TODO())
    
}


