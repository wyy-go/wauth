# wauth

![GitHub Repo stars](https://img.shields.io/github/stars/wyy-go/wauth?style=social)
![](https://img.shields.io/badge/license-MIT-green)
![GitHub](https://img.shields.io/github/license/wyy-go/wauth)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/wyy-go/wauth)
![GitHub CI Status](https://img.shields.io/github/workflow/status/wyy-go/wauth/ci?label=CI)
[![Go Report Card](https://goreportcard.com/badge/github.com/wyy-go/wauth)](https://goreportcard.com/report/github.com/wyy-go/wauth)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/wyy-go/wauth?tab=doc)
[![codecov](https://codecov.io/gh/wyy-go/wauth/branch/main/graph/badge.svg)](https://codecov.io/gh/wyy-go/wauth)


wauth is an authorization middleware for [Gin](https://github.com/gin-gonic/gin), it's based on [https://github.com/casbin/casbin](https://github.com/casbin/casbin).

## Installation

```bash
go get github.com/wyy-go/wauth
```

## Simple Example

```Go
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
```

## Documentation

The authorization determines a request based on ``{subject, object, action}``, which means what ``subject`` can perform what ``action`` on what ``object``. In this plugin, the meanings are:

1. ``subject``: the logged-on user name
2. ``object``: the URL path for the web resource like "dataset1/item1"
3. ``action``: HTTP method like GET, POST, PUT, DELETE, or the high-level actions you defined like "read-file", "write-blog"

For how to write authorization policy and other details, please refer to [the Casbin's documentation](https://github.com/casbin/casbin).

## Getting Help

- [Casbin](https://github.com/casbin/casbin)
- [Gin-authz](https://github.com/gin-contrib/authz)

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
