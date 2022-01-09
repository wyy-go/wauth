// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package wauth

import (
	"context"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type ctxAuthKey struct{}

type Option func(*Options)

type Options struct {
	errorFn     func(*gin.Context, error)
	forbiddenFn func(*gin.Context)
	skipAuthFn  func(*gin.Context) bool
	subjectFn   func(*gin.Context) string
}

func WithErrorFn(fn func(*gin.Context, error)) Option {
	return func(opts *Options) {
		if fn != nil {
			opts.errorFn = fn
		}
	}
}

func WithForbiddenFn(fn func(*gin.Context)) Option {
	return func(opts *Options) {
		if fn != nil {
			opts.forbiddenFn = fn
		}
	}
}

func WithSkipAuthFn(fn func(*gin.Context) bool) Option {
	return func(opts *Options) {
		if fn != nil {
			opts.skipAuthFn = fn
		}
	}
}

func WithSubjectFn(fn func(*gin.Context) string) Option {
	return func(opts *Options) {
		if fn != nil {
			opts.subjectFn = fn
		}
	}
}

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func NewAuthorizer(e casbin.IEnforcer, opts ...Option) gin.HandlerFunc {
	options := Options{
		func(c *gin.Context, err error) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  http.StatusText(http.StatusInternalServerError),
			})
		},
		func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  http.StatusText(http.StatusForbidden),
			})
		},
		func(c *gin.Context) bool { return false },
		Subject,
	}
	for _, opt := range opts {
		opt(&options)
	}

	a := &BasicAuthorizer{enforcer: e, options: options}

	return func(c *gin.Context) {

		if !a.options.skipAuthFn(c) {
			allowed, err := a.enforcer.Enforce(a.options.subjectFn(c), c.Request.URL.Path, c.Request.Method)
			if err != nil {
				a.options.errorFn(c, err)
				return
			}
			if !allowed {
				a.options.forbiddenFn(c)
				return
			}
		}

		c.Next()
	}
}

// BasicAuthorizer stores the casbin handler
type BasicAuthorizer struct {
	enforcer casbin.IEnforcer
	options  Options
}

func Subject(c *gin.Context) string {
	val, _ := c.Request.Context().Value(ctxAuthKey{}).(string)
	return val
}

func CtxWithSubject(c *gin.Context, subject string) {
	ctx := context.WithValue(c.Request.Context(), ctxAuthKey{}, subject)
	c.Request = c.Request.WithContext(ctx)
}
