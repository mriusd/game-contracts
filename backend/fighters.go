package main

import (
	"math/big"
	"sync"

    "github.com/gorilla/websocket"	
    "log"
    "strconv"

    "time"
    "fmt"
    "math/rand"
)

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
    OwnerAddress            string          `json:"ownerAddress"`
    Coordinates             Coordinate      `json:"coordinates"`
    MovementSpeed           int64           `json:"movementSpeed"` // squares per minute
    IsClosed                bool            `json:"isClosed"` 
    AttackDistance          int64           `json:"attackDistance"`
    Skill                   int64           `json:"skill"`
    SpawnCoords             Coordinate      `json:"spawnCoords"`

    Conn 					*websocket.Conn
    ConnMutex               sync.Mutex
}

var Fighters = make(map[string]*Fighter)
var FightersMutex sync.RWMutex

var FighterAttributesCache = make(map[string]FighterAttributes)
var FighterAttributesCacheMutex sync.RWMutex


func getFighterSafely(id string) *Fighter {
    FightersMutex.Lock()
    defer FightersMutex.Unlock()
    return Fighters[id]
}


func getHealthRegenerationRate(atts FighterAttributes) (float64) {

    vitality := atts.Vitality.Int64()
    healthRegenBonus := 0

    regenRate := (float64(vitality)/float64(HealthRegenerationDivider) + float64(healthRegenBonus))/5000;
    log.Printf("[getHealthRegenerationRate] vitality=", vitality," regenRate=", regenRate)
    return regenRate
}

func getNpcHealth(fighter *Fighter) int64 {
	if !fighter.IsDead {
		return getHealth(fighter);
	} else {
		now := time.Now().UnixNano() / 1e6
		elapsedTimeMs := now - fighter.LastDmgTimestamp

		if elapsedTimeMs >= 5000 {
			fmt.Println("[getNpcHealth] At least 5 seconds have passed since TimeOfDeath.")
            FightersMutex.Lock()
			fighter.IsDead = false;
			fighter.HealthAfterLastDmg = fighter.MaxHealth;
            FightersMutex.Unlock()
			return fighter.MaxHealth;
		} else {
			return 0;
		}	
	}
	
}

func getHealth(fighter *Fighter) int64 {
	maxHealth := fighter.MaxHealth;
	lastDmgTimestamp := fighter.LastDmgTimestamp;
	healthAfterLastDmg := fighter.HealthAfterLastDmg;

    healthRegenRate := fighter.HpRegenerationRate;
    currentTime := time.Now().UnixNano() / int64(time.Millisecond);

    health := float64(healthAfterLastDmg) + (float64((currentTime - lastDmgTimestamp)) * healthRegenRate)

    //log.Printf("[getHealth] currentTime=", currentTime," maxHealth=", maxHealth," lastDmgTimestamp=",lastDmgTimestamp," healthAfterLastDmg=",healthAfterLastDmg," healthRegenRate=", healthRegenRate, " health=", health)

    return min(maxHealth, int64(health))
}

func initiateFighterRoutine(fighter *Fighter) {
    log.Printf("[initiateFighterRoutine] fighter=", fighter)
    speed := fighter.MovementSpeed

    msPerHit := 60000 / speed
    delay := time.Duration(msPerHit) * time.Millisecond

    for {
        time.Sleep(delay)

        // Check if the connection is closed
        fighter.ConnMutex.Lock()
        isClosed := fighter.IsClosed
        fighter.ConnMutex.Unlock()

        if isClosed {
            // Break the loop if the connection is closed
            log.Printf("[initiateFighterRoutine] Connection closed, stopping the loop for fighter:", fighter.ID)
            break
        }

        pingFighter(fighter)
    }
}


func authFighter(conn *websocket.Conn, playerId int64, ownerAddess string, locationKey string) {
    log.Printf("[authFighter] playerId=%v ownerAddess=%v locationKey=%v", playerId, ownerAddess, locationKey)
    strId := strconv.Itoa(int(playerId));
    stats := getFighterStats(playerId);

    location := decodeLocation(locationKey);
    town := location[0]

    centerCoord := Coordinate{X: 5, Y: 5}
    emptySquares := getEmptySquares(centerCoord, 5, town)

    // if len(emptySquares) == 0 {
    //     return // No empty squares available to spawn the NPC
    // }

    rand.Seed(time.Now().UnixNano())
    spawnCoord := emptySquares[rand.Intn(len(emptySquares))]
        

    fighter := &Fighter{
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
        OwnerAddress: ownerAddess,
        MovementSpeed: 60,
        Coordinates: spawnCoord,
    }

    FightersMutex.Lock()
    Fighters[strId] = fighter
    FightersMutex.Unlock()

    fighterAttributes, err := getFighterAttributes(fighter)
    if err != nil {
        return;
    }
    
    
    if Population[town] == nil {
        Population[town] = make([]*Fighter, 0)
    }
    Population[town] = append(Population[town], fighter)

    log.Printf("[authFighter] ", Population[town])
    FightersMutex.Lock()
    fighter.Location = town
    fighter.HpRegenerationRate = getHealthRegenerationRate(fighterAttributes);
    FightersMutex.Unlock()

    go initiateFighterRoutine(fighter)
}

func findFighterByConn(conn *websocket.Conn) *Fighter {
    FightersMutex.Lock()
    defer FightersMutex.Unlock()
    for _, fighter := range Fighters {
        // Use ConnMutex to avoid data race while reading the Conn field
        fighter.ConnMutex.Lock()
        isMatchingConnection := fighter.Conn == conn
        fighter.ConnMutex.Unlock()

        if isMatchingConnection {
            return fighter
        }
    }
    return nil // Return nil if no matching Fighter is found
}

func addDamageToFighter(fighterID string, hitterID *big.Int, damage *big.Int) {
    found := false
    fighter := getFighterSafely(fighterID);

    //log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v fighter=%v", fighterID, hitterID, damage, fighter)

    damageReceived := fighter.DamageReceived

    // Check if there's already damage from the hitter
    for i, dmg := range damageReceived {
        if dmg.FighterId.Cmp(hitterID) == 0 {
            found = true
            //log.Printf("[addDamageToFighter] Damage found ")
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
    fighter.ConnMutex.Lock()
    fighter.DamageReceived = damageReceived
    fighter.ConnMutex.Unlock()

    //log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v fighter=%v", fighterID, hitterID, damage, fighter)
}