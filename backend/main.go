package main

import (
    "context"
    "log"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "time"
    "fmt"
    //"strings"
    "math/big"
    

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/bson"

    "github.com/gorilla/websocket"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum"
	//"github.com/onrik/ethrpc"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//"github.com/ethereum/go-ethereum/event"

	"math/rand"
	"os"
	"github.com/joho/godotenv"
)

var client *mongo.Client = ConnectToDB()

var BlockedDamageReduction = 0.1
var HeadAtackDamageIncrease = 1.3
var BodyAtackDamageIncrease = 1.0
var LegsAtackDamageIncrease = 0.7

var FighterAttributesContract string
var BattleContract string
var ItemsContract string
var PrivateKey string


var RpcClinetAddress 			= "http://127.0.0.1:7545"
var RPCNetworkID				= big.NewInt(1337)


type Battle struct {
	ID        				primitive.ObjectID `bson:"_id,omitempty"`
    Opponent1   	int64	`bson:"opponent1"`
    Opponent2   	int64	`bson:"opponent2"`
    Health1     	int64	`bson:"health1"`
    Health2     	int64	`bson:"health2"`
    Mana1       	int64	`bson:"mana1"`
    Mana2       	int64	`bson:"mana2"`
    Damage1     	int64	`bson:"damage1"`
    Damage2     	int64	`bson:"damage2"`
    ManaUsed1   	int64	`bson:"manaused1"`
    ManaUsed2   	int64	`bson:"manaused2"`
    Winner      	int64	`bson:"winner"`
    Purse       	int64	`bson:"purse"`
    Exp1        	int64	`bson:"exp1"`
    Exp2        	int64	`bson:"exp2"`
    LastRound1 		int64	`bson:"lastround1"` 	
    LastRound2 		int64	`bson:"lastround2"`  
    Closed 			int64 	`bson:"closed"`  	
    MaxHP1 			int64 	`bson:"maxHp1"`  	
    MaxHP2 			int64 	`bson:"maxHp2"`  	
    MaxMana1 		int64 	`bson:"maxMana1"`  	
    MaxMana2		int64 	`bson:"maxMana2"`   
    IsClosed 		bool 	`bson:"isClosed"`  	

    TxHash 			string `bson:"txhash"`  	
}


type Player struct {
	ID        				primitive.ObjectID `bson:"_id,omitempty"`
	PlayerID				int64   	`bson:"playerid"`
    Strength   				int64 	`bson:"strength"`
    Agility   				int64 	`bson:"agility"`
    Energy     				int64 	`bson:"energy"`
    Vitality     			int64 	`bson:"vitality"`
    Experience       		int64 	`bson:"experience"`
    HaelthAfterLastDmg  	int64 	`bson:"haelthAfterLastDmg"`
    LastDmgBlockNumber  	int64 	`bson:"lastDmgBlockNumber"`
    Class     				int64 	`bson:"class"`
    HealthRegenerationBonus int64 	`bson:"healthRegenerationBonus"`
    ManaAfterLastUse   		int64 	`bson:"manaAfterLastUse"`
    LastManaUseBlock      	int64 	`bson:"lastManaUseBlock"`
    ManaRegenerationBonus   int64 	`bson:"manaRegenerationBonus"`
    BaseHP        			int64 	`bson:"baseHP"`
    BaseMana        		int64 	`bson:"baseMana"`
    HpPerVitalityPoint 		int64 	`bson:"hpPerVitalityPoint"`
    ManaPerEnergyPoint 		int64  	`bson:"manaPerEnergyPoint"`
    HpIncreasePerLevel 		int64  	`bson:"hpIncreasePerLevel"`
    ManaIncreasePerLevel 	int64  	`bson:"manaIncreasePerLevel"`
}

type FighterAttributes struct {
	TokenID   				*big.Int `json:"TokenID"`
    Strength                *big.Int `json:"Strength"`
	Agility                 *big.Int `json:"Agility"`
	Energy                  *big.Int `json:"Energy"`
	Vitality                *big.Int `json:"Vitality"`
	Experience              *big.Int `json:"Experience"`
	Class                   *big.Int `json:"Class"`
	HpPerVitalityPoint      *big.Int `json:"HpPerVitalityPoint"`
	ManaPerEnergyPoint      *big.Int `json:"ManaPerEnergyPoint"`
	HpIncreasePerLevel      *big.Int `json:"HpIncreasePerLevel"`
	ManaIncreasePerLevel   	*big.Int `json:"manaIncreasePerLevel"`
	StatPointsPerLevel   	*big.Int `json:"statPointsPerLevel"`
	AttackSpeed   			*big.Int `json:"attackSpeed"`
	AgilityPointsPerSpeed   *big.Int `json:"agilityPointsPerSpeed"`
	IsNpc   				*big.Int `json:"isNpc"`

	HelmSlot   				*big.Int `json:"helmSlot"`
	ArmourSlot   			*big.Int `json:"armourSlot"`
	PantsSlot   			*big.Int `json:"pantsSlot"`
	GlovesSlot   			*big.Int `json:"glovesSlot"`
	BootsSlot   			*big.Int `json:"bootsSlot"`
	LeftHandSlot   			*big.Int `json:"leftHandSlot"`
	RightHandSlot   		*big.Int `json:"rightHandSlot"`
	LeftRingSlot   			*big.Int `json:"leftRingSlot"`
	RightRingSlot   		*big.Int `json:"rightRingSlot"`
	PendSlot   				*big.Int `json:"pendSlot"`
	WingsSlot   			*big.Int `json:"wingsSlot"`
}

type FighterStats struct {
	TokenID           		*big.Int `json:"TokenID"`
	MaxHealth               *big.Int `json:"maxHealth"`
	MaxMana                 *big.Int `json:"maxMana"`
	Level              		*big.Int `json:"level"`
	Exp      				*big.Int `json:"exp"`
	TotalStatPoints         *big.Int `json:"totalStatPoints"`
	MaxStatPoints           *big.Int `json:"maxStatPoints"`
}

