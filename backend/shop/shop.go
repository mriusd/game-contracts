package shop

import (
    "log"
    "io/ioutil"
    "fmt"
    "encoding/json"
    "time"

    "github.com/mriusd/game-contracts/inventory" 
    "github.com/mriusd/game-contracts/items" 
    "github.com/mriusd/game-contracts/fighters" 
)

type ShopItem struct {
    ItemName string `json:"itemName"`
    Level int `json:"level"`
    PackSize int `json:"packSize"`
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

func BuyItem(fighter *fighters.Fighter, shopName, hash string) error {
    shop, err := GetShop(shopName)
    if err != nil {
        return err
    }

    shopItem := shop.FindByHash(hash).GetAttributes()
    if shopItem == nil {
        return fmt.Errorf("[BuyItem] Item not found %v", hash)
    }

    itemAttributes := shopItem.GetItemAttributes()

    backpack := fighter.GetBackpack()
    availableGold := backpack.GetGold()
    price := itemAttributes.Price
    
    if availableGold < price {
        return fmt.Errorf("[BuyItem] Not enought gold")
    }

    itemParams := shopItem.GetItemParameters()
    if !backpack.IsEnoughSpace(itemParams.ItemWidth, itemParams.ItemHeight) {
        return fmt.Errorf("[BuyItem] Not enough space")
    }

    // Find the item by name
    item := &items.TokenAttributes{
        Name: shopItem.Name,
        ItemLevel: shopItem.ItemLevel,
        PackSize: shopItem.PackSize,

        CreatedAt: int(time.Now().UnixNano()),

        ItemAttributes: itemAttributes,
        ItemParameters: itemParams,
        ExcellentItemAttributes: items.ExcellentItemAttributes{},
    }

    itemHash, err := items.HashItemAttributes(item)
    if err != nil {
        return fmt.Errorf("[BuyItem] failed generating item hash item=%v", shopItem.Name)
    }

    backpack.SetGold(availableGold - price)
    backpack.AddItem(item, shopItem.PackSize, itemHash)
    return nil
}

























