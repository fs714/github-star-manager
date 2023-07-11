package middleware

import (
	"fmt"
	"time"

	"github.com/fs714/github-star-manager/pkg/config"
	"github.com/fs714/github-star-manager/pkg/utils/log"
	"github.com/gin-gonic/gin"
)

func LogWithSkipPath(skipPath []string) gin.HandlerFunc {
	skipPathMap := make(map[string]bool, len(skipPath))
	for _, path := range skipPath {
		skipPathMap[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		query := c.Request.URL.RawQuery
		c.Next()

		if _, ok := skipPathMap[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)

			statusCode := c.Writer.Status()
			clientIP := c.ClientIP()
			method := c.Request.Method
			userAgent := c.Request.UserAgent()
			errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			if log.ParseFormat(config.Config.Logging.Format) == log.JsonFormat {
				log.CurrentLog().ZapSugaredLogger.Infow("GIN",
					"StatusCode", statusCode,
					"Latency", latency,
					"ClientIP", clientIP,
					"Method", method,
					"Path", path,
					"Query", query,
					"UserAgent", userAgent,
					"ErrorMessage", errorMessage,
				)
			} else {
				var consoleString string
				if errorMessage == "" {
					consoleString = fmt.Sprintf("[GIN] | %3d | %13v | %15s | %-7s %#v",
						statusCode,
						latency,
						clientIP,
						method,
						path,
					)
				} else {
					consoleString = fmt.Sprintf("[GIN] | %3d | %13v | %15s |%-7s %#v\n%s",
						statusCode,
						latency,
						clientIP,
						method,
						path,
						errorMessage,
					)
				}

				log.CurrentLog().ZapSugaredLogger.Infof(consoleString)
			}
		}
	}
}
