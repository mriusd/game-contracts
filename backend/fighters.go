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

    "github.com/ethereum/go-ethereum/common"

    "errors"
    "regexp"
    "unicode/utf8"
)

type FighterAttributes struct {
    Name                    string `json:"Name"`
    Class                   string `json:"Class"`
	TokenID   				*big.Int `json:"TokenID"`
    Strength                *big.Int `json:"Strength"`
	Agility                 *big.Int `json:"Agility"`
	Energy                  *big.Int `json:"Energy"`
	Vitality                *big.Int `json:"Vitality"`
	Experience              *big.Int `json:"Experience"`
	HpPerVitalityPoint      *big.Int `json:"HpPerVitalityPoint"`
	ManaPerEnergyPoint      *big.Int `json:"ManaPerEnergyPoint"`
	HpIncreasePerLevel      *big.Int `json:"HpIncreasePerLevel"`
	ManaIncreasePerLevel   	*big.Int `json:"manaIncreasePerLevel"`
	StatPointsPerLevel   	*big.Int `json:"statPointsPerLevel"`
	AttackSpeed   			*big.Int `json:"attackSpeed"`
	AgilityPointsPerSpeed   *big.Int `json:"agilityPointsPerSpeed"`
    IsNpc                   *big.Int `json:"isNpc"`
	DropRarityLevel   	    *big.Int `json:"dropRarityLevel"`
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
	ID    					string  		    `json:"id"`
    MaxHealth     			int64 			    `json:"maxHealth"`
    Name           			string 			    `json:"name"`
    IsNpc         			bool 			    `json:"isNpc"`
    IsDead         			bool 			    `json:"isDead"`
    CanFight 				bool 			    `json:"canFight"`
    LastDmgTimestamp 		int64  			    `json:"lastDmgTimestamp"`
    HealthAfterLastDmg 		int64  			    `json:"healthAfterLastDmg"`

    TokenID                 int64               `json:"tokenId"`
    Location                string              `json:"location"`
    
    DamageReceived          []Damage            `json:"damageDealt"`
    OwnerAddress            string              `json:"ownerAddress"`
    Coordinates             Coordinate          `json:"coordinates"`
    MovementSpeed           int64               `json:"movementSpeed"` // squares per minute
    Skill                   int64               `json:"skill"`
    SpawnCoords             Coordinate          `json:"spawnCoords"`
    

    LastMoveTimestamp       int64               `json:"lastMoveTimestamp"` // milliseconds


    // Fighter stats
    Strength                int64               `json:"strength"`
    Agility                 int64               `json:"agility"`
    Energy                  int64               `json:"energy"`
    Vitality                int64               `json:"vitality"`


    // Fighter dynamic paramters
    CurrentHealth           int64               `json:"currentHealth"`
    CurrentMana             int64               `json:"currentMana"`


    // Fighter parameters with equipped items
    Damage                  int64               `json:"damage"`
    Defence                 int64               `json:"defence"`
    AttackSpeed             int64               `json:"attackSpeed"` 
    HpRegenerationRate      float64             `json:"hpRegenerationRate"`
    HpRegenerationBonus     float64             `json:"hpRegenerationBonus"`

    // Damage type rates
    CriticalDmgRate         int64               `json:"criticalDmgRate"`
    ExcellentDmgRate        int64               `json:"excellentDmgRate"`
    DoubleDmgRate           int64               `json:"doubleDmgRate"`
    IgnoreDefRate           int64               `json:"ignoreDefRate"`


    Level                   int64               `json:"level"`
    Experience              int64               `json:"experience"`

    Direction               Direction           `json:"direction"`

    Skills                  map[int64]*Skill    `json:"skills"`
    Backpack                *Inventory           `json:"-"`
    Vault                   *Inventory           `json:"-"`
    Equipment               map[int64]*InventorySlot `json:"equipment"`

    LastChatMsg             string              `json:"lastChatMessage"`
    LastChatMsgTimestamp    int64               `json:"lastChatMsgTimestamp"`

    Credits                 int64               `json:"credits"`

    Mutex                   sync.RWMutex        `json:"-"`
}


type FighterCreatedEvent struct {
    TokenID         *big.Int            `json:"tokenId"`
    Owner           common.Address      `json:"owner"`
    
    FighterClass    string            `json:"fighterClass"`
    Name            string              `json:"name"`
}

var Fighters = make(map[string]*Fighter)
var FightersMutex sync.RWMutex

var FighterAttributesCache = make(map[int64]FighterAttributes)
var FighterAttributesCacheMutex sync.RWMutex

