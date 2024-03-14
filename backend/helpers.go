package main

import (
	"math/rand"
	"time"
	"encoding/json"
	"github.com/joho/godotenv"
    "strconv"
    "log"
    "os"
    "fmt"
)

var PrivateKey string

var FightersContract string
var FightersHelperContract string

var BattleContract string
var BattleHelperContract string

var ItemsContract string
var ItemsHelperContract string

var MoneyContract string
var MoneyHelperContract string

var BackpackContract string
var BackpackHelperContract string

var TradeContract string
var TradeHelperContract string

var DropContract string
var DropHelperContract string

var CreditsContract string
var CreditsHelperContract string




func randomValueWithinRange(value int, percentage float64) int {
	rand.Seed(time.Now().UnixNano())
	min := float64(value) * (1.0 - percentage)
	max := float64(value) * (1.0 + percentage)
	return int(min + rand.Float64()*(max-min))
}


func decodeJson(jsonStr []byte) map[string]interface{} {
    type Message struct {
        Data json.RawMessage `json:"data"`
    }
	var msg Message
	// Decode JSON into Message struct
    if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
        panic(err)
    }

    // Decode the raw data based on its structure
    var data map[string]interface{}
    if err := json.Unmarshal(msg.Data, &data); err != nil {
        panic(err)
    }

    return data;
}


func convertIdToString(id int) string {
    return strconv.Itoa(int(id))
}

func convertIntToString(id int) string {
    return strconv.Itoa(int(id))
}


func loadEnv() {
    envFilePath := "../.env"
    err := godotenv.Load(envFilePath)
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    FightersContract = os.Getenv("FIGHTERS_CONTRACT")
    FightersHelperContract = os.Getenv("FIGHTERS_HELPER_CONTRACT")
    fmt.Println("FightersContract:", FightersContract)
    fmt.Println("FightersHelperContract:", FightersHelperContract)

    BattleContract = os.Getenv("BATTLE_CONTRACT")
    BattleHelperContract = os.Getenv("BATTLE_HELPER_CONTRACT")
    fmt.Println("BattleContract:", BattleContract)
    fmt.Println("BattleHelperContract:", BattleHelperContract)

    ItemsContract = os.Getenv("ITEMS_CONTRACT")
    ItemsHelperContract = os.Getenv("ITEMS_HELPER_CONTRACT")
    fmt.Println("ItemsContract:", ItemsContract)
    fmt.Println("ItemsHelperContract:", ItemsHelperContract)

    MoneyContract = os.Getenv("MONEY_CONTRACT")
    MoneyHelperContract = os.Getenv("MONEY_HELPER_CONTRACT")
    fmt.Println("MoneyContract:", MoneyContract)
    fmt.Println("MoneyHelperContract:", MoneyHelperContract)

    BackpackContract = os.Getenv("BACKPACK_CONTRACT")
    BackpackHelperContract = os.Getenv("BACKPACK_HELPER_CONTRACT")
    fmt.Println("BackpackContract:", BackpackContract)
    fmt.Println("BackpackHelperContract:", BackpackHelperContract)

    TradeContract = os.Getenv("TRADE_CONTRACT")
    TradeHelperContract = os.Getenv("TRADE_HELPER_CONTRACT")
    fmt.Println("TradeContract:", TradeContract)
    fmt.Println("TradeHelperContract:", TradeHelperContract)

    DropContract = os.Getenv("DROP_CONTRACT")
    DropHelperContract = os.Getenv("DROP_HELPER_CONTRACT")
    fmt.Println("DropContract:", DropContract)
    fmt.Println("DropHelperContract:", DropHelperContract)
    
    CreditsContract = os.Getenv("CREDITS_CONTRACT")
    CreditsHelperContract = os.Getenv("CREDITS_HELPER_CONTRACT")
    fmt.Println("CreditsContract:", CreditsContract)
    fmt.Println("CreditsHelperContract:", CreditsHelperContract)

    PrivateKey = os.Getenv("PRIVATE_KEY")
    Environment = os.Getenv("ENVIRONMENT")

    
    GAS_PRICE = getIntEnvProp("GAS_PRICE")
    
    
    
    
    
    
    
    

    fmt.Println("Environment:", Environment)
}



func getIntEnvProp(prop string) int {
    propStr := os.Getenv(prop)

    propInt, err := strconv.Atoi(propStr)
    if err != nil {
        log.Printf("[getintEnvProp] Error parsing %s: %v", prop, err)
        return 0  // Return a default value or handle the error as appropriate
    }

    return propInt
}



func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}