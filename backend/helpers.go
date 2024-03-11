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

    "math/big"
    "math"
    "encoding/hex"

    "crypto/ecdsa"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"

)

var PRIVATE_KEY *ecdsa.PrivateKey

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


var ADMIN_ADDRESS common.Address



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

func convertIntToString(id int64) string {
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

    // Load your private key
    PRIVATE_KEY, err = crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
    if err != nil {
        log.Fatalf("[loadEnv] Failed to load private key: %v", err)
        return
    }
    ADMIN_ADDRESS = crypto.PubkeyToAddress(PRIVATE_KEY.PublicKey)

    Environment = os.Getenv("ENVIRONMENT")

    
    GAS_PRICE = getInt64EnvProp("GAS_PRICE")
    
    
    
    
    
    
    
    

    fmt.Println("Environment:", Environment)
}

func ZeroInt() *big.Int {
    return big.NewInt(0)
}

// bigIntToString safely converts a *big.Int to a string.
func bigIntToString(num *big.Int) string {
    if num == nil {
        return "0"
    }
    return num.String()
}

func addZerosToAmount(amount *big.Int, decimals int) *big.Int {
    if decimals == 18 {
        return amount
    }
    zerosToAdd := 18 - decimals
    multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(zerosToAdd)), nil)
    return new(big.Int).Mul(amount, multiplier)
}

func removeZerosFromAmount(amount *big.Int, decimals int) *big.Int {
    if decimals == 18 {
        return amount
    }
    zerosToRemove := 18 - decimals
    divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(zerosToRemove)), nil)
    return new(big.Int).Div(amount, divisor)
}

func bigIntToBigFloat(i *big.Int) *big.Float {
    return new(big.Float).SetInt(i)
}

func bigFloatToString(f *big.Float) string {
    return f.Text('f', -1)
}

func multiplyBigIntByBigFloat(bigIntVal *big.Int, bigFloatVal *big.Float) *big.Int {
    // Convert big.Int to big.Float
    bigIntAsFloat := new(big.Float).SetInt(bigIntVal)

    // Multiply
    result := new(big.Float).Mul(bigIntAsFloat, bigFloatVal)

    // Convert the result back to big.Int
    resultBigInt := new(big.Int)
    result.Int(resultBigInt) // This will truncate the fractional part

    return resultBigInt
}

func stringToBigFloat(str string) (*big.Float, bool) {
    f := new(big.Float)
    _, ok := f.SetString(str)
    return f, ok
}

// min returns the minimum of a and b
func minBigInt(a, b *big.Int) *big.Int {
    // Create a new big.Int to store the result
    result := new(big.Int)

    // Compare a and b, and set result to the smaller one
    if a.Cmp(b) <= 0 {
        result.Set(a)
    } else {
        result.Set(b)
    }

    return result
}

func maxBigInt(a, b *big.Int) *big.Int {
    // Create a new big.Int to store the result
    result := new(big.Int)

    // Compare a and b, and set result to the smaller one
    if a.Cmp(b) >= 0 {
        result.Set(a)
    } else {
        result.Set(b)
    }

    return result
}

func stringToBigInt(str string) (*big.Int) {
    bigInt := new(big.Int)
    bigInt.SetString(str, 10) // 10 is the base for decimal
    return bigInt
}

// float64ToBigFloat converts a float64 to a big.Float.
func float64ToBigFloat(value float64) *big.Float {
    bigFloat := new(big.Float)
    bigFloat.SetFloat64(value)
    return bigFloat
}

// float64PriceToBigIntE8 converts a float64 price to a *big.Int and multiplies it by 1e8.
func float64PriceToBigIntE8(price float64) *big.Int {
    // Multiply the price by 1e8 using float64 arithmetic
    priceScaled := price * 1e8

    // Convert the scaled price to *big.Int
    result := new(big.Int)
    result.SetInt64(int64(priceScaled))

    return result
}

func bigFloatToBigInt(value *big.Float) (*big.Int, bool) {
    if value == nil {
        return nil, false
    }

    // Create a big.Int to hold the converted value
    bigInt := new(big.Int)

    // Convert and round to the nearest integer
    bigInt, accuracy := value.Int(bigInt)

    // Check if the conversion was exact
    if accuracy != big.Exact {
        //log.Println("Warning: Conversion from big.Float to big.Int was not exact.")
    }

    return bigInt, true
}


func formatAmount(amount *big.Int, amountDec int) string {
    // Divide amount by 1e8
    divisor := big.NewInt(1e8)
    amountFloat := new(big.Float).SetInt(amount)
    divisorFloat := new(big.Float).SetInt(divisor)
    amountFloat.Quo(amountFloat, divisorFloat)

    // Format and return the string
    return fmt.Sprintf("%.*f", amountDec, amountFloat)
}

func formatAmountE18(amount *big.Int, amountDec int) string {
    // Divide amount by 1e18
    divisor := big.NewInt(1e18)
    amountFloat := new(big.Float).SetInt(amount)
    divisorFloat := new(big.Float).SetInt(divisor)
    amountFloat.Quo(amountFloat, divisorFloat)

    // Format and return the string
    return fmt.Sprintf("%.*f", amountDec, amountFloat)
}

