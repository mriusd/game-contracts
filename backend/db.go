package main

import (
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"

    "context"
    "log"
    "time"
    "encoding/json"
)

var client *mongo.Client = ConnectToDB()

// func recordItemToDB(item ItemAttributes) {
// 	coll := client.Database("game").Collection("items")

// 	// create a filter for the query
// 	filter := bson.M{"tokenId": item.TokenId}

// 	// create an update document containing the values to update
// 	update := bson.M{
// 	    "$set": item,
// 	}

// 	// create options for the query
// 	opts := options.Update().SetUpsert(true)

// 	// execute the update query
// 	result, err := coll.UpdateOne(context.Background(), filter, update, opts)
// 	if err != nil {
// 	    log.Printf("[recordItemToDB] error: %v\n", err, result)
// 	}
// }

func updateFighterDB(fighter *Fighter) {
    collection := client.Database("game").Collection("fighters")

    // Marshalling the fighter object to JSON
    jsonFighter, err := json.Marshal(fighter)
    if err != nil {
        log.Fatal(err)
    }

    filter := bson.D{{"fighterID", fighter.ID}}
    update := bson.D{
        {"$set", bson.D{
            {"fighterID", fighter.ID},
            {"atts", string(jsonFighter)},
        }},
    }

    upsert := true
    opt := options.UpdateOptions{
        Upsert: &upsert,
    }

    _, err = collection.UpdateOne(context.Background(), filter, update, &opt)
    if err != nil {
        log.Fatal(err)
    }
}

func retrieveFighterFromDB(fighterID string) (*Fighter, error) {
    log.Printf("[retrieveFighterFromDB] fighterID=%v", fighterID)
    collection := client.Database("game").Collection("fighters")

    filter := bson.D{{"fighterID", fighterID}}
    var result struct {
        FighterID string `bson:"fighterID"`
        Atts      string `bson:"atts"`
    }

    err := collection.FindOne(context.Background(), filter).Decode(&result)
    if err != nil {
        return nil, err
    }

    var fighter Fighter
    err = json.Unmarshal([]byte(result.Atts), &fighter)
    if err != nil {
        return nil, err
    }
    log.Printf("[retrieveFighterFromDB] fighter=%v", fighter)
    return &fighter, nil
}

func getBackpackFromDB(fighter *Fighter) (bool) {
    collection := client.Database("game").Collection("backpacks")

    fighter.Mutex.RLock()
    filter := bson.M{"fighterId": fighter.TokenID}
    fighter.Mutex.RUnlock()

    var result bson.M
    err := collection.FindOne(context.Background(), filter).Decode(&result)

    if err != nil {
        log.Printf("[getBackpackFromDB] Error getting backpack from database: %v", err)
        return false
    }

    backpackStr, ok := result["backpack"].(string)
    if !ok {
        log.Printf("[getBackpackFromDB] Error asserting backpack as string")
        return false
    }

    var backpack Backpack
    err = json.Unmarshal([]byte(backpackStr), &backpack)
    if err != nil {
        log.Printf("[getBackpackFromDB] Error unmarshaling backpack: %v", err)
        return false
    }

    equipmentStr, ok := result["equipment"].(string)
    if !ok {
        log.Printf("[getBackpackFromDB] Error asserting equipment as string")
        return false
    }

    var equipment map[int64]*BackpackSlot
    err = json.Unmarshal([]byte(equipmentStr), &equipment)
    if err != nil {
        log.Printf("[getBackpackFromDB] Error unmarshaling equipment: %v", err)
        return false
    }

    log.Printf("[getBackpackFromDB] backpack=%v equipment=%v", backpack, equipment)
    fighter.Mutex.Lock()
    fighter.Backpack = &backpack
    fighter.Equipment = equipment
    fighter.Mutex.Unlock()

    return true;
}

func saveBackpackToDB(fighter *Fighter) error {
    log.Printf("[saveBackpackToDB] fighter=%v", fighter)
    collection := client.Database("game").Collection("backpacks")

    fighter.Mutex.RLock()
    backpackJSON, err := json.Marshal(fighter.Backpack)
    if err != nil {
        log.Printf("[saveBackpackToDB] Error marshaling backpack: %v", err)
        return err
    }

    backpackStr := string(backpackJSON)

    equipmentJSON, err := json.Marshal(fighter.Equipment)
    if err != nil {
        log.Printf("[saveBackpackToDB] Error marshaling backpack: %v", err)
        return err
    }
    filter := bson.M{"fighterId": fighter.TokenID}
    equipmentStr := string(equipmentJSON)
    fighter.Mutex.RUnlock()

    
    update := bson.M{"$set": bson.M{"backpack": backpackStr, "equipment": equipmentStr}}
    opts := options.Update().SetUpsert(true)

    _, err = collection.UpdateOne(context.Background(), filter, update, opts)
    if err != nil {
        log.Printf("[saveBackpackToDB] Error updating backpack in database: %v", err)
        return err
    }

    return nil
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

func getItemAttributesFromDB(itemId int64) (ItemAttributes, bool) {
    collection := client.Database("game").Collection("items")

    var itemWithAttributes struct {
        Attributes string `bson:"attributes"`
    }

    filter := bson.M{"tokenId": itemId}
    err := collection.FindOne(context.Background(), filter).Decode(&itemWithAttributes)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return ItemAttributes{}, false
        }
        log.Fatal("[getItemAttributesFromDB] ", err)
    }

    var item ItemAttributes
    err = json.Unmarshal([]byte(itemWithAttributes.Attributes), &item)
    if err != nil {
        log.Fatal("[getItemAttributesFromDB] JSON unmarshal error: ", err)
    }

    return item, true
}

func removeItemFromDB(itemId int64) (bool, error) {
    log.Printf("[removeItemFromDB] itemId=%v", itemId)
    collection := client.Database("game").Collection("items")

    filter := bson.M{"tokenId": itemId}
    result, err := collection.DeleteOne(context.Background(), filter)

    if err != nil {
        log.Printf("[removeItemFromDB] Error: %v", err)
        return false, err
    }

    if result.DeletedCount == 0 {
        log.Printf("[removeItemFromDB] Item with tokenId %d not found", itemId)
        return false, nil
    }

    log.Printf("[removeItemFromDB] Successfully removed item with tokenId %d", itemId)
    return true, nil
}

func saveItemAttributesToDB(item ItemAttributes) {
    log.Printf("[saveItemAttributesToDB] item=%v", item)
    collection := client.Database("game").Collection("items")

    jsonData, _ := json.Marshal(item)
    filter := bson.M{"tokenId": item.TokenId.Int64()}
    update := bson.M{"$set": bson.M{"attributes": string(jsonData)}}
    opts := options.Update().SetUpsert(true)
    _, err := collection.UpdateOne(context.Background(), filter, update, opts)

    if err != nil {
        log.Fatal("[saveItemAttributesToDB] ", err)
    }
}





