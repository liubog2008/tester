package data

import (
	"io/ioutil"
	"path/filepath"
)

const defaultRefDir = "common"

type refManager struct {
	refs map[string][]byte
}

func NewRefManager() RefManager {
	return &refManager{
		refs: map[string][]byte{},
	}
}

// TODO(liubog2008): recycle references after all tests which
//   use same data dir are finished
func (r *refManager) Ref(dir, ref string) ([]byte, error) {
	path := filepath.Join(dir, defaultRefDir, ref)

	body, ok := r.refs[path]
	if !ok {
		body, err := openRef(path)
		if err != nil {
			return nil, err
		}

		r.refs[path] = body

		return body, nil
	}

	return body, nil
}

func openRef(path string) ([]byte, error) {
	body, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return body, nil
}
