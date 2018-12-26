package costexp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
)

type Repository struct {
	Profile  string     `json:"profile"`
	Internal RecordList `json:"internal"`
}

func NewRepository(path string) (*Repository, error) {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	var repo Repository
	if err := json.Unmarshal(read, &repo); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return &repo, nil
}

func (r *Repository) SelectAll() RecordList {
	return r.Internal
}

func (r *Repository) AccountID() []string {
	selected := make(map[string]bool)
	for i := range r.Internal {
		selected[r.Internal[i].AccountID] = true
	}

	out := []string{}
	for k := range selected {
		out = append(out, k)
	}

	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