func validateFighterName(name string) error {
    if utf8.RuneCountInString(name) > 13 {
        return errors.New("Name too long")
    }

    // Check if name contains only A-Z, a-z, 0-9
    matched, err := regexp.MatchString(`^[a-zA-Z0-9]*$`, name)
    if err != nil {
        return err
    }

    if !matched {
        return errors.New("Name contains invalid characters")
    }

    return nil
}


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
            fighter.Mutex.Lock()
			fighter.IsDead = false;
			fighter.HealthAfterLastDmg = fighter.MaxHealth;
            fighter.Mutex.Unlock()
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

    fighter.Mutex.Lock()
    fighter.CurrentHealth = min(maxHealth, int64(health))
    fighter.Mutex.Unlock()

    return fighter.CurrentHealth
}

func initiateFighterRoutine(conn *websocket.Conn, fighter *Fighter) {
    log.Printf("[initiateFighterRoutine] fighter=", fighter.ID)
    speed := fighter.MovementSpeed

    msPerHit := 60000 / speed
    delay := time.Duration(msPerHit) * time.Millisecond

    for {
        conn, _ := findConnectionByFighter(fighter)

        if conn == nil {
            log.Printf("[initiateFighterRoutine] Connection closed, stopping the loop for fighter:", fighter.ID)
            removeFighterFromPopulation(fighter)
            return;
        }

       

        // Check if LastChatMsg is not empty and if 30 seconds passed from the LastChatMsgTimestamp
        currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)
        if fighter.LastChatMsg != "" && (currentTimeMillis - fighter.LastChatMsgTimestamp > 30 * 1000) {
            // Lock the fighter struct to prevent concurrent write
            fighter.Mutex.Lock()
            fighter.LastChatMsg = ""
            fighter.Mutex.Unlock()
        }

        pingFighter(fighter)
        time.Sleep(delay)
    }
}



func authFighter(conn *websocket.Conn, playerId int64, ownerAddess string, locationKey string) {
    log.Printf("[authFighter] playerId=%v ownerAddess=%v locationKey=%v conn=%v", playerId, ownerAddess, locationKey, conn)

    if playerId == 0 {
        log.Printf("[authFighter] Player id cannot be zero")
        return
    }

    strId := strconv.Itoa(int(playerId))
    stats := getFighterStats(playerId)
    atts, _ := getFighterAttributes(playerId)

    location := decodeLocation(locationKey)
    town := location[0]  


    if _, ok := Connections[conn]; ok {
        // "ok" is true if "conn" is a key in the map
        // "connValue" is the value assigned to the key "conn"
        removeFighterFromPopulation(Connections[conn].Fighter)
    }
    
    if fighter, ok := Fighters[strId]; ok {
        log.Printf("[authFighter] Fighter already exists, only update the Conn value")

        oldConn, _ := findConnectionByFighter(fighter)
        if oldConn != nil {
            ConnectionsMutex.Lock()
            delete(Connections, oldConn)
            ConnectionsMutex.Unlock()
        }

        PopulationMutex.Lock()
        if Population[town] == nil {
            Population[town] = make([]*Fighter, 0)
        }
        Population[town] = append(Population[town], fighter)
        PopulationMutex.Unlock() 

        Connections[conn].Mutex.Lock()
        Connections[conn].Fighter = Fighters[strId]
        Connections[conn].OwnerAddress = common.HexToAddress(ownerAddess)
        Connections[conn].Mutex.Unlock()

        go initiateFighterRoutine(conn, fighter)
    } else {        
        centerCoord := Coordinate{X: 5, Y: 5}
        emptySquares := getEmptySquares(centerCoord, 5, town)

        rand.Seed(time.Now().UnixNano())
        spawnCoord := emptySquares[rand.Intn(len(emptySquares))]


        fighter, err := retrieveFighterFromDB(strId)

        if err == nil {
            fighter.Backpack = NewInventory(8, 8) 
            fighter.Vault = NewInventory(8, 16) 
            fighter.Equipment = make(map[int64]*InventorySlot)

        } else {
            log.Printf("[authFighter] err=%v", err)

            fighter = &Fighter{
                ID: strId,
                TokenID: playerId,
                MaxHealth: stats.MaxHealth.Int64(),
                CurrentHealth: stats.MaxHealth.Int64(),
                Name: atts.Name,
                IsNpc: false,
                CanFight: true,
                LastDmgTimestamp: 0,
                HealthAfterLastDmg: 0,
                OwnerAddress: ownerAddess,
                MovementSpeed: 270,
                Coordinates: spawnCoord,
                Backpack: NewInventory(8, 8),
                Vault: NewInventory(8, 16), 
                Location: town,
                Strength: atts.Strength.Int64(),
                Agility: atts.Agility.Int64(),
                Energy: atts.Energy.Int64(),
                Vitality: atts.Vitality.Int64(),
                HpRegenerationRate: getHealthRegenerationRate(atts),
                Level: stats.Level.Int64(),
                Experience: atts.Experience.Int64(),
                Direction: Direction{Dx: 0, Dy: 1},
                Skills: Skills,
                Equipment: make(map[int64]*InventorySlot),
            }
        }        

        FightersMutex.Lock()
        Fighters[strId] = fighter
        FightersMutex.Unlock()

        getBackpackFromDB(fighter)
        updateFighterParams(fighter)

        PopulationMutex.Lock()
        if Population[town] == nil {
            Population[town] = make([]*Fighter, 0)
        }
        Population[town] = append(Population[town], fighter)
        PopulationMutex.Unlock() 

        Connections[conn].Mutex.Lock()
        Connections[conn].Fighter = Fighters[strId]
        Connections[conn].OwnerAddress = common.HexToAddress(ownerAddess)
        Connections[conn].Mutex.Unlock()       

        
        
        go initiateFighterRoutine(conn, fighter)
    }

    FaucetCredits(conn)

    

    getFighterItems(playerId)
}


