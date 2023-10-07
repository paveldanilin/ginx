package tests

import (
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
		ginx.ResponseContentType(gin.MIMEJSON), // <- overwrite controller ContentType
	)

	// [GET] /users
	controller.GET("/users", loadUsers, ginx.ResponseContentType(gin.MIMEPlain))

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
