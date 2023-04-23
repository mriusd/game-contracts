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
    "strings"
    "strconv"
    "errors"
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

type Damage struct {
    FighterId        *big.Int
    Damage           *big.Int
    //MedianPartyLevel *big.Int
}

type Fighter struct {
	ID    					string  		`json:"id"`
    MaxHealth     			int64 			`json:"maxHealth"`
    Name           			string 			`json:"name"`
    IsNpc         			bool 			`json:"isNpc"`
    IsDead         			bool 			`json:"isDead"`
    CanFight 				bool 			`json:"canFight"`
    LastDmgTimestamp 		int64  			`json:"lastDmgTimestamp"`
    HealthAfterLastDmg 		int64  			`json:"healthAfterLastDmg"`
    HpRegenerationRate 		float64 		`json:"hpRegenerationRate"`
    HpRegenerationBonus 	float64 		`json:"hpRegenerationBonus"`
    TokenID                 int64           `json:"tokenId"`
    Location                string          `json:"location"`
    AttackSpeed             int64           `json:"attackSpeed"` 
    DamageReceived          []Damage        `json:"damageDealt"`

    Conn 					*websocket.Conn

}

type Coordinate struct {
    X int64
    Y int64
}

var npcs []NPC;
var Fighters = make(map[string]Fighter)
var Battles = make(map[primitive.ObjectID]Battle)

var uniqueNpcIdCounter int64 = 1000

var fighterAttributesCache = make(map[string]FighterAttributes)






var Population = make(map[string]map[Coordinate][]Fighter)


var HealthRegenerationDivider = 8;
var ManaRegenerationDivider = 8;
var AgilityPerDefence = 4;
var StrengthPerDamage = 8;
var EnergyPerDamage = 8;
var MaxExperience = 291342500;
var ExperienceDivider = 5;

type WsMessage struct {
    Type string  `json:"type"`
    Data Fighter `json:"data"`
}

func emitNpcSpawnMessage(npc Fighter) {
    for _, fighter := range Fighters {
        if !fighter.IsNpc && fighter.Location == npc.Location {
            err := sendSpawnNpcMessage(fighter.Conn, npc)
            if err != nil {
                log.Printf("Error sending spawn_npc message to fighter %s: %v", fighter.ID, err)
            }
        }
    }
}

func sendSpawnNpcMessage(conn *websocket.Conn, npc Fighter) error {
    log.Printf("[sendSpawnNpcMessage] ", npc)
    if conn == nil {
        return errors.New("WebSocket connection is nil")
    }

    type jsonResponse struct {
        Action string `json:"action"`
        Npc Fighter `json:"npc"`
    }

    jsonResp := jsonResponse{
        Action: "spawn_npc",
        Npc: npc,
    }

    messageJSON, err := json.Marshal(jsonResp)
    if err != nil {
        return err
    }

    return conn.WriteMessage(websocket.TextMessage, messageJSON)
}

func addDamageToFighter(fighterID string, hitterID *big.Int, damage *big.Int) {
    found := false
    fighter := Fighters[fighterID];

    //log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v fighter=%v", fighterID, hitterID, damage, fighter)

    damageReceived := fighter.DamageReceived

    // Check if there's already damage from the hitter
    for i, dmg := range damageReceived {
        if dmg.FighterId.Cmp(hitterID) == 0 {
            found = true
            log.Printf("[addDamageToFighter] Damage found ")
            // Add the new damage to the existing damage
            damageReceived[i].Damage = new(big.Int).Add(damageReceived[i].Damage, damage)
            break
        }
    }

    // If no existing damage from the hitter, create a new Damage object and append it to DamageReceived
    if !found {
        newDamage := Damage{
            FighterId: hitterID,
            Damage:    damage,
        }
        damageReceived = append(damageReceived, newDamage)
    }

    fighter.DamageReceived = damageReceived
    Fighters[fighterID] = fighter;

    //log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v fighter=%v", fighterID, hitterID, damage, fighter)
}

