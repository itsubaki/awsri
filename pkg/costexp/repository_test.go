package costexp

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func TestSerialize(t *testing.T) {
	date := []*costexplorer.DateInterval{
		{
			Start: aws.String("2018-01-01"),
			End:   aws.String("2018-02-01"),
		},
		{
			Start: aws.String("2018-02-01"),
			End:   aws.String("2018-03-01"),
		},
		{
			Start: aws.String("2018-03-01"),
			End:   aws.String("2018-04-01"),
		},
		{
			Start: aws.String("2018-04-01"),
			End:   aws.String("2018-05-01"),
		},
		{
			Start: aws.String("2018-05-01"),
			End:   aws.String("2018-06-01"),
		},
		{
			Start: aws.String("2018-06-01"),
			End:   aws.String("2018-07-01"),
		},
		{
			Start: aws.String("2018-07-01"),
			End:   aws.String("2018-08-01"),
		},
		{
			Start: aws.String("2018-08-01"),
			End:   aws.String("2018-09-01"),
		},
		{
			Start: aws.String("2018-09-01"),
			End:   aws.String("2018-10-01"),
		},
		{
			Start: aws.String("2018-10-01"),
			End:   aws.String("2018-11-01"),
		},
		{
			Start: aws.String("2018-11-01"),
			End:   aws.String("2018-12-01"),
		},
		{
			Start: aws.String("2018-12-01"),
			End:   aws.String("2019-01-01"),
		},
	}

	for i := range date {
		repo, err := NewRepository("example", []*costexplorer.DateInterval{date[i]})
		if err != nil {
			t.Errorf("new repository: %v", err)
		}

		start := *date[i].Start
		path := fmt.Sprintf("/var/tmp/hermes/costexp/example_%s.out", start[:7])
		if err := repo.Write(path); err != nil {
			t.Errorf("write file: %v", err)
		}
	}
}

func TestMergedRepository(t *testing.T) {
	dir := "/var/tmp/hermes/costexp"
	path := []string{
		fmt.Sprintf("%s/%s.out", dir, "example_2018-01"),
		fmt.Sprintf("%s/%s.out", dir, "example_2018-02"),
		fmt.Sprintf("%s/%s.out", dir, "example_2018-03"),
		fmt.Sprintf("%s/%s.out", dir, "example_2018-04"),
		fmt.Sprintf("%s/%s.out", dir, "example_2018-05"),
		fmt.Sprintf("%s/%s.out", dir, "example_2018-06"),
		fmt.Sprintf("%s/%s.out", dir, "example_2018-07"),
		fmt.Sprintf("%s/%s.out", dir, "example_2018-08"),
		fmt.Sprintf("%s/%s.out", dir, "example_2018-09"),
	}

	repo := &Repository{
		Profile: "example",
	}

	for _, p := range path {
		tmp, err := Read(p)
		if err != nil {
			t.Errorf("read file: %v", err)
		}

		repo.Internal = append(repo.Internal, tmp.Internal...)
	}

	if len(repo.SelectAll()) < 1 {
		t.Errorf("invalid repository")
	}

	for _, ID := range repo.AccountID() {
		if len(ID) != 12 {
			t.Errorf("invalid AWS AccountID")
		}
	}
}

func TestRepository(t *testing.T) {
	path := fmt.Sprintf("/var/tmp/hermes/costexp/%s.out", "example_2018-09")
	repo, err := Read(path)
	if err != nil {
		t.Errorf("read file: %v", err)
	}

	if len(repo.SelectAll()) < 1 {
		t.Errorf("repository is empty")
	}

	if repo.Profile != "example" {
		t.Errorf("invalid profile")
	}

}
