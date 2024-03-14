// db.go

package db

import (
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"

    "context"
    "log"
    "time"
    "fmt"
)

var Client *mongo.Client = ConnectToDB()



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

func RemoveItemFromDB(itemId int) (bool, error) {
    log.Printf("[removeItemFromDB] itemId=%v", itemId)
    collection := Client.Database("game").Collection("items")

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

func GetNextSequenceValue(sequenceName string) (int, error) {
    sequencesCollection := Client.Database("game").Collection("sequences")
    filter := bson.M{"_id": sequenceName}
    update := bson.M{"$inc": bson.M{"seq": 1}}
    options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

    var updatedDoc struct {
        Seq int `bson:"seq"`
    }

    err := sequencesCollection.FindOneAndUpdate(context.Background(), filter, update, options).Decode(&updatedDoc)
    if err != nil {
        return 0, fmt.Errorf("getNextSequenceValue: %w", err)
    }

    return updatedDoc.Seq, nil
}







