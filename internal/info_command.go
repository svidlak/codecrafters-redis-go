package commands

import (
	"strconv"
	"strings"

	config "github.com/codecrafters-io/redis-starter-go"
)

func InfoCommand(splitMsg []string) string {
	infoParam := splitMsg[4]
	if infoParam == "replication" {
		info := []string{
			"role:" + config.Configs.ClusterType,
			"master_replid:" + config.Configs.ID,
			"master_repl_offset:" + strconv.Itoa(config.Configs.Offset),
		}

		joined := strings.Join(info, "\n")

		return "$" + strconv.Itoa(len(joined)) + "\r\n" + joined
	}
	return ""
}
