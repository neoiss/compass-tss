package constants

import (
	"os"
)

const (
	DefaultDir = ".compass-tss"
)

var (
	DefaultHome = os.ExpandEnv("$HOME/") + DefaultDir
)
