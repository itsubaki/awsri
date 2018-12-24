package forecast

import (
	"fmt"
	"os"
	"testing"

	"github.com/itsubaki/awsri/internal/costexp"
)

func TestMergedRepository(t *testing.T) {
	dir := fmt.Sprintf(
		"%s/%s",
		os.Getenv("GOPATH"),
		"src/github.com/itsubaki/awsri/internal/_serialized/costexp",
	)

	path := []string{
		fmt.Sprintf("%s/%s", dir, "example_2017-12.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-01.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-02.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-03.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-04.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-05.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-06.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-07.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-08.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-09.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-10.out"),
		fmt.Sprintf("%s/%s", dir, "example_2018-11.out"),
	}

	repo := &costexp.Repository{
		Profile: "example",
	}

	for _, p := range path {
		p, err := costexp.NewRepository(p)
		if err != nil {
			t.Errorf("new costexp repository: %v", err)
		}
		repo.Internal = append(repo.Internal, p.Internal...)
	}

	for _, r := range repo.SelectAll().Sort() {
		fmt.Println(r)
	}
}
