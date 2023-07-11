package jsondb

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRepositoriesUpdateNameIndexes(t *testing.T) {
	repos := GenerateRepositories()

	repos.UpdateNameIndexes()

	nameIndexesJson, err := json.MarshalIndent(repos.NameIndexes, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(nameIndexesJson))
}

func TestRepositoriesUpdateTagMap(t *testing.T) {
	repos := GenerateRepositories()

	repos.UpdateTagMap()

	tagMapJson, err := json.MarshalIndent(repos.TagMap, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(tagMapJson))
}

func TestRepositoriesAdd(t *testing.T) {
	repos := GenerateRepositories()

	repoLinuxEbpf03 := Repository{
		Name:       "linux_ebpf_03",
		Url:        "http://linux_ebpf_03.com",
		Language:   "golang",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"linux", "ebpf"},
	}

	repoCloudNetwork01 := Repository{
		Name:       "cloud_network_01",
		Url:        "http://cloud_network_01.com",
		Language:   "golang",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"cloud", "network"},
	}

	repos.Add(repoLinuxEbpf03.Tags, &repoLinuxEbpf03)
	repos.Add(repoCloudNetwork01.Tags, &repoCloudNetwork01)

	reposJson, err := json.MarshalIndent(&repos, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(reposJson))

	nameIndexesJson, err := json.MarshalIndent(repos.NameIndexes, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(nameIndexesJson))

	tagMapJson, err := json.MarshalIndent(repos.TagMap, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(tagMapJson))
}

func TestRepositoriesGet(t *testing.T) {
	repos := GenerateRepositories()

	srs := repos.Get([]string{"linux", "ebpf"})
	srsJson, err := json.MarshalIndent(srs, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(srsJson))
}

func TestRepositoriesGetAllByPath(t *testing.T) {
	repos := GenerateRepositories()

	allRepos := repos.GetAllRepositoryByPath([]string{"linux"})
	allReposJson, err := json.MarshalIndent(allRepos, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(allReposJson))
}

func TestRepositoriesGetAllByTag(t *testing.T) {
	repos := GenerateRepositories()

	allRepos := repos.GetAllRepositoryByTag("linux")
	allReposJson, err := json.MarshalIndent(allRepos, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(allReposJson))
}

func TestRepositoriesGetAllTags(t *testing.T) {
	repos := GenerateRepositories()

	allTags := repos.GetAllTag()
	allTagsJson, err := json.MarshalIndent(allTags, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(allTagsJson))
}

func TestRepositoriesGetByNameWithIndexes(t *testing.T) {
	repos := GenerateRepositories()

	repo := repos.getRepositoryByNameWithIndexes(&RepositoryNameIndex{
		Path:  []string{"linux", "ebpf"},
		Index: 0,
	})
	RepoJson, err := json.MarshalIndent(repo, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(RepoJson))
}

func TestRepositoriesGetByNameWithoutIndexes(t *testing.T) {
	repos := GenerateRepositories()

	path, idx, repo := repos.getRepositoryByNameWithoutIndexes("linux_ebpf_01")
	RepoJson, err := json.MarshalIndent(repo, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(path)
	fmt.Println(idx)
	fmt.Println(string(RepoJson))
}

func TestRepositoryUpdate(t *testing.T) {
	repos := GenerateRepositories()

	_, _, repo := repos.GetRepositoryByName("linux_ebpf_01")
	preTags := repo.Tags
	newTags := append(preTags, "golang")
	newRepo := *repo
	newRepo.Tags = newTags

	repos.Update(&newRepo)

	reposJson, err := json.MarshalIndent(&repos, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(reposJson))

	nameIndexesJson, err := json.MarshalIndent(repos.NameIndexes, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(nameIndexesJson))

	tagMapJson, err := json.MarshalIndent(repos.TagMap, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(tagMapJson))
}

func TestRepositoryDelete(t *testing.T) {
	repos := GenerateRepositories()

	_, _, repo := repos.GetRepositoryByName("linux_ebpf_01")

	repos.Delete(repo.Name)

	reposJson, err := json.MarshalIndent(&repos, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(reposJson))

	nameIndexesJson, err := json.MarshalIndent(repos.NameIndexes, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(nameIndexesJson))

	tagMapJson, err := json.MarshalIndent(repos.TagMap, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(tagMapJson))
}

func GenerateRepositories() Repositories {
	repoLinux01 := Repository{
		Name:       "linux01",
		Url:        "http://linux01.com",
		Language:   "c",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"linux"},
	}

	repoLinux02 := Repository{
		Name:       "linux02",
		Url:        "http://linux02.com",
		Language:   "c",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"linux"},
	}

	repoLinuxEbpf01 := Repository{
		Name:       "linux_ebpf_01",
		Url:        "http://linux_ebpf_01.com",
		Language:   "golang",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"linux", "ebpf"},
	}

	repoLinuxEbpf02 := Repository{
		Name:       "linux_ebpf_02",
		Url:        "http://linux_ebpf_02.com",
		Language:   "golang",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"linux", "ebpf"},
	}

	repoLinuxProxy01 := Repository{
		Name:       "linux_proxy_01",
		Url:        "http://linux_proxy_01.com",
		Language:   "golang",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"linux", "proxy"},
	}

	repoLinuxProxy02 := Repository{
		Name:       "linux_proxy_02",
		Url:        "http://linux_proxy_02.com",
		Language:   "golang",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"linux", "proxy"},
	}

	repoAi01 := Repository{
		Name:       "ai_01",
		Url:        "http://ai_01.com",
		Language:   "python",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"ai"},
	}

	repoAiPicture01 := Repository{
		Name:       "ai_picture_01",
		Url:        "http://ai_picture_01.com",
		Language:   "python",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"ai", "picture"},
	}

	repoAiSound01 := Repository{
		Name:       "ai_sound_01",
		Url:        "http://ai_sound_01.com",
		Language:   "python",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"ai", "sound"},
	}

	repoTool01 := Repository{
		Name:       "tool_01",
		Url:        "http://tool_01.com",
		Language:   "c++",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{"tool"},
	}

	repoUnclassified01 := Repository{
		Name:       "unclassified_01",
		Url:        "http://unclassified_01.com",
		Language:   "c++",
		StarsCount: 100,
		ForksCount: 200,
		CreatedAt:  123,
		UpdatedAt:  456,
		PushedAt:   789,
		Tags:       []string{},
	}

	repos := NewRepositories()
	repos.Add(repoLinux01.Tags, &repoLinux01)
	repos.Add(repoLinux02.Tags, &repoLinux02)
	repos.Add(repoLinuxEbpf01.Tags, &repoLinuxEbpf01)
	repos.Add(repoLinuxEbpf02.Tags, &repoLinuxEbpf02)
	repos.Add(repoLinuxProxy01.Tags, &repoLinuxProxy01)
	repos.Add(repoLinuxProxy02.Tags, &repoLinuxProxy02)
	repos.Add(repoAi01.Tags, &repoAi01)
	repos.Add(repoAiPicture01.Tags, &repoAiPicture01)
	repos.Add(repoAiSound01.Tags, &repoAiSound01)
	repos.Add(repoTool01.Tags, &repoTool01)
	repos.Add(repoUnclassified01.Tags, &repoUnclassified01)

	return *repos
}
