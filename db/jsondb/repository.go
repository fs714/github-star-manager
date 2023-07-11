package jsondb

import (
	"errors"
	"sync"
)

type Repository struct {
	Name        string
	Url         string
	Language    string
	StarsCount  int
	ForksCount  int
	Description string
	CreatedAt   int64
	UpdatedAt   int64
	PushedAt    int64
	Tags        []string
}

type RepositoryNameIndex struct {
	Path  []string
	Index int
}

type Repositories struct {
	Repositories    []*Repository
	SubRepositories map[string]*Repositories
	NameIndexes     map[string]*RepositoryNameIndex `json:"-"`
	TagMap          map[string][]*Repository        `json:"-"`
	sync.RWMutex    `json:"-"`
}

func NewRepositories() *Repositories {
	return &Repositories{
		Repositories:    make([]*Repository, 0),
		SubRepositories: make(map[string]*Repositories),
		NameIndexes:     make(map[string]*RepositoryNameIndex),
		TagMap:          make(map[string][]*Repository),
	}
}

func (rs *Repositories) UpdateNameIndexes() {
	rs.RLock()
	defer rs.RUnlock()

	rs.updateNameIndexes([]string{}, rs.NameIndexes)
}

func (rs *Repositories) updateNameIndexes(prePath []string, nameIndexes map[string]*RepositoryNameIndex) {
	for idx, r := range rs.Repositories {
		nameIndexes[r.Name] = &RepositoryNameIndex{
			Path:  prePath,
			Index: idx,
		}
	}

	if len(rs.SubRepositories) > 0 {
		for k, v := range rs.SubRepositories {
			v.updateNameIndexes(append(prePath, k), nameIndexes)
		}
	}
}

func (rs *Repositories) addRepoToNameIndexes(path []string, index int, repo *Repository) {
	rs.NameIndexes[repo.Name] = &RepositoryNameIndex{
		Path:  path,
		Index: index,
	}
}

func (rs *Repositories) delRepoFromNameIndexes(name string) {
	delete(rs.NameIndexes, name)
}

func (rs *Repositories) UpdateTagMap() {
	rs.RLock()
	defer rs.RUnlock()

	rs.updateTagMap(rs.TagMap)
}

func (rs *Repositories) updateTagMap(tagMap map[string][]*Repository) {
	for _, r := range rs.Repositories {
		for _, t := range r.Tags {
			tagMap[t] = append(tagMap[t], r)
		}
	}

	if len(rs.SubRepositories) > 0 {
		for _, v := range rs.SubRepositories {
			v.updateTagMap(tagMap)
		}
	}
}

func (rs *Repositories) addRepoToTagMap(repo *Repository) {
	for _, t := range repo.Tags {
		rs.TagMap[t] = append(rs.TagMap[t], repo)
	}
}

func (rs *Repositories) delRepoFromTagMap(repo *Repository) {
	for _, t := range repo.Tags {
		newRepoList := make([]*Repository, 0)
		for _, r := range rs.TagMap[t] {
			if r.Name != repo.Name {
				newRepoList = append(newRepoList, r)
			}
		}
		rs.TagMap[t] = newRepoList
	}
}

func (rs *Repositories) Add(path []string, repo *Repository) {
	rs.Lock()
	defer rs.Unlock()

	index := rs.add(path, repo)

	rs.addRepoToNameIndexes(path, index, repo)
	rs.addRepoToTagMap(repo)
}

func (rs *Repositories) add(path []string, repo *Repository) int {
	var index int

	if len(path) == 0 {
		rs.Repositories = append(rs.Repositories, repo)
		index = len(rs.Repositories) - 1
	} else {
		if _, ok := rs.SubRepositories[path[0]]; ok {
			index = rs.SubRepositories[path[0]].add(path[1:], repo)
		} else {
			rs.SubRepositories[path[0]] = NewRepositories()
			index = rs.SubRepositories[path[0]].add(path[1:], repo)
		}
	}

	return index
}

