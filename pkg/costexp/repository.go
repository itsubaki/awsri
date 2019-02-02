package costexp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/itsubaki/hermes/internal/costexp"
)

type Repository struct {
	Date     []*Date
	Internal RecordList `json:"internal"`
}

func New(date []*Date) (*Repository, error) {
	repo := NewRepository(date)
	return repo, repo.Fetch()
}

func NewRepository(date []*Date) *Repository {
	return &Repository{
		Date: date,
	}
}

func (repo *Repository) Fetch() error {
	return repo.FetchWithClient(http.DefaultClient)
}

func (repo *Repository) FetchWithClient(client *http.Client) error {
	cli := costexp.New()
	cli.Client.Config.WithHTTPClient(client)
	for i := range repo.Date {
		q, err := cli.GetUsageQuantity(&costexp.Date{
			Start: repo.Date[i].Start,
			End:   repo.Date[i].End,
		})
		if err != nil {
			return fmt.Errorf("get usage quantity: %v", err)
		}

		for _, qq := range q {
			repo.Internal = append(repo.Internal, &Record{
				AccountID:      qq.AccountID,
				Description:    qq.Description,
				Date:           qq.Date,
				Region:         qq.Region,
				UsageType:      qq.UsageType,
				Platform:       qq.Platform,
				CacheEngine:    qq.CacheEngine,
				DatabaseEngine: qq.DatabaseEngine,
				InstanceHour:   qq.InstanceHour,
				InstanceNum:    qq.InstanceNum,
			})
		}
	}

	return nil
}

func Read(path string) (*Repository, error) {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	repo := &Repository{}
	if err := repo.Deserialize(read); err != nil {
		return nil, fmt.Errorf("new repository: %v", err)
	}

	return repo, nil
}

func (repo *Repository) Write(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}

	bytes, err := repo.Serialize()
	if err != nil {
		return fmt.Errorf("serialize: %v", err)
	}

	if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}

func (repo *Repository) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(repo)
	if err != nil {
		return []byte{}, fmt.Errorf("marshal: %v", err)
	}

	return bytes, nil
}

func (repo *Repository) Deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, repo); err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	return nil
}

func (repo *Repository) SelectAll() RecordList {
	return repo.Internal
}

func (repo *Repository) AccountID() []string {
	selected := make(map[string]bool)
	for i := range repo.Internal {
		selected[repo.Internal[i].AccountID] = true
	}

	out := []string{}
	for k := range selected {
		out = append(out, k)
	}

	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
