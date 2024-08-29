// fighters.broadcast.go

package fighters

import (
    "log"
	"encoding/json"
	"errors"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func GetJsonSerializedFighters(accountId primitive.ObjectID) (json.RawMessage, error) {
    userFighters := GetUserFighters(accountId)
    if userFighters == nil {
        userFighters = []*Fighter{}
    }
    
    log.Printf("[GetJsonSerializedFighters] userFighters=%v", userFighters)
    
    type jsonResponse struct {
        Action string `json:"action"`
        Fighters []*Fighter `json:"fighters"`        
    }

    jsonResp := jsonResponse{
        Action: "user_fighters",
        Fighters: userFighters,
    }


    jr, err := json.Marshal(jsonResp)
    if err != nil {        
        return json.RawMessage{}, errors.New("Failed serializing fighter list")
    }
    
    return jr, nil
}