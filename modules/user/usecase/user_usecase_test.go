package usecase_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/cyruzin/puppet_master/domain"
// 	"github.com/cyruzin/puppet_master/domain/mocks"
// 	"github.com/cyruzin/puppet_master/modules/user/usecase"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestFetch(t *testing.T) {
// 	mockUserRepo := new(mocks.UserRepository)
// 	mockUser := &domain.User{
// 		Name:     "Homer Simpson",
// 		Email:    "homer@simpsons.org",
// 		Password: "123",
// 	}

// 	mockListUser := make([]*domain.User, 0)
// 	mockListUser = append(mockListUser, mockUser)

// 	t.Run("success", func(t *testing.T) {
// 		mockUserRepo.On("Fetch",
// 			mock.Anything, mock.AnythingOfType("string"),
// 			mock.Anything, mock.AnythingOfType("string"),
// 			mock.Anything, mock.AnythingOfType("string")).
// 			Return(mockListUser, nil).Once()

// 		u := usecase.NewUserUsecase(mockUserRepo)
// 		list, err := u.Fetch(context.TODO())

// 		assert.Equal(t, "Homer Simpson", list[0].Name)
// 		assert.NoError(t, err)
// 		assert.Len(t, list, len(mockListUser))

// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("error-failed", func(t *testing.T) {
// 		mockUserRepo.On("Fetch",
// 			mock.Anything, mock.AnythingOfType("string"),
// 			mock.Anything, mock.AnythingOfType("string"),
// 			mock.Anything, mock.AnythingOfType("string")).
// 			Return(nil, errors.New("Unexpected error")).Once()

// 		u := usecase.NewUserUsecase(mockUserRepo)
// 		_, err := u.Fetch(context.TODO())

// 		assert.NotNil(t, err)

// 		mockUserRepo.AssertExpectations(t)
// 	})
// }

// func TestGetByID(t *testing.T) {
// 	mockUserRepo := new(mocks.UserRepository)
// 	mockUser := &domain.User{
// 		Name:     "Homer Simpson",
// 		Email:    "homer@simpsons.org",
// 		Password: "123",
// 	}

// 	t.Run("success", func(t *testing.T) {
// 		mockUserRepo.On("GetByID",
// 			mock.Anything, mock.AnythingOfType("int64")).
// 			Return(mockUser, nil).Once()

// 		u := usecase.NewUserUsecase(mockUserRepo)
// 		user, err := u.GetByID(context.TODO(), 1)

// 		assert.Equal(t, "Homer Simpson", user.Name)
// 		assert.NoError(t, err)
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("failure", func(t *testing.T) {
// 		mockUserRepo.On("GetByID",
// 			mock.Anything, mock.AnythingOfType("int64")).
// 			Return(nil, errors.New("Unexpected error")).Once()

// 		u := usecase.NewUserUsecase(mockUserRepo)
// 		_, err := u.GetByID(context.TODO(), 1)

// 		assert.NotNil(t, err)

// 		mockUserRepo.AssertExpectations(t)
// 	})
// }

// func TestStore(t *testing.T) {
// 	mockUserRepo := new(mocks.UserRepository)
// 	mockUser := &domain.User{
// 		Name:     "Homer Simpson",
// 		Email:    "homer@simpsons.org",
// 		Password: "123",
// 	}

// 	t.Run("success", func(t *testing.T) {
// 		mockUserRepo.On("Store",
// 			mock.Anything, mock.AnythingOfType("*domain.User")).
// 			Return(nil).Once()

// 		u := usecase.NewUserUsecase(mockUserRepo)
// 		err := u.Store(context.TODO(), mockUser)

// 		assert.NoError(t, err)
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("failure", func(t *testing.T) {
// 		mockUserRepo.On("Store",
// 			mock.Anything, mock.AnythingOfType("*domain.User")).
// 			Return(errors.New("Unexpected error")).Once()

// 		u := usecase.NewUserUsecase(mockUserRepo)
// 		err := u.Store(context.TODO(), mockUser)

// 		assert.NotNil(t, err)

// 		mockUserRepo.AssertExpectations(t)
// 	})
// }

// func TestUpdate(t *testing.T) {
// 	mockUserRepo := new(mocks.UserRepository)
// 	mockUser := &domain.User{
// 		Name:     "Homer Simpson",
// 		Email:    "homer@simpsons.org",
// 		Password: "123",
// 	}

// 	t.Run("success", func(t *testing.T) {
// 		mockUserRepo.On("Update",
// 			mock.Anything, mock.AnythingOfType("*domain.User")).
// 			Return(nil).Once()

// 		u := usecase.NewUserUsecase(mockUserRepo)
// 		err := u.Update(context.TODO(), mockUser)

// 		assert.NoError(t, err)
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("failure", func(t *testing.T) {
// 		mockUserRepo.On("Update",
// 			mock.Anything, mock.AnythingOfType("*domain.User")).
// 			Return(errors.New("Unexpected error")).Once()

// 		u := usecase.NewUserUsecase(mockUserRepo)
// 		err := u.Update(context.TODO(), mockUser)

// 		assert.NotNil(t, err)

// 		mockUserRepo.AssertExpectations(t)
// 	})
// }

// func TestDelete(t *testing.T) {
// 	mockUserRepo := new(mocks.UserRepository)

// 	t.Run("success", func(t *testing.T) {
// 		mockUserRepo.On("Delete",
// 			mock.Anything, mock.AnythingOfType("int64")).
// 			Return(nil).Once()

// 		u := usecase.NewUserUsecase(mockUserRepo)
// 		err := u.Delete(context.TODO(), 1)

// 		assert.NoError(t, err)
// 		mockUserRepo.AssertExpectations(t)
// 	})

// 	t.Run("failure", func(t *testing.T) {
// 		mockUserRepo.On("Delete",
// 			mock.Anything, mock.AnythingOfType("int64")).
// 			Return(errors.New("Unexpected error")).Once()

// 		u := usecase.NewUserUsecase(mockUserRepo)
// 		err := u.Delete(context.TODO(), 1)

// 		assert.NotNil(t, err)

// 		mockUserRepo.AssertExpectations(t)
// 	})
// }
