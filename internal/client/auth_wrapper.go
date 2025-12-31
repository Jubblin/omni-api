package client

import (
	"context"
	"fmt"

	"github.com/siderolabs/omni/client/pkg/client"
)

// ServiceAccount represents a service account
type ServiceAccount struct {
	ID          string
	Name        string
	Description string
	Roles       []string
}

// CreateServiceAccountRequest represents a request to create a service account
type CreateServiceAccountRequest struct {
	Name        string
	Description string
	Roles       []string
}

// AuthService defines the interface for Auth operations
type AuthService interface {
	ListServiceAccounts(ctx context.Context) ([]*ServiceAccount, error)
	GetServiceAccount(ctx context.Context, id string) (*ServiceAccount, error)
	CreateServiceAccount(ctx context.Context, req *CreateServiceAccountRequest) (*ServiceAccount, error)
	DeleteServiceAccount(ctx context.Context, id string) error
}

// authService implements AuthService
type authService struct {
	client interface{} // Will be *auth.Client once we know the type
}

// NewAuthService creates a new AuthService wrapper
func NewAuthService(c *client.Client) AuthService {
	return &authService{
		client: c.Auth(),
	}
}

// Implementation methods - these will call the actual Auth client API
// TODO: Implement once actual API methods are known

func (a *authService) ListServiceAccounts(ctx context.Context) ([]*ServiceAccount, error) {
	return nil, fmt.Errorf("ListServiceAccounts not yet implemented - Auth API integration needed")
}

func (a *authService) GetServiceAccount(ctx context.Context, id string) (*ServiceAccount, error) {
	return nil, fmt.Errorf("GetServiceAccount not yet implemented - Auth API integration needed")
}

func (a *authService) CreateServiceAccount(ctx context.Context, req *CreateServiceAccountRequest) (*ServiceAccount, error) {
	return nil, fmt.Errorf("CreateServiceAccount not yet implemented - Auth API integration needed")
}

func (a *authService) DeleteServiceAccount(ctx context.Context, id string) error {
	return fmt.Errorf("DeleteServiceAccount not yet implemented - Auth API integration needed")
}
