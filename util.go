package expe

import (
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
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

func PortFromBatSock(socket, batcmd string) uint16 {
	regexStr := `^.*:(?P<Port>\d+)$`
	r := regexp.MustCompile(regexStr)
	capture := r.FindStringSubmatch(socket)

	if capture == nil {
		log.WithFields(log.Fields{
			"socket endpoint":  socket,
			"extraction regex": regexStr,
			"batsim command":   batcmd,
		}).Fatal("Cannot retrieve port from batsim command")
	}

	port, err := strconv.Atoi(capture[1])
	if err != nil {
		log.WithFields(log.Fields{
			"socket endpoint":  socket,
			"extraction regex": regexStr,
			"captured string":  capture[1],
		}).Fatal("Cannot convert string to int")
	}

	return uint16(port)
}