type ItemAttributes struct {
    Name                                  		string      `json:"name" bson:"name"`
    TokenId                               		int64    	`json:"tokenId" bson:"tokenId"`
    ItemLevel                            		int64       `json:"itemLevel" bson:"itemLevel"`
    MaxLevel                              		int64       `json:"maxLevel" bson:"maxLevel"`
    Durability                            		int64       `json:"durability" bson:"durability"`
    ClassRequired                         		int64       `json:"classRequired" bson:"classRequired"`
    StrengthRequired                      		int64       `json:"strengthRequired" bson:"strengthRequired"`
    AgilityRequired                       		int64       `json:"agilityRequired" bson:"agilityRequired"`
    EnergyRequired                        		int64       `json:"energyRequired" bson:"energyRequired"`
    VitalityRequired                      		int64       `json:"vitalityRequired" bson:"vitalityRequired"`
    ItemWidth                             		int64       `json:"itemWidth" bson:"itemWidth"`
    ItemHeight                            		int64       `json:"itemHeight" bson:"itemHeight"`
    AcceptableSlot1                       		int64       `json:"acceptableSlot1" bson:"acceptableSlot1"`
    AcceptableSlot2                       		int64       `json:"acceptableSlot2" bson:"acceptableSlot2"`
    PhysicalDamage                        		int64       `json:"physicalDamage" bson:"physicalDamage"`
    MagicDamage                           		int64       `json:"magicDamage" bson:"magicDamage"`
    Defense                               		int64       `json:"defense" bson:"defense"`
    AttackSpeed                           		int64       `json:"attackSpeed" bson:"attackSpeed"`
    DefenseSuccessRate                    		int64       `json:"defenseSuccessRate" bson:"defenseSuccessRate"`
    AdditionalDamage                      		int64       `json:"additionalDamage" bson:"additionalDamage"`
    AdditionalDefense                     		int64       `json:"additionalDefense" bson:"additionalDefense"`
    IncreasedExperienceGain               		int64       `json:"increasedExperienceGain" bson:"increasedExperienceGain"`
    DamageIncrease                        		int64       `json:"damageIncrease" bson:"damageIncrease"`
    DefenseSuccessRateIncrease            		int64       `json:"defenseSuccessRateIncrease" bson:"defenseSuccessRateIncrease"`
    LifeAfterMonsterIncrease              		int64       `json:"lifeAfterMonsterIncrease" bson:"lifeAfterMonsterIncrease"`
    ManaAfterMonsterIncrease              		int64       `json:"manaAfterMonsterIncrease" bson:"manaAfterMonsterIncrease"`
    ZenAfterMonsterIncrease               		int64       `json:"zenAfterMonsterIncrease" bson:"zenAfterMonsterIncrease"`
    DoubleDamageProbabilityIncrease      		int64       `json:"doubleDamageProbabilityIncrease" bson:"doubleDamageProbabilityIncrease"`
    ExcellentDamageProbabilityIncrease    		int64       `json:"excellentDamageProbabilityIncrease" bson:"excellentDamageProbabilityIncrease"`
    IgnoreOpponentsDefenseRateIncrease    		int64       `json:"ignoreOpponentsDefenseRateIncrease" bson:"ignoreOpponentsDefenseRateIncrease"`
    ReflectDamage                         		int64       `json:"reflectDamage" bson:"reflectDamage"`
    MaxLifeIncrease                       		int64       `json:"maxLifeIncrease" bson:"maxLifeIncrease"`
    MaxManaIncrease                       		int64       `json:"maxManaIncrease" bson:"maxManaIncrease"`
    ExcellentDamageRateIncrease           		int64       `json:"excellentDamageRateIncrease" bson:"excellentDamageRateIncrease"`
    DoubleDamageRateIncrease             		int64    	`json:"doubleDamageRateIncrease" bson:"doubleDamageRateIncrease"`
    IgnoreOpponentsDefenseSuccessRateIncrease 	int64    	`json:"ignoreOpponentsDefenseSuccessRateIncrease" bson:"ignoreOpponentsDefenseSuccessRateIncrease"`
    AttackDamageIncrease                 		int64    	`json:"attackDamageIncrease" bson:"attackDamageIncrease"`
    IsAncient    								int64    	`json:"isAncient" bson:"isAncient"`
    ReflectDamageRateIncrease            		int64    	`json:"reflectDamageRateIncrease" bson:"reflectDamageRateIncrease"`
    DecreaseDamageRateIncrease           		int64    	`json:"decreaseDamageRateIncrease" bson:"decreaseDamageRateIncrease"`
    HPRecoveryRateIncrease               		int64    	`json:"hpRecoveryRateIncrease" bson:"hpRecoveryRateIncrease"`
    MPRecoveryRateIncrease               		int64    	`json:"mpRecoveryRateIncrease" bson:"mpRecoveryRateIncrease"`
    HPIncrease                           		int64    	`json:"hpIncrease" bson:"hpIncrease"`
    MPIncrease                           		int64    	`json:"mpIncrease" bson:"mpIncrease"`
    IncreaseDefenseRate                  		int64    	`json:"increaseDefenseRate" bson:"increaseDefenseRate"`
    IncreaseStrength                     		int64    	`json:"increaseStrength" bson:"increaseStrength"`
    IncreaseAgility                      		int64    	`json:"increaseAgility" bson:"increaseAgility"`
    IncreaseEnergy                       		int64    	`json:"increaseEnergy" bson:"increaseEnergy"`
    IncreaseVitality                     		int64    	`json:"increaseVitality" bson:"increaseVitality"`
    AttackSpeedIncrease                  		int64    	`json:"attackSpeedIncrease" bson:"attackSpeedIncrease"`
    ItemRarityLevel                     		int64    	`json:"itemRarityLevel" bson:"itemRarityLevel"`
    ItemAttributesId                     		int64    	`json:"itemAttributesId" bson:"itemAttributesId"`

    Luck              							bool        `json:"luck" bson:"luck"`
    Skill             							bool        `json:"skill" bson:"skill"`
    IsBox             							bool        `json:"isBox" bson:"isBox"`
    IsWeapon          							bool        `json:"isWeapon" bson:"isWeapon"`
    IsArmour          							bool        `json:"isArmour" bson:"isArmour"`
    IsJewel           							bool        `json:"isJewel" bson:"isJewel"`
    IsMisc            							bool        `json:"isMisc" bson:"isMisc"`
    IsConsumable      							bool        `json:"isConsumable" bson:"isConsumable"`
    InShop            							bool        `json:"inShop" bson:"inShop"`
}



