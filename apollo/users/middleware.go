package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"apollo/redis"
)

const cookieName = "exchange_userCookie"

func AuthRequired() func(c *gin.Context) {
	return func(c *gin.Context) {
		cookieValue, err := c.Cookie(cookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No authorization"})
			return
		}

		userId, err := redis.GetClient().Get(cookieValue).Result();
		if err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authorization expired"})
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.Set("user_id", userId)
		}
	}
}
