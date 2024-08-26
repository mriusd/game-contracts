// fighters.broadcast.go

package fighters

import (
	"encoding/json"
	"errors"
)

func GetJsonSerializedFighters(accountId int) (json.RawMessage, error) {
	userFighters := GetUserFighters(accountId)

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