// fighterModel.go

package fighters

import (
	"sync"	
    "math"
    "context"
    "log"
    "fmt"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"


    "github.com/mriusd/game-contracts/battle"
    "github.com/mriusd/game-contracts/maps"
    "github.com/mriusd/game-contracts/inventory"
    "github.com/mriusd/game-contracts/skill"
    "github.com/mriusd/game-contracts/db"
)

type Fighter struct {
	ID    					string  		        `json:"id" bson:"-"`
    Class                   string                  `json:"class" bson:"class"`
    MaxHealth     			int 			        `json:"maxHealth" bson:"-"`
    Name           			string 			        `json:"name" bson:"name"`
    IsNpc         			bool 			        `json:"isNpc" bson:"-"`
    IsDead         			bool 			        `json:"isDead" bson:"-"`
    CanFight 				bool 			        `json:"canFight" bson:"-"`
    
    HealthAfterLastDmg 		int  			        `json:"healthAfterLastDmg" bson:"healthAfterLastDmg"`

    TokenID                 int                   `json:"tokenId" bson:"tokenId"`
    Location                string                  `json:"location" bson:"location"`
    
    DamageReceived          []battle.Damage         `json:"damageDealt" bson:"-"`
    OwnerAddress            string                  `json:"ownerAddress" bson:"owner_address"`
    Coordinates             maps.Coordinate         `json:"coordinates" bson:"coordinates"`
    MovementSpeed           int                     `json:"movementSpeed" bson:"-"` // squares per minute
    Skill                   int                     `json:"skill" bson:"skill"`
    SpawnCoords             maps.Coordinate         `json:"spawnCoords" bson:"-"`
    




    // Fighter stats
    Strength                int                     `json:"strength" bson:"strength"`
    Agility                 int                     `json:"agility" bson:"agility"`
    Energy                  int                     `json:"energy" bson:"energy"`
    Vitality                int                     `json:"vitality" bson:"vitality"`
    AvailableStats          int                     `json:"available_stats" bson:"-"`


    // Fighter dynamic paramters
    CurrentHealth           int                     `json:"currentHealth" bson:"-"`
    CurrentMana             int                     `json:"currentMana" bson:"-"`


    // Fighter parameters with equipped items
    Damage                  int                     `json:"damage" bson:"-"`
    Defence                 int                     `json:"defence" bson:"-"`
    AttackSpeed             int                     `json:"attackSpeed" bson:"-"` 
    HpRegenerationRate      float64                 `json:"hpRegenerationRate" bson:"-"`
    HpRegenerationBonus     float64                 `json:"hpRegenerationBonus" bson:"-"`

    // Damage type rates
    CriticalDmgRate         int                     `json:"criticalDmgRate" bson:"-"`
    ExcellentDmgRate        int                     `json:"excellentDmgRate" bson:"-"`
    DoubleDmgRate           int                     `json:"doubleDmgRate" bson:"-"`
    IgnoreDefRate           int                     `json:"ignoreDefRate" bson:"-"`


    Level                   int                     `json:"level" bson:"level"`
    LevelProgress           int                     `json:"level_progress" bson:"level_progress"`
    Experience              int                     `json:"experience" bson:"experience"`

    Direction               maps.Direction          `json:"direction" bson:"-"`

    Skills                  map[int]skill.Skill     `json:"skills" bson:"skills"`
    Backpack                *inventory.Inventory    `json:"backpack" bson:"-"`
    Vault                   *inventory.Inventory    `json:"-" bson:"-"`
    Equipment               *inventory.Equipment    `json:"equipment" bson:"-"`

    LastChatMsg             string                  `json:"lastChatMessage" bson:"-"`

    LastDmgTimestamp        int                     `json:"lastDmgTimestamp" bson:"lastDmgTimestamp"`
    LastMoveTimestamp       int                     `json:"lastMoveTimestamp" bson:"-"` // milliseconds
    LastChatMsgTimestamp    int                     `json:"lastChatMsgTimestamp" bson:"-"`

    Credits                 int                     `json:"credits" bson:"-"`

    sync.RWMutex                                    `json:"-" bson:"-"`
}


