package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sort"
	"strings"
)

func loadFileDescriptors(filePaths ...string) (map[string]*desc.FileDescriptor, error) {
	fds := []*dpb.FileDescriptorProto{}
	for _, filePath := range filePaths {
		fs, err := readFileDescriptorSet(filePath)
		if err != nil {
			return nil, err
		}
		fds = append(fds, fs.GetFile()...)
	}
	return desc.CreateFileDescriptors(fds)
}

func readFileDescriptorSet(filePath string) (*dpb.FileDescriptorSet, error) {
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fs := &dpb.FileDescriptorSet{}
	if err := proto.Unmarshal(in, fs); err != nil {
		return nil, err
	}
	return fs, nil
}

var outputFormatFuncs = map[string]formatFunc{
	"yaml": yaml.Marshal,
	"yml":  yaml.Marshal,
	"json": json.Marshal,
	"summary": func(i interface{}) ([]byte, error) {
		return printSummaryTable(i.([]lint.Response))
	},
}

type formatFunc func(interface{}) ([]byte, error)

func getOutputFormatFunc(formatType string) formatFunc {
	if f, found := outputFormatFuncs[strings.ToLower(formatType)]; found {
		return f
	}
	return yaml.Marshal
}

// printSummaryTable returns a summary table of violation counts.
func printSummaryTable(responses []lint.Response) ([]byte, error) {
	s := createSummary(responses)

	data := []summary{}
	for ruleID, fileViolations := range s {
		totalViolations := 0
		for _, count := range fileViolations {
			totalViolations += count
		}
		data = append(data, summary{ruleID, totalViolations, len(fileViolations)})
	}
	sort.SliceStable(data, func(i, j int) bool { return data[i].violations < data[j].violations })

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Rule", "Total Violations", "Violated Files"})
	table.SetCaption(true, fmt.Sprintf("Linted %d proto files", len(responses)))
	for _, d := range data {
		table.Append([]string{
			d.ruleID,
			fmt.Sprintf("%d", d.violations),
			fmt.Sprintf("%d", d.files),
		})
	}
	table.Render()

	return buf.Bytes(), nil
}

func createSummary(responses []lint.Response) map[string]map[string]int {
	summary := make(map[string]map[string]int)
	for _, r := range responses {
		filePath := string(r.FilePath)
		for _, p := range r.Problems {
			ruleID := string(p.RuleID)
			if summary[ruleID] == nil {
				summary[ruleID] = make(map[string]int)
			}
			summary[ruleID][filePath]++
		}
	}
	return summary
}

type summary struct {
	ruleID     string
	violations int
	files      int
}
