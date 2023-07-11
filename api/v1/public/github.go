package public

import (
	"net/http"

	"github.com/fs714/github-star-manager/db/jsondb"
	"github.com/fs714/github-star-manager/pkg/github_api"
	"github.com/fs714/github-star-manager/pkg/utils/code"
	"github.com/fs714/github-star-manager/pkg/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func SyncFromGithub(c *gin.Context) {
	msg, err := doSyncFromGithub(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": code.RespCommonError,
			"msg":    msg,
			"data":   "",
		})

		log.Errorf("failed to sync from github:\n%+v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": code.RespOk,
		"msg":    "",
		"data":   "",
	})
}

func doSyncFromGithub(c *gin.Context) (string, error) {
	var msg string

	var postData = struct {
		User string
	}{}
	err := c.ShouldBindJSON(&postData)
	if err != nil {
		msg = "failed to bind post json to struct"
		err = errors.Wrap(err, msg)
		return msg, err
	}

	repos, err := github_api.GetStarredRepos(postData.User)
	if err != nil {
		msg = "failed to get starred repos from github"
		err = errors.Wrap(err, msg)
		return msg, err
	}

	newRepos := jsondb.NewRepositories()
	for _, repo := range repos {
		path, _, r := jsondb.Jsondb.GetAllRepositoryByName(*repo.Repository.FullName)
		if r != nil {
			if repo.Repository.FullName != nil {
				r.Name = *repo.Repository.FullName
			}

			if repo.Repository.HTMLURL != nil {
				r.Url = *repo.Repository.HTMLURL
			}

			if repo.Repository.Language != nil {
				r.Language = *repo.Repository.Language
			}

			if repo.Repository.StargazersCount != nil {
				r.StarsCount = *repo.Repository.StargazersCount
			}

			if repo.Repository.ForksCount != nil {
				r.ForksCount = *repo.Repository.ForksCount
			}

			if repo.Repository.Description != nil {
				r.Description = *repo.Repository.Description
			}

			r.CreatedAt = repo.Repository.CreatedAt.Unix()
			r.UpdatedAt = repo.Repository.UpdatedAt.Unix()
			r.PushedAt = repo.Repository.PushedAt.Unix()

			newRepos.Add(path, r)
		} else {
			r := &jsondb.Repository{
				CreatedAt: repo.Repository.CreatedAt.Unix(),
				UpdatedAt: repo.Repository.UpdatedAt.Unix(),
				PushedAt:  repo.Repository.PushedAt.Unix(),
			}

			if repo.Repository.FullName != nil {
				r.Name = *repo.Repository.FullName
			}

			if repo.Repository.HTMLURL != nil {
				r.Url = *repo.Repository.HTMLURL
			}

			if repo.Repository.Language != nil {
				r.Language = *repo.Repository.Language
			}

			if repo.Repository.StargazersCount != nil {
				r.StarsCount = *repo.Repository.StargazersCount
			}

			if repo.Repository.ForksCount != nil {
				r.ForksCount = *repo.Repository.ForksCount
			}

			if repo.Repository.Description != nil {
				r.Description = *repo.Repository.Description
			}

			newRepos.Add([]string{}, r)
		}
	}

	err = jsondb.Jsondb.LoadRepositories(newRepos)
	if err != nil {
		msg = "failed to load new repositories to db"
		err = errors.Wrap(err, msg)
		return msg, err
	}

	return msg, err
}
