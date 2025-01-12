package LogUtil

import (
	"fmt"
	"log"
)

var LogPrefixStorage = ""

const ModuleLogPrefix = "ESCAPE-API"

// Designed to be used in conjunction with some variant of SetLogPrefixFor...Package().
// Defer this function to make sure the log prefix is properly reset. See endpoints/creation.go for examples on use
func EnsureLogPrefixIsReset() {
	log.Println()
	log.SetPrefix(LogPrefixStorage)
}

// NOTE: ONLY USE THIS IS CONJUNCTION WITH A DEFERRED EnsureLogPrefixIsReset()!!
// Sets the prefix for logging based on the given package prefix
func SetLogPrefix(modulePrefix string, packagePrefix string) {
	LogPrefixStorage = log.Prefix()
	log.SetPrefix(fmt.Sprintf("%s/%s: ", modulePrefix, packagePrefix))
	log.Println()
}

// Logs the error in the format of "[funcLogPrefix] ERROR! [err]"
func LogError(funcLogPrefix string, err error) {
	log.Printf("%s ERROR! %s", funcLogPrefix, err)
}
