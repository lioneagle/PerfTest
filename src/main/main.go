package main

import (
	"core"
	"flag"
	"fmt"
	"perf"
)

type RunConfig struct {
	statFileName   string
	dataFileName   string
	outputFileName string
}

func (this *RunConfig) Parse() {
	flag.StringVar(&this.statFileName, "stat", "", "perf stat file")
	flag.StringVar(&this.dataFileName, "data", "", "perf data file")
	flag.StringVar(&this.outputFileName, "output", "output.csv", "output csv file")

	flag.Parse()
}

func (this *RunConfig) Check() bool {
	ok, _ := core.PathOrFileIsExist(this.statFileName)
	if !ok {
		fmt.Printf("ERROR: perf stat file \"%s\" is not exist\n", this.statFileName)
		return false
	}
	ok, _ = core.PathOrFileIsExist(this.dataFileName)
	if !ok {
		fmt.Printf("ERROR: perf output file \"%s\" is not exist\n", this.dataFileName)
		return false
	}

	return true
}

func main() {

	runConfig := &RunConfig{}
	runConfig.Parse()
	if !runConfig.Check() {
		return
	}

	perfData := perf.NewPerfData()

	if !perfData.ParseAll(runConfig.statFileName, runConfig.dataFileName) {
		return
	}

	perfData.WriteDataToCsvFile(runConfig.outputFileName)
}
