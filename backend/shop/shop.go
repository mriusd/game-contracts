package shop

import (
    "log"
    "io/ioutil"
    "fmt"
    "encoding/json"

    "github.com/mriusd/game-contracts/inventory" 
    "github.com/mriusd/game-contracts/items" 
)

type ShopItem struct {
    ItemName string `json:"itemName"`
    Level int `json:"level"`
    PackSize int `json:"packSize"`
    Price int `json:"price"`
}

var Shops = make(map[string]*inventory.Inventory)

func LoadShops() {
    log.Printf("[LoadShops]")
    file, err := ioutil.ReadFile("./shop/shops.json")
    if err != nil {
        log.Fatalf("failed to read file: %v", err)
    }

    var jsonShops map[string][]ShopItem
    err = json.Unmarshal(file, &jsonShops)
    if err != nil {
        log.Fatalf("failed to unmarshal JSON: %v", err)
    }

    
    for shopName, shopItems := range jsonShops {
        shop := inventory.NewInventory(0, "shop")

        for _, shopItem := range shopItems {
            itemAtts, exists := items.BaseItemAttributes[shopItem.ItemName]
            if !exists {
                log.Fatalf("[LoadShops] Error Populating shop, item not found itemName=%v ", shopItem.ItemName)
            }

            itemParams, _ := items.BaseItemParameters[shopItem.ItemName]

            // Find the item by name
            item := &items.TokenAttributes{
                Name: shopItem.ItemName,
                ItemLevel: shopItem.Level,
                PackSize: shopItem.PackSize,

                ItemAttributes: itemAtts,
                ItemParameters: itemParams,
                ExcellentItemAttributes: items.ExcellentItemAttributes{},
            }

            itemHash, err := items.HashItemAttributes(item)
            if err != nil {
                log.Fatalf("[LoadShops] failed generating item hash item=%v", shopItem.ItemName)
            }

            shop.AddItem(item, 1, itemHash)
        }

        Shops[shopName] = shop

        log.Printf("[LoadShops] Shop %v loaded %v", shopName, Shops[shopName]) 
    }    
}


func GetShop (shopName string) (*inventory.Inventory, error) {
    shop, exists := Shops[shopName]
    if !exists {
        return nil, fmt.Errorf("[GetShop] Shop not found: %v", shopName)
    }

    return shop, nil
}