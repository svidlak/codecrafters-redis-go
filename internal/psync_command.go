package commands

import (
	"fmt"

	config "github.com/codecrafters-io/redis-starter-go"
)

func PsyncCommand(splitMsg []string) string {
	fmt.Print(config.Configs.ID)
	return "+FULLRESYNC " + config.Configs.ID + " 0"
}