func (i *Fighter) SetStrength(v int) {
    i.Lock()
    i.Strength = v
    i.Unlock()

    i.RecordToDB()
}

func (i *Fighter) SetAgility(v int) {
    i.Lock()
    i.Agility = v
    i.Unlock()

    i.RecordToDB()
}

func (i *Fighter) SetEnergy(v int) {
    i.Lock()
    i.Energy = v
    i.Unlock()

    i.RecordToDB()
}

func (i *Fighter) SetVitality(v int) {
    i.Lock()
    i.Vitality = v
    i.Unlock()

    i.RecordToDB()
}


func (i *Fighter) GetVault() *inventory.Inventory {
    i.RLock()
    defer i.RUnlock()

    return i.Vault
}

func (i *Fighter) GetLevel() int {
    i.RLock()
    defer i.RUnlock()

    return i.Level
}

func (i *Fighter) GetBackpack() *inventory.Inventory {
    i.RLock()
    defer i.RUnlock()

    return i.Backpack
}


func (i *Fighter) GetName() string {
    i.RLock()
    defer i.RUnlock()

    return i.Name
}


func (i *Fighter) GetOwnerAddress() string {
    i.RLock()
    defer i.RUnlock()

    return i.OwnerAddress
}

func (i *Fighter) GetID() string {
    i.RLock()
    defer i.RUnlock()

    return i.ID 
}

func (i *Fighter) GetCoordinates() maps.Coordinate {
    i.RLock()
    defer i.RUnlock()

    return i.Coordinates
}

func (i *Fighter) GetDamageReceived() []battle.Damage {
    i.RLock()
    defer i.RUnlock()

    return i.DamageReceived
}

func (i *Fighter) GetHealthAfterLastDmg() int {
    i.RLock()
    defer i.RUnlock()

    return i.HealthAfterLastDmg
}

func (i *Fighter) GetLastChatMsgTimestamp() int {
    i.RLock()
    defer i.RUnlock()

    return i.LastChatMsgTimestamp
}

func (i *Fighter) GetHpRegenerationRate() float64 {
    i.RLock()
    defer i.RUnlock()

    return i.HpRegenerationRate
}

func (i *Fighter) GetMovementSpeed() int {
    i.RLock()
    defer i.RUnlock()

    return i.MovementSpeed
}

func (i *Fighter) GetLastDmgTimestamp() int {
    i.RLock()
    defer i.RUnlock()

    return i.LastDmgTimestamp
}

func (i *Fighter) GetMaxHealth() int {
    i.RLock()
    defer i.RUnlock()

    return i.MaxHealth
}

func (i *Fighter) CalcMaxHealth() int {
    i.RLock()
    defer i.RUnlock()

    classStats := getClassStats(i.Class)
    hpPerVitalityPoint := classStats.HpPerVitalityPoint
    increasePerLevel := classStats.HpIncreasePerLevel

    return increasePerLevel * i.GetLevel() + hpPerVitalityPoint * i.GetVitality()
}

func (i *Fighter) GetDefence() int {
    i.RLock()
    defer i.RUnlock()

    return i.Defence
}

func (i *Fighter) GetDamage() int {
    i.RLock()
    defer i.RUnlock()

    return i.Damage
}

func (i *Fighter) GetStrength() int {
    i.RLock()
    defer i.RUnlock()

    return i.Strength + getClassStats(i.Class).BaseStrength
}

func (i *Fighter) GetAgility() int {
    i.RLock()
    defer i.RUnlock()

    return i.Agility + getClassStats(i.Class).BaseAgility
}

func (i *Fighter) GetEnergy() int {
    i.RLock()
    defer i.RUnlock()

    return i.Energy + getClassStats(i.Class).BaseEnergy
}

