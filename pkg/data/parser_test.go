package data

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONParser(t *testing.T) {
	cases := []struct {
		desc     string
		in       []byte
		expected []TestCaseData
	}{
		{
			desc: "normal case",
			in: []byte(`
[
	{
		"description": "xxx",
		"labels": {
			"aaa": "bbb"
		},
		"data": {"xxx":"yyy"}
	},
	{
		"description": "yyy",
		"labels": {
			"ccc": "ddd"
		},
		"data": {"mmm":"nnn"}
	}
]
			`),
			expected: []TestCaseData{
				{
					Description: "xxx",
					Labels: map[string]string{
						"aaa": "bbb",
					},
					Data: map[string]json.RawMessage{
						"xxx": json.RawMessage(`"yyy"`),
					},
				},
				{
					Description: "yyy",
					Labels: map[string]string{
						"ccc": "ddd",
					},
					Data: map[string]json.RawMessage{
						"mmm": json.RawMessage(`"nnn"`),
					},
				},
			},
		},
		{
			desc: "labels is unset",
			in: []byte(`
[{
	"description": "xxx",
	"data": {"xxx":"yyy"}
}]
			`),
			expected: []TestCaseData{
				{
					Description: "xxx",
					Data: map[string]json.RawMessage{
						"xxx": json.RawMessage(`"yyy"`),
					},
				},
			},
		},
		{
			desc: "data is unset",
			in: []byte(`
[{
	"description": "xxx",
	"labels": {
		"aaa": "bbb"
	}
}]
			`),
			expected: []TestCaseData{
				{
					Description: "xxx",
					Labels: map[string]string{
						"aaa": "bbb",
					},
				},
			},
		},
		{
			desc: "data is empty",
			in: []byte(`
[{
	"description": "xxx",
	"labels": {
		"aaa": "bbb"
	},
	"data": {}
}]
			`),
			expected: []TestCaseData{
				{
					Description: "xxx",
					Labels: map[string]string{
						"aaa": "bbb",
					},
					Data: map[string]json.RawMessage{},
				},
			},
		},
	}
	for _, c := range cases {
		ds, err := JSONParser().Parse(c.in)
		assert.NoError(t, err, c.desc)
		assert.Equal(t, c.expected, ds, c.desc)
	}
}

func TestYAMLParser(t *testing.T) {
	cases := []struct {
		desc     string
		in       []byte
		expected []TestCaseData
	}{
		{
			desc: "normal case",
			in: []byte(`
- description: xxx
  labels:
    aaa: bbb
  data:
    xxx: yyy
`,
			),
			expected: []TestCaseData{
				{
					Description: "xxx",
					Labels: map[string]string{
						"aaa": "bbb",
					},
					Data: map[string]json.RawMessage{
						"xxx": json.RawMessage(`"yyy"`),
					},
				},
			},
		},
		{
			desc: "labels is unset",
			in: []byte(`
- description: xxx
  data:
    xxx: yyy
`,
			),
			expected: []TestCaseData{
				{
					Description: "xxx",
					Data: map[string]json.RawMessage{
						"xxx": json.RawMessage(`"yyy"`),
					},
				},
			},
		},
		{
			desc: "data is unset",
			in: []byte(`
- description: xxx
  labels:
    aaa: bbb
`,
			),
			expected: []TestCaseData{
				{
					Description: "xxx",
					Labels: map[string]string{
						"aaa": "bbb",
					},
				},
			},
		},
	}
	for _, c := range cases {
		ds, err := YAMLParser().Parse(c.in)
		assert.NoError(t, err, c.desc)
		assert.Equal(t, c.expected, ds, c.desc)
	}
}
