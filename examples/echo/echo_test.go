package echo

import (
	"testing"

	"github.com/liubog2008/tester/pkg/data"
	"github.com/liubog2008/tester/pkg/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEcho(t *testing.T) {
	tester.Test(t, testEcho)
}

type EchoTestCase struct {
	In       string `json:"in"`
	Expected string `json:"expected"`
}

func testEcho(t *testing.T, tc data.TestCase) {
	c := EchoTestCase{}
	require.NoError(t, tc.Unmarshal(&c))
	assert.Equal(t, c.Expected, Echo(c.In), tc.Description())
}
