package filesystem

import (
	"os"
	"runtime"
)

func getHomeDirectory() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		if home == "" {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
		return home
	}

	home := os.Getenv("HOME")
	if home == "" {
		home = "/"
	}
	return home
}
