package main

import (
	"log"
	"github.com/gorilla/websocket"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum"
	"context"
	"math/big"
    "encoding/json"
)

var DailyCreditsFaucet = int64(150);

func FaucetCredits(conn *websocket.Conn) {
    ownerAddress, err := getOwnerAddressByConn(conn);

    if err != nil {
        log.Printf("[FaucetCredits] err=%v", err)
        return
    }    

    log.Printf("[FaucetCredits] ownerAddress=%v", ownerAddress);

    // faucetCredits(address playerAddress, uint256 amount) e
    // Load contract ABI from file
    amount := big.NewInt(DailyCreditsFaucet)
    contractABI := loadABI("CreditsHelper");

    data, err := contractABI.Pack("faucetCredits", ownerAddress, amount)
    if err != nil {
        log.Printf("[FaucetCredits] Failed to encode function arguments: %v", err)
    }

    sendBlockchainTransaction(
        nil, 
        "CreditsHelper", 
        CreditsHelperContract, 
        data, 
        "CreditsHelper",
        CreditsHelperContract,
        "Faucet", 
        Coordinate{X: 0, Y: 0}, 
        common.Hash{},
        conn,
    )
}

func getUserCredits(conn *websocket.Conn) int64 {

	log.Printf("[getUserCredits]")

    ownerAddress, err := getOwnerAddressByConn(conn);

    if err != nil {
        log.Printf("[getUserCredits] err=%v", err)
        return 0
    }    
    
    // Connect to the Ethereum network using an Ethereum client
    rpcClient := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(CreditsHelperContract)
    contractABI := loadABI("CreditsHelper")

    //log.Printf("contractABI: ", contractABI);
    callData, err := contractABI.Pack("balanceOf", ownerAddress)
    if err != nil {
        log.Printf("[getUserCredits] Failed to pack call data: %v", err)
    }

    // log.Printf("callData: %v ", callData);
    // log.Printf("fighterID: %v ", fighterID);

    // Call the contract using the Ethereum client
    result, err := rpcClient.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
        Gas: 3000000,
    }, nil)
    if err != nil {
        if err.Error()[:36] == "VM Exception while processing transaction" {
            reason, err := abi.UnpackRevert(result)
            if err != nil {
                log.Printf("[getUserCredits] Failed to decode revert reason: %v", err)
            }
            log.Printf("[getUserCredits] Revert reason: %v", reason)
        } else {
            log.Printf("[getUserCredits] Failed to call %v", err)
        }
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var creditsResp  []interface{};

    //log.Printf("result: %v ", result);



    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    creditsResp, err = contractABI.Unpack("balanceOf", result)
    if err != nil {
        log.Printf("[getUserCredits] Failed to unpack error: %v", err)
    }

    jsonatts, err := json.Marshal(creditsResp[0])

    var credits big.Int
    json.Unmarshal(jsonatts, &credits)
    if err != nil {
        log.Printf("[getUserCredits] Failed to call contract: %v", err)
    }
   	log.Printf("[getUserCredits] credits: %v", credits.Int64())
    return credits.Int64()
}
