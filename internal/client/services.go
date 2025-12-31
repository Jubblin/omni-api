package client

import (
	"github.com/siderolabs/omni/client/pkg/client"
)

// Services wraps all available Omni client services
type Services struct {
	client *client.Client
}

// NewServices creates a new Services wrapper
func NewServices(c *client.Client) *Services {
	return &Services{client: c}
}

// State returns the Omni State interface for read operations
func (s *Services) State() interface{} {
	return s.client.Omni().State()
}

// Management returns the Management service client for write operations
func (s *Services) Management() interface{} {
	return s.client.Management()
}

// Talos returns the Talos service client
func (s *Services) Talos() interface{} {
	return s.client.Talos()
}

// Auth returns the Auth service client
func (s *Services) Auth() interface{} {
	return s.client.Auth()
}

// OIDC returns the OIDC service client
func (s *Services) OIDC() interface{} {
	return s.client.OIDC()
}

// Omni returns the Omni service client
func (s *Services) Omni() interface{} {
	return s.client.Omni()
}
