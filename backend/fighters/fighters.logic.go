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
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mriusd/game-contracts/db"
	"github.com/mriusd/game-contracts/maps"


)



type ClassAttributes struct {
    BaseStrength int
    BaseAgility int
    BaseEnergy int
    BaseVitality int

    HpPerVitalityPoint int
    ManaPerEnergyPoint int
    HpIncreasePerLevel int
    ManaIncreasePerLevel int
    StatPointsPerLevel int
    AttackSpeed int
    AgilityPointsPerSpeed int      
}

func CreateFighter(ownerAddress, name, class  string) (*Fighter, error) {
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

	//stats := getClassStats(class)

	// record fighter to db
	fighter := &Fighter{
		Name: name,
		Class: class,
		OwnerAddress: ownerAddress,
		Location: "lorencia",
		Coordinates: maps.Coordinate{X: 10, Y: 10},
	}

	err = RecordFighterToDB(fighter)
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


func RecordFighterToDB(fighter *Fighter) error {
    // Assuming db.GetNextSequenceValue generates the next sequence value for Fighter ID
    nextID, err := db.GetNextSequenceValue("fighter")
    if err != nil {
        return err
    }

    fighter.TokenID = nextID

    collection := db.Client.Database("game").Collection("fighters")
    _, err = collection.InsertOne(context.Background(), fighter)
    if err != nil {
        return fmt.Errorf("RecordFighterToDB: %w", err)
    }

    log.Printf("[RecordFighterToDB] New fighter recorded with TokenID: %v", nextID)
    return nil
}


func GetUserFighters(ownerAddress string) []*Fighter {
    var fighters []*Fighter // Slice to store the result

    // Assuming `db.Client` is your MongoDB client instance and `fighters` is your collection
    collection := db.Client.Database("game").Collection("fighters")

    // Define the filter to match documents by ownerAddress
    filter := bson.M{"owner_address": ownerAddress}

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
        fighters = append(fighters, &fighter)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("Cursor error: %v", err)
        return nil
    }

    return fighters
}

func GetFighter(tokenId int64) (*Fighter, error) {
    // Assuming `db.Client` is your MongoDB client instance and `fighters` is your collection
    collection := db.Client.Database("game").Collection("fighters")

    var fighter Fighter
    // Create a filter to find the document by tokenId
    filter := bson.M{"tokenId": tokenId}

    // Retrieve the document
    err := collection.FindOne(context.Background(), filter).Decode(&fighter)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("No fighter found with TokenID: %d", tokenId)
        } else {
            return nil, fmt.Errorf("Error retrieving fighter with TokenID %d: %v", tokenId, err)
        }       
    }

    return &fighter, nil
}


