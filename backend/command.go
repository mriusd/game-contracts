package main

import (
	"log"
	"strings"
	"strconv"
)

type ParsedCommand struct {
	Name       string
	Attributes []string
}

func handleCommand(fighter *Fighter, command string) {
	if !strings.HasPrefix(command, "/") {
		log.Printf("User sent chat message: %v", command)
		return
	}

	parsedCommand := parseCommand(command)
	executeCommand(parsedCommand, fighter)
}

func parseCommand(command string) ParsedCommand {
	words := strings.Fields(command)
	name := strings.TrimPrefix(words[0], "/")
	attributes := words[1:]

	return ParsedCommand{Name: name, Attributes: attributes}
}

func executeCommand(parsedCommand ParsedCommand, fighter *Fighter) {
	switch parsedCommand.Name {
	case "slide":
		// Handle command1
		log.Printf("[executeCommand:move] Attributes:", parsedCommand.Attributes)

		if (len(parsedCommand.Attributes) < 2) {
			log.Printf("[executeCommand:slide] Attributes:", parsedCommand.Attributes)
			sendErrorMessage(fighter, "Invalid coords")
			return
		}

		coords := Coordinate{}
		x, err1 := strconv.ParseInt(parsedCommand.Attributes[0], 10, 64)
		y, err2 := strconv.ParseInt(parsedCommand.Attributes[1], 10, 64)

		if err1 != nil || err2 != nil {
			log.Printf("[executeCommand:slide] Attributes:", parsedCommand.Attributes)
			sendErrorMessage(fighter, "Invalid coords")
			return
		}

		coords.X = x
		coords.Y = y

		moveFighter(fighter, coords)

	case "spawn":
		// Handle command1
		log.Printf("[executeCommand:spawn] Attributes:", parsedCommand.Attributes)

		if (len(parsedCommand.Attributes) < 1) {
			log.Printf("[executeCommand:spawn] Attributes:", parsedCommand.Attributes)
			sendErrorMessage(fighter, "Invalid npcId")
			return
		}

		npcId, err := strconv.ParseInt(parsedCommand.Attributes[0], 10, 64)

		if err != nil {
			log.Printf("[executeCommand:spawn] Invalid npcId:", parsedCommand.Attributes)
			sendErrorMessage(fighter, "Invalid npcId")
			return
		}

		npc := findNpcById(npcId)
		if npc == nil {
			log.Printf("[executeCommand:spawn] Invalid npcId:", parsedCommand.Attributes)
			sendErrorMessage(fighter, "Invalid npcId")
			return
		}
		
		spawnNPC(npcId, []string{"lorencia", strconv.FormatInt(fighter.Coordinates.X, 10), strconv.FormatInt(fighter.Coordinates.Y, 10)})

	default:
		log.Printf("[executeCommand] uknown command:", parsedCommand.Name)
	}
}