type NPC struct {
    ID               int64    `json:"id"`
    Name             string   `json:"name"`
    Level            int64      `json:"level"`
    Strength         int64      `json:"strength"`
    Agility          int64      `json:"agility"`
    Energy           int64      `json:"energy"`
    Vitality         int64      `json:"vitality"`
    AttackSpeed      int64      `json:"attackSpeed"`
    DropRarityLevel  int64      `json:"dropRarityLevel"`
    RespawnLocations [][]string `json:"respawnLocations"`
    CanFight         bool     `json:"canFight"`
    MaxHealth        int64      `json:"maxHealth"`
}


type Fighter struct {
	ID    					int64  			`json:"id"`
    MaxHealth     			int64 			`json:"maxHealth"`
    Name           			string 			`json:"name"`
    IsNpc         			bool 			`json:"isNpc"`
    IsDead         			bool 			`json:"isDead"`
    CanFight 				bool 			`json:"canFight"`
    LastDmgTimestamp 		int64  			`json:"lastDmgTimestamp"`
    HealthAfterLastDmg 		int64  			`json:"healthAfterLastDmg"`
    HpRegenerationRate 		float64 		`json:"hpRegenerationRate"`
    HpRegenerationBonus 	float64 		`json:"hpRegenerationBonus"`
    InBattle 				bool 			`json:"inBattle"`
    Conn 					*websocket.Conn
}

type RecordMoveMsg struct {
    PlayerID 	int64 	`json:"playerID"`
    OpponentID  int64 	`json:"opponentID`
    Skill		int64 	`json:"skill"`
}


var npcs map[string]map[int64][]NPC;
var Fighters = make(map[int64]Fighter)
var Battles = make(map[primitive.ObjectID]Battle)


var Population = make(map[string]map[int64][]Fighter)


var HealthRegenerationDivider = 8;
var ManaRegenerationDivider = 8;
var AgilityPerDefence = 4;
var StrengthPerDamage = 8;
var EnergyPerDamage = 8;
var MaxExperience = 291342500;
var ExperienceDivider = 5;

func findNpcById(id int64) *NPC {
	for _, zone := range npcs {
		for _, npcList := range zone {
			for _, npc := range npcList {
				if npc.ID == id {
					return &npc
				}
			}
		}
	}
	return nil
}

func initMap() {
	Population["lorencia"] = make(map[int64][]Fighter);

	log.Printf("Map Initialized")
}

func authFighter(conn *websocket.Conn, playerId int64) {
	fighter := Fighters[playerId];


	// Check if the fighter is not already in the map
	found := false
	var f Fighter;
	for i :=0; i<len(Population["lorencia"][0]); i++ {
		f = Population["lorencia"][0][i]
		if f.ID == fighter.ID {
			found = true
			fighter.Conn = conn
			Population["lorencia"][0][i] = fighter;
			break
		}
	}

	if !found {
		Population["lorencia"][0] = append(Population["lorencia"][0], fighter)
	} 

	log.Printf("[authFighter] ", Population["lorencia"][0])
	Fighters[playerId] = fighter;
}

func initiateNpcRoutine(npcId int64) {
	attributes := getFighterAttributes(npcId)
	npc := findNpcById(npcId)
	Fighters[npcId] = Fighter{
		ID: npcId,
		MaxHealth: npc.MaxHealth, 
		Name: npc.Name,
	    IsNpc: true,
	    CanFight: npc.CanFight,
	    HpRegenerationRate: 0,
	    HpRegenerationBonus: 0,
	    LastDmgTimestamp: 0,
	    HealthAfterLastDmg: npc.MaxHealth,
	}

	fmt.Println("[initiateNpcRoutine] npcId=", npcId)

	

	speed := attributes.AttackSpeed;

	msPerHit := 60000/speed.Int64();
	delay := time.Duration(msPerHit) * time.Millisecond;
	var fighter Fighter;

	for i := 0; i<1000; i++  {
		if Fighters[npcId].IsDead {
			fighter = Fighters[npcId]
			now := time.Now().UnixNano() / 1e6
			elapsedTimeMs := now - fighter.LastDmgTimestamp

			if elapsedTimeMs >= 5000 {
				fmt.Println("At least 5 seconds have passed since TimeOfDeath.")
				fighter.IsDead = false;
				fighter.HealthAfterLastDmg = fighter.MaxHealth;
			} 

			Fighters[npcId] = fighter;

		}

		if len(Population["lorencia"][0]) > 0 {
			rand.Seed(time.Now().UnixNano()) // Seed the random number generator
			randomIndex := rand.Intn(len(Population["lorencia"][0]))
			randomFighter := Population["lorencia"][0][randomIndex]

			data := RecordMoveMsg{
				OpponentID: randomFighter.ID,
				PlayerID: npcId,
				Skill: 0,
			}

			rawMessage, err := json.Marshal(data);
			if err != nil {
				fmt.Println("[initiateNpcRoutine] Error marshaling data:", err)
				return
			}

			recordMove(randomFighter.Conn, rawMessage)
		}
		
		time.Sleep(delay);
	}
}

func loadNPCs() {
    // Open the JSON file
    file, err := os.Open("../npcList.json")
    if err != nil {
        log.Printf("[loadNPCs] error= ", err)
    }
    defer file.Close()

    // Read the JSON data
    data, err := ioutil.ReadAll(file)
    if err != nil {
        log.Printf("[loadNPCs] error= ", err)
    }

    // Unmarshal the JSON data into a slice of NPCs
    var npcs []NPC
    err = json.Unmarshal(data, &npcs)
    if err != nil {
        log.Printf("[loadNPCs] error= ", err)
    }

    // Set default values and initiate NPC routines
    for i, npc := range npcs {
        npcs[i].CanFight = true
        npcs[i].MaxHealth = npc.Vitality
        go initiateNpcRoutine(npc.ID)
    }

    log.Printf("NPCs Loaded")
}

