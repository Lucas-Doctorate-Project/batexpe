package batexpe

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func CreateDirIfNeeded(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("Cannot create directory (err=%s)", err.Error())
		}
	}
	return nil
}

func max(x, y int) (maxVal int) {
	if x > y {
		return x
	} else {
		return y
	}
}

func PortFromBatSock(socket string) (port uint16, err error) {
	regexStr := `^.*:(?P<Port>\d+)$`
	r := regexp.MustCompile(regexStr)
	capture := r.FindStringSubmatch(socket)

	if capture == nil {
		return 0, fmt.Errorf("Cannot extract port with regex '%s': No match.",
			regexStr)
	}

	iport, err := strconv.Atoi(capture[1])
	if err != nil {
		return 0, fmt.Errorf("Cannot convert port '%s' to int", capture[1])
	}

	return uint16(iport), nil
}

func PreviewFile(filename string, maxLines int64) (preview string, err error) {
	// Retrieve the file length
	wcCmd := exec.Command("wc")
	wcCmd.Args = []string{wcCmd.Args[0], "-l", filename}

	wcOut, err := wcCmd.Output()
	if err != nil {
		return "", fmt.Errorf("Cannot call 'wc -l %s'", filename)
	}

	wcR := regexp.MustCompile(`(?m)^(\d+)\s+.*$`)
	wcCap := wcR.FindStringSubmatch(string(wcOut))

	if wcCap == nil {
		return "", fmt.Errorf("Cannot retrieve number of lines in "+
			"wc output '%s'", string(wcOut))
	}

	nbLines, _ := strconv.ParseInt(wcCap[1], 10, 32)

	if nbLines <= maxLines {
		// Retrieve the whole file content
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", fmt.Errorf("Cannot read file")
		}
		return string(content), nil
	} else {
		// Only retrieve the first and last lines
		// First lines
		headCmd := exec.Command("head")
		headCmd.Args = []string{headCmd.Args[0], "-n",
			strconv.Itoa(int(maxLines / 2)), filename}

		headOut, err := headCmd.Output()
		if err != nil {
			return "", fmt.Errorf("Cannot call 'head -n %s %s'",
				string(maxLines/2), filename)
		}

		// Last lines
		tailCmd := exec.Command("tail")
		tailCmd.Args = []string{tailCmd.Args[0], "-n",
			strconv.Itoa(int(maxLines / 2)), filename}

		tailOut, err := tailCmd.Output()
		if err != nil {
			return "", fmt.Errorf("Cannot call 'tail -n %s %s'",
				string(maxLines/2), filename)
		}

		return fmt.Sprintf("%s...\n...\n"+
			"... (truncated... whole log in '%s')\n...\n...\n%s",
			string(headOut), filename, string(tailOut)), nil
	}
}

func IsBatsimOrBatschedRunning() (bool, error) {
	// This function directly searches for batsim or batsched processes.
	rBatsim := regexp.MustCompile(`(?m)^\S*\bbatsim .*$`)
	rBatsched := regexp.MustCompile(`(?m)^\S*\bbatsched.*$`)

	psCmd := exec.Command("ps")
	psCmd.Args = []string{psCmd.Args[0], "-e", "-o", "command"}

	outBuf, err := psCmd.Output()
	if err != nil {
		log.WithFields(log.Fields{
			"err":     err,
			"command": psCmd,
		}).Error("Cannot list running processes via ps")
		return false, err
	}

	batsimMatches := rBatsim.FindAllString(string(outBuf), -1)
	batschedMatches := rBatsched.FindAllString(string(outBuf), -1)
	return len(batsimMatches) > 0 || len(batschedMatches) > 0, nil
}
