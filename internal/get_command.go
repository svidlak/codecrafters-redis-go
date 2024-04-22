package commands

func GetCommand(splitMsg []string) string {
	key := splitMsg[4]
	val, ok := Set[key]
	if ok {
		return "+" + val
	}
	return "$-1"
}
