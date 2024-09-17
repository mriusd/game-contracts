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


func (i *TokenAttributes) UpgradeItemOption(chance float64) error {
	itemAttributes := i.GetItemAttributes()

	if itemAttributes.IsArmour {
		addOption := i.GetAdditionalDefense()
		if addOption == MAX_ADDITIONAL_OPTION {
			return fmt.Errorf("[UpgradeItemOption] Max additional option reached")
		}

		if rand.Float64() <= chance {
			err := i.IncreaseAdditionalDefense()
			if err != nil {
				return err
			}
		} else {
			i.DecreaseAdditionalDefense()
		}
	} else if itemAttributes.IsWeapon || itemAttributes.IsWings {
		addOption := i.GetAdditionalDamage()
		if addOption == MAX_ADDITIONAL_OPTION {
			return fmt.Errorf("[UpgradeItemOption] Max additional option reached")
		}

		if rand.Float64() <= chance {
			err := i.IncreaseAdditionalDamage()
			if err != nil {
				return err
			}
		} else {
			i.DecreaseAdditionalDamage()
		}
	} else {
		return fmt.Errorf("[UpgradeItemOption] Item type not upgradeable")
	}

	

	return nil
}