// damage.go

package battle 

import (
	"math/big"
)

type Damage struct {
    FighterId        *big.Int
    Damage           *big.Int
}

type DamageType struct {
    IsCritical          bool `json:"isCritical"`
    IsExcellent         bool `json:"isExcellent"`
    IsDouble            bool `json:"isDouble"`
    IsIgnoreDefence     bool `json:"isIgnoreDefence"`
}