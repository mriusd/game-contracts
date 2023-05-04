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


func randomValueWithinRange(value int64, percentage float64) int64 {
	rand.Seed(time.Now().UnixNano())
	min := float64(value) * (1.0 - percentage)
	max := float64(value) * (1.0 + percentage)
	return int64(min + rand.Float64()*(max-min))
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

func convertIdToString(id int64) string {
    return strconv.Itoa(int(id))
}



func loadEnv() {
    envFilePath := "../.env"
    err := godotenv.Load(envFilePath)
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    FighterAttributesContract = os.Getenv("FIGHTER_ATTRIBUTES_CONTRACT")
    BattleContract = os.Getenv("BATTLE_CONTRACT")
    ItemsContract = os.Getenv("ITEMS_CONTRACT")
    MoneyContract = os.Getenv("MONEY_CONTRACT")
    BackpackContract = os.Getenv("BACKPACK_CONTRACT")
    PrivateKey = os.Getenv("PRIVATE_KEY")

    fmt.Println("FighterAttributesContract:", FighterAttributesContract)
    fmt.Println("BattleContract:", BattleContract)
    fmt.Println("ItemsContract:", ItemsContract)
    fmt.Println("MoneyContract:", MoneyContract)
    fmt.Println("BackpackContract:", BackpackContract)
}

func min(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}

func max(a, b int64) int64 {
    if a > b {
        return a
    }
    return b
}