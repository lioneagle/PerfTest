package perf

import (
	"core"
	//"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPerfDataParseFuncDataFile(t *testing.T) {
	prefix := "TestPerfDataParseFuncDataFile"

	perfData := NewPerfData()

	input := os.Args[len(os.Args)-1] + "\\src\\testdata\\data1.txt"
	standard := os.Args[len(os.Args)-1] + "\\src\\testdata\\test_standard\\test_data_1.csv"
	output := os.Args[len(os.Args)-1] + "\\src\\testdata\\test_output\\test_data_1.csv"

	if !perfData.ParseFuncDataFile(input) {
		t.Errorf("%s: parse file failed\n", prefix)
		return
	}

	if !perfData.WriteDataToCsvFile(output) {
		t.Errorf("%s: write file failed\n", prefix)
		return
	}

	if !core.FileEqual(standard, output) {
		t.Errorf("%s: ouput file \"%s\" is not equal standard file \"%s\"", prefix, filepath.Base(output), filepath.Base(standard))
	}
}

func TestPerfDataParseStatDataFile(t *testing.T) {
	prefix := "TestPerfDataParseFuncDataFile"

	perfData := NewPerfData()

	input := os.Args[len(os.Args)-1] + "\\src\\testdata\\stat1.txt"
	standard := os.Args[len(os.Args)-1] + "\\src\\testdata\\test_standard\\test_stat_1.csv"
	output := os.Args[len(os.Args)-1] + "\\src\\testdata\\test_output\\test_stat_1.csv"

	if !perfData.ParseStatDataFile(input) {
		t.Errorf("%s: parse file failed\n", prefix)
		return
	}

	if !perfData.WriteStatToCsvFile(output) {
		t.Errorf("%s: write file failed\n", prefix)
		return
	}

	if !core.FileEqual(standard, output) {
		t.Errorf("%s: ouput file \"%s\" is not equal standard file \"%s\"", prefix, filepath.Base(output), filepath.Base(standard))
	}
}
