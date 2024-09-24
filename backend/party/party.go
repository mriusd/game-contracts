// party.go

package party

import (
	"context"
	"log"
	"sync"

	"github.com/mriusd/game-contracts/fighter"
)

var MAX_PARTY_MEMBERS = 5

type Party struct {
	Map []*fighters.Fighter 
	Leader *fighter.Fighters
	sync.RWMutex
}

type SafePartiesMap struct {
	Map []*Party
	sync.RWMutex
}

var PartyMap SafePartiesMap

type PartyRequest struct {
	Maker *fighter.Fighters
	Taker *fighter.Fighters
	CreatedAt time.Time
}

type SafePartyRequestMap struct {
	Map []*PartyRequest 
	sync.RWMutex
}

var PartyRequestMap SafePartyRequestMap


func (i *SafePartyRequestMap) FindRequestByMaker(maker *fighters.Fighter) *PartyRequest {
	i.RLock()
	defer i.RUnlock()

	for _, partyRequest := range i.Map {
		if partyRequest.Maker == maker {
			return partyRequest
		}
	}

	return nil
}

func (i *SafePartyRequestMap) CreateRequest(maker, taker  *fighters.Fighter) {
	i.Lock()
	defer i.Unlock()

	i.Map = append(i.Map, &PartyRequest{
		Maker: maker,
		Taker: taker,
	})
}

func (i *SafePartyRequestMap) RemoveRequest(request *PartyRequest) {
	i.Lock()
	defer i.Unlock()

	// Find the index of the request to remove
	for idx, req := range i.Map {
		if req == request {
			// Remove the request from the slice
			i.Map = append(i.Map[:idx], i.Map[idx+1:]...)
			break
		}
	}
}



func (i *SafePartyRequestMap) FindRequestByTaker(taker *fighters.Fighter) *PartyRequest {
	i.RLock()
	defer i.RUnlock()

	for _, partyRequest := range i.Map {
		if partyRequest.Taker == taker {
			return partyRequest
		}
	}

	return nil
}

func (i *Party) GetLeader() *fighters.Fighter {
	i.RLock()
	defer i.RUnlock()

    return i.Leader
}

func (i *Party) FindFighter(fighter *fighters.Fighter) bool {
	i.RLock()
	defer i.RUnlock()

	for _, val := range i.Map {
        if val == fighter {
        	return true
        }
    }

    return false
}

func (i *Party) AddMember(fighter *fighters.Fighter) error {
	i.Lock()
	defer i.Unlock()


	if len(makerParty.Map) >= MAX_PARTY_MEMBERS {
		return fmt.Errorf("Party full")
	}

	i.Map = append(i.Map, fighter)

	return nil
}



func (i *SafePartiesMap) FindParty(fighter *fighters.Fighter) *Party {
	i.RLock()
	defer i.RUnlock()

	for _, val := range i.Map {
        if val.FindFighter(fighter) {
        	return val
        }
    }

    return nil
}


func (i *SafePartiesMap) CreateParty(leader *fighter.Fighter) (*Party, error) {
	existingParty := i.FindParty(leader)
	if existingParty != nil {
		return nil, fmt.Errorf("[CreateParty] Player already in party")
	}

	// Create a new Party with the leader
	newParty := &Party{
		Map:    []*fighter.Fighter{leader},
		Leader: leader,
	}

	// Lock the SafePartiesMap for writing
	i.Lock()
	defer i.Unlock()

	// Add the new Party to the PartyMap
	i.Map = append(i.Map, newParty)

	return newParty, nil
}


func AcceptRequest(taker, maker *fighter.Fighter) error {
	// Find the party request where maker sent a request to taker
	request := PartyRequestMap.FindRequestByMaker(maker)
	if request == nil || request.Taker != taker {
		return fmt.Errorf("No party request found")
	}

	// Check if the taker is already in a party
	if PartyMap.FindParty(taker) != nil {
		return fmt.Errorf("Taker is already in a party")
	}

	// Find maker's party
	makerParty := PartyMap.FindParty(maker)

	if makerParty == nil {
		// Maker is not in a party; create a new party with maker as leader
		newParty, err := PartyMap.CreateParty(maker)
		if err != nil {
			return err
		}

		err = newParty.AddMember(taker)
		if err != nil {
			return err
		}

	} else {
		err = makerParty.AddMember(taker)
		if err != nil {
			return err
		}
	}

	// Remove the party request from the map
	PartyRequestMap.RemoveRequest(request)

	return nil
}



func RequestParty (maker, taker *fighters.Fighter) error {
	if SafePartiesMap.FindParty(taker) != nil {
		return fmt.Errorf("Player already in party")
	}

	makerParty := SafePartiesMap.FindParty(maker)
	if makerParty != nil && makerParty.GetLeader() != maker {
		return fmt.Errorf("You are not the party leader")
	}

	takerRequestTaker := PartyRequestMap.FindRequestByTaker(taker)
	takerRequestMaker := PartyRequestMap.FindRequestByMaker(taker)

	if takerRequestTaker != nil && takerRequestMaker != nil {
		return fmt.Errorf("Player has a party request already")		
	}


	makerRequestTaker := PartyRequestMap.FindRequestByTaker(maker)
	if makerRequestTaker != nil {
		PartyRequestMap.RemoveRequest(makerRequestTaker)
	}

	makerRequestMaker := PartyRequestMap.FindRequestByMaker(maker)
		if makerRequestMaker != nil {
		PartyRequestMap.RemoveRequest(makerRequestMaker)
	}

	PartyRequestMap.CreateRequest(maker, taker)

	return nil
}





