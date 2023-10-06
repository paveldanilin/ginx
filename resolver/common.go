package resolver

type Scope string

const (
	// ScopePath variable will be resolved by a request path i.e.: /api/user/:name.
	ScopePath Scope = "path"

	// ScopeQuery variable will be resolved by a request query string i.e. /api/user?name=John.
	ScopeQuery Scope = "query"

	// ScopeHeader variable will be resolved by a request headers.
	ScopeHeader Scope = "header"
)

type priority struct {
	value int
}

func (p priority) Priority() int {
	return p.value
}
