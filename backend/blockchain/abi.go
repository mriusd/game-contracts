// abi.go

package blockchain

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
    "io/ioutil"
    "fmt"
    "strings"
)

func LoadABI(contract string) (abi.ABI) {
    // Read the contract ABI file
    raw, err := ioutil.ReadFile("../build/contracts/" + contract + ".json")
    if err != nil {
        panic(fmt.Sprintf("Error reading ABI file: %v", err))
    }

    // Unmarshal the ABI JSON into the contractABI object
    var contractABIContent struct {
        ABI json.RawMessage `json:"abi"`
    }

    err = json.Unmarshal(raw, &contractABIContent)
    if err != nil {
        panic(fmt.Sprintf("Error unmarshalling ABI JSON: %v", err))
    }

    // Use the abi.JSON function to parse the ABI directly
    parsedABI, err := abi.JSON(strings.NewReader(string(contractABIContent.ABI)))
    if err != nil {
        panic(fmt.Sprintf("Error parsing ABI JSON: %v", err))
    }

    return parsedABI
}