func (i *Fighter) GetVitality() int {
    i.RLock()
    defer i.RUnlock()

    return i.Vitality + getClassStats(i.Class).BaseVitality
}

func (i *Fighter) GetIgnoreDefRate() int {
    i.RLock()
    defer i.RUnlock()

    return i.IgnoreDefRate
}

func (i *Fighter) GetExcellentDmgRate() int {
    i.RLock()
    defer i.RUnlock()

    return i.ExcellentDmgRate
}

func (i *Fighter) GetCriticalDmgRate() int {
    i.RLock()
    defer i.RUnlock()

    return i.CriticalDmgRate
}

func (i *Fighter) GetDoubleDmgRate() int {
    i.RLock()
    defer i.RUnlock()

    return i.DoubleDmgRate
}

func (i *Fighter) GetSkill() int {
    i.RLock()
    defer i.RUnlock()

    return i.Skill
}

func (i *Fighter) GetTokenID() int {
    i.RLock()
    defer i.RUnlock()

    return i.TokenID
}

func (i *Fighter) GetLocation() string {
    i.RLock()
    defer i.RUnlock()

    return i.Location
}

func (i *Fighter) GetIsDead() bool {
    i.RLock()
    defer i.RUnlock()

    return i.IsDead
}

func (i *Fighter) GetIsNpc() bool {
    i.RLock()
    defer i.RUnlock()

    return i.IsNpc
}

func (i *Fighter) GetDirection() maps.Direction {
    i.RLock()
    defer i.RUnlock()

    return i.Direction
}


func (i *Fighter) GetSpawnCoords() maps.Coordinate {
    i.RLock()
    defer i.RUnlock()

    return i.SpawnCoords
}

func (i *Fighter) GetEquipment() *inventory.Equipment {
    i.RLock()
    defer i.RUnlock()

    return i.Equipment
}

func (i *Fighter) SetDirection(v maps.Direction) {
    i.Lock()
    defer i.Unlock()

    i.Direction = v
}

func (i *Fighter) SetCoordinates(v maps.Coordinate) {
    i.Lock()
    defer i.Unlock()

    i.Coordinates = v
}

func (i *Fighter) SetCurrentHealth(v int) {
    i.Lock()
    defer i.Unlock()

    i.CurrentHealth = v
}

func (i *Fighter) SetIsDead(v bool) {
    i.Lock()
    defer i.Unlock()

    i.IsDead = v
}

func (i *Fighter) SetLastDmgTimestamp(v int) {
    i.Lock()
    defer i.Unlock()

    i.LastDmgTimestamp = v
}

func (i *Fighter) SetHealthAfterLastDmg(v int) {
    i.Lock()
    defer i.Unlock()

    i.HealthAfterLastDmg = v
}

func (i *Fighter) SetLastChatMsg(v string) {
    i.Lock()
    defer i.Unlock()

    i.LastChatMsg = v
}

func (i *Fighter) SetEquipment(v *inventory.Equipment) {
    i.Lock()
    defer i.Unlock()

    i.Equipment = v
}

func (i *Fighter) AddExperience(v int) {    
    i.Lock()
    i.Experience += v    
    i.Unlock()

    lvl := i.CalcLevel()

    i.Lock()
    i.Level = lvl
    i.Unlock()

    progress := i.CalcLevelProgress()

    i.Lock()
    i.LevelProgress = progress
    i.Unlock()

    availableStats := i.GetAvailableStats()
    i.Lock()
    i.AvailableStats = availableStats
    i.Unlock()

    i.RecordToDB()
}


type ClassAttributes struct {
    BaseStrength            int
    BaseAgility             int
    BaseEnergy              int
    BaseVitality            int

    HpPerVitalityPoint      int
    ManaPerEnergyPoint      int
    HpIncreasePerLevel      int
    ManaIncreasePerLevel    int
    StatPointsPerLevel      int
    AttackSpeed             int
    AgilityPointsPerSpeed   int      
}

