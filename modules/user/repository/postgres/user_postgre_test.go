package postgre_test

import (
	"context"
	"testing"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	postgreRepository "github.com/cyruzin/puppet_master/modules/user/repository/postgres"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	mockUsers := []domain.User{
		{
			ID:         1,
			Name:       "Homer Simpson",
			Email:      "homer@simpsons.org",
			Password:   "123",
			SuperAdmin: true, // of course :p
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         2,
			Name:       "Bart Simpson",
			Email:      "bart@simpsons.org",
			Password:   "456",
			SuperAdmin: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	rows := sqlxmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password",
		"superadmin",
		"created_at",
		"updated_at",
	}).AddRow(
		mockUsers[0].ID,
		mockUsers[0].Name,
		mockUsers[0].Email,
		mockUsers[0].Password,
		mockUsers[0].SuperAdmin,
		mockUsers[0].CreatedAt,
		mockUsers[0].UpdatedAt,
	).AddRow(
		mockUsers[1].ID,
		mockUsers[1].Name,
		mockUsers[1].Email,
		mockUsers[1].Password,
		mockUsers[1].SuperAdmin,
		mockUsers[1].CreatedAt,
		mockUsers[1].UpdatedAt,
	)

	query := "SELECT \\* FROM users"

	mock.ExpectQuery(query).WillReturnRows(rows)

	userRepo := postgreRepository.NewPostgreUserRepository(db)

	usersList, err := userRepo.Fetch(context.Background())

	assert.NoError(t, err)
	assert.Len(t, usersList, 2)
	assert.Equal(t, usersList[0].Email, "homer@simpsons.org")
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	rows := sqlxmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password",
		"superadmin",
		"created_at",
		"updated_at",
	}).
		AddRow(1, "Homer Simpson", "homer@simpsons.org", "123", true, time.Now(), time.Now())

	query := "SELECT \\* FROM users WHERE id = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)

	userRepo := postgreRepository.NewPostgreUserRepository(db)

	currentUser, err := userRepo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, currentUser)
	assert.Equal(t, "Homer Simpson", currentUser.Name)
}
