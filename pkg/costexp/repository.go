package costexp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
