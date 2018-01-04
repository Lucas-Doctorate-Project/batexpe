package batexpe

import (
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
)

type Experiment struct {
	Batcmd            string  `json:"batcmd"`
	OutputDir         string  `json:"output-dir"`
	Schedcmd          string  `json:"schedcmd"`
	SimulationTimeout float64 `json:"simulation-timeout"`
	ReadyTimeout      float64 `json:"ready-timeout"`
	SuccessTimeout    float64 `json:"success-timeout"`
	FailureTimeout    float64 `json:"failure-timeout"`
}

func readStringFromDict(data map[string]interface{}, key string, yam string) string {
	ret := ""
	if val, ok := data[key]; ok {
		if val == nil {
			ret = ""
		} else {
			switch val.(type) {
			case string:
				ret = val.(string)
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

	return ret
}

func readFloat64FromDict(data map[string]interface{}, key string, yam string) float64 {
	ret := 0.0
	if val, ok := data[key]; ok {
		switch val.(type) {
		case int:
			ret = float64(val.(int))
		case float64:
			ret = val.(float64)
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

	return ret
}

func FromYaml(str string) Experiment {
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

	var expe Experiment
	expe.Batcmd = readStringFromDict(data, "batcmd", str)
	expe.OutputDir = readStringFromDict(data, "output-dir", str)
	expe.Schedcmd = readStringFromDict(data, "schedcmd", str)
	expe.SimulationTimeout = readFloat64FromDict(data, "simulation-timeout", str)
	expe.ReadyTimeout = readFloat64FromDict(data, "ready-timeout", str)
	expe.SuccessTimeout = readFloat64FromDict(data, "success-timeout", str)
	expe.FailureTimeout = readFloat64FromDict(data, "failure-timeout", str)

	log.WithFields(log.Fields{
		"expe": expe,
	}).Debug("dict->expe")

	return expe
}

func ToYaml(exp Experiment) string {
	byt, err := yaml.Marshal(exp)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"expe": exp,
		}).Fatal("Could not convert expe to yaml")
	}

	yam := string(byt)

	log.WithFields(log.Fields{
		"yaml": yam,
		"expe": exp,
	}).Debug("expe -> yaml")

	return yam
}
