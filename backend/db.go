package main

import (
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"

    "context"
    "log"
    "time"
)

var client *mongo.Client = ConnectToDB()

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

    var item ItemAttributes
    filter := bson.M{"tokenId": itemId}
    err := collection.FindOne(context.Background(), filter).Decode(&item)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return ItemAttributes{}, false
        }
        log.Fatal("[getItemAttributesFromDB] ",err)
    }

    return item, true
}

func saveItemAttributesToDB(item ItemAttributes) {
    collection := client.Database("game").Collection("items")

    filter := bson.M{"tokenId": item.TokenId}
    update := bson.M{"$set": item}

    opts := options.Update().SetUpsert(true)
    _, err := collection.UpdateOne(context.Background(), filter, update, opts)
    if err != nil {
        log.Fatal("[printAllItemsInDB] ",err)
    }
}