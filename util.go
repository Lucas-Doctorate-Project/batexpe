package expe

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func CreateDirIfNeeded(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": dir,
			}).Fatal("Cannot create directory")
		}
	}
}

func max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}
