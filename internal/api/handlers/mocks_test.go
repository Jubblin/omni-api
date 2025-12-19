package handlers

import (
	"context"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/stretchr/testify/mock"
)

type MockState struct {
	mock.Mock
	state.State
}

func (m *MockState) Get(ctx context.Context, p resource.Pointer, opts ...state.GetOption) (resource.Resource, error) {
	args := m.Called(ctx, p, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(resource.Resource), args.Error(1)
}

func (m *MockState) List(ctx context.Context, k resource.Kind, opts ...state.ListOption) (resource.List, error) {
	args := m.Called(ctx, k, opts)
	return args.Get(0).(resource.List), args.Error(1)
}