func getNextUniqueNpcId() string {
    uniqueNpcIdCounter++
    return "npc_" + strconv.Itoa(int(uniqueNpcIdCounter))
}

func findNpcById(id int64) *NPC {
    for _, npc := range npcs {
        if npc.ID == id {
            return &npc
        }
    }
    return nil
}

func convertIdToString(id int64) string {
    return strconv.Itoa(int(id))
}

func authFighter(conn *websocket.Conn, playerId int64, locationKey string) {
    log.Printf("[authFighter] playerId=%v locationKey=%v", playerId, locationKey)
    strId := strconv.Itoa(int(playerId));
    stats := getFighterStats(playerId);
        

    fighter := Fighter{
        ID: strId,
        TokenID: playerId,
        MaxHealth: stats.MaxHealth.Int64(), 
        Name: "",
        IsNpc: false,
        CanFight: true,
        HpRegenerationRate: 0,
        LastDmgTimestamp: 0,
        HealthAfterLastDmg: 0,
        Conn: conn,
    }

    Fighters[strId] = fighter;

    fighterAttributes := getFighterAttributes(strId)
    
    

    location := decodeLocation(locationKey);
    town := location[0]
    x, _ := strconv.ParseInt(location[1], 10, 64)
    y, _ := strconv.ParseInt(location[2], 10, 64)

    coord := Coordinate{X: x, Y: y}

    // Remove the fighter from any other population
    for key, zone := range Population {        
        for coord, npcList := range zone {
            for i := 0; i < len(npcList); i++ {
                if npcList[i].ID == fighter.ID {
                    Population[key][coord] = append(npcList[:i], npcList[i+1:]...)
                    break
                }
            }
        }
    }

    // Check if the fighter is not already in the map
    found := false
    var f Fighter
    for coord, npcList := range Population[town] {
        for i := 0; i < len(npcList); i++ {
            f = Fighters[npcList[i].ID]
            if f.ID == fighter.ID {
                found = true
                fighter.Conn = conn
                fighter.Location = locationKey
                Population[locationKey][coord][i] = fighter
                break
            }
        }
        if found {
            break
        }
    }

    if !found {
        fighter.Location = locationKey
        if Population[town] == nil {
            Population[town] = make(map[Coordinate][]Fighter)
        }
        Population[town][coord] = append(Population[town][coord], fighter)
    }

    log.Printf("[authFighter] ", Population[town])
    fighter.HpRegenerationRate = getHealthRegenerationRate(fighterAttributes);
    Fighters[strId] = fighter;



}

func decodeLocation(locationHash string) []string {
    return strings.Split(locationHash, "_")
}

func spawnNPC(npcId int64, location []string) Fighter {
    
    npc := findNpcById(npcId)
    log.Printf("[spawnNPC] %v %v", npcId, npc)
    locationHash := strings.Join(location, "_")

    uniqueNpcId := getNextUniqueNpcId()

    Fighters[uniqueNpcId] = Fighter{
        ID: uniqueNpcId,
        MaxHealth: npc.MaxHealth, 
        Name: npc.Name,
        IsNpc: true,
        CanFight: npc.CanFight,
        HpRegenerationRate: 0,
        HpRegenerationBonus: 0,
        LastDmgTimestamp: 0,
        HealthAfterLastDmg: npc.MaxHealth,
        TokenID: npcId,
        Location: locationHash,
        AttackSpeed: npc.AttackSpeed,
    }

    fighter := Fighters[uniqueNpcId];

    emitNpcSpawnMessage(fighter);

    zone := location[0]
    x, _ := strconv.ParseInt(location[1], 10, 64)
    y, _ := strconv.ParseInt(location[2], 10, 64)

    coord := Coordinate{X: x, Y: y}

    if _, exists := Population[zone]; !exists {
        Population[zone] = make(map[Coordinate][]Fighter)
    }

    Population[zone][coord] = append(Population[zone][coord], fighter)

    return fighter;
}

