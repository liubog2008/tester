package tester

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/liubog2008/tester/pkg/data"
)

const (
	defaultTestDataDir = "testdata"
)

var defaultRefManager data.RefManager

// Tester defines interface contains main test logic
// Tester should be implemented by struct which contains test case data
type Tester interface {
	// Test is the main test logic
	Test(t *testing.T)
}

// Options defines options of test
type Options struct {
	// Selector select labels of test case
	Selector map[string]string

	// Parallel runs test parallel
	Parallel bool
}

// Test parses testdata and runs test case one by one
func Test(t *testing.T, h Tester) {
	t.Helper()

	internalTest(t, h, nil, defaultTestDataDir, false)
}

// TestWithOptions select matched testdata and runs them
func TestWithOptions(t *testing.T, h Tester, options Options) {
	t.Helper()

	internalTest(t, h, options.Selector, defaultTestDataDir, options.Parallel)
}

func internalTest(t *testing.T, h Tester, selector map[string]string, dataDir string, parallel bool) {
	t.Helper()

	ptr := reflect.ValueOf(h)
	if ptr.Kind() != reflect.Ptr || ptr.IsNil() {
		t.Fatalf("Tester is not pointer of struct or Tester is nil")
	}

	v := ptr.Elem()
	vt := v.Type()
	zeroValue := reflect.Zero(vt)

	_, file, _, ok := runtime.Caller(2)
	if !ok {
		t.Fatalf("can't find caller of tester")
	}

	dir := filepath.Dir(file)

	if !filepath.IsAbs(dataDir) {
		dataDir = filepath.Join(dir, dataDir)
	}

	if defaultRefManager == nil {
		defaultRefManager = data.NewRefManager()
	}

	tcl, err := data.NewTestCaseList(dataDir, t.Name(), defaultRefManager)
	if err != nil {
		t.Fatalf("can't parse test data from files: %v", err)
	}

	tcs := tcl.Select(selector)

	// TODO(liubog2008): try to run tests parallel
	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.Description(), func(t *testing.T) {
			handler := h
			if parallel {
				t.Parallel()

				nh, ok := reflect.New(vt).Interface().(Tester)
				if !ok {
					t.Errorf("can't convert interface to Tester")
				}
				handler = nh
			} else {
				v.Set(zeroValue)
			}
			if err := tc.Unmarshal(handler); err != nil {
				t.Errorf("can't unmarshal data to test case: %v", err)
			}

			handler.Test(t)
		})
	}
}