func findFighterByConn(conn *websocket.Conn) *Fighter {
    //log.Printf("[findFighterByConn] conn=%v", conn)
    ConnectionsMutex.Lock()
    defer ConnectionsMutex.Unlock()

    return Connections[conn].Fighter
}

func addDamageToFighter(fighterID string, hitterID *big.Int, damage *big.Int) {
    found := false
    fighter := getFighterSafely(fighterID);

    log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v fighter=%v", fighterID, hitterID, damage, fighter)

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
    fighter.Mutex.Lock()
    fighter.DamageReceived = damageReceived
    fighter.Mutex.Unlock()

    //log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v fighter=%v", fighterID, hitterID, damage, fighter)
}

func updateFighterParams(fighter *Fighter) {

    fighter.Mutex.RLock()
    equipment := fighter.Equipment
    fighter.Mutex.RUnlock()


    defence := fighter.Agility/4
    damage := fighter.Strength/4 + fighter.Energy/4

    criticalDmg     := int64(0)
    excellentDmg    := int64(0)
    doubleDmg       := int64(0)
    ignoreDef       := int64(0)

    for _, item := range equipment {
        // Perform your logic with the current item and slot
        defence += item.Attributes.ItemParameters.Defense + item.Attributes.AdditionalDefense.Int64()
        damage += item.Attributes.ItemParameters.MinPhysicalDamage + item.Attributes.AdditionalDamage.Int64()

        if item.Attributes.Luck {
            criticalDmg += 5
        }

        if item.Attributes.ExcellentItemAttributes.ExcellentDamageProbabilityIncrease.Int64() > 0 {
            excellentDmg += 10
        }

        if item.Attributes.ExcellentItemAttributes.DoubleDamageChance.Int64() > 0 {
            doubleDmg += 10
        }

        if item.Attributes.ExcellentItemAttributes.IgnoreOpponentDefenseChance.Int64() > 0 {
            ignoreDef += item.Attributes.ExcellentItemAttributes.IgnoreOpponentDefenseChance.Int64()
        }
    }


    fighter.Mutex.Lock()
    fighter.Damage = damage
    fighter.Defence = defence

    fighter.CriticalDmgRate = criticalDmg
    fighter.ExcellentDmgRate = excellentDmg
    fighter.DoubleDmgRate = doubleDmg
    fighter.IgnoreDefRate = ignoreDef

    fighter.Mutex.Unlock()

    //pingFighter(fighter)

    log.Printf("[updateFighterParams] fighter=%v", fighter)
}


func applyConsumable(fighter *Fighter, item TokenAttributes) {
    switch item.Name {
        case "Small Healing Potion":
            go graduallyIncreaseHp(fighter, 100, 5)
            break

        case "Small Mana Potion":
            go graduallyIncreaseMana(fighter, 100, 5)
            break

        default:
            log.Printf("[applyConsumable] Unknown consumable=%v", item)
            break
    }
}


func graduallyIncreaseHp(fighter *Fighter, hp int64, chunks int64) {
    // Calculate how much to increase HP by each chunk
    hpIncrease := hp / chunks

    for i := int64(0); i < chunks; i++ {
        // Lock the mutex before updating fighter's HP
        fighter.Mutex.Lock()
        fighter.CurrentHealth += hpIncrease
        fighter.Mutex.Unlock()

        // Print HP for debugging purposes, remove in production code
        fmt.Println("[graduallyIncreaseHp] HP after chunk", i+1, ":", fighter.CurrentHealth)

        // Sleep for one second
        time.Sleep(1 * time.Second)
    }
}

func graduallyIncreaseMana(fighter *Fighter, mana int64, chunks int64) {
    // Calculate how much to increase HP by each chunk
    manaIncrease := mana / chunks

    for i := int64(0); i < chunks; i++ {
        // Lock the mutex before updating fighter's HP
        fighter.Mutex.Lock()
        fighter.CurrentMana += manaIncrease
        fighter.Mutex.Unlock()

        // Print HP for debugging purposes, remove in production code
        fmt.Println("[graduallyIncreaseMana] MP after chunk", i+1, ":", fighter.CurrentMana)

        // Sleep for one second
        time.Sleep(1 * time.Second)
    }
}









