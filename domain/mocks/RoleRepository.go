// Code generated by mockery v2.7.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/cyruzin/puppet_master/domain"
	mock "github.com/stretchr/testify/mock"
)

// RoleRepository is an autogenerated mock type for the RoleRepository type
type RoleRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *RoleRepository) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fetch provides a mock function with given fields: ctx
func (_m *RoleRepository) Fetch(ctx context.Context) ([]*domain.Role, error) {
	ret := _m.Called(ctx)

	var r0 []*domain.Role
	if rf, ok := ret.Get(0).(func(context.Context) []*domain.Role); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Role)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *RoleRepository) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.Role
	if rf, ok := ret.Get(0).(func(context.Context, int64) *domain.Role); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Role)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, role
func (_m *RoleRepository) Store(ctx context.Context, role *domain.Role) error {
	ret := _m.Called(ctx, role)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Role) error); ok {
		r0 = rf(ctx, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, role
func (_m *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	ret := _m.Called(ctx, role)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Role) error); ok {
		r0 = rf(ctx, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}