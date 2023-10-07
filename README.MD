## GINx

```go
blogController := ginx.NewController(gin.Default())

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
```