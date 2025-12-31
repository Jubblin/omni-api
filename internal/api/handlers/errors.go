package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// handleManagementError maps gRPC errors from Management service to HTTP status codes
func handleManagementError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	code := status.Code(err)
	switch code {
	case codes.NotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case codes.AlreadyExists:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case codes.InvalidArgument:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case codes.PermissionDenied:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case codes.Unauthenticated:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case codes.Unavailable:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	case codes.DeadlineExceeded:
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "operation timeout"})
	case codes.FailedPrecondition:
		c.JSON(http.StatusPreconditionFailed, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// handleTalosError maps gRPC errors from Talos service to HTTP status codes
func handleTalosError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	code := status.Code(err)
	switch code {
	case codes.NotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
	case codes.Unavailable:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "machine unreachable"})
	case codes.DeadlineExceeded:
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "operation timeout"})
	case codes.InvalidArgument:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