func formatPrice(price *big.Int, decimals int) string {
    // Divide price by 1e8
    divisor := big.NewInt(1e8)
    priceFloat := new(big.Float).SetInt(price)
    divisorFloat := new(big.Float).SetInt(divisor)
    priceFloat.Quo(priceFloat, divisorFloat)

    // Format and return the string
    return fmt.Sprintf("%.*f", decimals, priceFloat)
}

func formatPriceToBigFloat(price *big.Int, decimals int) *big.Float {
    // Divide price by 1e8
    divisor := big.NewInt(1e8)
    priceFloat := new(big.Float).SetInt(price)
    divisorFloat := new(big.Float).SetInt(divisor)
    priceFloat.Quo(priceFloat, divisorFloat)

    // Format and return the string
    return priceFloat
}

func formatPriceToFloat64(price *big.Int, decimals int) float64 {
    // Divide price by 1e8
    divisor := big.NewInt(1e8)
    priceFloat := new(big.Float).SetInt(price)
    divisorFloat := new(big.Float).SetInt(divisor)
    priceFloat.Quo(priceFloat, divisorFloat)

    // Convert to float64
    formattedPrice, _ := priceFloat.Float64()

    return formattedPrice
}

func formatInt(value *big.Int) string {
    // Define the divisor as 1e8
    divisor := big.NewInt(1e8)

    // Convert the big.Int to big.Float for division
    valueFloat := new(big.Float).SetInt(value)
    divisorFloat := new(big.Float).SetInt(divisor)

    // Perform the division
    resultFloat := new(big.Float).Quo(valueFloat, divisorFloat)

    // Convert the result to a string rounded to the nearest whole number
    resultString := resultFloat.Text('f', 0) // 'f' format with 0 decimal places

    return resultString
}

func getFloatEnvProp(prop string) float64 {
    propStr := os.Getenv(prop)

    propFloat, err := strconv.ParseFloat(propStr, 64)
    if err != nil {
        log.Printf("Error parsing BASE_FUNDING_RATE: %v", err)
        return 0.0  // Return a default value or handle the error as appropriate
    }

    return propFloat
}

func getInt64EnvProp(prop string) int64 {
    propStr := os.Getenv(prop)

    propInt, err := strconv.ParseInt(propStr, 10, 64)
    if err != nil {
        log.Printf("[getInt64EnvProp] Error parsing %s: %v", prop, err)
        return 0  // Return a default value or handle the error as appropriate
    }

    return propInt
}

func baseFeeToBigInt(value float64) *big.Int {
    // Scale the float64 value
    scaleFactor := math.Pow(10, 16)
    scaledValue := value * scaleFactor

    // Convert the scaled value to *big.Int
    bigIntValue := new(big.Int)
    bigIntValue.SetUint64(uint64(scaledValue))

    return bigIntValue
}

func boolToBigInt(b bool) *big.Int {
    if b {
        return big.NewInt(1)
    }
    return big.NewInt(0)
}

func to32ByteSlice(str string) [32]byte {
    var arr [32]byte
    copy(arr[:], []byte(str))
    return arr
}

func to32ByteArray(b common.Hash) [32]byte {
    var arr [32]byte
    copy(arr[:], b.Bytes())
    return arr
}


// Define the maxBigFloat function used in the code above
func maxBigFloat(a, b *big.Float) *big.Float {
    if a.Cmp(b) > 0 {
        return a
    }
    return b
}


func stringToFloat64(s string) float64 {
    f, err := strconv.ParseFloat(s, 64)
    if err != nil {
        return 0.0 // Handle the error according to your use case
    }
    return f
}


func multiplyBigIntByFloat64(bigIntValue *big.Int, floatValue float64) *big.Int {
    // Convert bigIntValue to big.Float
    bigFloatValue := new(big.Float).SetInt(bigIntValue)

    // Convert floatValue to big.Float
    floatMultiplier := new(big.Float).SetFloat64(floatValue)

    // Multiply big.Float values
    result := new(big.Float).Mul(bigFloatValue, floatMultiplier)

    // Convert result back to big.Int
    resultInt := new(big.Int)
    result.Int(resultInt) // Note: this will truncate the decimal part

    return resultInt
}


func absBigInt(x *big.Int) *big.Int {
    if x.Sign() < 0 {
        return new(big.Int).Neg(x)
    }
    return new(big.Int).Set(x)
}

func hexStringToByteArray(hexStr string) ([32]byte, error) {
    var byteArray [32]byte

    // Remove the '0x' prefix if present
    if len(hexStr) > 1 && hexStr[:2] == "0x" {
        hexStr = hexStr[2:]
    }

    // The length of the hex string should be 64 characters (32 bytes)
    if len(hexStr) != 64 {
        return byteArray, fmt.Errorf("hex string length should be 64 characters (32 bytes), got %v characters", len(hexStr))
    }

    // Decode the string into a byte slice
    bytes, err := hex.DecodeString(hexStr)
    if err != nil {
        return byteArray, err
    }

    // Copy the bytes into the array
    copy(byteArray[:], bytes)

    return byteArray, nil
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