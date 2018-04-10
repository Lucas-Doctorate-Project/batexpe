package batexpe

import (
	"fmt"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
)

// Stores info on one Batsim simulation instance
type Experiment struct {
	Batcmd            string  `json:"batcmd"`
	OutputDir         string  `json:"output-dir"`
	Schedcmd          string  `json:"schedcmd"`
	SimulationTimeout float64 `json:"simulation-timeout"`
	ReadyTimeout      float64 `json:"ready-timeout"`
	SuccessTimeout    float64 `json:"success-timeout"`
	FailureTimeout    float64 `json:"failure-timeout"`
}

func readStringFromDict(data map[string]interface{}, key string, yam string) (strRead string, err error) {
	if val, ok := data[key]; ok {
		if val == nil {
			strRead = ""
		} else {
			switch val.(type) {
			case string:
				strRead = val.(string)
			default:
				log.WithFields(log.Fields{
					"yaml": yam,
					"key":  key,
					"map":  data,
				}).Error("Invalid yaml: field is not a string")
				return "", fmt.Errorf("Invalid yaml: field is not a string")
			}
		}
	} else {
		log.WithFields(log.Fields{
			"yaml": yam,
			"key":  key,
			"map":  data,
		}).Error("Invalid yaml: missing field")
		return "", fmt.Errorf("Invalid yaml: missing field")
	}

	return strRead, nil
}

func readFloat64FromDict(data map[string]interface{}, key string, yam string) (fltRead float64, err error) {
	if val, ok := data[key]; ok {
		switch val.(type) {
		case float64:
			fltRead = val.(float64)
		default:
			log.WithFields(log.Fields{
				"yaml": yam,
				"key":  key,
				"map":  data,
			}).Error("Invalid yaml: field is not a float64")
			return -1, fmt.Errorf("Invalid yaml: field is not a float64")
		}
	} else {
		log.WithFields(log.Fields{
			"yaml": yam,
			"key":  key,
			"map":  data,
		}).Error("Invalid yaml: missing field")
		return -1, fmt.Errorf("Invalid yaml: missing field")
	}

	return fltRead, nil
}

func FromYaml(str string) (exp Experiment, convertErr error) {
	byt := []byte(str)

	var data map[string]interface{}
	err := yaml.Unmarshal(byt, &data)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"yaml": str,
		}).Error("Cannot yaml -> dict")
		return exp, fmt.Errorf("Cannot yaml -> dict")
	}

	log.WithFields(log.Fields{
		"yaml": str,
		"dict": data,
	}).Debug("yaml -> dict")

	var err1, err2, err3, err4, err5, err6, err7 error

	exp.Batcmd, err1 = readStringFromDict(data, "batcmd", str)
	exp.OutputDir, err2 = readStringFromDict(data, "output-dir", str)
	exp.Schedcmd, err3 = readStringFromDict(data, "schedcmd", str)
	exp.SimulationTimeout, err4 = readFloat64FromDict(data,
		"simulation-timeout", str)
	exp.ReadyTimeout, err5 = readFloat64FromDict(data, "ready-timeout", str)
	exp.SuccessTimeout, err6 = readFloat64FromDict(data, "success-timeout",
		str)
	exp.FailureTimeout, err7 = readFloat64FromDict(data, "failure-timeout",
		str)

	if (err1 != nil) || (err2 != nil) || (err3 != nil) || (err4 != nil) ||
		(err5 != nil) || (err6 != nil) || (err7 != nil) {
		return exp, fmt.Errorf("Invalid yaml")
	}

	log.WithFields(log.Fields{
		"expe": exp,
	}).Debug("dict->expe")

	return exp, nil
}

func ToYaml(exp Experiment) (yam string, err error) {
	byt, err := yaml.Marshal(exp)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"expe": exp,
		}).Error("Could not convert expe to yaml")
		return "", fmt.Errorf("Could not convert expe to yaml")
	}

	yam = string(byt)

	log.WithFields(log.Fields{
		"yaml": yam,
		"expe": exp,
	}).Debug("expe -> yaml")

	return yam, nil
}
