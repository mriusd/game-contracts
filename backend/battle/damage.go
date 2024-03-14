// damage.go

package battle 

import (
	
)

type Damage struct {
    FighterId        string
    Damage           int
}

type DamageType struct {
    IsCritical          bool `json:"isCritical"`
    IsExcellent         bool `json:"isExcellent"`
    IsDouble            bool `json:"isDouble"`
    IsIgnoreDefence     bool `json:"isIgnoreDefence"`
}