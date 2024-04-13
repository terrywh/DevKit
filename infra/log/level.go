package log

import "strings"

const TRACE = 0
const DEBUG = 10
const INFO = 20
const WARN = 30
const ERROR = 40
const FATAL = 50

func LevelFromString(name string) int {
	switch strings.ToUpper(name) {
	case "TRACE":
		return TRACE
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return WARN
	}
}
