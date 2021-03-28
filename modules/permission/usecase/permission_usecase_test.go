package usecase_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/cyruzin/puppet_master/domain"
// 	"github.com/cyruzin/puppet_master/domain/mocks"
// 	"github.com/cyruzin/puppet_master/modules/permission/usecase"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestFetch(t *testing.T) {
// 	mockPermissionRepo := new(mocks.PermissionRepository)
// 	mockPermission := &domain.Permission{
// 		Name: "read articles",
// 	}

// 	mockListPermission := make([]*domain.Permission, 0)
// 	mockListPermission = append(mockListPermission, mockPermission)

// 	t.Run("success", func(t *testing.T) {
// 		mockPermissionRepo.On("Fetch",
// 			mock.Anything, mock.AnythingOfType("string"),
// 			mock.Anything, mock.AnythingOfType("string"),
// 			mock.Anything, mock.AnythingOfType("string")).
// 			Return(mockListPermission, nil).Once()

// 		u := usecase.NewPermissionUsecase(mockPermissionRepo)
// 		list, err := u.Fetch(context.TODO())

// 		assert.Equal(t, "read articles", list[0].Name)
// 		assert.NoError(t, err)
// 		assert.Len(t, list, len(mockListPermission))

// 		mockPermissionRepo.AssertExpectations(t)
// 	})

// 	t.Run("error-failed", func(t *testing.T) {
// 		mockPermissionRepo.On("Fetch",
// 			mock.Anything, mock.AnythingOfType("string"),
// 			mock.Anything, mock.AnythingOfType("string"),
// 			mock.Anything, mock.AnythingOfType("string")).
// 			Return(nil, errors.New("Unexpected error")).Once()

// 		u := usecase.NewPermissionUsecase(mockPermissionRepo)
// 		_, err := u.Fetch(context.TODO())

// 		assert.NotNil(t, err)

// 		mockPermissionRepo.AssertExpectations(t)
// 	})
// }

// func TestGetByID(t *testing.T) {
// 	mockPermissionRepo := new(mocks.PermissionRepository)
// 	mockPermission := &domain.Permission{
// 		Name: "read articles",
// 	}

// 	t.Run("success", func(t *testing.T) {
// 		mockPermissionRepo.On("GetByID",
// 			mock.Anything, mock.AnythingOfType("int64")).
// 			Return(mockPermission, nil).Once()

// 		u := usecase.NewPermissionUsecase(mockPermissionRepo)
// 		permission, err := u.GetByID(context.TODO(), 1)

// 		assert.Equal(t, "read articles", permission.Name)
// 		assert.NoError(t, err)
// 		mockPermissionRepo.AssertExpectations(t)
// 	})

// 	t.Run("failure", func(t *testing.T) {
// 		mockPermissionRepo.On("GetByID",
// 			mock.Anything, mock.AnythingOfType("int64")).
// 			Return(nil, errors.New("Unexpected error")).Once()

// 		u := usecase.NewPermissionUsecase(mockPermissionRepo)
// 		_, err := u.GetByID(context.TODO(), 1)

// 		assert.NotNil(t, err)

// 		mockPermissionRepo.AssertExpectations(t)
// 	})
// }

// func TestStore(t *testing.T) {
// 	mockPermissionRepo := new(mocks.PermissionRepository)
// 	mockPermission := &domain.Permission{
// 		Name: "read articles",
// 	}

// 	t.Run("success", func(t *testing.T) {
// 		mockPermissionRepo.On("Store",
// 			mock.Anything, mock.AnythingOfType("*domain.Permission")).
// 			Return(nil).Once()

// 		u := usecase.NewPermissionUsecase(mockPermissionRepo)
// 		err := u.Store(context.TODO(), mockPermission)

// 		assert.NoError(t, err)
// 		mockPermissionRepo.AssertExpectations(t)
// 	})

// 	t.Run("failure", func(t *testing.T) {
// 		mockPermissionRepo.On("Store",
// 			mock.Anything, mock.AnythingOfType("*domain.Permission")).
// 			Return(errors.New("Unexpected error")).Once()

// 		u := usecase.NewPermissionUsecase(mockPermissionRepo)
// 		err := u.Store(context.TODO(), mockPermission)

// 		assert.NotNil(t, err)

// 		mockPermissionRepo.AssertExpectations(t)
// 	})
// }

// func TestUpdate(t *testing.T) {
// 	mockPermissionRepo := new(mocks.PermissionRepository)
// 	mockPermission := &domain.Permission{
// 		Name: "read articles",
// 	}

// 	t.Run("success", func(t *testing.T) {
// 		mockPermissionRepo.On("Update",
// 			mock.Anything, mock.AnythingOfType("*domain.Permission")).
// 			Return(nil).Once()

// 		u := usecase.NewPermissionUsecase(mockPermissionRepo)
// 		err := u.Update(context.TODO(), mockPermission)

// 		assert.NoError(t, err)
// 		mockPermissionRepo.AssertExpectations(t)
// 	})

// 	t.Run("failure", func(t *testing.T) {
// 		mockPermissionRepo.On("Update",
// 			mock.Anything, mock.AnythingOfType("*domain.Permission")).
// 			Return(errors.New("Unexpected error")).Once()

// 		u := usecase.NewPermissionUsecase(mockPermissionRepo)
// 		err := u.Update(context.TODO(), mockPermission)

// 		assert.NotNil(t, err)

// 		mockPermissionRepo.AssertExpectations(t)
// 	})
// }

// func TestDelete(t *testing.T) {
// 	mockPermissionRepo := new(mocks.PermissionRepository)

// 	t.Run("success", func(t *testing.T) {
// 		mockPermissionRepo.On("Delete",
// 			mock.Anything, mock.AnythingOfType("int64")).
// 			Return(nil).Once()

// 		u := usecase.NewPermissionUsecase(mockPermissionRepo)
// 		err := u.Delete(context.TODO(), 1)

// 		assert.NoError(t, err)
// 		mockPermissionRepo.AssertExpectations(t)
// 	})

// 	t.Run("failure", func(t *testing.T) {
// 		mockPermissionRepo.On("Delete",
// 			mock.Anything, mock.AnythingOfType("int64")).
// 			Return(errors.New("Unexpected error")).Once()

// 		u := usecase.NewPermissionUsecase(mockPermissionRepo)
// 		err := u.Delete(context.TODO(), 1)

// 		assert.NotNil(t, err)

// 		mockPermissionRepo.AssertExpectations(t)
// 	})
// }
