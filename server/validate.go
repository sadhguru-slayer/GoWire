package main

import (
	"errors"
	"strings"
)

func validateCommand(parts []string) error {

	if len(parts) == 0 {
		return errors.New("ERR unknown command ''")
	}

	switch strings.ToUpper(parts[0]) {

	case "SET":
		if len(parts) != 3 {
			return errors.New("SET requires key and value")
		}
	case "GET":
		if len(parts) != 2 {
			return errors.New("GET requires key")
		}
	case "DEL":
		if len(parts) != 2 {
			return errors.New("DEL requires key")
		}
	case "EXISTS":
		if len(parts) != 2 {
			return errors.New("EXISTS requires key")
		}
	case "HSET":
		if len(parts) != 4 {
			return errors.New("HSET requires key, field and value")
		}

	case "HGET":
		if len(parts) != 3 {
			return errors.New("HGET requires key and field")
		}

	case "HDEL":
		if len(parts) != 3 {
			return errors.New("HDEL requires key and field")
		}

	case "HEXISTS":
		if len(parts) != 3 {
			return errors.New("HEXISTS requires key and field")
		}
	default:
		return errors.New("ERR unknown command " + parts[0])
	}
	return nil
}
