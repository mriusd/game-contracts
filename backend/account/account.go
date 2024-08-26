// account.go
package account

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"
	"sync"

	"github.com/google/uuid"
	"github.com/mriusd/game-contracts/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account struct {
	ID       int    `json:"id" bson:"id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Session struct {
	AccountID int
	ExpiresAt time.Time
}

var sessionStore = struct {
	sync.RWMutex
	sessions map[string]Session
}{sessions: make(map[string]Session)}

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

// Login attempts to find an account with the given email and password
func Login(email, password string) (*Account, error) {
	// Convert the email to lowercase
	email = strings.ToLower(email)

	// Access the accounts collection
	collection := db.Client.Database("game").Collection("accounts")

	// Find the account by email
	var foundAccount Account
	filter := bson.M{"email": email}
	err := collection.FindOne(context.Background(), filter).Decode(&foundAccount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to find account: %v", err)
	}

	// Encode the provided password
	encodedPassword := encodePassword(password)

	// Compare the encoded passwords
	if foundAccount.Password != encodedPassword {
		return nil, fmt.Errorf("invalid email or password")
	}

	return &foundAccount, nil
}

// CreateSession generates a new session ID for the account and stores it
func CreateSession(accountID int) (string, error) {
	// Generate a unique session ID
	sessionId := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour) // Session expires in 24 hours

	// Store the session
	sessionStore.Lock()
	sessionStore.sessions[sessionId] = Session{
		AccountID: accountID,
		ExpiresAt: expiresAt,
	}
	sessionStore.Unlock()

	return sessionId, nil
}

// ValidateSession checks if the provided session ID is valid and returns the associated session
func ValidateSession(sessionId string) (*Session, error) {
	sessionStore.RLock()
	session, exists := sessionStore.sessions[sessionId]
	sessionStore.RUnlock()

	if !exists || session.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("invalid or expired session")
	}

	return &session, nil
}

func CheckSessionExpired(session *Session) (error) {
	if session.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("Session expired. Please login.")
	}

	return nil
}


