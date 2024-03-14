// drop.loader.go

package drop

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type DropParams struct {
    WeaponsDropRate int `json:"weaponsDropRate"`
    ArmoursDropRate int `json:"armoursDropRate"`
    JewelsDropRate  int `json:"jewelsDropRate"`
    MiscDropRate    int `json:"miscDropRate"`
    BoxDropRate     int `json:"boxDropRate"`

    LuckDropRate 	int `json:"luckDropRate"`
    SkillDropRate   int `json:"SkillDropRate"`

    ExcDropRate   	int `json:"excDropRate"`
    BoxId         	int `json:"boxId"`

    MinItemLevel 	int `json:"minItemLevel"`
    MaxItemLevel 	int `json:"maxItemLevel"`
    MaxAddPoints 	int `json:"maxAddPoints"`
}

var DropParamsMobMap = make(map[int] DropParams)
var DropParamsBoxMap = make(map[int] DropParams)

func LoadDropParamsMob() {
	log.Printf("[LoadDropParamsMob]")
	file, err := ioutil.ReadFile("./drop/dropparams_mobs.json")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	var params map[int] DropParams
	err = json.Unmarshal(file, &params)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	DropParamsMobMap = params
}

func LoadDropParamsBox() {
	log.Printf("[LoadDropParamsBox]")
	file, err := ioutil.ReadFile("./drop/dropparams_mobs.json")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	var params map[int] DropParams
	err = json.Unmarshal(file, &params)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	DropParamsBoxMap = params
}