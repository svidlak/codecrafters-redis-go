package commands

func PsyncCommand(splitMsg []string) string {
	return "+FULLRESYNC <REPL_ID> 0"
}
