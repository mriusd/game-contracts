package main

import (
	"log"
	"strings"
	"strconv"
	"time"
	// "math/big"

	"github.com/mriusd/game-contracts/maps"
	"github.com/mriusd/game-contracts/items"
	"github.com/mriusd/game-contracts/fighters"
	"github.com/mriusd/game-contracts/drop"
)

type ParsedCommand struct {
	Name       string
	Attributes []string
}

func handleCommand(fighter *fighters.Fighter, command string) {
	// if !strings.HasPrefix(command, "/") {
	// 	log.Printf("User sent chat message: %v", command)
	// 	return
	// }

	parsedCommand := parseCommand(command)
	executeCommand(parsedCommand, fighter)
}

func parseCommand(command string) ParsedCommand {
    if !strings.HasPrefix(command, "/") {
        // This is a chat message
        return ParsedCommand{Name: "chat", Attributes: []string{command}}
    }

    words := strings.Fields(command)
    name := strings.TrimPrefix(words[0], "/")
    attributes := words[1:]

    return ParsedCommand{Name: name, Attributes: attributes}
}


func executeCommand(parsedCommand ParsedCommand, fighter *fighters.Fighter) {
	switch parsedCommand.Name {
		case "slide":
			// Handle command1
			log.Printf("[executeCommand:move] Attributes: %v", parsedCommand.Attributes)

			if (len(parsedCommand.Attributes) < 2) {
				log.Printf("[executeCommand:slide] Attributes: %v", parsedCommand.Attributes)
				sendErrorMessage(fighter, "Invalid coords")
				return
			}

			coords := maps.Coordinate{}
			x, err1 := strconv.ParseInt(parsedCommand.Attributes[0], 10, 64)
			y, err2 := strconv.ParseInt(parsedCommand.Attributes[1], 10, 64)

			if err1 != nil || err2 != nil {
				log.Printf("[executeCommand:slide] Attributes: %v", parsedCommand.Attributes)
				sendErrorMessage(fighter, "Invalid coords")
				return
			}

			coords.X = int(x)
			coords.Y = int(y)

			moveFighter(fighter, coords)

		case "spawn":
			// Handle command1
			log.Printf("[executeCommand:spawn] Attributes: %v", parsedCommand.Attributes)

			if (len(parsedCommand.Attributes) < 1) {
				log.Printf("[executeCommand:spawn] Attributes: %v", parsedCommand.Attributes)
				sendErrorMessage(fighter, "Invalid npcName")
				return
			}
			
			npcName :=  parsedCommand.Attributes[0]

			npc := FindNpcByName(npcName)
			if npc == nil {
				log.Printf("[executeCommand:spawn] Invalid npcId: %v", parsedCommand.Attributes)
				sendErrorMessage(fighter, "Invalid npcName")
				return
			}
			
			SpawnNPC(npc.ID, []string{"lorencia", strconv.FormatInt(int64(fighter.Coordinates.X), 10), strconv.FormatInt(int64(fighter.Coordinates.Y), 10)})


		case "make":
		    // Handle the make command
		    log.Printf("[executeCommand:make] Attributes: %v", parsedCommand.Attributes)

		    if len(parsedCommand.Attributes) < 1 {
		        log.Printf("[executeCommand:make] Invalid item name: %v", parsedCommand.Attributes)
		        sendErrorMessage(fighter, "Invalid item name")
		        return
		    }

		    itemWords := []string{}
		    level := int(0)
		    additionalPoints := int(0)
		    luck := false
		    excellent := false

		    for _, attr := range parsedCommand.Attributes {
		        if strings.HasPrefix(attr, "+") {
		            num, err := strconv.Atoi(attr[1:])
		            if err == nil {
		                switch {
		                case level == 0:
		                    level = num
		                case additionalPoints == 0:
		                    additionalPoints = num
		                default:
		                    log.Printf("[executeCommand:make] Ignoring extra '+' parameter: %s", attr)
		                }
		            } else {
		                log.Printf("[executeCommand:make] Error parsing number from attribute: %s", attr)
		            }
		        } else if strings.ToLower(attr) == "l" {
		            luck = true
		        } else if strings.ToLower(attr) == "exc" {
		            excellent = true
		        } else {
		            itemWords = append(itemWords, attr)
		        }
		    }

		    itemName := strings.Join(itemWords, " ")
		    generateItem(fighter, itemName, level, additionalPoints, luck, excellent)


		case "makeset":
		    // Handle the makeset command
		    log.Printf("[executeCommand:makeset] Attributes: %v", parsedCommand.Attributes)

		    if len(parsedCommand.Attributes) < 2 {
		        log.Printf("[executeCommand:makeset] Invalid command parameters: %v", parsedCommand.Attributes)
		        sendErrorMessage(fighter, "Invalid command parameters")
		        return
		    }

		    setName := ""
		    level := int(0)
		    excellent := false

		    for i, attr := range parsedCommand.Attributes {
		        if i == 0 {
		            setName = attr
		        } else if strings.HasPrefix(attr, "+") {
		            num, err := strconv.Atoi(attr[1:])
		            if err == nil {
		                level = num
		            } else {
		                log.Printf("[executeCommand:makeset] Error parsing number from attribute: %s", attr)
		            }
		        } else if strings.ToLower(attr) == "exc" {
		            excellent = true
		        } else {
		            log.Printf("[executeCommand:makeset] Ignoring unrecognized attribute: %s", attr)
		        }
		    }

		    generateItem(fighter, setName+" Helm", level, 0, false, excellent)
		    generateItem(fighter, setName+" Armour", level, 0, false, excellent)
		    generateItem(fighter, setName+" Pants", level, 0, false, excellent)
		    generateItem(fighter, setName+" Boots", level, 0, false, excellent)
		    generateItem(fighter, setName+" Gloves", level, 0, false, excellent)

		case "chat":
	        // Handle chat messages
	        chatMessage := strings.Join(parsedCommand.Attributes, " ")
	        log.Printf("[executeCommand:chat] %s: %s", fighter.Name, chatMessage)

	        fighter.Lock()
	        fighter.LastChatMsg = chatMessage
	        fighter.LastChatMsgTimestamp = int(time.Now().UnixNano()) / int(time.Millisecond)
	        fighter.Unlock()
	        broadcastChatMsg(fighter.Location, fighter.Name, chatMessage, "local")
	    

		default:
			log.Printf("[executeCommand] uknown command: %v", parsedCommand.Name)
	}
}


