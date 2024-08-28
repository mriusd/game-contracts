// http.go

package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/mriusd/game-contracts/account"
)


type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
    // Handle CORS preflight request
    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.WriteHeader(http.StatusOK)
        return
    }

    // Handle the actual POST request
    if r.Method == http.MethodPost {
        // Your existing registration logic here
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        // Process the registration and send the response
        // return
    }

    // Only accept POST requests
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Parse the JSON request body
    var req RegisterRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    log.Printf("[handleRegister] email=%v", req.Email)

    // Create the account
    acc, err := account.CreateAccount(req.Email, req.Password)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to create account: %v", err), http.StatusInternalServerError)
        return
    }

    // Respond with the created account's ID
    //w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Account created successfully",
        "id":      acc.ID,
    })
}


func handleLogin(w http.ResponseWriter, r *http.Request) {
    // Handle CORS preflight request
    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.WriteHeader(http.StatusOK)
        return
    }

    // Handle the actual POST request
    if r.Method == http.MethodPost {
        // Your existing registration logic here
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        // // Process the registration and send the response
        // return
    }

    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var req RegisterRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    acc, err := account.Login(req.Email, req.Password)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to login: %v", err), http.StatusUnauthorized)
        return
    }

    // Create a new session
    sessionId, err := account.CreateSession(acc.ID)
    if err != nil {
        http.Error(w, "Failed to create session", http.StatusInternalServerError)
        return
    }

    log.Printf("[handleLogin] email=%v", req.Email)


    // Respond with the session ID
    //w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":   "Login successful",
        "accountId": acc.ID,
        "sessionId": sessionId,
    })
}


