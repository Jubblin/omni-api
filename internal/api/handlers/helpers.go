package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// buildURL constructs a full URL from the gin context and a path
func buildURL(c *gin.Context, path string) string {
	req := c.Request
	scheme := "http"
	
	// Check for HTTPS or forwarded protocol
	if req.TLS != nil {
		scheme = "https"
	} else if proto := req.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	if host == "" {
		host = "localhost:8080" // fallback
	}
	
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	
	return scheme + "://" + host + path
}
