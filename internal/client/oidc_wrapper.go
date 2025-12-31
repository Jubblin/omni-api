package client

import (
	"context"
	"fmt"

	"github.com/siderolabs/omni/client/pkg/client"
)

// OIDCProvider represents an OIDC provider configuration
type OIDCProvider struct {
	ID        string
	Name      string
	IssuerURL string
	ClientID  string
	Scopes    []string
	Enabled   bool
}

// CreateOIDCProviderRequest represents a request to create an OIDC provider
type CreateOIDCProviderRequest struct {
	Name      string
	IssuerURL string
	ClientID  string
	Scopes    []string
	Enabled   bool
}

// UpdateOIDCProviderRequest represents a request to update an OIDC provider
type UpdateOIDCProviderRequest struct {
	Name      string
	IssuerURL string
	ClientID  string
	Scopes    []string
	Enabled   *bool
}

// OIDCService defines the interface for OIDC operations
type OIDCService interface {
	ListOIDCProviders(ctx context.Context) ([]*OIDCProvider, error)
	GetOIDCProvider(ctx context.Context, id string) (*OIDCProvider, error)
	CreateOIDCProvider(ctx context.Context, req *CreateOIDCProviderRequest) (*OIDCProvider, error)
	UpdateOIDCProvider(ctx context.Context, id string, req *UpdateOIDCProviderRequest) (*OIDCProvider, error)
	DeleteOIDCProvider(ctx context.Context, id string) error
}

// oidcService implements OIDCService
type oidcService struct {
	client interface{} // Will be *oidc.Client once we know the type
}

// NewOIDCService creates a new OIDCService wrapper
func NewOIDCService(c *client.Client) OIDCService {
	return &oidcService{
		client: c.OIDC(),
	}
}

// Implementation methods - these will call the actual OIDC client API
// TODO: Implement once actual API methods are known

func (o *oidcService) ListOIDCProviders(ctx context.Context) ([]*OIDCProvider, error) {
	return nil, fmt.Errorf("ListOIDCProviders not yet implemented - OIDC API integration needed")
}

func (o *oidcService) GetOIDCProvider(ctx context.Context, id string) (*OIDCProvider, error) {
	return nil, fmt.Errorf("GetOIDCProvider not yet implemented - OIDC API integration needed")
}

func (o *oidcService) CreateOIDCProvider(ctx context.Context, req *CreateOIDCProviderRequest) (*OIDCProvider, error) {
	return nil, fmt.Errorf("CreateOIDCProvider not yet implemented - OIDC API integration needed")
}

func (o *oidcService) UpdateOIDCProvider(ctx context.Context, id string, req *UpdateOIDCProviderRequest) (*OIDCProvider, error) {
	return nil, fmt.Errorf("UpdateOIDCProvider not yet implemented - OIDC API integration needed")
}

func (o *oidcService) DeleteOIDCProvider(ctx context.Context, id string) error {
	return fmt.Errorf("DeleteOIDCProvider not yet implemented - OIDC API integration needed")
}