func getClassStats(class string) ClassAttributes {
    switch class {
        case "Warrior": 
            return ClassAttributes{
                BaseStrength: 42, 
                BaseAgility: 21,  
                BaseEnergy: 5, 
                BaseVitality: 20,  
                HpPerVitalityPoint: 5,  
                ManaPerEnergyPoint: 5, 
                HpIncreasePerLevel: 5, 
                ManaIncreasePerLevel: 1, 
                StatPointsPerLevel: 5, 
                AttackSpeed: 27, 
                AgilityPointsPerSpeed: 7,
            }

        case "Wizard":
            return ClassAttributes{
                BaseStrength: 15, 
                BaseAgility: 20,  
                BaseEnergy: 50, 
                BaseVitality: 20,  
                HpPerVitalityPoint: 3,  
                ManaPerEnergyPoint: 10, 
                HpIncreasePerLevel: 3, 
                ManaIncreasePerLevel: 5, 
                StatPointsPerLevel: 5, 
                AttackSpeed: 16, 
                AgilityPointsPerSpeed: 5,
            }
    }

    return ClassAttributes{}
}




func (i *Fighter) CalcLevel() int {
    i.RLock()
    defer i.RUnlock()
    return (sqrtint((5 * i.Experience) + 125) - 5) / 10;
}

func (i *Fighter) GetAvailableStats() int {
    i.RLock()
    defer i.RUnlock()
    
    classStats := getClassStats(i.Class)
    earnedStats := i.CalcLevel() * classStats.StatPointsPerLevel
    usedStats := i.Strength + i.Agility + i.Energy + i.Vitality 

    return earnedStats - usedStats
}

func (i *Fighter) CalcLevelProgress() int {
    i.RLock()
    defer i.RUnlock()

    currentLevel := (sqrtint((5 * i.Experience) + 125) - 5) / 10

    // Correct the formulas for currentLevelExp and nextLevelExp
    nextLevelExp := ((10*(currentLevel+1) + 5) * (10*(currentLevel+1) + 5) - 125) / 5
    currentLevelExp := ((10*currentLevel + 5) * (10*currentLevel + 5) - 125) / 5

    log.Printf("[CalcLevelProgress] currentLevel=%v nextLevelExp=%v currentLevelExp=%v", currentLevel, nextLevelExp, currentLevelExp)

    progress := ((i.Experience - currentLevelExp) * 100) / (nextLevelExp - currentLevelExp)
    return progress
}


func sqrtint(x int) int {
    // Convert the int to float64 to use math.Sqrt
    result := math.Sqrt(float64(x))
    // Convert back to int. This step truncates the decimal part.
    return int(result)
}


func GetFromDB(fighterId int) (*Fighter, error) {
    filter := bson.M{"tokenId": fighterId}
    collection := db.Client.Database("game").Collection("fighters")

    var fighter Fighter
    err := collection.FindOne(context.Background(), filter).Decode(&fighter)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("[Fighters: GetFromDB] Fighter not found in db %v", fighterId) 
        }
        return nil, fmt.Errorf("[Fighters: GetFromDB] fighterId=%v err=%v", fighterId, err)
    }

    log.Printf("[Fighters: GetFromDB] Fighter found=%v", fighter)
    return &fighter, nil
}


func (i *Fighter) RecordToDB() error {
    i.RLock()
    copy := *i 
    i.RUnlock()

    filter := bson.M{"tokenId": copy.TokenID}
    update := bson.M{"$set": copy}
    options := options.Update().SetUpsert(true)

    collection := db.Client.Database("game").Collection("fighters")
    _, err := collection.UpdateOne(context.Background(), filter, update, options)
    if err != nil {
        log.Printf("[Fighters: RecordToDB]: %w", err)
        return fmt.Errorf("[Fighters: RecordToDB]: %w", err)
    }

    log.Printf("[Fighters: RecordToDB] Fighter Updated")
    return nil
}



















