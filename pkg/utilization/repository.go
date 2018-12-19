package utilization

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/itsubaki/awsri/internal/costviz"
)

type Repository struct {
	AccountID string             `json:"account_id"`
	Internal  costviz.RecordList `json:"internal"`
}

func NewRepository(path string) (*Repository, error) {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	var records costviz.RecordList
	if err := json.Unmarshal(read, &records); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return &Repository{
		Internal: records,
	}, nil
}

func (r *Repository) SelectAll() costviz.RecordList {
	return r.Internal
}
