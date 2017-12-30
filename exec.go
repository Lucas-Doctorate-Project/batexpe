package expe

// import (
// 	log "github.com/sirupsen/logrus"
// 	"os"
// )

func PrepareOutput(exp Experiment) {
	// Create output directory if needed
	CreateDirIfNeeded(exp.OutputDir)
	CreateDirIfNeeded(exp.OutputDir + "/log")
}
