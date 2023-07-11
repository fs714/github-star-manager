package public

import (
	"net/http"

	"github.com/fs714/github-star-manager/db/jsondb"
	"github.com/fs714/github-star-manager/pkg/utils/code"
	"github.com/gin-gonic/gin"
)

func GetRepos(c *gin.Context) {
	repos := jsondb.Jsondb.GetAllRepositoryByPath([]string{})
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status": code.RespCommonError,
	// 		"msg":    msg,
	// 		"data":   "",
	// 	})

	// 	log.Errorf("failed to get repos:\n%+v", err)
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"status": code.RespOk,
		"msg":    "",
		"data":   repos,
	})
}