func getNPCs(town string, location int64) []NPC {
    return npcs[town][location]
}

func recordBattleOnChain(player, opponent int64) (string) {
	if !Fighters[opponent].IsNpc { return "" }
	log.Printf("[recordBattleOnChain] Recording")

	// Connect to the Ethereum network
	client := getRpcClient();

	// Load your private key
	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Load contract ABI from file
	contractABI := loadABI("Battle");

	// Set contract address
	contractAddress := common.HexToAddress(BattleContract)

	// Prepare function arguments
	vals1 := [7]*big.Int{}
	vals1[0] = big.NewInt(player)
	vals1[1] = big.NewInt(0)
	vals1[2] = big.NewInt(0)
	vals1[3] = big.NewInt(Fighters[opponent].MaxHealth)
	vals1[4] = big.NewInt(0)
	vals1[5] = big.NewInt(0)
	vals1[6] = big.NewInt(0)
	

	// ...
	vals2 := [7]*big.Int{}
	vals2[0] = big.NewInt(opponent)
	vals2[1] = big.NewInt(0)
	vals2[2] = big.NewInt(0)
	vals2[3] = big.NewInt(0)
	vals2[4] = big.NewInt(Fighters[opponent].MaxHealth)
	vals2[5] = big.NewInt(0)
	vals2[6] = big.NewInt(0)

	type Damage struct {
	    FighterId        *big.Int
	    Damage           *big.Int
	    MedianPartyLevel *big.Int
	}

	killedFighter := big.NewInt(opponent)
	damageDealt := []Damage{
	    {
	        FighterId:        big.NewInt(1),
	        Damage:           big.NewInt(100),
	        MedianPartyLevel: big.NewInt(0),
	    },
	}
	battleNonce := big.NewInt(time.Now().UnixNano() / int64(time.Millisecond))

	damageDealtInterface := make([]interface{}, len(damageDealt))
	for i, d := range damageDealt {
	    damageDealtInterface[i] = d
	}

	log.Printf("[recordBattleOnChain] battleNonce=", battleNonce)

	// Prepare transaction options
	nonce, err := client.NonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey), nil)
	if err != nil {
		log.Printf("Failed to retrieve nonce: %v", err)
	}
	gasLimit := uint64(500000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("Failed to retrieve gas price: %v", err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	// Encode function arguments
	data, err := contractABI.Pack("recordKill", killedFighter, damageDealtInterface, battleNonce)
	if err != nil {
		log.Printf("Failed to encode function arguments: %v", err)
	}

	// Create transaction and sign it
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(RPCNetworkID), privateKey)
	if err != nil {
		log.Printf("Failed to sign transaction: %v", err)
	}

	// Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Printf("Failed to send transaction: %v", err)
	}

	fmt.Println("[recordBattleOnChain] Transaction hash:", signedTx.Hash().Hex())

	//updateString("battles", battle.ID, "txhash", signedTx.Hash().Hex())

	return signedTx.Hash().Hex();
}

func getRpcClient() *ethclient.Client {
	// Connect to the Ethereum network using an Ethereum client
	client, err := ethclient.Dial(RpcClinetAddress)
	if err != nil {
		log.Fatalf("[getRpcClient] Failed to connect to Ethereum network: %v", err)
	}

	return client
}


func getFighterAttributes(fighterID int64) FighterAttributes {

	// Connect to the Ethereum network using an Ethereum client
    rpcClient := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(FighterAttributesContract)
    contractABI := loadABI("FighterAttributes")

    //log.Printf("contractABI: ", contractABI.Methods[method.Name]);

    // Prepare the call to the getTokenAttributes function
    tokenID := big.NewInt(fighterID)
    callData, err := contractABI.Pack("getTokenAttributes", tokenID)
    if err != nil {
        log.Fatalf("[getFighterAttributes] Failed to pack call data: %v", err)
    }

    log.Printf("callData: %v ", callData);

    // chainId, err := client.ChainID(context.Background());
    //  if err != nil {
    //     log.Fatalf("[getFighterAttributes] Failed to get chain id: %v", err)
    // }
    // log.Printf("ChainId:  ", chainId);

    // Call the contract using the Ethereum client
    result, err := rpcClient.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
    }, nil)
    if err != nil {
        if err.Error()[:36] == "VM Exception while processing transaction" {
            reason, err := abi.UnpackRevert(result)
            if err != nil {
                log.Fatalf("[getFighterAttributes] Failed to decode revert reason: %v", err)
            }
            log.Fatalf("[getFighterAttributes] Revert reason: %v", reason)
        } else {
            log.Fatalf("[getFighterAttributes] Failed to call contract: %v", fighterID, err)
        }
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var attributes []interface{};


    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    attributes, err = contractABI.Unpack("getTokenAttributes", result)
    if err != nil {
        log.Printf("[getFighterAttributes] Failed to unpack error: %v", err)
    }

    jsonatts, err := json.Marshal(attributes[0])

    var fighter FighterAttributes
    json.Unmarshal(jsonatts, &fighter)
    if err != nil {
        log.Fatalf("[getFighterAttributes] Failed to call contract: %v", err)
    }
   	log.Printf("[getFighterAttributes] fighter: %v", fighter)
   	return fighter;
}

func getItemAttributes(itemId int64) ItemAttributes {
	//log.Printf("[getItemAttributes] itemId: %v", itemId)
	if itemId == 0 {
		return ItemAttributes{};
	}

	// Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(ItemsContract)
    contractABI := loadABI("Items")

    // Prepare the call to the getTokenAttributes function
    tokenID := big.NewInt(itemId)
    callData, err := contractABI.Pack("getTokenAttributes", tokenID)
    if err != nil {
        log.Fatalf("[getItemAttributes] Failed to pack call data: %v", err)
    }

    // Call the contract using the Ethereum client
    result, err := client.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
    }, nil)
    if err != nil {
        log.Fatalf("[getItemAttributes] Failed to call contract: %v", err)
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var attributes []interface{};


    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    attributes, err = contractABI.Unpack("getTokenAttributes", result)
    if err != nil {
        log.Printf("[getItemAttributes] Failed to unpack error: %v", err)
    }

    jsonatts, err := json.Marshal(attributes[0])

    var item ItemAttributes
    json.Unmarshal(jsonatts, &item)
    if err != nil {
        log.Fatalf("[getItemAttributes] Failed to call contract: %v", err)
    }

    //log.Printf("[getItemAttributes] item: %v", item)
   	
   	return item;
}

