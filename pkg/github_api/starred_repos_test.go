package github_api

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetStarredRepos(t *testing.T) {
	repos, err := GetStarredRepos("fs714")
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.FailNow()
	}

	reposJson, err := json.MarshalIndent(repos, "", "  ")
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.FailNow()
	}

	fmt.Println(string(reposJson))
}
