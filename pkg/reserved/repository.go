package reserved

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Repository struct {
	Region   []string   `json:"region"`
	Internal RecordList `json:"internal"`
}

func New(region []string) (*Repository, error) {
	repo := NewRepository(region)
	return repo, repo.Fetch()
}

func NewRepository(region []string) *Repository {
	return &Repository{
		Region: region,
	}
}

func (repo *Repository) Fetch() error {
	return repo.FetchWithClient(http.DefaultClient)
}

func (repo *Repository) FetchWithClient(client *http.Client) error {
	for _, f := range NewFetchReservedRecordList() {
		if err := repo.fetchWithClient(client, f); err != nil {
			return fmt.Errorf("fetch with client: %v", err)
		}
	}
	return nil
}

func (repo *Repository) fetchWithClient(client *http.Client, fetch FetchReservedRecordList) error {
	for _, r := range repo.Region {
		ses, err := session.NewSession(
			&aws.Config{
				Region:     aws.String(r),
				HTTPClient: client,
			},
		)
		if err != nil {
			return fmt.Errorf("new session (region=%s): %v", r, err)
		}

		rs, err := fetch(ses, r)
		if err != nil {
			return fmt.Errorf("query: %v", err)
		}
		repo.Internal = append(repo.Internal, rs...)
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

func (repo *Repository) FindByInstanceType(tipe string) RecordList {
	out := RecordList{}
	for i := range repo.Internal {
		if repo.Internal[i].InstanceType == tipe {
			out = append(out, repo.Internal[i])
		}
	}

	return out
}

func Download(region []string, dir string) error {
	path := fmt.Sprintf("%s/reserved.out", dir)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}

	repo := NewRepository(region)
	if err := repo.Fetch(); err != nil {
		return fmt.Errorf("fetch reservation: %v", err)
	}

	if err := repo.Write(path); err != nil {
		return fmt.Errorf("write reservation: %v", err)
	}

	fmt.Printf("write: %v\n", path)
	return nil
}