func recordItemToDB(item ItemAttributes) {
	coll := client.Database("game").Collection("items")

	// create a filter for the query
	filter := bson.M{"tokenId": item.TokenId}

	// create an update document containing the values to update
	update := bson.M{
	    "$set": item,
	}

	// create options for the query
	opts := options.Update().SetUpsert(true)

	// execute the update query
	result, err := coll.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
	    log.Printf("[recordItemToDB] error: %v\n", err, result)
	}
}

func getFighterItems(conn *websocket.Conn, data json.RawMessage, w http.ResponseWriter, r *http.Request)  {

	type ItemReqData struct {
		UserAddress string 
		FighterId int64 
	}

	var reqData ItemReqData
	err := json.Unmarshal(data, &reqData)
    if err != nil {
        log.Printf("websocket unmarshal error: %v", err)
        return
    }

    

	// Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(ItemsContract)
    contractABI := loadABI("Items")

    // Prepare the call to the getTokenAttributes function
    tokenID := big.NewInt(reqData.FighterId)
    callData, err := contractABI.Pack("getFighterItems", common.HexToAddress(reqData.UserAddress), tokenID)
    if err != nil {
        log.Fatalf("[getFighterItems] Failed to pack call data: %v", err)
    }


    //log.Printf("[getFighterItems] FreqData: %v", reqData)

    // Call the contract using the Ethereum client
    result, err := client.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
    }, nil)
    if err != nil {
        log.Fatalf("[getFighterItems] Failed to call contract: %v", err)
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var attributes []interface{};

	//log.Printf("[getFighterItems] Result", result)
    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    attributes, err = contractABI.Unpack("getFighterItems", result)
    if err != nil {
        log.Printf("[getFighterItems] Failed to unpack error: %v", err)
    }

    var items []ItemAttributes

	//log.Printf("[getFighterItems] attributes: %v", attributes)

	for _, v := range attributes {
		attrs, ok := v.([][2]*big.Int)
		if !ok {
	        // handle invalid attribute format
	        log.Printf("[getFighterItems] Error iterating attributes")
	    }

	    for _, attr := range attrs {

	    	itemId := attr[0]
	    	// lastUpdBlock := attr[1]

	    	// log.Printf("[getFighterItems] itemId=",itemId ," lastUpdBlock=", lastUpdBlock)

	    	// get item attributes
	    	itemAttributes := getItemAttributes(itemId.Int64());

	    	recordItemToDB(itemAttributes);

	    	items = append(items, itemAttributes);

	    }
	}

	var jsonatts []byte
	if len(items) == 0 {
		jsonatts = nil;
	} else {
		jsonatts, err = json.Marshal(items)
	}	

	stats := getFighterStats(reqData.FighterId);

    jsonstats, err := json.Marshal(stats)
    //log.Print("[getFighter] jsonstats: %s", stats)


    fighterAttributes := getFighterAttributes(reqData.FighterId);
    jsonfighteratts, err := json.Marshal(fighterAttributes)


    //log.Print("[getFighter] jsonatts: %s", fighter)


    _, ok := Fighters[reqData.FighterId]
	if !ok {
		Fighters[reqData.FighterId] = Fighter{
			ID: reqData.FighterId,
			MaxHealth: stats.MaxHealth.Int64(), 
			Name: "",
		    IsNpc: false,
		    CanFight: true,
		    HpRegenerationRate: getHealthRegenerationRate(fighterAttributes),
		    LastDmgTimestamp: 0,
		    HealthAfterLastDmg: 0,
		    Conn: conn,
		}
	}

	authFighter(conn, reqData.FighterId);

    log.Printf("[getFighterItems] fighterAttributes: %v", fighterAttributes)
    type Equipment struct {
		Helm 		ItemAttributes `json:"helm"`
		Armour 		ItemAttributes `json:"armour"`
		Pants 		ItemAttributes `json:"pants"`
		Gloves 		ItemAttributes `json:"gloves"`
		Boots 		ItemAttributes `json:"boots"`
		LeftHand 	ItemAttributes `json:"leftHand"`
		RightHand 	ItemAttributes `json:"rightHand"`
		LeftRing 	ItemAttributes `json:"leftRing"`
		RightRing 	ItemAttributes `json:"rightRing"`
		Pendant 	ItemAttributes `json:"pendant"`
		Wings 		ItemAttributes `json:"wings"`
	}

	equipment := Equipment{
		Helm: 		getItemAttributes(fighterAttributes.HelmSlot.Int64()),
		Armour: 	getItemAttributes(fighterAttributes.ArmourSlot.Int64()),
		Pants: 		getItemAttributes(fighterAttributes.PantsSlot.Int64()),
		Gloves: 	getItemAttributes(fighterAttributes.GlovesSlot.Int64()),
		Boots: 		getItemAttributes(fighterAttributes.BootsSlot.Int64()),
		LeftHand: 	getItemAttributes(fighterAttributes.LeftHandSlot.Int64()),
		RightHand: 	getItemAttributes(fighterAttributes.RightHandSlot.Int64()),
		LeftRing: 	getItemAttributes(fighterAttributes.LeftRingSlot.Int64()),
		RightRing: 	getItemAttributes(fighterAttributes.RightRingSlot.Int64()),
		Pendant: 	getItemAttributes(fighterAttributes.PendSlot.Int64()),
		Wings: 		getItemAttributes(fighterAttributes.WingsSlot.Int64()),
	}

	jsonequip, err := json.Marshal(equipment)
    //log.Print("[getFighter] jsonstats: %s", equipment)


    

    npcs := getNPCs("lorencia", 0);
    log.Print("[getFighterItems] npcs: ", npcs)
    var npcatts []Fighter
    var npcStats FighterStats
    for _, npc := range npcs {
    	_, ok = Fighters[npc.ID]
    	npcStats = getFighterStats(npc.ID);
    	if !ok {
    		Fighters[npc.ID] = Fighter{
    			ID: npc.ID,
    			MaxHealth: npcStats.MaxHealth.Int64(), 
    			Name: npc.Name,
			    IsNpc: true,
			    CanFight: npc.CanFight,
			    HpRegenerationRate: 0,
			    HpRegenerationBonus: 0,
			    LastDmgTimestamp: 0,
			    HealthAfterLastDmg: npcStats.MaxHealth.Int64(),
    		}
    	} 
	    npcatts = append(npcatts, Fighters[npc.ID]);
	}

    jsonnpcs, err := json.Marshal(npcatts)


    jsonfighter, err := json.Marshal(Fighters[reqData.FighterId])

    type jsonResponse struct {
		Action string `json:"action"`
		Items string `json:"items"`
		Attributes string `json:"attributes"`
		Equipment string `json:"equipment"`
		Stats string `json:"stats"`
		NPCs string `json:"npcs"`
		Fighter string `json:"fighter"`
	}

    jsonResp := jsonResponse{
    	Action: "fighter_items",
    	Items: string(jsonatts),
    	Attributes: string(jsonfighteratts),
    	Equipment: string(jsonequip),
    	Stats: string(jsonstats),
    	NPCs: string(jsonnpcs),
    	Fighter: string(jsonfighter),
    }

    //log.Print("[getFighterItems] jsonResp: ", jsonResp)    

    response, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[getFighterItems] error: ", err)
        return
    }
    respond(conn, response)
}

