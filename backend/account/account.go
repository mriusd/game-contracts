// account.go
package account

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/mriusd/game-contracts/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Account struct {
	ID       int    `json:"id" bson:"id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (a *Account) RecordToDB() error {
	// Create a copy of the Account object
	copyOfAccount := *a

	// Convert the email to lowercase
	copyOfAccount.Email = strings.ToLower(copyOfAccount.Email)

	// Encode the password
	copyOfAccount.Password = encodePassword(copyOfAccount.Password)

	// Access the collection
	collection := db.Client.Database("game").Collection("accounts")

	// Check if the email is already in use
	filter := bson.M{"email": copyOfAccount.Email}
	err := collection.FindOne(context.Background(), filter).Err()
	if err == nil {
		log.Printf("[Account: RecordToDB] Email already in use: %s", copyOfAccount.Email)
		return fmt.Errorf("email already in use: %s", copyOfAccount.Email)
	}

	// Ensure the account does not have an ID, which means it's a new account
	if copyOfAccount.ID != 0 {
		return fmt.Errorf("id not zero: %d", copyOfAccount.ID)
	}

	// Insert the document
	result, err := collection.InsertOne(context.Background(), copyOfAccount)
	if err != nil {
		log.Printf("[Account: RecordToDB]: %v", err)
		return fmt.Errorf("[Account: RecordToDB]: %v", err)
	}

	// Set the ID to the inserted document's ID
	copyOfAccount.ID = result.InsertedID.(int)

	log.Printf("[Account: RecordToDB] Account Inserted with ID: %d", copyOfAccount.ID)
	return nil
}

func encodePassword(password string) string {
	hashedPassword := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashedPassword[:])
}


func CreateAccount(email, password string) (*Account, error) {
	// Create the Account struct
	account := &Account{
		Email:    email,
		Password: password,
	}

	// Call RecordToDB to save the account to the database
	err := account.RecordToDB()
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %v", err)
	}

	return account, nil
}


