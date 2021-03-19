package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/domain/mocks"
	"github.com/cyruzin/puppet_master/modules/role/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetch(t *testing.T) {
	mockRoleRepo := new(mocks.RoleRepository)
	mockRole := &domain.Role{
		Name: "Admin",
	}

	mockListRole := make([]*domain.Role, 0)
	mockListRole = append(mockListRole, mockRole)

	t.Run("success", func(t *testing.T) {
		mockRoleRepo.On("Fetch",
			mock.Anything, mock.AnythingOfType("string"),
			mock.Anything, mock.AnythingOfType("string"),
			mock.Anything, mock.AnythingOfType("string")).
			Return(mockListRole, nil).Once()

		u := usecase.NewRoleUsecase(mockRoleRepo)
		list, err := u.Fetch(context.TODO())

		assert.Equal(t, "Admin", list[0].Name)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListRole))

		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockRoleRepo.On("Fetch",
			mock.Anything, mock.AnythingOfType("string"),
			mock.Anything, mock.AnythingOfType("string"),
			mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("Unexpected error")).Once()

		u := usecase.NewRoleUsecase(mockRoleRepo)
		_, err := u.Fetch(context.TODO())

		assert.NotNil(t, err)

		mockRoleRepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockRoleRepo := new(mocks.RoleRepository)
	mockRole := &domain.Role{
		Name: "Admin",
	}

	t.Run("success", func(t *testing.T) {
		mockRoleRepo.On("GetByID",
			mock.Anything, mock.AnythingOfType("int64")).
			Return(mockRole, nil).Once()

		u := usecase.NewRoleUsecase(mockRoleRepo)
		permission, err := u.GetByID(context.TODO(), 1)

		assert.Equal(t, "Admin", permission.Name)
		assert.NoError(t, err)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockRoleRepo.On("GetByID",
			mock.Anything, mock.AnythingOfType("int64")).
			Return(nil, errors.New("Unexpected error")).Once()

		u := usecase.NewRoleUsecase(mockRoleRepo)
		_, err := u.GetByID(context.TODO(), 1)

		assert.NotNil(t, err)

		mockRoleRepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	mockRoleRepo := new(mocks.RoleRepository)
	mockRole := &domain.Role{
		Name: "Admin",
	}

	t.Run("success", func(t *testing.T) {
		mockRoleRepo.On("Store",
			mock.Anything, mock.AnythingOfType("*domain.Role")).
			Return(nil).Once()

		u := usecase.NewRoleUsecase(mockRoleRepo)
		err := u.Store(context.TODO(), mockRole)

		assert.NoError(t, err)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockRoleRepo.On("Store",
			mock.Anything, mock.AnythingOfType("*domain.Role")).
			Return(errors.New("Unexpected error")).Once()

		u := usecase.NewRoleUsecase(mockRoleRepo)
		err := u.Store(context.TODO(), mockRole)

		assert.NotNil(t, err)

		mockRoleRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockRoleRepo := new(mocks.RoleRepository)
	mockRole := &domain.Role{
		Name: "Admin",
	}

	t.Run("success", func(t *testing.T) {
		mockRoleRepo.On("Update",
			mock.Anything, mock.AnythingOfType("*domain.Role")).
			Return(nil).Once()

		u := usecase.NewRoleUsecase(mockRoleRepo)
		err := u.Update(context.TODO(), mockRole)

		assert.NoError(t, err)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockRoleRepo.On("Update",
			mock.Anything, mock.AnythingOfType("*domain.Role")).
			Return(errors.New("Unexpected error")).Once()

		u := usecase.NewRoleUsecase(mockRoleRepo)
		err := u.Update(context.TODO(), mockRole)

		assert.NotNil(t, err)

		mockRoleRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockRoleRepo := new(mocks.RoleRepository)

	t.Run("success", func(t *testing.T) {
		mockRoleRepo.On("Delete",
			mock.Anything, mock.AnythingOfType("int64")).
			Return(nil).Once()

		u := usecase.NewRoleUsecase(mockRoleRepo)
		err := u.Delete(context.TODO(), 1)

		assert.NoError(t, err)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockRoleRepo.On("Delete",
			mock.Anything, mock.AnythingOfType("int64")).
			Return(errors.New("Unexpected error")).Once()

		u := usecase.NewRoleUsecase(mockRoleRepo)
		err := u.Delete(context.TODO(), 1)

		assert.NotNil(t, err)

		mockRoleRepo.AssertExpectations(t)
	})
}
