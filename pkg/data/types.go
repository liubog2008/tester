package data

import "encoding/json"

// TestCaseFileFormat defines format of test case
type TestCaseFileFormat string

const (
	// JSONFormat defines json format
	JSONFormat TestCaseFileFormat = ".json"
	// YAMLFormat defines yaml format
	YAMLFormat TestCaseFileFormat = ".yaml"
)

// TestCase defines a test case interface
type TestCase interface {
	// Match test whether case matches the labels
	Match(labels map[string]string) bool
	// Description returns description of test case
	Description() string
	// Unmarshal decode case data into obj
	Unmarshal(obj interface{}) error
}

// TestCaseList dfeines list of test case
type TestCaseList interface {
	// Select select cases by labels
	Select(labels map[string]string) []TestCase
}

// TestCaseData defines data of test case
type TestCaseData struct {
	// Description is case description
	Description string `json:"description"`
	// Labels defines labels of case
	Labels map[string]string `json:"labels,omitempty"`
	// Data defines raw custom defined data of case
	Data map[string]json.RawMessage `json:"data,omitempty"`
	// References defines references of external docs
	References []Reference `json:"refs,omitempty"`
}

// Reference defines ref of external docs
type Reference struct {
	Name string `json:"name"`
	Ref  string `json:"ref"`
}

// TestCaseParser defines a parser to parse test case from bytes
type TestCaseParser interface {
	// Parse parse bytes to test cases
	Parse(body []byte) ([]TestCaseData, error)
}

// RefManager defines manager to manage refs of data
type RefManager interface {
	// Ref returns data referenced
	Ref(dir, ref string) ([]byte, error)
}
