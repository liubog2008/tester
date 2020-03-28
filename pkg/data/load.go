package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ghodss/yaml"
)

// testCase implements TestCase
type testCase struct {
	dataDir string
	format  TestCaseFileFormat
	data    TestCaseData

	refManager RefManager
}

// Description implements data.TestCase
func (c *testCase) Description() string {
	return c.data.Description
}

// Unmarshal implements data.TestCase
func (c *testCase) Unmarshal(obj interface{}) error {
	switch c.format {
	case JSONFormat:
		merged, err := c.merge()
		if err != nil {
			return err
		}

		body, err := json.Marshal(merged)
		if err != nil {
			return err
		}

		return json.Unmarshal(body, obj)
	case YAMLFormat:
		merged, err := c.merge()
		if err != nil {
			return err
		}

		body, err := yaml.Marshal(merged)
		if err != nil {
			return err
		}

		return yaml.Unmarshal(body, obj)
	}

	return fmt.Errorf("unrecognized format: %v, only support json and yaml now", c.format)
}

func (c *testCase) merge() (map[string]json.RawMessage, error) {
	merged := map[string]json.RawMessage{}
	for k, v := range c.data.Data {
		merged[k] = v
	}
	for _, ref := range c.data.References {
		body, err := c.refManager.Ref(c.dataDir, ref.Ref)
		if err != nil {
			return nil, fmt.Errorf("invalid reference %s named %s: %v", ref.Ref, ref.Name, err)
		}
		_, exist := merged[ref.Name]
		if exist {
			return nil, fmt.Errorf("%s has been defined in data", ref.Name)
		}
		var raw json.RawMessage
		switch c.format {
		case JSONFormat:
			if err := json.Unmarshal(body, &raw); err != nil {
				return nil, fmt.Errorf("cant' unmarshal ref %s: %v", ref.Name, err)
			}
		case YAMLFormat:
			if err := yaml.Unmarshal(body, &raw); err != nil {
				return nil, fmt.Errorf("cant' unmarshal ref %s: %v", ref.Name, err)
			}
		}
		merged[ref.Name] = raw
	}
	return merged, nil
}

// Match implements data.TestCase
func (c *testCase) Match(labels map[string]string) bool {
	return contains(c.data.Labels, labels)
}

type testCaseList struct {
	items []TestCase
}

// NewTestCaseList parses multiple files and returns list of test case
func NewTestCaseList(dataDir, callerName string, refManager RefManager) (TestCaseList, error) {
	files, err := findFiles(dataDir, callerName)
	if err != nil {
		return nil, err
	}

	var parser TestCaseParser
	cl := &testCaseList{}
	for _, file := range files {
		ext := TestCaseFileFormat(filepath.Ext(file))
		switch ext {
		case JSONFormat:
			parser = JSONParser()
		case YAMLFormat:
			parser = YAMLParser()
		default:
			return nil, fmt.Errorf("can't find parser for %v", ext)
		}
		body, err := ioutil.ReadFile(filepath.Clean(file))
		if err != nil {
			return nil, err
		}
		cs, err := parser.Parse(body)
		if err != nil {
			return nil, err
		}
		for _, c := range cs {
			cl.items = append(cl.items, &testCase{
				dataDir:    dataDir,
				format:     ext,
				data:       c,
				refManager: refManager,
			})
		}
	}

	return cl, nil
}

func (cl *testCaseList) Select(labels map[string]string) []TestCase {
	cs := []TestCase{}
	for _, item := range cl.items {
		if item.Match(labels) {
			cs = append(cs, item)
		}
	}
	return cs
}

// if all KVs in b are also in a, return true
func contains(a, b map[string]string) bool {
	for k, v := range b {
		if a == nil {
			return false
		}
		av, ok := a[k]
		if !ok {
			return false
		}
		if av != v {
			return false
		}
	}
	return true
}

// findFile finds testdata from pkg dir
// e.g.
//      dataDir: /go/github.com/src/xxx/yyy/testdata
//   callerName: TestEcho
// There are two optional way to find testdata files
//   1. First find file with json or yaml extension, it will be
//      /go/github.com/src/xxx/yyy/testdata/TestEcho.[yaml|json]
//   2. if not find, all files in dir
//      /go/github.com/src/xxx/yyy/testdata/TestEcho
//      will be returned
func findFiles(dataDir, callerName string) ([]string, error) {
	filePrefix := filepath.Join(dataDir, callerName)
	matchedFiles, err := filepath.Glob(filePrefix + ".*")

	if err != nil {
		return nil, err
	}

	if len(matchedFiles) > 1 {
		return nil, fmt.Errorf("find more than one matched file: %v", matchedFiles[0])
	}

	if len(matchedFiles) == 1 {
		return matchedFiles, nil
	}

	all := filepath.Join(filePrefix, "*")

	filesInDir, err := filepath.Glob(all)
	if err != nil {
		return nil, err
	}

	if len(filesInDir) == 0 {
		return nil, fmt.Errorf("can't find any files in %v", all)
	}

	return filesInDir, nil
}
