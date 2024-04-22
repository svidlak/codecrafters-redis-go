package commands

import (
	"strconv"
	"time"
)

func SetCommand(splitMsg []string) string {
	key := splitMsg[4]
	val := splitMsg[6]

	if len(splitMsg) >= 10 {
		expirationTime, err := strconv.Atoi(splitMsg[10])
		if err != nil {
			return "-ERROR"
		}
		go time.AfterFunc(time.Duration(expirationTime)*time.Millisecond, func() {
			delete(Set, key)
		})
	}
	Set[key] = val
	return "+OK"
}
