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
    TokenID                 *big.Int `json:"TokenID"`
	BirthBlock   			*big.Int `json:"BirthBlock"`
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
    BirthBlock              int64               `json:"birthBlock"`
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

    sync.RWMutex
}

func (i *Fighter) gName() string {
    i.RLock()
    i.RUnlock()

    return i.Name
}

func (i *Fighter) gID() string {
    i.RLock()
    i.RUnlock()

    return i.ID
}

func (i *Fighter) gCoordinates() Coordinate {
    i.RLock()
    i.RUnlock()

    return i.Coordinates
}

func (i *Fighter) gDamageReceived() []Damage {
    i.RLock()
    i.RUnlock()

    return i.DamageReceived
}

func (i *Fighter) gMovementSpeed() int64 {
    i.RLock()
    i.RUnlock()

    return i.MovementSpeed
}

func (i *Fighter) gLastDmgTimestamp() int64 {
    i.RLock()
    i.RUnlock()

    return i.LastDmgTimestamp
}

func (i *Fighter) gMaxHealth() int64 {
    i.RLock()
    i.RUnlock()

    return i.MaxHealth
}

func (i *Fighter) gSkill() int64 {
    i.RLock()
    i.RUnlock()

    return i.Skill
}

func (i *Fighter) gTokenID() int64 {
    i.RLock()
    i.RUnlock()

    return i.TokenID
}

func (i *Fighter) gLocation() string {
    i.RLock()
    i.RUnlock()

    return i.Location
}

func (i *Fighter) gIsDead() bool {
    i.RLock()
    i.RUnlock()

    return i.IsDead
}

func (i *Fighter) gIsNpc() bool {
    i.RLock()
    i.RUnlock()

    return i.IsNpc
}


func (i *Fighter) gSpawnCoords() Coordinate {
    i.RLock()
    i.RUnlock()

    return i.SpawnCoords
}

func (i *Fighter) sDirection(v Direction) {
    i.Lock()
    defer i.Unlock()

    i.Direction = v
}

func (i *Fighter) sCoordinates(v Coordinate) {
    i.Lock()
    defer i.Unlock()

    i.Coordinates = v
}



type FighterCreatedEvent struct {
    TokenID         *big.Int            `json:"tokenId"`
    Owner           common.Address      `json:"owner"`
    
    FighterClass    string            `json:"fighterClass"`
    Name            string              `json:"name"`
}


type SafeFightersMap struct {
    Map map[string]*Fighter
    sync.RWMutex
}

var FightersMap = &SafeFightersMap{Map: make(map[string]*Fighter)}

var FighterAttributesCache = make(map[int64]FighterAttributes)
var FighterAttributesCacheMutex sync.RWMutex


func (i *SafeFightersMap) gMap() map[string]*Fighter {
    i.RLock()
    defer i.RUnlock()

    copy := make(map[string]*Fighter, len(i.Map))
    for key, val := range i.Map {
        copy[key] = val
    }
    return copy
}


func (i *SafeFightersMap) Find(id string) *Fighter {
    i.RLock()
    defer i.RUnlock()

    fighter, exists := FightersMap.Map[id]
    if exists {
        return fighter
    }

    return nil
}

func (i *SafeFightersMap) Add(id string, f *Fighter) {
    i.Lock()
    defer i.Unlock()

    FightersMap.Map[id] = f
}



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


// func getFighterSafely(id string) *Fighter {
//     FightersMutex.RLock()
//     defer FightersMutex.RUnlock()

//     fighter, exists := Fighters[id]
//     if exists {
//         return fighter
//     }

//     log.Printf("[getFighterSafely] Fighter not found id=%v", id)
//     return nil
// }


func getHealthRegenerationRate(atts FighterAttributes) (float64) {

    vitality := atts.Vitality.Int64()
    healthRegenBonus := 0

    regenRate := (float64(vitality)/float64(HealthRegenerationDivider) + float64(healthRegenBonus))/5000;
    log.Printf("[getHealthRegenerationRate] vitality=%v regenRate=%v", vitality, regenRate)
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
            fighter.Lock()
			fighter.IsDead = false;
			fighter.HealthAfterLastDmg = fighter.MaxHealth;
            fighter.Unlock()
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

    fighter.Lock()
    fighter.CurrentHealth = min(maxHealth, int64(health))
    fighter.Unlock()

    return fighter.CurrentHealth
}

func initiateFighterRoutine(conn *websocket.Conn, fighter *Fighter) {
    log.Printf("[initiateFighterRoutine] fighter=%v", fighter.ID)
    speed := fighter.MovementSpeed

    msPerHit := 60000 / speed
    delay := time.Duration(msPerHit) * time.Millisecond

    for {
        conn, _ := findConnectionByFighter(fighter)

        if conn == nil {
            log.Printf("[initiateFighterRoutine] Connection closed, stopping the loop for fighter: %v", fighter.ID)
            //removeFighterFromPopulation(fighter)
            PopulationMap.Remove(fighter)
            return;
        }

       

        // Check if LastChatMsg is not empty and if 30 seconds passed from the LastChatMsgTimestamp
        currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)
        if fighter.LastChatMsg != "" && (currentTimeMillis - fighter.LastChatMsgTimestamp > 30 * 1000) {
            // Lock the fighter struct to prevent concurrent write
            fighter.Lock()
            fighter.LastChatMsg = ""
            fighter.Unlock()
        }

        pingFighter(fighter)
        time.Sleep(delay)
    }
}



