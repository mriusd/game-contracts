package main

import (
    "log"
    "os"
    "strconv"

    "github.com/mriusd/game-contracts/items"
)

type PriceList struct {
    WeaponBasePrice       int64
    ArmourBasePrice       int64
    WingBasePrice         int64

    JocPrice              int64
    JosPrice              int64
    JobPrice              int64
    JolPrice              int64

    RarityMultiplierPct   int64
    LevelMultiplierPct    int64
    AddPointsMultiplierPct int64
    LuckMultiplierPct     int64
    ExceMultiplierPct     int64

    BuySellMultiplier     int64
}

var ShopPriceList PriceList;

func CalculateItemSellingPrice(item *items.TokenAttributes) int64 {
    var basePrice int64

    if item.ItemAttributes.IsWeapon {
        basePrice = ShopPriceList.WeaponBasePrice
    } else if item.ItemAttributes.IsArmour {
        basePrice = ShopPriceList.ArmourBasePrice
    } else if item.ItemAttributes.IsJewel {
        switch item.Name {
        case "Jewel of Chaos":
            basePrice = ShopPriceList.JocPrice
        case "Jewel of Soul":
            basePrice = ShopPriceList.JosPrice
        case "Jewel of Bless":
            basePrice = ShopPriceList.JobPrice
        case "Jewel of Life":
            basePrice = ShopPriceList.JolPrice
        }
    } else if item.ItemAttributes.IsWings {
        basePrice = ShopPriceList.WingBasePrice
    }

    rarityMultiplier := 1 + (ShopPriceList.RarityMultiplierPct*item.ItemAttributes.ItemRarityLevel.Int64())/100
    levelMultiplier := 1 + (ShopPriceList.LevelMultiplierPct*item.ItemLevel.Int64())/100
    optionMultiplier := 1 + (ShopPriceList.AddPointsMultiplierPct*(item.AdditionalDamage.Int64()+item.AdditionalDefense.Int64())/4)/100
    luckMultiplier := int64(1)
    if item.Luck {
        luckMultiplier = 1 + (ShopPriceList.LuckMultiplierPct)/100
    }
    excellentMultiplier := int64(1)
    if item.ExcellentItemAttributes.IsExcellent {
        excellentMultiplier = 1 + (ShopPriceList.ExceMultiplierPct)/100
    }

    totalPrice := (basePrice * rarityMultiplier * levelMultiplier * optionMultiplier * luckMultiplier * excellentMultiplier) 

    return totalPrice
}


func loadShopPriceList() {
	var err error

	ShopPriceList.WeaponBasePrice, err = strconv.ParseInt(os.Getenv("PL_WEAPON_BASE_PRICE"), 10, 64)
	ShopPriceList.ArmourBasePrice, err = strconv.ParseInt(os.Getenv("PL_ARMOUR_BASE_PRICE"), 10, 64)
	ShopPriceList.WingBasePrice, err = strconv.ParseInt(os.Getenv("PL_WING_BASE_PRICE"), 10, 64)
	ShopPriceList.JocPrice, err = strconv.ParseInt(os.Getenv("PL_JOC_PRICE"), 10, 64)
	ShopPriceList.JosPrice, err = strconv.ParseInt(os.Getenv("PL_JOS_PRICE"), 10, 64)
	ShopPriceList.JobPrice, err = strconv.ParseInt(os.Getenv("PL_JOB_PRICE"), 10, 64)
	ShopPriceList.JolPrice, err = strconv.ParseInt(os.Getenv("PL_JOL_PRICE"), 10, 64)
	ShopPriceList.RarityMultiplierPct, err = strconv.ParseInt(os.Getenv("PL_RARITY_MULTIPLIER_PCT"), 10, 64)
	ShopPriceList.LevelMultiplierPct, err = strconv.ParseInt(os.Getenv("PL_LEVEL_MULTIPLIER_PCT"), 10, 64)
	ShopPriceList.AddPointsMultiplierPct, err = strconv.ParseInt(os.Getenv("PL_ADDPOINTS_MULTIPLIER_PCT"), 10, 64)
	ShopPriceList.LuckMultiplierPct, err = strconv.ParseInt(os.Getenv("PL_LUCK_MULTIPLIER_PCT"), 10, 64)
	ShopPriceList.ExceMultiplierPct, err = strconv.ParseInt(os.Getenv("PL_EXCE_MULTIPLIER_PCT"), 10, 64)
	ShopPriceList.BuySellMultiplier, err = strconv.ParseInt(os.Getenv("PL_BUYSELL_MULTIPLIER_PCT"), 10, 64)

	if err != nil {
		log.Fatalf("[loadShopPriceList] Failed to load price list=%v", err)
	}

	log.Printf("[loadShopPriceList] ShopPriceList=%v", ShopPriceList)
}