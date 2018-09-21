package controller

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/sysu-go-online/service-end/model"
)

func handleMapCommand(command []string) (*PortMapping, error) {
	if len(command) <= 0 || command[0] != "map" {
		return nil, errors.New("Invalid map command")
	}
	mapping := PortMapping{}
	// store which flag is un used
	isUsed := make([]bool, len(command))
	for i := 0; i < len(isUsed); i++ {
		isUsed[i] = false
	}

	// scan command
	for i := 1; i < len(command); i++ {
		if command[i][0] == '"' {
			if i != len(command)-1 {
				return nil, errors.New("Invalid command")
			}
			break
		}
		if isUsed[i] {
			continue
		}
		switch command[i] {
		case "-p":
			// parse port number
			if i == len(command)-1 {
				return nil, errors.New("can not get port number")
			}
			next := command[i+1]
			port, err := strconv.Atoi(next)
			if err != nil || port <= 0 || port >= 65535 {
				return nil, errors.New("Invalid port number")
			}
			mapping.Port = port
			isUsed[i] = true
			isUsed[i+1] = true
		default:
			return nil, fmt.Errorf("Can not parse %s", command[i])
		}
	}

	if mapping.Port == 0 {
		return nil, errors.New("Can not get port number")
	}

	// distribute domain name
	cnt := 0
	for {
		if cnt == 5 {
			return nil, errors.New("Can not get suitable domain name")
		}
		uuid := generateUUID()
		if has, err := model.IsUUIDExists(uuid, DomainNameRedisClient); err == nil {
			if has {
				cnt++
				continue
			} else {
				mapping.DomainName = uuid
			}
		} else {
			return nil, err
		}
	}
	// parse command
	userCommand := command[len(command)-1]
	if userCommand[0] == '"' && userCommand[len(userCommand)-1] == '"' {
		userCommand = userCommand[1 : len(userCommand)-1]
		mapping.Command = userCommand
		return &mapping, nil
	}
	return nil, errors.New("Invalid user command")
}