func initiateNpcRoutine(npcId string) {
	
	//fmt.Println("[initiateNpcRoutine] npcId=", fighter.ID)	
    fighter := Fighters[npcId]
	speed := fighter.AttackSpeed;

	msPerHit := 60000/speed;
	delay := time.Duration(msPerHit) * time.Millisecond;

    location := decodeLocation(fighter.Location);

    zone := location[0]
    x, _ := strconv.ParseInt(location[1], 10, 64)
    y, _ := strconv.ParseInt(location[2], 10, 64)

    coord := Coordinate{X: x, Y: y}

	for {
        fighter = Fighters[npcId]
        now := time.Now().UnixNano() / 1e6
        elapsedTimeMs := now - fighter.LastDmgTimestamp
		if fighter.IsDead && elapsedTimeMs >= 5000  {
            fmt.Println("At least 5 seconds have passed since TimeOfDeath.")
            fighter.IsDead = false;
            fighter.HealthAfterLastDmg = fighter.MaxHealth;
            fighter.DamageReceived = []Damage{};
            Fighters[fighter.ID] = fighter;
            emitNpcSpawnMessage(fighter);
		} else if len(Population[zone][coord]) > 0 {
			rand.Seed(time.Now().UnixNano()) // Seed the random number generator

            nonNpcFighters := []Fighter{}
            for _, fighter := range Population[zone][coord] {
                if !fighter.IsNpc {
                    nonNpcFighters = append(nonNpcFighters, fighter)
                }
            }


            if len(nonNpcFighters) > 0 {
    			randomIndex := rand.Intn(len(nonNpcFighters))
                randomFighter := nonNpcFighters[randomIndex]

    			data := RecordHitMsg{
    				OpponentID: randomFighter.ID,
    				PlayerID: fighter.ID,
    				Skill: 0,
    			}

    			rawMessage, err := json.Marshal(data);
    			if err != nil {
    				fmt.Println("[initiateNpcRoutine] Error marshaling data:", err)
    				return
    			}

    			ProcessHit(randomFighter.Conn, rawMessage)
            }
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
    err = json.Unmarshal(data, &npcs)
    if err != nil {
        log.Printf("[loadNPCs] error= ", err)
    }

    log.Printf("[loadNPCs] %v", npcs)

    // Set default values and initiate NPC routines
    for i, npc := range npcs {
        npcs[i].CanFight = true
        npcs[i].MaxHealth = npc.Vitality

        // Iterate through respawn locations
        for _, location := range npc.RespawnLocations {
            fighter := spawnNPC(npc.ID, location)
            go initiateNpcRoutine(fighter.ID)
        }

        
    }

    log.Printf("NPCs Loaded", npcs )
}

func getNPCs(locationHash string) []Fighter {
    location := decodeLocation(locationHash);

    zone := location[0]
    x, _ := strconv.ParseInt(location[1], 10, 64)
    y, _ := strconv.ParseInt(location[2], 10, 64)

    coord := Coordinate{X: x, Y: y}
    npcFighters := []Fighter{}
    for _, fighter := range Population[zone][coord] {
        if fighter.IsNpc {
            fighter.LastDmgTimestamp = Fighters[fighter.ID].LastDmgTimestamp
            fighter.HealthAfterLastDmg = Fighters[fighter.ID].HealthAfterLastDmg
            npcFighters = append(npcFighters, fighter)
        }
    }

    return npcFighters
}

func recordBattleOnChain(player, opponent string) (string) {
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



    type DamageTuple struct {
        FighterId        *big.Int
        Damage           *big.Int
    }

	killedFighter := big.NewInt(Fighters[opponent].TokenID)
    damageDealt := Fighters[opponent].DamageReceived
	// damageDealt := []Damage{
	//     {
	//         FighterId:        big.NewInt(Fighters[player].TokenID),
	//         Damage:           big.NewInt(damage),
	//         MedianPartyLevel: big.NewInt(1),
	//     },
	// }
	battleNonce := big.NewInt(time.Now().UnixNano() / int64(time.Millisecond))

    log.Printf("[recordBattleOnChain] damageDealt %v", damageDealt)

	damageDealtTuples := make([]DamageTuple, len(damageDealt))
    for i, d := range damageDealt {
        damageDealtTuples[i] = DamageTuple{
            FighterId:        d.FighterId,
            Damage:           d.Damage,
        }
    }

	log.Printf("[recordBattleOnChain] battleNonce=", battleNonce)

	// Prepare transaction options
	nonce, err := client.NonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey), nil)
	if err != nil {
		log.Printf("Failed to retrieve nonce: %v", err)
	}
	gasLimit := uint64(5000000)
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
	data, err := contractABI.Pack("recordKill", killedFighter, damageDealtTuples, battleNonce)
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

func getFighterAttributes(id string) FighterAttributes {

    //log.Printf("[getFighterAttributes] id: %v", id)
    fighterID := Fighters[id].TokenID
    atts, ok := fighterAttributesCache[id];
    if Fighters[id].IsNpc && ok {
        return atts
    }

	// Connect to the Ethereum network using an Ethereum client
    rpcClient := getRpcClient();

    // Define the contract address and ABI
    contractAddress := common.HexToAddress(FighterAttributesContract)
    contractABI := loadABI("FighterAttributes")

    //log.Printf("contractABI: ", contractABI.Methods[method.Name]);

    // Prepare the call to the getTokenAttributes function
    tokenId := big.NewInt(fighterID) 
    callData, err := contractABI.Pack("getTokenAttributes", tokenId)
    if err != nil {
        log.Fatalf("[getFighterAttributes] Failed to pack call data: %v", err)
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

    //log.Printf("result: %v ", result);

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

    fighterAttributesCache[id] = fighter
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

    _, ok := Fighters[convertIdToString(reqData.FighterId)]
    if !ok {
        log.Printf("[getFighterItems] Fighter not authenticated", )
        return;
    }

    log.Printf("[getFighterItems] reqData: %v", reqData)
    

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

	log.Printf("[getFighterItems] attributes: %v", attributes)

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

    log.Printf("[getFighterItems] items: %v", items)

	var jsonatts []byte
	if len(items) == 0 {
		jsonatts = nil;
	} else {
		jsonatts, err = json.Marshal(items)
	}	

	stats := getFighterStats(reqData.FighterId);

    jsonstats, err := json.Marshal(stats)
    //log.Print("[getFighter] jsonstats: %s", stats)


    fighterAttributes := getFighterAttributes(convertIdToString(reqData.FighterId));
    jsonfighteratts, err := json.Marshal(fighterAttributes)


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


    

    npcatts := getNPCs("lorencia_0_0");
    log.Print("[getFighterItems] npcs: ", npcs)
    jsonnpcs, err := json.Marshal(npcatts)


    jsonfighter, err := json.Marshal(Fighters[convertIdToString(reqData.FighterId)])

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
    contractABI := loadABI("FighterAttributes")

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

func getHealthRegenerationRate(atts FighterAttributes) (float64) {

    vitality := atts.Vitality.Int64()
    healthRegenBonus := 0

    regenRate := (float64(vitality)/float64(HealthRegenerationDivider) + float64(healthRegenBonus))/5000;
    log.Printf("[getHealthRegenerationRate] vitality=", vitality," regenRate=", regenRate)
    return regenRate
}

func getHealth(id string) int64 {
	maxHealth := Fighters[id].MaxHealth;
	lastDmgTimestamp := Fighters[id].LastDmgTimestamp;
	healthAfterLastDmg := Fighters[id].HealthAfterLastDmg;

    healthRegenRate := Fighters[id].HpRegenerationRate;
    currentTime := time.Now().UnixNano() / int64(time.Millisecond);

    health := float64(healthAfterLastDmg) + (float64((currentTime - lastDmgTimestamp)) * healthRegenRate)

    //log.Printf("[getHealth] currentTime=", currentTime," maxHealth=", maxHealth," lastDmgTimestamp=",lastDmgTimestamp," healthAfterLastDmg=",healthAfterLastDmg," healthRegenRate=", healthRegenRate, " health=", health)

    return min(maxHealth, int64(health))
}

func getNpcHealth(id string) int64 {
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

type RecordHitMsg struct {
    PlayerID    string  `json:"playerID"`
    OpponentID  string  `json:"opponentID`
    Skill       int64   `json:"skill"`
}

//// !!!!
func ProcessHit(conn *websocket.Conn, data json.RawMessage) {
    var hitData RecordHitMsg
    err := json.Unmarshal(data, &hitData)
    if err != nil {
        log.Printf("websocket unmarshal error: %v", err)
        return
    }

    playerFighter, ok := Fighters[hitData.PlayerID]
    if !ok {
        log.Printf("[ProcessHit] Unknown Player");
        return;
    }

    opponentFighter, ok := Fighters[hitData.OpponentID]
    if !ok {
        log.Printf("[ProcessHit] Unknown Opponent");
        return;
    }


    playerId := playerFighter.TokenID
    opponentId := opponentFighter.TokenID
    

    stats1 := getFighterAttributes(hitData.PlayerID);
    stats2 := getFighterAttributes(hitData.OpponentID);

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
    npcHealth := getNpcHealth(hitData.OpponentID)

   	// Update battle 
	damage = float64(min(npcHealth, max(0, dmg1 - def2 - itemDefence2)));
	oppNewHealth = max(0, npcHealth - int64(damage));    	


   	if (opponentFighter.IsNpc) {
   		if opponentFighter.IsDead {
			now := time.Now().UnixNano() / 1e6
			elapsedTimeMs := now - opponentFighter.LastDmgTimestamp

			if elapsedTimeMs >= 5000 {
				fmt.Println("At least 5 seconds have passed since TimeOfDeath.")
				opponentFighter.IsDead = false;
                opponentFighter.DamageReceived = []Damage{};
				opponentFighter.HealthAfterLastDmg = opponentFighter.MaxHealth;
			} else {
				log.Printf("[ProcessHit] NPC Dead playerId=", playerId, "opponentId=", opponentId)
				return
			}
   			
   		} else if oppNewHealth == 0 {
            opponentFighter.IsDead = true;
   			
   		}
   	}
   	opponentFighter.LastDmgTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
   	opponentFighter.HealthAfterLastDmg = oppNewHealth
    Fighters[hitData.OpponentID] = opponentFighter

    if damage > 0 {
        addDamageToFighter(hitData.OpponentID, big.NewInt(playerId), big.NewInt(int64(damage)))
    }	

   	if (oppNewHealth == 0) {
   		recordBattleOnChain(hitData.PlayerID, hitData.OpponentID)
   	}

    
   	
   	type jsonResponse struct {
		Action string `json:"action"`
    	Damage int64 `json:"damage"`
    	Opponent string `json:"opponent"`
    	Player string `json:"player"`
    	OpponentNewHealth int64 `json:"opponentHealth"`
    	LastDmgTimestamp int64 `json:"lastDmgTimestamp"`
    	HealthAfterLastDmg int64 `json:"healthAfterLastDmg"`
    }

    jsonResp := jsonResponse{
    	Action: "damage_dealt",
    	Damage: int64(damage),
		Opponent: hitData.OpponentID,
		Player: hitData.PlayerID,
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
        	authFighter(conn, 1, "lorencia_0_0");
                //createNewBattle(conn, msg.Data, w, r)
            continue
            
        case "recordMove":
            ProcessHit(conn, msg.Data)
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


func lastBlockNumber() (uint64, error) {
    client := getRpcClient()
    header, err := client.HeaderByNumber(context.Background(), nil)
    if err != nil {
        return 0, err
    }

    var blockNumber = header.Number.Uint64()
    log.Printf("[lastBlockNumber] %v", blockNumber)

    return blockNumber, nil
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
	//initMap()    	
    lastBlockNumber()   
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
