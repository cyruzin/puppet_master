package postgre_test

// import (
// 	"context"
// 	"regexp"
// 	"testing"
// 	"time"

// 	"github.com/cyruzin/puppet_master/domain"
// 	postgreRepository "github.com/cyruzin/puppet_master/modules/role/repository/postgres"
// 	"github.com/stretchr/testify/assert"
// 	sqlxmock "github.com/zhashkevych/go-sqlxmock"
// )

// func TestFetch(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	mockRoles := []domain.Role{
// 		{
// 			ID:          1,
// 			Name:        "Admin",
// 			Description: "This is the admin role",
// 			CreatedAt:   time.Now(),
// 			UpdatedAt:   time.Now(),
// 		},
// 		{
// 			ID:          2,
// 			Name:        "Manager",
// 			Description: "This is the manager role",
// 			CreatedAt:   time.Now(),
// 			UpdatedAt:   time.Now(),
// 		},
// 	}

// 	rows := sqlxmock.NewRows([]string{
// 		"id",
// 		"name",
// 		"description",
// 		"created_at",
// 		"updated_at",
// 	}).AddRow(
// 		mockRoles[0].ID,
// 		mockRoles[0].Name,
// 		mockRoles[0].Description,
// 		mockRoles[0].CreatedAt,
// 		mockRoles[0].UpdatedAt,
// 	).AddRow(
// 		mockRoles[1].ID,
// 		mockRoles[1].Name,
// 		mockRoles[1].Description,
// 		mockRoles[1].CreatedAt,
// 		mockRoles[1].UpdatedAt,
// 	)

// 	query := "SELECT \\* FROM roles"
// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	rolesList, err := roleRepo.Fetch(context.TODO())
// 	assert.NoError(t, err)
// 	assert.Len(t, rolesList, 2)
// 	assert.Equal(t, rolesList[0].Name, "Admin")
// }

// func TestFetchFailure(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	rows := sqlxmock.NewRows([]string{
// 		"id",
// 		"name",
// 		"description",
// 		"created_at",
// 		"updated_at",
// 	}).AddRow("", "", "", "", "")
// 	query := "SELECT \\* FROM roles"
// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	_, err = roleRepo.Fetch(context.TODO())
// 	assert.NotNil(t, err)
// }

// func TestGetByID(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	rows := sqlxmock.NewRows([]string{
// 		"id",
// 		"name",
// 		"description",
// 		"created_at",
// 		"updated_at",
// 	}).
// 		AddRow(1, "Admin", "This is the admin role", time.Now(), time.Now())

// 	query := "SELECT \\* FROM roles WHERE id = \\$1"
// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	currentRole, err := roleRepo.GetByID(context.TODO(), 1)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, currentRole)
// 	assert.Equal(t, "Admin", currentRole.Name)
// }

// func TestGetByIDFailure(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	rows := sqlxmock.NewRows([]string{
// 		"id",
// 		"name",
// 		"description",
// 		"created_at",
// 		"updated_at",
// 	}).
// 		AddRow("", "", "", "", "")
// 	query := "SELECT \\* FROM roles WHERE id = \\$1"
// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	_, err = roleRepo.GetByID(context.TODO(), 0)
// 	assert.NotNil(t, err)
// }

// func TestStore(t *testing.T) {
// 	now := time.Now()
// 	role := &domain.Role{
// 		ID:          12,
// 		Name:        "Admin",
// 		Description: "This is the admin role",
// 		CreatedAt:   now,
// 		UpdatedAt:   now,
// 	}

// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := `
// 	  INSERT INTO roles (
// 		name,
// 		description,
// 		created_at,
// 		updated_at
// 		)
// 		VALUES ($1, $2, $3, $4)
// 		`

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
// 		role.Name,
// 		role.Description,
// 		role.CreatedAt,
// 		role.UpdatedAt,
// 	).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	err = roleRepo.Store(context.TODO(), role)
// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(12), role.ID)
// }

// func TestStoreFailure(t *testing.T) {
// 	role := &domain.Role{}

// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := `
// 	  INSERT INTO roles (
// 		name,
// 		description,
// 		created_at,
// 		updated_at
// 		)
// 		VALUES (?, ?, ?, ?)
// 		`

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
// 		"", "", "", "",
// 	).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	err = roleRepo.Store(context.TODO(), role)
// 	assert.NotNil(t, err)
// }

// func TestUpdate(t *testing.T) {
// 	now := time.Now()
// 	role := &domain.Role{
// 		ID:          12,
// 		Name:        "Admin",
// 		Description: "This is the admin role",
// 		UpdatedAt:   now,
// 	}

// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := `
// 		UPDATE roles
// 		SET
// 		name = $1,
// 		description = $2,
// 		updated_at = $3
// 		WHERE id = $4
// 	`

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
// 		role.Name,
// 		role.Description,
// 		role.UpdatedAt,
// 		role.ID,
// 	).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	err = roleRepo.Update(context.TODO(), role)
// 	assert.NoError(t, err)
// }

// func TestUpdateFailure(t *testing.T) {
// 	role := &domain.Role{}

// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := `
// 		UPDATE roles
// 		SET
// 		name = $1,
// 		description = $2,
// 		updated_at = $3
// 		WHERE id = $4
// 	`

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
// 		"", "", "",
// 	).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	err = roleRepo.Update(context.TODO(), role)
// 	assert.NotNil(t, err)
// }

// func TestUpdateRowsAffected(t *testing.T) {
// 	now := time.Now()
// 	role := &domain.Role{
// 		ID:          12,
// 		Name:        "Admin",
// 		Description: "This is the admin role",
// 		UpdatedAt:   now,
// 	}

// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := `
// 		UPDATE roles
// 		SET
// 		name = $1,
// 		description = $2,
// 		updated_at = $3
// 		WHERE id = $4
// 	`

// 	mock.ExpectBegin()

// 	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
// 		role.Name,
// 		role.Description,
// 		role.UpdatedAt,
// 		role.ID,
// 	).WillReturnResult(sqlxmock.NewResult(12, 0))

// 	mock.ExpectCommit()
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	err = roleRepo.Update(context.TODO(), role)
// 	assert.NotNil(t, err)
// }

// func TestDelete(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := "DELETE FROM roles WHERE id = \\$1"
// 	mock.ExpectBegin()
// 	mock.ExpectExec(query).WithArgs(12).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	err = roleRepo.Delete(context.TODO(), int64(12))
// 	assert.NoError(t, err)
// }

// func TestDeleteFailure(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := "DELETE FROM roles WHERE id = \\$1"
// 	mock.ExpectBegin()
// 	mock.ExpectExec(query).WithArgs(0).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	err = roleRepo.Delete(context.TODO(), int64(12))
// 	assert.NotNil(t, err)
// }

// func TestDeleteRowsAffected(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := "DELETE FROM roles WHERE id = \\$1"
// 	mock.ExpectBegin()
// 	mock.ExpectExec(query).WithArgs(12).WillReturnResult(sqlxmock.NewResult(12, 0))
// 	mock.ExpectCommit()
// 	roleRepo := postgreRepository.NewPostgreRoleRepository(db)
// 	err = roleRepo.Delete(context.TODO(), int64(12))
// 	assert.NotNil(t, err)
// }
