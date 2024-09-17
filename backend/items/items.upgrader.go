// items.upgrader.go

package items

import (
	"fmt"
	"math/rand"
)

func (i *TokenAttributes) UpgradeItemLevel(chance float64) error {
	if i.GetItemLevel() >= 9 {
		return fmt.Errorf("[UpgradeItemLevel] Item level must be less than 9")
	}

	if rand.Float64() <= chance {
		err := i.IncreaseItemLevel()
		if err != nil {
			return err
		}
	} else {
		i.DecreaseItemLevel()
	}

	return nil
}