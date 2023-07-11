package jsondb

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/fs714/github-star-manager/pkg/config"
	"github.com/pkg/errors"
)

var Jsondb JsonConfig

func InitJsondbFromConfig() (err error) {
	err = InitJsondb(config.Config.Database.Path)
	if err != nil {
		errors.WithMessage(err, "failed to init db from config")
		return
	}

	return
}

func InitJsondb(path string) (err error) {
	Jsondb = JsonConfig{
		Path:         path,
		Common:       &Common{},
		Repositories: NewRepositories(),
	}

	if _, err := os.Stat(path); err == nil {
		err = Jsondb.Read()
		if err != nil {
			return err
		}

		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		return Jsondb.Write()
	} else {
		return errors.Wrap(err, "stat error")
	}
}

type JsonConfig struct {
	Path         string
	Common       *Common
	Repositories *Repositories
	sync.RWMutex
}

func (j *JsonConfig) Read() error {
	data, err := os.ReadFile(j.Path)
	if err != nil {
		return errors.Wrap(err, "read file error")
	}

	err = json.Unmarshal(data, j)
	if err != nil {
		return errors.Wrap(err, "unmarshal error")
	}

	return nil
}

func (j *JsonConfig) Write() error {
	data, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshal error")
	}

	err = os.WriteFile(j.Path, data, 0644)
	if err != nil {
		return errors.Wrap(err, "write file error")
	}

	return nil
}

func (j *JsonConfig) GetGithubToken() string {
	j.RLock()
	defer j.RUnlock()

	return j.Common.GithubToken
}

func (j *JsonConfig) UpdateGithubToken(token string) error {
	j.Lock()
	defer j.Unlock()

	j.Common.GithubToken = token

	return j.Write()
}

func (j *JsonConfig) LoadRepositories(repos *Repositories) error {
	j.Lock()
	defer j.Unlock()

	j.Repositories = repos

	return j.Write()
}

func (j *JsonConfig) AddRepository(path []string, repo *Repository) error {
	j.Repositories.Add(path, repo)

	j.Lock()
	defer j.Unlock()

	return j.Write()
}

func (j *JsonConfig) GetRepositories(path []string) *Repositories {
	return j.Repositories.Get(path)
}

func (j *JsonConfig) GetAllRepositoryByPath(path []string) []*Repository {
	return j.Repositories.GetAllRepositoryByPath(path)
}

func (j *JsonConfig) GetAllRepositoryByTag(tag string) []*Repository {
	return j.Repositories.GetAllRepositoryByTag(tag)
}

func (j *JsonConfig) GetAllTag() []string {
	return j.Repositories.GetAllTag()
}

func (j *JsonConfig) GetAllRepositoryByName(name string) ([]string, int, *Repository) {
	return j.Repositories.GetRepositoryByName(name)
}

func (j *JsonConfig) UpdateRepository(repo *Repository) error {
	err := j.Repositories.Update(repo)
	if err != nil {
		return err
	}

	j.Lock()
	defer j.Unlock()

	return j.Write()
}

func (j *JsonConfig) DeleteRepository(name string) error {
	j.Repositories.Delete(name)

	j.Lock()
	defer j.Unlock()
	return j.Write()
}
