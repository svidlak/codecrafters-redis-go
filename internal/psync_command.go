package commands

import (
	config "github.com/codecrafters-io/redis-starter-go"
)

func PsyncCommand(splitMsg []string) string {
	return "+FULLRESYNC " + config.Configs.ID + " 0"
}
