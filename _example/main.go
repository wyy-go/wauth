package main

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"github.com/wyy-go/wauth"
)

func main() {
	e, err := casbin.NewEnforcer("../casbin_model.conf", "../casbin_policy.csv")
	if err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(
		func(c *gin.Context) {
			wauth.CtxWithSubject(c, "alice")
		},
		wauth.NewAuthorizer(e),
	)
	router.GET("/dataset1/resource1", func(c *gin.Context) {
		c.String(http.StatusOK, "alice own this resource")
	})
	router.GET("/dataset2/resource1", func(c *gin.Context) {
		c.String(http.StatusOK, "alice do not own this resource")
	})
	router.Run(":8080")
}
