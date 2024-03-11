// market.go

package main

import (
	"sync"
	"math/big"

	"github.com/mriusd/game-contracts/items"
)


type MarketItem struct {
	Item 			items.TokenAttributes		`json:"item"`
	PriceGold 		*big.Int 			`json:"price_gold"`
	PriceCredits 	*big.Int 			`json:"price_credits"`
	PriceChaos 		int 				`json:"price_chaos"`
	PriceSoul 		int 				`json:"price_soul"`
	PriceBless 		int 				`json:"price_bless"`
	PriceLife 		int 				`json:"price_life"`

	sync.RWMutex
}


type SafeMarketMap struct {
	Items []*MarketItem
	sync.RWMutex
}