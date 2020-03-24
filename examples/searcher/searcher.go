// Package searcher value of key
package searcher

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Searcher struct {
	store map[string]string
}

func New(kvs []KV) *Searcher {
	s := Searcher{
		store: map[string]string{},
	}
	for _, kv := range kvs {
		s.store[kv.Key] = kv.Value
	}
	return &s
}

func (s *Searcher) Search(key string) (string, bool) {
	v, ok := s.store[key]
	return v, ok
}