func generateItem(fighter *fighters.Fighter, itemName string, level, additionalPoints int, luck, excellent bool) {
    log.Printf("[generateItem] itemName=%v", itemName)

    itemAtts, exists := items.BaseItemAttributes[itemName]
    if !exists {
    	log.Printf("[generateItem] Error Generating item not found itemName=%v ", itemName)
		sendErrorMsgToFighter(fighter, "SYSTEM" , "Item not found")
		return;
    }

    itemParams, _ := items.BaseItemParameters[itemName]

    // Find the item by name
	item := items.TokenAttributes{
		Name: itemName,
		ItemLevel: level,
		PackSize: 1,
		Luck: luck,

		ItemAttributes: itemAtts,
		ItemParameters: itemParams,
		ExcellentItemAttributes: items.ExcellentItemAttributes{},
	}

	// type TokenAttributes struct {
	// 	Name            		string 					`json:"name" bson:"name"`
	// 	TokenId         		int 					`json:"tokenId" bson:"token_id"`
	// 	ItemLevel       		int 					`json:"itemLevel" bson:"item_level"`
	// 	AdditionalDamage 		int 					`json:"additionalDamage" bson:"additional_damage"`
	// 	AdditionalDefense 		int 					`json:"additionalDefense" bson:"additional_defence"`
	// 	FighterId       		int 					`json:"fighterId" bson:"fighter_id"`
	// 	PackSize        		int 					`json:"packSize" bson:"pack_size"`
	// 	Luck            		bool   					`json:"luck" bson:"luck"`
	// 	Skill           		bool   					`json:"skill" bson:"skill"`

	// 	CreatedAt				int 					`json:"createdAt" bson:"created_at"`

	// 	ItemAttributes  		ItemAttributes 			`json:"itemAttributes" bson:"-"`
	// 	ItemParameters 			ItemParameters 			`json:"itemParameters" bson:"-"`
	// 	ExcellentItemAttributes ExcellentItemAttributes `json:"excellentItemAttributes" bson:"excellent_item_attributes"`

	// 	sync.RWMutex									`json:"-" bson:"-"`
	// }



	log.Printf("[generateItem] item=%v", item)
	
    
    if item.ItemAttributes.IsWeapon {
    	item.AdditionalDamage = additionalPoints
    	item.Skill = true
    } 

    if item.ItemAttributes.IsArmour {
    	item.AdditionalDefense = additionalPoints
    } 
    
    item.Luck = luck

    
	if excellent {
    	item.ExcellentItemAttributes.IncreaseAttackSpeedPoints = 1
    	item.ExcellentItemAttributes.ManaAfterMonsterIncrease = 1
    	item.ExcellentItemAttributes.LifeAfterMonsterIncrease = 1
    	item.ExcellentItemAttributes.GoldAfterMonsterIncrease = 1
    	item.ExcellentItemAttributes.ReflectDamagePercent = 1
    	item.ExcellentItemAttributes.RestoreHPChance = 1
    	item.ExcellentItemAttributes.RestoreMPChance = 1
    	item.ExcellentItemAttributes.DoubleDamageChance = 1
    	item.ExcellentItemAttributes.IgnoreOpponentDefenseChance = 1
    	item.ExcellentItemAttributes.ExcellentDamageProbabilityIncrease = 1
    	item.ExcellentItemAttributes.AttackSpeedIncrease = 1
    	item.ExcellentItemAttributes.AttackLvl20 = 1
    	item.ExcellentItemAttributes.AttackIncreasePercent = 1
    	item.ExcellentItemAttributes.DefenseSuccessRateIncrease = 1
    	item.ExcellentItemAttributes.ReflectDamage = 1
    	item.ExcellentItemAttributes.MaxLifeIncrease = 1
    	item.ExcellentItemAttributes.MaxManaIncrease = 1
    	item.ExcellentItemAttributes.DecreaseDamageRateIncrease = 1
    	item.ExcellentItemAttributes.HpRecoveryRateIncrease = 1
    	item.ExcellentItemAttributes.MpRecoveryRateIncrease = 1
    }    


    drop.MakeItem(&item, fighter, fighter.GetLocation(), fighter.GetCoordinates()) 

    broadcastDropMessage()  
}