func getFighterStats(fighterID int64) FighterStats {

	// Connect to the Ethereum network using an Ethereum client
    client := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(FighterAttributesContract)
    contractABI := loadABI("FighteAttributes")

    // Prepare the call to the getTokenAttributes function
    tokenID := big.NewInt(fighterID)
    callData, err := contractABI.Pack("getFighterStats", tokenID)
    if err != nil {
        log.Fatalf("[getFighterStats] Failed to pack call data: %v", err)
    }

    // Call the contract using the Ethereum client
    result, err := client.CallContract(context.Background(), ethereum.CallMsg{
        To:   &contractAddress,
        Data: callData,
    }, nil)
    if err != nil {
        log.Fatalf("[getFighterStats] Failed to call contract fighterID=",fighterID," error=", err)
    }

    // Unpack the result into the attributes struct
    //var attributes []FighterAttributes
    var attributes []interface{};


    //err = contractABI.UnpackIntoInterface(&attributes, "getTokenAttributes", result)
    //attributes, err = contractABI.UnmarshalJSON("getTokenAttributes", result)
    attributes, err = contractABI.Unpack("getFighterStats", result)
    if err != nil {
        log.Printf("[getFighterStats] Failed to unpack error: %v", err)
    }

    jsonatts, err := json.Marshal(attributes[0])

    var fighter FighterStats
    json.Unmarshal(jsonatts, &fighter)
    if err != nil {
        log.Fatalf("[getFighterStats] Failed to call contract fighterID ", fighterID, "error=", err)
    }

   	return fighter;
}

func loadABI(contract string) (abi.ABI) {
    contractABI := abi.ABI{
		Methods: make(map[string]abi.Method),
		Events:  make(map[string]abi.Event),
	}

    type ABIEntry struct {
		Type            string          `json:"type"`
		Name            string          `json:"name"`
		Inputs          []abi.Argument  `json:"inputs"`
		Outputs         []abi.Argument  `json:"outputs"`
		StateMutability string          `json:"stateMutability"`
		Constant        bool            `json:"constant"`
		Payable         bool            `json:"payable"`
		Anonymous       bool            `json:"anonymous"`
	}


    type ContractABI struct {
		ContractName string     `json:"contractName"`
		ABI          []ABIEntry `json:"abi"`
	}


    // Read the contract ABI file
    raw, err := ioutil.ReadFile("../build/contracts/"+contract+".json")
    if err != nil {
        return contractABI
    }

    // Unmarshal the ABI JSON into the contractABI object
    var contractABIContent ContractABI
	err = json.Unmarshal(raw, &contractABIContent)
	if err != nil {
		fmt.Printf("Error unmarshalling ABI JSON: %v\n", err)
		return contractABI
	}

	// Loop through contractABIContent.ABI and create an abi.ABI object
	for _, entry := range contractABIContent.ABI {
		switch entry.Type {
		case "function":
			method := abi.Method{
				Name:            entry.Name,
				Inputs:          entry.Inputs,
				Outputs:         entry.Outputs,
				StateMutability: entry.StateMutability,
			}

			contractABI.Methods[method.Name] = method
		case "event":
			event := abi.Event{
				Name:      entry.Name,
				Anonymous: entry.Anonymous,
				Inputs:    entry.Inputs,
			}
			contractABI.Events[event.Name] = event
		}
	}

    return contractABI
}

func getHealthRegenerationRate(atts FighterAttributes) (float64) {

    vitality := atts.Vitality.Int64()
    healthRegenBonus := 0

    regenRate := (float64(vitality)/float64(HealthRegenerationDivider) + float64(healthRegenBonus))/5000;
    log.Printf("[getHealthRegenerationRate] vitality=", vitality," regenRate=", regenRate)
    return regenRate
}

func getHealth(id int64) int64 {
	maxHealth := Fighters[id].MaxHealth;
	lastDmgTimestamp := Fighters[id].LastDmgTimestamp;
	healthAfterLastDmg := Fighters[id].HealthAfterLastDmg;

    healthRegenRate := Fighters[id].HpRegenerationRate;
    currentTime := time.Now().UnixNano() / int64(time.Millisecond);

    health := float64(healthAfterLastDmg) + (float64((currentTime - lastDmgTimestamp)) * healthRegenRate)

    //log.Printf("[getHealth] currentTime=", currentTime," maxHealth=", maxHealth," lastDmgTimestamp=",lastDmgTimestamp," healthAfterLastDmg=",healthAfterLastDmg," healthRegenRate=", healthRegenRate, " health=", health)

    return min(maxHealth, int64(health))
}

