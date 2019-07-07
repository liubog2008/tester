# Tester

Tester is a common library to help write testing code.

There are some features:
- data driven
- label selector

## Get Started

```go
func TestEcho(t *testing.T) {
	tester.Test(t, testEcho)
}

type EchoTestCase struct {
	In       string `json:"in"`
	Expected string `json:"expected"`
}

func testEcho(t *testing.T, tc data.TestCase) {
	c := EchoTestCase{}
	require.NoError(t, tc.Load(&c))
	assert.Equal(t, c.Expected, Echo(c.In), tc.Description())
}
```
See `examples/echo` for more detail.

## Future Work

- fuzzy case auto generation (not implemented)
