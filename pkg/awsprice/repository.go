package awsprice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Repository struct {
	Region   string     `json:"region"`
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

func (r *Repository) FindByInstanceType(tipe string) RecordList {
	out := RecordList{}
	for i := range r.Internal {
		if r.Internal[i].InstanceType == tipe {
			out = append(out, r.Internal[i])
		}
	}

	return out
}

func (r *Repository) FindByUsageType(tipe string) RecordList {
	out := RecordList{}
	for i := range r.Internal {
		if r.Internal[i].UsageType == tipe {
			out = append(out, r.Internal[i])
		}
	}

	return out
}
