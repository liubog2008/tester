package searcher

import (
	"testing"

	"github.com/liubog2008/tester/pkg/tester"
	"github.com/stretchr/testify/assert"
)

func TestSearcher(t *testing.T) {
	tester.Test(t, new(SearcherCase))
}

type SearcherCase struct {
	KVs []KV `json:"kvs"`

	Key string `json:"key"`

	Found bool `json:"found"`

	Value string `json:"value"`
}

func (c *SearcherCase) Test(t *testing.T) {
	s := New(c.KVs)
	v, ok := s.Search(c.Key)
	assert.Equal(t, c.Found, ok)
	assert.Equal(t, c.Value, v)
}
