package main

import (
	"log"
	"strings"
	"strconv"
	"time"
	"math/big"

	"github.com/mriusd/game-contracts/maps"
	"github.com/mriusd/game-contracts/items"
	"github.com/mriusd/game-contracts/fighters"
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

			coords.X = x
			coords.Y = y

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

			npc := findNpcByName(npcName)
			if npc == nil {
				log.Printf("[executeCommand:spawn] Invalid npcId: %v", parsedCommand.Attributes)
				sendErrorMessage(fighter, "Invalid npcName")
				return
			}
			
			spawnNPC(npc.ID, []string{"lorencia", strconv.FormatInt(fighter.Coordinates.X, 10), strconv.FormatInt(fighter.Coordinates.Y, 10)})


		case "make":
		    // Handle the make command
		    log.Printf("[executeCommand:make] Attributes: %v", parsedCommand.Attributes)

		    if len(parsedCommand.Attributes) < 1 {
		        log.Printf("[executeCommand:make] Invalid item name: %v", parsedCommand.Attributes)
		        sendErrorMessage(fighter, "Invalid item name")
		        return
		    }

		    itemWords := []string{}
		    level := int64(0)
		    additionalPoints := int64(0)
		    luck := false
		    excellent := false

		    for _, attr := range parsedCommand.Attributes {
		        if strings.HasPrefix(attr, "+") {
		            num, err := strconv.ParseInt(attr[1:], 10, 64)
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
		    level := int64(0)
		    excellent := false

		    for i, attr := range parsedCommand.Attributes {
		        if i == 0 {
		            setName = attr
		        } else if strings.HasPrefix(attr, "+") {
		            num, err := strconv.ParseInt(attr[1:], 10, 64)
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

		    generateItem(fighter, setName+" helm", level, 0, false, excellent)
		    generateItem(fighter, setName+" armour", level, 0, false, excellent)
		    generateItem(fighter, setName+" pants", level, 0, false, excellent)
		    generateItem(fighter, setName+" boots", level, 0, false, excellent)
		    generateItem(fighter, setName+" gloves", level, 0, false, excellent)

		 case "chat":
	        // Handle chat messages
	        chatMessage := strings.Join(parsedCommand.Attributes, " ")
	        log.Printf("[executeCommand:chat] %s: %s", fighter.Name, chatMessage)

	        fighter.Lock()
	        fighter.LastChatMsg = chatMessage
	        fighter.LastChatMsgTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
	        fighter.Unlock()
	        broadcastChatMsg(fighter.Location, fighter.Name, chatMessage, "local")
	    

		default:
			log.Printf("[executeCommand] uknown command: %v", parsedCommand.Name)
	}
}


func generateItem(fighter *fighters.Fighter, itemName string, level, additionalPoints int64, luck, excellent bool) {
    log.Printf("[generateItem] itemName=%v", itemName)

    // Find the item by name
	item, err := items.GenerateSolidityItem(strings.ToLower(itemName))

	if err != nil {
		log.Printf("[generateItem] Error Generating item itemName=%v error=%v", itemName, err)
		sendErrorMsgToFighter(fighter, "SYSTEM" , "Item not found")
		return;
	}

	log.Printf("[generateItem] item=%v", item)
	
    // Update item attributes based on the drop command
    item.ItemLevel = big.NewInt(level)

    if item.IsWeapon {
    	item.AdditionalDamage = big.NewInt(additionalPoints)
    	item.Skill = true
    } 

    if item.IsArmour {
    	item.AdditionalDefense = big.NewInt(additionalPoints)
    } 
    
    item.Luck = luck

    
	if excellent {
    	item.IncreaseAttackSpeedPoints = big.NewInt(1)
    	item.ManaAfterMonsterIncrease = big.NewInt(1)
    	item.LifeAfterMonsterIncrease = big.NewInt(1)
    	item.GoldAfterMonsterIncrease = big.NewInt(1)
    	item.ReflectDamagePercent = big.NewInt(1)
    	item.RestoreHPChance = big.NewInt(1)
    	item.RestoreMPChance = big.NewInt(1)
    	item.DoubleDamageChance = big.NewInt(1)
    	item.IgnoreOpponentDefenseChance = big.NewInt(1)
    	item.ExcellentDamageProbabilityIncrease = big.NewInt(1)
    	item.AttackSpeedIncrease = big.NewInt(1)
    	item.AttackLvl20 = big.NewInt(1)
    	item.AttackIncreasePercent = big.NewInt(1)
    	item.DefenseSuccessRateIncrease = big.NewInt(1)
    	item.ReflectDamage = big.NewInt(1)
    	item.MaxLifeIncrease = big.NewInt(1)
    	item.MaxManaIncrease = big.NewInt(1)
    	item.DecreaseDamageRateIncrease = big.NewInt(1)
    	item.HpRecoveryRateIncrease = big.NewInt(1)
    	item.MpRecoveryRateIncrease = big.NewInt(1)
    }    


    MakeItem(fighter, &item)   
}