func authFighter(conn *websocket.Conn, playerId int64, ownerAddess string, locationKey string) {
    log.Printf("[authFighter] playerId=%v ownerAddess=%v locationKey=%v conn=%v", playerId, ownerAddess, locationKey, conn)

    if playerId == 0 {
        log.Printf("[authFighter] Player id cannot be zero")
        sendErrorMsgToConn(conn, "SYSTEM", "Invalid player id")
        return
    }

    strId := strconv.Itoa(int(playerId))
    stats := getFighterStats(playerId)
    atts, _ := getFighterAttributes(playerId)

    location := decodeLocation(locationKey)
    town := location[0]  

    connection := ConnectionsMap.Find(conn)
    if connection != nil {
        //removeFighterFromPopulation(connection.gFighter())
        PopulationMap.Remove(connection.gFighter())
    }
    
    fighter := FightersMap.Find(strId)
    if fighter != nil {
        log.Printf("[authFighter] Fighter already exists, only update the Conn value")

        oldConn, _ := findConnectionByFighter(fighter)
        if oldConn != nil {
            ConnectionsMap.Remove(oldConn)
        }

        // PopulationMutex.Lock()
        // if Population[town] == nil {
        //     Population[town] = make([]*Fighter, 0)
        // }
        // Population[town] = append(Population[town], fighter)
        // PopulationMutex.Unlock() 

        PopulationMap.Add(town, fighter)

        ConnectionsMap.AddWithValues(conn, fighter, common.HexToAddress(ownerAddess))

        go initiateFighterRoutine(conn, fighter)
        getFighterItems(fighter)
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
                BirthBlock: atts.BirthBlock.Int64(),
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


        FightersMap.Add(strId, fighter)
        // FightersMutex.Lock()
        // Fighters[strId] = fighter
        // FightersMutex.Unlock()

        getBackpackFromDB(fighter)
        updateFighterParams(fighter)

        // PopulationMutex.Lock()
        // if Population[town] == nil {
        //     Population[town] = make([]*Fighter, 0)
        // }
        // Population[town] = append(Population[town], fighter)
        // PopulationMutex.Unlock()  

        PopulationMap.Add(town, fighter)    

        ConnectionsMap.AddWithValues(conn, fighter, common.HexToAddress(ownerAddess))

        
        
        go initiateFighterRoutine(conn, fighter)
        getFighterItems(fighter)
    }

    //FaucetCredits(conn)

    

    
}


func findFighterByConn(conn *websocket.Conn) (*Fighter, error) {
    connection := ConnectionsMap.Find(conn)

    if connection == nil {
        return nil, fmt.Errorf("connection not found")
    }

    fighter := connection.gFighter()
    if fighter == nil {
        return nil, fmt.Errorf("connection has no fighter")
    }

    return fighter, nil
}


func addDamageToFighter(fighterID string, hitterID *big.Int, damage *big.Int) {
    found := false
    fighter := FightersMap.Find(fighterID);

    log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v", fighterID, hitterID, damage)

    damageReceived := fighter.gDamageReceived()

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

    fighter.Lock()
    fighter.DamageReceived = damageReceived
    fighter.Unlock()

    //log.Printf("[addDamageToFighter] fighterID=%v hitterID=%v damage=%v fighter=%v", fighterID, hitterID, damage, fighter)
}

func updateFighterParams(fighter *Fighter) {

    fighter.RLock()
    equipment := fighter.Equipment
    fighter.RUnlock()


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


    fighter.Lock()
    fighter.Damage = damage
    fighter.Defence = defence

    fighter.CriticalDmgRate = criticalDmg
    fighter.ExcellentDmgRate = excellentDmg
    fighter.DoubleDmgRate = doubleDmg
    fighter.IgnoreDefRate = ignoreDef

    fighter.Unlock()

    //pingFighter(fighter)

    log.Printf("[updateFighterParams] fighter=%v", fighter)
}


func applyConsumable(fighter *Fighter, item *TokenAttributes) {
    switch item.gName() {
        case "Small Healing Potion":
            go graduallyIncreaseHp(fighter, 100, 5)
            break

        case "Small Mana Potion":
            go graduallyIncreaseMana(fighter, 100, 5)
            break

        default:
            log.Printf("[applyConsumable] Unknown consumable=%v", item.gName())
            break
    }
}


func graduallyIncreaseHp(fighter *Fighter, hp int64, chunks int64) {
    // Calculate how much to increase HP by each chunk
    hpIncrease := hp / chunks

    for i := int64(0); i < chunks; i++ {
        // Lock the mutex before updating fighter's HP
        fighter.Lock()
        fighter.CurrentHealth += hpIncrease
        fighter.Unlock()

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
        fighter.Lock()
        fighter.CurrentMana += manaIncrease
        fighter.Unlock()

        // Print HP for debugging purposes, remove in production code
        fmt.Println("[graduallyIncreaseMana] MP after chunk", i+1, ":", fighter.CurrentMana)

        // Sleep for one second
        time.Sleep(1 * time.Second)
    }
}









