# Tester

Tester is a common library to help write testing code.

There are some features:
- data driven
- label selector

## Get Started

```go
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
```
See `examples/echo` for more detail.

## Future Work

- fuzzy case auto generation (not implemented)