func (rs *Repositories) Get(path []string) *Repositories {
	rs.RLock()
	defer rs.RUnlock()

	if len(path) == 0 {
		return rs
	} else {
		if _, ok := rs.SubRepositories[path[0]]; ok {
			return rs.SubRepositories[path[0]].Get(path[1:])
		} else {
			return nil
		}
	}
}

func (rs *Repositories) GetAllRepositoryByPath(path []string) []*Repository {
	repos := make([]*Repository, 0)
	rsTemp := rs.Get(path)

	rs.RLock()
	defer rs.RUnlock()

	if rsTemp != nil {
		repos = append(repos, rsTemp.Repositories...)
		for _, srs := range rsTemp.SubRepositories {
			repos = append(repos, srs.GetAllRepositoryByPath([]string{})...)
		}
	}

	return repos
}

func (rs *Repositories) GetAllRepositoryByTag(tag string) []*Repository {
	rs.RLock()
	defer rs.RUnlock()

	if _, ok := rs.TagMap[tag]; ok {
		return rs.TagMap[tag]
	} else {
		return make([]*Repository, 0)
	}
}

func (rs *Repositories) GetAllTag() []string {
	rs.RLock()
	defer rs.RUnlock()

	tags := make([]string, 0)
	for k := range rs.TagMap {
		tags = append(tags, k)
	}

	return tags
}

func (rs *Repositories) GetRepositoryByName(name string) ([]string, int, *Repository) {
	rs.RLock()
	defer rs.RUnlock()

	if nameIndex, ok := rs.NameIndexes[name]; ok {
		return nameIndex.Path, nameIndex.Index, rs.getRepositoryByNameWithIndexes(nameIndex)
	} else {
		return rs.getRepositoryByNameWithoutIndexes(name)
	}
}

func (rs *Repositories) getRepositoryByNameWithIndexes(indexes *RepositoryNameIndex) *Repository {
	repos := rs.Get(indexes.Path)
	if repos == nil {
		return nil
	}

	return repos.Repositories[indexes.Index]
}

func (rs *Repositories) getRepositoryByNameWithoutIndexes(name string) ([]string, int, *Repository) {
	var path []string
	var index int
	var repo *Repository

	found := false
	for idx, r := range rs.Repositories {
		if r.Name == name {
			repo = r
			index = idx
			found = true
			break
		}
	}

	if found {
		return path, index, repo
	} else {
		for k, v := range rs.SubRepositories {
			p, i, r := v.getRepositoryByNameWithoutIndexes(name)
			if r != nil {
				path = append(path, k)
				path = append(path, p...)
				index = i
				repo = r
				found = true
				break
			}
		}

		if found {
			return path, index, repo
		}
	}

	return path, index, repo
}

func (rs *Repositories) Update(repo *Repository) error {
	path, idx, existRepo := rs.GetRepositoryByName(repo.Name)
	if existRepo == nil {
		return errors.New("repository not found")
	}

	repos := rs.Get(path)
	if repos == nil {
		return errors.New("path not found")
	}

	rs.Lock()
	defer rs.Unlock()

	repos.Repositories[idx] = repo

	rs.delRepoFromTagMap(repo)
	rs.addRepoToTagMap(repo)

	return nil
}

func (rs *Repositories) Delete(name string) {
	path, _, repo := rs.GetRepositoryByName(name)
	if repo == nil {
		return
	}

	repos := rs.Get(path)
	if repos == nil {
		return
	}

	rs.Lock()
	defer rs.Unlock()

	UpdatedRepoList := make([]*Repository, 0)
	for _, r := range repos.Repositories {
		if r.Name != name {
			UpdatedRepoList = append(UpdatedRepoList, r)
		}
	}

	repos.Repositories = UpdatedRepoList

	rs.delRepoFromNameIndexes(name)
	rs.delRepoFromTagMap(repo)
}
