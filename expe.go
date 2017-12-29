package expe

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
)

type Experiment struct {
	Batcmd            string `json:"batcmd"`
	OutputDir         string `json:"output-dir"`
	Schedcmd          string `json:"schedcmd"`
	SimulationTimeout int    `json:"simulation-timeout"`
	SocketTimeout     int    `json:"socket-timeout"`
	SuccessTimeout    int    `json:"success-timeout"`
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
	expe.Batcmd = data["batcmd"].(string)
	expe.OutputDir = data["output-dir"].(string)
	expe.Schedcmd = data["schedcmd"].(string)
	expe.SimulationTimeout = int(data["simulation-timeout"].(float64))
	expe.SocketTimeout = int(data["socket-timeout"].(float64))
	expe.SuccessTimeout = int(data["success-timeout"].(float64))
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
