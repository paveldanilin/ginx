package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paveldanilin/ginx"
	"github.com/paveldanilin/ginx/requestbody"
	"github.com/paveldanilin/ginx/resolver"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var controller *ginx.Controller

func loadUsers(req *http.Request) []user {
	if req.Header.Get("token") == "1234567890" {
		return []user{{Login: "user1"}, {Login: "user2"}}
	}
	return []user{}
}

func init() {
	controller = ginx.NewController(gin.Default())
	controller.ContentType = gin.MIMEXML

	controller.Use(ginx.ErrorInterceptorFunc(func(e ginx.Error) {
		// Decorate error
		e.Response().SetBody(fmt.Sprintf("[[%s]]", e.Error().Error()))
	}))
	controller.Use(resolver.HttpRequest()) // injects *http.Request
	controller.Use(resolver.Struct())      // injects an instance of the given struct (populated from the income request)

	// /hello/:user
	controller.GET("/hello/:user", func(userName string) string {
		return "Hello, " + userName + "!"
	}, resolver.Path("user", 1))

	// [GET] /posts
	controller.GET("/posts", func(page int) ([]blogPost, int) {
		if page == 33 {
			return []blogPost{
				{Title: "First post", Content: "Hello, world"},
				{Title: "Monday", Content: "This is monday"},
			}, 200
		}
		return []blogPost{}, 404
	},
		resolver.Query("page", 1),
		ginx.Produce(gin.MIMEJSON), // <- overwrite controller ContentType
	)

	// [GET] /users
	controller.GET("/users", loadUsers, ginx.Produce(gin.MIMEPlain))

	// [POST] /users
	controller.POST("/users", func(body requestbody.JSONData) string {
		userID := int(body["id"].(float64))
		return fmt.Sprintf("<%s:%d>", body["username"], userID)
	})

	// [POST] /orders
	controller.POST("/orders", func(o order) string {
		return fmt.Sprintf("<%s:%d:%s>", o.Product, o.ID, o.Extra)
	})

	// [GET] /error
	controller.GET("/error", func() {
		panic(fmt.Errorf("this is error from the '/error'"))
	})

	// [GET] /resolve/PathVariable?var2=1234&var3=test
	controller.GET("/resolve/:var1",
		func(pathVar1 string, queryVar1 int, queryVar2, queryUnknownVar, queryDefaultVar, headerVar1 string) map[string]any {
			return map[string]any{
				"pathVar1":        pathVar1,
				"queryVar1":       queryVar1,
				"queryVar2":       queryVar2,
				"headerVar1":      headerVar1,
				"queryUnknownVar": queryUnknownVar,
				"queryDefaultVar": queryDefaultVar,
			}
		},
		resolver.Path("var1", 1),
		resolver.Query("var2", 2),
		resolver.Query("var3", 3),
		resolver.Query("var4", 4),
		resolver.QueryOrDefault("var4", 5, "default-value"),
		resolver.Header("token", 6),
		ginx.ProduceJSON(),
	)
}

func Test_GET(t *testing.T) {
	res := controller.Tester().GET("/hello/John", nil)

	assert.Equal(t, "Hello, John!", res.Body.String())
}

type blogPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func Test_GET_QueryResolver(t *testing.T) {
	res := controller.Tester().GET("/posts?page=33", nil)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "[{\"title\":\"First post\",\"content\":\"Hello, world\"},{\"title\":\"Monday\",\"content\":\"This is monday\"}]", res.Body.String())
}

func Test_GET_RequestResolver(t *testing.T) {
	res := controller.Tester().GET("/users", map[string]string{"token": "1234567890"})

	assert.Equal(t, "[{user1} {user2}]", res.Body.String())
}

func Test_POST_JsonBodyMap(t *testing.T) {
	res := controller.Tester().POSTJson("/users", map[string]any{
		"username": "Root",
		"id":       12345,
	})

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "<Root:12345>", res.Body.String())
}

func Test_POST_JsonBodyStruct(t *testing.T) {
	res := controller.Tester().POSTJson(
		"/orders?extra=promotion",
		order{ID: 123, Name: "MyORDER", Product: "coca-cola"})

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "<coca-cola:123:promotion>", res.Body.String())
}

func Test_ErrorIntercept(t *testing.T) {
	res := controller.Tester().GET("/error", nil)

	assert.Equal(t, "[[this is error from the '/error']]", res.Body.String())
}

func Test_GET_ResolveQueryPathHeader(t *testing.T) {
	res := controller.Tester().GET("/resolve/PathVariable?var2=1234&var3=test", map[string]string{"token": "abcde1234567890"})

	assert.Equal(t, "application/json; charset=utf-8", res.Header().Get("content-type"))

	var responseData map[string]any
	err := json.Unmarshal(res.Body.Bytes(), &responseData)

	assert.NoError(t, err)

	assert.Equal(t, "abcde1234567890", responseData["headerVar1"])
	assert.Equal(t, "PathVariable", responseData["pathVar1"])
	assert.Equal(t, "", responseData["queryUnknownVar"])
	assert.Equal(t, float64(1234), responseData["queryVar1"])
	assert.Equal(t, "test", responseData["queryVar2"])
	assert.Equal(t, "default-value", responseData["queryDefaultVar"])
}
