package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paveldanilin/ginx"
	"github.com/paveldanilin/ginx/requestbody"
	"github.com/paveldanilin/ginx/resolver"
	"math/rand"
	"sync"
	"time"
)

var storage = map[string][]BlogPost{}
var storageMu sync.Mutex

type BlogPost struct {
	requestbody.JSON
	ID          int       `json:"id,omitempty"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PublishDate time.Time `json:"publish_date,omitempty"`
}

func getUserFeed(username string) []BlogPost {
	storageMu.Lock()
	defer storageMu.Unlock()

	if _, exists := storage[username]; exists {
		return storage[username]
	}

	return []BlogPost{}
}

func createBlogPost(username string, blogPost BlogPost) {
	storageMu.Lock()
	defer storageMu.Unlock()

	if _, exists := storage[username]; !exists {
		storage[username] = []BlogPost{}
	}

	blogPost.ID = rand.Int()
	blogPost.Author = username
	blogPost.PublishDate = time.Now()

	storage[username] = append(storage[username], blogPost)
}

func main() {
	r := gin.Default()

	blogController := ginx.NewController(r)

	// Controller produces JSON response.
	blogController.ContentType = gin.MIMEJSON

	// Controller resolves user function handler struct arguments.
	blogController.Use(resolver.Struct())

	blogController.GET(
		"/blog/:username/feed",
		getUserFeed,
		// Resolve the first argument as a path variable and inject value into user handler.
		resolver.Path("username", 1))

	blogController.POST("/blog/:username",
		createBlogPost,
		resolver.Path("username", 1))

	_ = r.Run()
}
