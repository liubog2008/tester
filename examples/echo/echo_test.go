package echo

import (
	"testing"

	"github.com/liubog2008/tester/pkg/tester"
	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	tester.Test(t, new(EchoTestCase))
}

type EchoTestCase struct {
	In       string `json:"in"`
	Expected string `json:"expected"`
}

func (c *EchoTestCase) Test(t *testing.T) {
	assert.Equal(t, c.Expected, Echo(c.In))
}
