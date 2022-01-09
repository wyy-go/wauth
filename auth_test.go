package wauth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func testAuthzRequest(t *testing.T, router *gin.Engine, user string, path string, method string, code int) {
	r, _ := http.NewRequest(method, path, nil)
	r.SetBasicAuth(user, "123")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != code {
		t.Errorf("%s, %s, %s: %d, supposed to be %d", user, path, method, w.Code, code)
	}
}

func TestBasic(t *testing.T) {
	router := gin.New()
	e, _ := casbin.NewEnforcer("casbin_model.conf", "casbin_policy.csv")
	router.Use(func(context *gin.Context) {
		CtxWithSubject(context, "alice")
	})
	router.Use(NewAuthorizer(e, WithSubjectFn(Subject)))
	router.Any("/*anypath", func(c *gin.Context) {
		c.Status(200)
	})

	testAuthzRequest(t, router, "alice", "/dataset1/resource1", "GET", 200)
	testAuthzRequest(t, router, "alice", "/dataset1/resource1", "POST", 200)
	testAuthzRequest(t, router, "alice", "/dataset1/resource2", "GET", 200)
	testAuthzRequest(t, router, "alice", "/dataset1/resource2", "POST", 403)
}

func TestPathWildcard(t *testing.T) {
	router := gin.New()
	e, _ := casbin.NewEnforcer("casbin_model.conf", "casbin_policy.csv")
	router.Use(func(context *gin.Context) {
		CtxWithSubject(context, "bob")
	})
	router.Use(NewAuthorizer(e, WithSubjectFn(Subject)))
	router.Any("/*anypath", func(c *gin.Context) {
		c.Status(200)
	})

	testAuthzRequest(t, router, "bob", "/dataset2/resource1", "GET", 200)
	testAuthzRequest(t, router, "bob", "/dataset2/resource1", "POST", 200)
	testAuthzRequest(t, router, "bob", "/dataset2/resource1", "DELETE", 200)
	testAuthzRequest(t, router, "bob", "/dataset2/resource2", "GET", 200)
	testAuthzRequest(t, router, "bob", "/dataset2/resource2", "POST", 403)
	testAuthzRequest(t, router, "bob", "/dataset2/resource2", "DELETE", 403)

	testAuthzRequest(t, router, "bob", "/dataset2/folder1/item1", "GET", 403)
	testAuthzRequest(t, router, "bob", "/dataset2/folder1/item1", "POST", 200)
	testAuthzRequest(t, router, "bob", "/dataset2/folder1/item1", "DELETE", 403)
	testAuthzRequest(t, router, "bob", "/dataset2/folder1/item2", "GET", 403)
	testAuthzRequest(t, router, "bob", "/dataset2/folder1/item2", "POST", 200)
	testAuthzRequest(t, router, "bob", "/dataset2/folder1/item2", "DELETE", 403)
}

func TestRBAC(t *testing.T) {
	router := gin.New()
	e, _ := casbin.NewEnforcer("casbin_model.conf", "casbin_policy.csv")
	router.Use(func(context *gin.Context) {
		CtxWithSubject(context, "cathy")
	})

	router.Use(NewAuthorizer(e,
		WithSubjectFn(Subject),
		WithErrorFn(func(context *gin.Context, err error) {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  http.StatusText(http.StatusInternalServerError),
			})
		}),
		WithForbiddenFn(func(context *gin.Context) {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  http.StatusText(http.StatusForbidden),
			})
		})),
	)
	router.Any("/*anypath", func(c *gin.Context) {
		c.Status(200)
	})

	// cathy can access all /dataset1/* resources via all methods because it has the dataset1_admin role.
	testAuthzRequest(t, router, "cathy", "/dataset1/item", "GET", 200)
	testAuthzRequest(t, router, "cathy", "/dataset1/item", "POST", 200)
	testAuthzRequest(t, router, "cathy", "/dataset1/item", "DELETE", 200)
	testAuthzRequest(t, router, "cathy", "/dataset2/item", "GET", 403)
	testAuthzRequest(t, router, "cathy", "/dataset2/item", "POST", 403)
	testAuthzRequest(t, router, "cathy", "/dataset2/item", "DELETE", 403)

	// delete all roles on user cathy, so cathy cannot access any resources now.
	_, err := e.DeleteRolesForUser("cathy")
	if err != nil {
		t.Errorf("got error %v", err)
	}

	testAuthzRequest(t, router, "cathy", "/dataset1/item", "GET", 403)
	testAuthzRequest(t, router, "cathy", "/dataset1/item", "POST", 403)
	testAuthzRequest(t, router, "cathy", "/dataset1/item", "DELETE", 403)
	testAuthzRequest(t, router, "cathy", "/dataset2/item", "GET", 403)
	testAuthzRequest(t, router, "cathy", "/dataset2/item", "POST", 403)
	testAuthzRequest(t, router, "cathy", "/dataset2/item", "DELETE", 403)
}
