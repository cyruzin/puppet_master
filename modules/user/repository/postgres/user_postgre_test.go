package postgre_test

import (
	"context"
	"regexp"
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
			ID:        1,
			Name:      "Homer Simpson",
			Email:     "homer@simpsons.org",
			Password:  "123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "Bart Simpson",
			Email:     "bart@simpsons.org",
			Password:  "456",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlxmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password",
		"created_at",
		"updated_at",
	}).AddRow(
		mockUsers[0].ID,
		mockUsers[0].Name,
		mockUsers[0].Email,
		mockUsers[0].Password,
		mockUsers[0].CreatedAt,
		mockUsers[0].UpdatedAt,
	).AddRow(
		mockUsers[1].ID,
		mockUsers[1].Name,
		mockUsers[1].Email,
		mockUsers[1].Password,
		mockUsers[1].CreatedAt,
		mockUsers[1].UpdatedAt,
	)

	query := "SELECT \\* FROM users"
	mock.ExpectQuery(query).WillReturnRows(rows)
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	usersList, err := userRepo.Fetch(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, usersList, 2)
	assert.Equal(t, usersList[0].Email, "homer@simpsons.org")
}

func TestFetchFailure(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlxmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password",
		"created_at",
		"updated_at",
	}).AddRow("", "", "", "", "", "")
	query := "SELECT \\* FROM users"
	mock.ExpectQuery(query).WillReturnRows(rows)
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	_, err = userRepo.Fetch(context.TODO())
	assert.NotNil(t, err)
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
		"created_at",
		"updated_at",
	}).
		AddRow(1, "Homer Simpson", "homer@simpsons.org", "123", time.Now(), time.Now())

	query := "SELECT \\* FROM users WHERE id = \\$1"
	mock.ExpectQuery(query).WillReturnRows(rows)
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	currentUser, err := userRepo.GetByID(context.TODO(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, currentUser)
	assert.Equal(t, "Homer Simpson", currentUser.Name)
}

func TestGetByIDFailure(t *testing.T) {
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
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "", "")
	query := "SELECT \\* FROM users WHERE id = \\$1"
	mock.ExpectQuery(query).WillReturnRows(rows)
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	_, err = userRepo.GetByID(context.TODO(), 0)
	assert.NotNil(t, err)
}

func TestStore(t *testing.T) {
	now := time.Now()
	user := &domain.User{
		ID:        12,
		Name:      "Homer Simpson",
		Email:     "homer@simpsons.org",
		Password:  "123",
		CreatedAt: now,
		UpdatedAt: now,
	}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	query := `
	  INSERT INTO users ( 
		name, 
		email, 
		password,
		created_at, 
		updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).WillReturnResult(sqlxmock.NewResult(12, 1))
	mock.ExpectCommit()
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	err = userRepo.Store(context.TODO(), user)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), user.ID)
}

func TestStoreFailure(t *testing.T) {
	user := &domain.User{}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	query := `
	  INSERT INTO users ( 
		name, 
		email, 
		password,
		created_at, 
		updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		"", "", "", "", "",
	).WillReturnResult(sqlxmock.NewResult(12, 1))
	mock.ExpectCommit()
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	err = userRepo.Store(context.TODO(), user)
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	user := &domain.User{
		ID:        12,
		Name:      "Homer Simpson",
		Email:     "homer@simpsons.org",
		Password:  "123",
		UpdatedAt: now,
	}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	query := `
		UPDATE users
		SET 
		name = $1, 
		email = $2, 
		password = $3, 
		updated_at = $4
		WHERE id = $5
	`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		user.Name,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.ID,
	).WillReturnResult(sqlxmock.NewResult(12, 1))
	mock.ExpectCommit()
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	err = userRepo.Update(context.TODO(), user)
	assert.NoError(t, err)
}

func TestUpdateFailure(t *testing.T) {
	user := &domain.User{}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	query := `
		UPDATE users
		SET 
		name = $1, 
		email = $2, 
		password = $3, 
		updated_at = $4
		WHERE id = $5
	`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		"", "", "", "", "", "",
	).WillReturnResult(sqlxmock.NewResult(12, 1))
	mock.ExpectCommit()
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	err = userRepo.Update(context.TODO(), user)
	assert.NotNil(t, err)
}

func TestUpdateRowsAffected(t *testing.T) {
	now := time.Now()
	user := &domain.User{
		ID:        12,
		Name:      "Homer Simpson",
		Email:     "homer@simpsons.org",
		Password:  "123",
		UpdatedAt: now,
	}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	query := `
		UPDATE users
		SET 
		name = $1, 
		email = $2, 
		password = $3, 
		updated_at = $4
		WHERE id = $5
	`

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		user.Name,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.ID,
	).WillReturnResult(sqlxmock.NewResult(12, 0))

	mock.ExpectCommit()
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	err = userRepo.Update(context.TODO(), user)
	assert.NotNil(t, err)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	query := "DELETE FROM users WHERE id = \\$1"
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(12).WillReturnResult(sqlxmock.NewResult(12, 1))
	mock.ExpectCommit()
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	err = userRepo.Delete(context.TODO(), int64(12))
	assert.NoError(t, err)
}

func TestDeleteFailure(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	query := "DELETE FROM users WHERE id = \\$1"
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(0).WillReturnResult(sqlxmock.NewResult(12, 1))
	mock.ExpectCommit()
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	err = userRepo.Delete(context.TODO(), int64(12))
	assert.NotNil(t, err)
}

func TestDeleteRowsAffected(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	query := "DELETE FROM users WHERE id = \\$1"
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(12).WillReturnResult(sqlxmock.NewResult(12, 0))
	mock.ExpectCommit()
	userRepo := postgreRepository.NewPostgreUserRepository(db)
	err = userRepo.Delete(context.TODO(), int64(12))
	assert.NotNil(t, err)
}
