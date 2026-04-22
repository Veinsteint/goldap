package model

// RoleCasbin represents a role permission rule
type RoleCasbin struct {
	Keyword string `json:"keyword"` // Role keyword
	Path    string `json:"path"`    // API path
	Method  string `json:"method"`  // HTTP method
}
