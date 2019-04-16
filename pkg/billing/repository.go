package billing

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
	c := costexp.New()
	c.Client.Config.WithHTTPClient(client)

	for i := range repo.Date {
		cost, err := c.GetCost(&costexp.Date{
			Start: repo.Date[i].Start,
			End:   repo.Date[i].End,
		})

		if err != nil {
			return fmt.Errorf("get cost: %v", err)
		}

		for _, c := range cost {
			repo.Internal = append(repo.Internal, &Record{
				AccountID:        c.AccountID,
				Description:      c.Description,
				Date:             c.Date,
				AmortizedCost:    c.AmortizedCost,
				BlendedCost:      c.BlendedCost,
				UnblendedCost:    c.UnblendedCost,
				NetAmortizedCost: c.NetAmortizedCost,
				NetUnblendedCost: c.NetUnblendedCost,
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

func (repo *Repository) Description() []string {
	selected := make(map[string]bool)
	for i := range repo.Internal {
		selected[repo.Internal[i].Description] = true
	}

	out := []string{}
	for k := range selected {
		out = append(out, k)
	}

	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func Download(dir string) error {
	path := fmt.Sprintf("%s/billing", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	date := GetCurrentDate()
	for i := range date {
		cache := fmt.Sprintf("%s/%s.out", path, date[i].YYYYMM())
		if _, err := os.Stat(cache); !os.IsNotExist(err) {
			continue
		}

		repo := NewRepository([]*Date{date[i]})
		if err := repo.Fetch(); err != nil {
			return fmt.Errorf("fetch billing (date=%s): %v", date[i], err)
		}

		if err := repo.Write(cache); err != nil {
			return fmt.Errorf("write billing (date=%s): %v", date[i], err)
		}

		fmt.Printf("write: %v\n", cache)
	}

	return nil
}
