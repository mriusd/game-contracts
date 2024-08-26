// fighters.logic.go

package fighters

import (
	"errors"
	"log"
	"fmt"	
	"unicode/utf8"
	"regexp"
	"context"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mriusd/game-contracts/db"
    "github.com/mriusd/game-contracts/maps"
    "github.com/mriusd/game-contracts/skill"
	"github.com/mriusd/game-contracts/inventory"
)


func (i *Fighter) AddStrength(v int) error {
    if v > i.AvailableStats {
        return errors.New("Not enought available stats")
    }

    i.SetStrength(i.GetNetStrength() + v)
    return nil
}

func (i *Fighter) AddAgility(v int) error {
    if v > i.AvailableStats {
        return errors.New("Not enought available stats")
    }

    i.SetAgility(i.GetNetAgility() + v)
    return nil
}

func (i *Fighter) AddEnergy(v int) error {
    if v > i.AvailableStats {
        return errors.New("Not enought available stats")
    }

    i.SetEnergy(i.GetNetEnergy() + v)
    return nil
}

func (i *Fighter) AddVitality(v int) error {
    if v > i.AvailableStats {
        return errors.New("Not enought available stats")
    }

    i.SetVitality(i.GetNetVitality() + v)
    return nil
}

func (i *Fighter) BindSkill(skill, key int) error {
    if key < 1 || key > 5 {
        return errors.New("Only key 1 to 5 can be used for bindings")
    }

    skills := i.GetSkills()

    skillObj, exists := skills[skill]
    if !exists {
        return errors.New("Skill not found for character")
    }

    skillBindings := i.GetSkillBindings()

    skillBindings[key] = &skillObj
    i.SetSkillBindings(skillBindings)

    return nil
}

func (i *Fighter) ConsumableBind(binding, key string) error {
    if binding != "hp" && binding != "mana" {
        return errors.New("Unknown binding")
    }

    if key != "Q" && key != "W" && key != "E" && key != "R" {
        return errors.New("Invalid key")
    }

    consumableBindings := i.GetConsumableBindings()

    consumableBindings[key] = binding
    i.SetConsumableBindings(consumableBindings)

    return nil
}

func CreateFighter(accountId int, name, class  string) (*Fighter, error) {
	err := validateFighterName(name)
	if err != nil {
		return nil, err
	}

	if !CheckNameAvailable(name) {
		return nil, fmt.Errorf("Name taken")
	}

	// check name already taken

	// validate class
	if !validateClass(class) {
		return nil, fmt.Errorf("Invalid class=%v", class)
	}

	stats := getClassStats(class)

    skillBindings := make(map[int]*skill.Skill)
    for i := 1; i <= 5; i++ {
        skillBindings[i] = nil
    }


    consumableBindings := map[string]string{
        "Q": "hp",
        "W": "mana",
        "E": "speed",
        "R": "cure",
    }

	// record fighter to db
	fighter := &Fighter{
		Name: name,
		Class: class,
		AccountID: accountId,
		Location: "lorencia",
		Coordinates: maps.Coordinate{X: 10, Y: 10},
        Strength: stats.BaseStrength,
        Agility: stats.BaseAgility,
        Energy: stats.BaseEnergy,
        Vitality: stats.BaseVitality,
        SkillBindings: skillBindings,
        ConsumableBindings: consumableBindings,
	}

	err = RecordNewFighterToDB(fighter)
	if err != nil {
		return nil, fmt.Errorf("Failed to create fighter, error=%v", err)
	}

	return fighter, nil
}


func validateFighterName(name string) error {
    if utf8.RuneCountInString(name) > 13 {
        return errors.New("Name too long")
    }

    // Check if name contains only A-Z, a-z, 0-9
    matched, err := regexp.MatchString(`^[a-zA-Z0-9]*$`, name)
    if err != nil {
        return fmt.Errorf("Error validating name: %v", err)
    }

    if !matched {
        return errors.New("Name contains forbidden characters")
    }

    return nil
}

func validateClass(class string) bool {
    return class == "Warrior" || class == "Wizard"
}




func CheckNameAvailable(name string) bool {
    // Assuming `db` is your MongoDB client instance and `fighters` is your collection
    collection := db.Client.Database("game").Collection("fighters")

    // Create a case-insensitive collation
    collation := options.Collation{Locale: "en", Strength: 2}

    // Count documents with the matching name, case-insensitive
    count, err := collection.CountDocuments(context.TODO(), bson.M{"name": name}, options.Count().SetCollation(&collation))
    if err != nil {
        // Handle error (e.g., log it, return false, throw an error, etc.)
        // For simplicity, returning false here, but you might want to handle this more gracefully
        return false
    }

    // Name is available if count is 0 (no documents found with the same name)
    return count == 0
}


func RecordNewFighterToDB(fighter *Fighter) error {
    // Assuming db.GetNextSequenceValue generates the next sequence value for Fighter ID
    nextID, err := db.GetNextSequenceValue("fighter")
    if err != nil {
        return err
    }

    nextID += 100000;

    fighter.TokenID = nextID

    collection := db.Client.Database("game").Collection("fighters")
    _, err = collection.InsertOne(context.Background(), fighter)
    if err != nil {
        return fmt.Errorf("RecordFighterToDB: %w", err)
    }

    log.Printf("[RecordFighterToDB] New fighter recorded with TokenID: %v", nextID)
    return nil
}


func GetUserFighters(accountId int) []*Fighter {
    var fighters []*Fighter // Slice to store the result

    // Assuming `db.Client` is your MongoDB client instance and `fighters` is your collection
    collection := db.Client.Database("game").Collection("fighters")

    // Define the filter to match documents by account_id
    filter := bson.M{"account_id": accountId}

    // Find returns a cursor for multiple documents
    cursor, err := collection.Find(context.Background(), filter)
    if err != nil {
        log.Printf("Error finding user fighters: %v", err)
        return nil
    }
    defer cursor.Close(context.Background())

    // Iterate through the cursor
    for cursor.Next(context.Background()) {
        var fighter Fighter
        err := cursor.Decode(&fighter)
        if err != nil {
            log.Printf("Error decoding fighter: %v", err)
            continue // Skip this iteration
        }        

        equipment, err := inventory.GetEquipmentFromDB(fighter.TokenID)
        // if err != nil {
        //     return nil, fmt.Errorf("Failed to load equipment: %v", err)
        // } 

        fighter.Equipment = equipment
        fighters = append(fighters, &fighter)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("Cursor error: %v", err)
        return nil
    }

    return fighters
}

// func GetFighter(tokenId int) (*Fighter, error) {
//     // Assuming `db.Client` is your MongoDB client instance and `fighters` is your collection
//     collection := db.Client.Database("game").Collection("fighters")

//     var fighter Fighter
//     // Create a filter to find the document by tokenId
//     filter := bson.M{"tokenId": tokenId}

//     // Retrieve the document
//     err := collection.FindOne(context.Background(), filter).Decode(&fighter)
//     if err != nil {
//         if err == mongo.ErrNoDocuments {
//             return nil, fmt.Errorf("No fighter found with TokenID: %d", tokenId)
//         } else {
//             return nil, fmt.Errorf("Error retrieving fighter with TokenID %d: %v", tokenId, err)
//         }       
//     }

//     return &fighter, nil
// }