func getNpcHealth(id int64) int64 {
	fighter := Fighters[id]
	if !fighter.IsDead {
		return getHealth(id);
	} else {
		now := time.Now().UnixNano() / 1e6
		elapsedTimeMs := now - fighter.LastDmgTimestamp

		if elapsedTimeMs >= 5000 {
			fmt.Println("At least 5 seconds have passed since TimeOfDeath.")
			fighter.IsDead = false;
			fighter.HealthAfterLastDmg = fighter.MaxHealth;
			Fighters[id] = fighter;
			return fighter.MaxHealth;
		} else {
			return 0;
		}	
	}
	
}

//// !!!!
func recordMove(conn *websocket.Conn, data json.RawMessage) {
	
	var moveData RecordMoveMsg
	err := json.Unmarshal(data, &moveData)
    if err != nil {
        log.Printf("websocket unmarshal error: %v", err)
        return
    }

	log.Print("[recordMove] playerID=", moveData.PlayerID, " opponentID=", moveData.OpponentID, " skill=", moveData.Skill);


	if moveData.OpponentID == 0 {
		log.Print("[recordMove] opponentId cannot be 0, playerId=", moveData.PlayerID);
		return
	}

	// battleRaw := getBattleRaw(moveData.BattleID)

	// if (battleRaw.Closed == 1) {
	// 	log.Print("[recordMove] Battle CloseD!  playerID=", moveData.PlayerID, " battleID=", moveData.BattleID);
	// 	return
	// }

	// coll := client.Database("game").Collection("moves")
	// newMove := Move{
	// 	BattleID: moveData.BattleID,
	// 	PlayerID: moveData.PlayerID,
	// 	Skill: moveData.Skill,
	// }
	// result, err := coll.InsertOne(context.Background(), newMove)
	

	ProcessHit(conn, moveData.PlayerID, moveData.OpponentID);   
}

func getEquippedItems(fighter FighterAttributes) []ItemAttributes {
	var items []ItemAttributes
	zero := big.NewInt(0)
	if fighter.HelmSlot.Cmp(zero) 		!= 0 { items = append(items, getItemAttributes(fighter.HelmSlot.Int64())) }
	if fighter.ArmourSlot.Cmp(zero) 		!= 0 { items = append(items, getItemAttributes(fighter.ArmourSlot.Int64())) }
	if fighter.PantsSlot.Cmp(zero) 		!= 0 { items = append(items, getItemAttributes(fighter.PantsSlot.Int64())) }
	if fighter.GlovesSlot.Cmp(zero) 		!= 0 { items = append(items, getItemAttributes(fighter.GlovesSlot.Int64())) }
	if fighter.BootsSlot.Cmp(zero) 		!= 0 { items = append(items, getItemAttributes(fighter.BootsSlot.Int64())) }
	if fighter.LeftHandSlot.Cmp(zero) 	!= 0 { items = append(items, getItemAttributes(fighter.LeftHandSlot.Int64())) }
	if fighter.RightHandSlot.Cmp(zero) 	!= 0 { items = append(items, getItemAttributes(fighter.RightHandSlot.Int64())) }
	if fighter.LeftRingSlot.Cmp(zero) 	!= 0 { items = append(items, getItemAttributes(fighter.LeftRingSlot.Int64())) }
	if fighter.RightRingSlot.Cmp(zero) 	!= 0 { items = append(items, getItemAttributes(fighter.RightRingSlot.Int64())) }
	if fighter.PendSlot.Cmp(zero) 		!= 0 { items = append(items, getItemAttributes(fighter.PendSlot.Int64())) }
	if fighter.WingsSlot.Cmp(zero) 		!= 0 { items = append(items, getItemAttributes(fighter.WingsSlot.Int64())) }
	return items;
}

func getTotalItemsDefence(items []ItemAttributes) int64 {
	var def = int64(0);
	for i := 0; i < len(items); i++ {
		def += items[i].Defense
	}

	return def;
}

func randomValueWithinRange(value int64, percentage float64) int64 {
	rand.Seed(time.Now().UnixNano())
	min := float64(value) * (1.0 - percentage)
	max := float64(value) * (1.0 + percentage)
	return int64(min + rand.Float64()*(max-min))
}

//// !!!!
func ProcessHit(conn *websocket.Conn, playerId, opponentId int64) {

    stats1 := getFighterAttributes(playerId);
    stats2 := getFighterAttributes(opponentId);

    //def1 := stats1.Agility.Int64()/4;
    def2 := stats2.Agility.Int64()/4;

    dmg1 := randomValueWithinRange((stats1.Strength.Int64()/4 + stats1.Energy.Int64()/4), 0.25);
    //dmg2 := randomValueWithinRange((stats2.Strength.Int64()/4 + stats2.Energy.Int64()/4), 0.25);


    //items1 := getEquippedItems(stats1);
    items2 := getEquippedItems(stats2);

    //itemDefence1 := getTotalItemsDefence(items1)
    itemDefence2 := getTotalItemsDefence(items2)

   	var damage float64;
   	var oppNewHealth int64;

   	// Update battle 
	damage = float64(max(0, dmg1 - def2 - itemDefence2));
	oppNewHealth = max(0, getNpcHealth(opponentId) - int64(damage));    	

   	fighter := Fighters[opponentId]
   	if (fighter.IsNpc) {
   		if fighter.IsDead {
			now := time.Now().UnixNano() / 1e6
			elapsedTimeMs := now - fighter.LastDmgTimestamp

			if elapsedTimeMs >= 5000 {
				fmt.Println("At least 5 seconds have passed since TimeOfDeath.")
				fighter.IsDead = false;
				fighter.HealthAfterLastDmg = fighter.MaxHealth;
			} else {
				log.Printf("[ProcessHit] NPC Dead playerId=", playerId, "opponentId=", opponentId)
				return
			}
   			
   		} else if oppNewHealth == 0 {
   			fighter.IsDead = true;
   		}
   	}
   	fighter.LastDmgTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
   	fighter.HealthAfterLastDmg = oppNewHealth
   	Fighters[opponentId] = fighter

   	if (oppNewHealth == 0) {
   		recordBattleOnChain(playerId, opponentId)
   	}
   	
   	type jsonResponse struct {
		Action string `json:"action"`
    	Damage int64 `json:"damage"`
    	Opponent int64 `json:"opponent"`
    	Player int64 `json:"player"`
    	OpponentNewHealth int64 `json:"opponentHealth"`
    	LastDmgTimestamp int64 `json:"lastDmgTimestamp"`
    	HealthAfterLastDmg int64 `json:"healthAfterLastDmg"`
    }

    jsonResp := jsonResponse{
    	Action: "damage_dealt",
    	Damage: int64(damage),
		Opponent: opponentId,
		Player: playerId,
		OpponentNewHealth: oppNewHealth,
		LastDmgTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
    }

    // Convert the struct to JSON
    response, err := json.Marshal(jsonResp)
    if err != nil {
        log.Print("[ProcessHit] error: ", err)
        return
    }

    respond(conn, response)


    log.Println("[ProcessHit] damage=", damage, "opponentId=", opponentId, "playerId=", playerId);
}

