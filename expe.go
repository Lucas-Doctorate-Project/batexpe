package batexpe

import (
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

func readStringFromDict(data map[string]interface{}, key string, yam string) (strRead string) {
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
				}).Fatal("Invalid yaml: field is not a string")
			}
		}
	} else {
		log.WithFields(log.Fields{
			"yaml": yam,
			"key":  key,
			"map":  data,
		}).Fatal("Invalid yaml: missing field")
	}

	return strRead
}

func readFloat64FromDict(data map[string]interface{}, key string, yam string) (fltRead float64) {
	if val, ok := data[key]; ok {
		switch val.(type) {
		case int:
			fltRead = float64(val.(int))
		case float64:
			fltRead = val.(float64)
		default:
			log.WithFields(log.Fields{
				"yaml": yam,
				"key":  key,
				"map":  data,
			}).Fatal("Invalid yaml: field is not a string")
		}
	} else {
		log.WithFields(log.Fields{
			"yaml": yam,
			"key":  key,
			"map":  data,
		}).Fatal("Invalid yaml: missing field")
	}

	return fltRead
}

func FromYaml(str string) (exp Experiment) {
	byt := []byte(str)

	var data map[string]interface{}
	err := yaml.Unmarshal(byt, &data)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"yaml": str,
		}).Fatal("Cannot yaml -> dict")
	}

	log.WithFields(log.Fields{
		"yaml": str,
		"dict": data,
	}).Debug("yaml -> dict")

	exp.Batcmd = readStringFromDict(data, "batcmd", str)
	exp.OutputDir = readStringFromDict(data, "output-dir", str)
	exp.Schedcmd = readStringFromDict(data, "schedcmd", str)
	exp.SimulationTimeout = readFloat64FromDict(data, "simulation-timeout", str)
	exp.ReadyTimeout = readFloat64FromDict(data, "ready-timeout", str)
	exp.SuccessTimeout = readFloat64FromDict(data, "success-timeout", str)
	exp.FailureTimeout = readFloat64FromDict(data, "failure-timeout", str)

	log.WithFields(log.Fields{
		"expe": exp,
	}).Debug("dict->expe")

	return exp
}

func ToYaml(exp Experiment) (yam string) {
	byt, err := yaml.Marshal(exp)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"expe": exp,
		}).Fatal("Could not convert expe to yaml")
	}

	yam = string(byt)

	log.WithFields(log.Fields{
		"yaml": yam,
		"expe": exp,
	}).Debug("expe -> yaml")

	return yam
}
