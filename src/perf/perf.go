package perf

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type FuncData struct {
	Name string
	Data map[string]float64
}

func NewFuncData() *FuncData {
	return &FuncData{Data: make(map[string]float64, 0)}
}

func (this *FuncData) AddData(name string, value float64) {
	this.Data[name] = value
}

func (this *FuncData) StringSlice(DataNames []string) (ret []string) {
	ret = append(ret, this.Name)
	for _, v := range DataNames {
		data, ok := this.Data[v]
		if !ok {
			fmt.Printf("ERROR: cannot get data \"%s\"\n", v)
		}
		ret = append(ret, fmt.Sprintf("%.2f%%", data))
	}
	return
}

type StatData struct {
	DataNames []string
	Data      map[string]float64
}

func NewStatData() *StatData {
	return &StatData{Data: make(map[string]float64, 0)}
}

func (this *StatData) AddData(name string, value float64) {
	_, ok := this.Data[name]
	if !ok {
		this.DataNames = append(this.DataNames, name)
	}
	this.Data[name] = value
}

func (this *StatData) WriteToCsvFile(filename string) bool {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Printf("ERROR: cannot open csv file %s to write\n", filename)
		return false
	}
	defer file.Close()

	w := csv.NewWriter(file)

	for _, v := range this.DataNames {
		data, ok := this.Data[v]
		if !ok {
			fmt.Printf("ERROR: cannot find data \"%s\"\n", v)
			return false
		}
		line := []string{v, fmt.Sprintf("%.2f", data)}
		w.Write(line)
	}

	w.Flush()
	return true
}

type PerfData struct {
	Funcs     map[string]*FuncData
	FuncNames []string
	DataNames []string

	statData *StatData
}

func NewPerfData() *PerfData {
	return &PerfData{Funcs: make(map[string]*FuncData, 0), statData: NewStatData()}
}

func (this *PerfData) AddFuncData(funcName, dataName string, value float64) {
	if !this.hasData(dataName) {
		this.DataNames = append(this.DataNames, dataName)
	}

	func1, ok := this.Funcs[funcName]
	if ok {
		func1.AddData(dataName, value)
	} else {
		func2 := NewFuncData()
		func2.Name = funcName
		func2.AddData(dataName, value)
		this.FuncNames = append(this.FuncNames, funcName)
		this.Funcs[funcName] = func2
	}
	return
}

func (this *PerfData) hasData(name string) bool {
	for _, v := range this.DataNames {
		if name == v {
			return true
		}
	}
	return false
}

func (this *PerfData) WriteDataToCsvFile(filename string) bool {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Printf("ERROR: cannot open csv file %s to write\n", filename)
		return false
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Write(append([]string{"FunctionName"}, this.DataNames...))

	for _, v := range this.FuncNames {
		func1, ok := this.Funcs[v]
		if !ok {
			fmt.Printf("ERROR: cannot find func \"%s\"\n", v)
			return false
		}
		line := []string{}
		line = append(line, func1.StringSlice(this.DataNames)...)
		w.Write(line)
	}

	w.Flush()
	return true
}

func (this *PerfData) WriteStatToCsvFile(filename string) bool {
	return this.statData.WriteToCsvFile(filename)
}

func (this *PerfData) ParseAll(statFileName, dataFileName string) bool {
	if !this.ParseStatDataFile(statFileName) {
		return false
	}
	return this.ParseFuncDataFile(dataFileName)
}

func (this *PerfData) ParseFuncDataFile(fileName string) bool {
	parser := &FuncDataParser{perfData: this}
	return parser.ParseFile(fileName)
}

func (this *PerfData) ParseStatDataFile(fileName string) bool {
	parser := &StatDataParser{statData: this.statData}
	return parser.ParseFile(fileName)
}

type FuncDataParser struct {
	perfData        *PerfData
	currentDataName string
}

func (this *FuncDataParser) ParseFile(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("ERROR: cannot open file %s\n", fileName)
		return false
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lineNo := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil && io.EOF != err {
			break
		}

		lineNo++
		if !this.ParseLine(line) {
			fmt.Printf("ERROR: parse failed at line %d\n", lineNo)
			return false
		}

		if io.EOF == err {
			break
		}
	}
	return true
}

func (this *FuncDataParser) ParseLine(line string) bool {
	line = strings.TrimSpace(line)

	if len(line) == 0 {
		return true
	}

	if strings.Contains(line, "Samples") {
		begin := strings.Index(line, "'")
		if begin < 0 {
			fmt.Printf("ERROR: no counter name after \"Samples\"\n")
			return false
		}
		begin++
		end := strings.Index(line[begin:], "'")
		if end < 0 {
			fmt.Printf("ERROR: no ' after counter name\n")
			return false
		}
		end += begin
		this.currentDataName = line[begin:end]
	} else if pos := strings.Index(line, "[.]"); pos > 0 {
		funcName := strings.TrimSpace(line[pos+3:])

		var children, self float64
		n, err := fmt.Sscanf(line, "%f%%  %f%%", &children, &self)
		if err != nil || n != 2 {
			fmt.Printf("ERROR: parse counter value failed\n")
			return false
		}

		this.perfData.AddFuncData(funcName, this.currentDataName+"(children)", children)
		this.perfData.AddFuncData(funcName, this.currentDataName+"(self)", self)
	}

	return true
}

type StatDataParser struct {
	statData *StatData
}

func (this *StatDataParser) ParseFile(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("ERROR: cannot open file %s\n", fileName)
		return false
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lineNo := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil && io.EOF != err {
			break
		}

		lineNo++
		if !this.ParseLine(line) {
			fmt.Printf("ERROR: parse failed at line %d\n", lineNo)
			return false
		}

		if io.EOF == err {
			break
		}
	}
	return true
}

func (this *StatDataParser) ParseLine(line string) bool {
	line = strings.TrimSpace(line)

	if len(line) == 0 {
		return true
	}

	var v1, v2 string

	n, err := fmt.Sscanf(line, "%s %s", &v1, &v2)
	if err != nil || n != 2 {
		fmt.Printf("ERROR: parse stat value failed\n")
		return false
	}

	if v2 == "seconds" {
		v2 = "total-time"
	} else {
		pos := strings.Index(v2, ":")
		if pos > 0 {
			v2 = v2[:pos]
		}
	}

	s1, err := strconv.ParseFloat(strings.Replace(v1, ",", "", -1), 10)
	if err != nil {
		fmt.Printf("ERROR: parse stat value1 failed\n")
		return false
	}

	this.statData.AddData(v2, s1)

	return true
}