func updateInt(coll string, objectID primitive.ObjectID, key string, value int64) {
	

	filter := bson.M{"_id": objectID};

	collection := client.Database("game").Collection(coll)
	update := bson.M{"$set": bson.M{key:value}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
	    // Handle error
	    log.Fatal(err)
	}

	//log.Println("[UpdateBattleRaw] battle ", objectID, " result=", result);
}

func updateString(coll string, objectID primitive.ObjectID, key string, value string) {
	

	filter := bson.M{"_id": objectID};

	collection := client.Database("game").Collection(coll)
	update := bson.M{"$set": bson.M{key:value}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
	    // Handle error
	    log.Fatal(err)
	}

	//log.Println("[UpdateBattleRaw] battle ", objectID, " result=", result);
}

func getBattleRaw(battleID string) Battle {
	collection := client.Database("game").Collection("battles")

	// define a filter to select the battle by ID
	objectID, err := primitive.ObjectIDFromHex(battleID)
	if err != nil{
	    log.Println("[getBattleRaw] Invalid id ", battleID);
	}

	filter := bson.M{"_id": objectID};

	// define a variable to hold the result
	var battle Battle

	// find the battle row that matches the filter
	err = collection.FindOne(context.Background(), filter).Decode(&battle)
	if err != nil {
	    log.Fatal(err)
	}

	//log.Print("[getBattleRaw] battleID=", battle.ID, " Opponent1 ", battle.Opponent1);
	return battle
}

func getPlayerRaw(playerID int64) Player {
	log.Println("[getPlayerRaw] playerID=[", playerID,"]");
	collection := client.Database("game").Collection("players")


	filter := bson.M{"playerid": playerID};

	// define a variable to hold the result
	var player Player

	// find the battle row that matches the filter
	err := collection.FindOne(context.TODO(), filter).Decode(&player)
	if err != nil {
	    log.Fatal(err)
	}

	//log.Print("[getPlayerRaw] playerID=", playerID, " Strength ", player.Strength);
	return player
}


func respond(conn *websocket.Conn, response json.RawMessage) {
	// w.Header().Set("Content-Type", "application/json")
    // w.Header().Set("Access-Control-Allow-Origin", "*")
    // w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

    err := conn.WriteMessage(websocket.TextMessage, response)

    if err != nil {
		log.Println("[respond] ", err)
        return
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    var msg struct {
        Type string `json:"type"`
        Data json.RawMessage `json:"data"`
    }


	var upgrader = websocket.Upgrader{
	    ReadBufferSize:  1024,
	    WriteBufferSize: 1024,
	    CheckOrigin: func(r *http.Request) bool {
	        // allow all connections by default
	        return true
	    },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()


		if err != nil {
	  		//log.Println("Failed to read message from WebSocket:", err)
	  		break
		}

		//log.Printf("message: ", message)

		err = json.Unmarshal(message, &msg)
        if err != nil {
            log.Printf("websocket unmarshal error: %v", err)
            continue
        }

        // log.Printf("Type: ", msg.Type)
        // log.Printf("Data: ", msg.Data)
        switch msg.Type {
        case "auth":
        	authFighter(conn, 1);
	            //createNewBattle(conn, msg.Data, w, r)
	            continue

	        case "startNewBattle":
	            //createNewBattle(conn, msg.Data, w, r)
	            continue
	            
	        case "recordMove":
	            recordMove(conn, msg.Data)
	            continue

	        case "getFighterItems":
	            getFighterItems(conn, msg.Data, w, r)
	            continue

	        

	        // case "getFighter":
	        //     getFighter(conn, msg.Data, w, r)
	        //     continue
	            
	        default:
	            log.Printf("unknown message type: %s", msg.Type)
	            continue
        }

		data := decodeJson(message);
		log.Printf("Received message: %s\n", data["action"].(string))

		// Handle the message here and send a response back to the client
		// response := "Hello, client!"
		// conn.WriteMessage(websocket.TextMessage, []byte(response))
	}
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

func ConnectToDB() *mongo.Client {
	// Set up MongoDB client options
    //connStr := "mongodb+srv://admin:sydeBlx2pDfiy0CP@cluster0.bwinsau.mongodb.net/?retryWrites=true&w=majority"
    connStr := "mongodb://localhost:27017"
    clientOptions := options.Client().ApplyURI(connStr)

    // Create a MongoDB client
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal("[ConnectToDB] ", err)
    }

    // Check the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal("[ConnectToDB] ", err)
    } else {
    	log.Print("Connected to MangoDB ");   
    }

    _, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	return client
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
    PrivateKey = os.Getenv("PRIVATE_KEY")

    fmt.Println("FighterAttributesContract:", FighterAttributesContract)
    fmt.Println("BattleContract:", BattleContract)
    fmt.Println("ItemsContract:", ItemsContract)
    
}

func main() {
	loadEnv()
	initMap()    
	loadNPCs()    
	//listenToBattleContractEvents()  	
	//subscribeToNewBlock()
    // Set up a web server
    // http.HandleFunc("/createNewBattle", createNewBattle)    
    // http.HandleFunc("/move", recordMove)    
    http.HandleFunc("/ws", handleWebSocket)

    // Start the server
    log.Fatal(http.ListenAndServe(":8080", nil))

    

    defer client.Disconnect(context.TODO())
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
