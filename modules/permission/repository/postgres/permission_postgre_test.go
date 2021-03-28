package postgre_test

// import (
// 	"context"
// 	"regexp"
// 	"testing"
// 	"time"

// 	"github.com/cyruzin/puppet_master/domain"
// 	postgreRepository "github.com/cyruzin/puppet_master/modules/permission/repository/postgres"
// 	"github.com/stretchr/testify/assert"
// 	sqlxmock "github.com/zhashkevych/go-sqlxmock"
// )

// func TestFetch(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	mockPermissions := []domain.Permission{
// 		{
// 			ID:          1,
// 			Name:        "create articles",
// 			Description: "Permission to create articles",
// 			CreatedAt:   time.Now(),
// 			UpdatedAt:   time.Now(),
// 		},
// 		{
// 			ID:          2,
// 			Name:        "edit articles",
// 			Description: "Permission to edit articles",
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
// 		mockPermissions[0].ID,
// 		mockPermissions[0].Name,
// 		mockPermissions[0].Description,
// 		mockPermissions[0].CreatedAt,
// 		mockPermissions[0].UpdatedAt,
// 	).AddRow(
// 		mockPermissions[1].ID,
// 		mockPermissions[1].Name,
// 		mockPermissions[1].Description,
// 		mockPermissions[1].CreatedAt,
// 		mockPermissions[1].UpdatedAt,
// 	)

// 	query := "SELECT \\* FROM permissions"
// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	permissionsList, err := permissionRepo.Fetch(context.TODO())
// 	assert.NoError(t, err)
// 	assert.Len(t, permissionsList, 2)
// 	assert.Equal(t, permissionsList[0].Name, "create articles")
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
// 	query := "SELECT \\* FROM permissions"
// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	_, err = permissionRepo.Fetch(context.TODO())
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
// 		AddRow(1, "create articles", "Permission to create articles", time.Now(), time.Now())

// 	query := "SELECT \\* FROM permissions WHERE id = \\$1"
// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	currentPermission, err := permissionRepo.GetByID(context.TODO(), 1)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, currentPermission)
// 	assert.Equal(t, "create articles", currentPermission.Name)
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
// 	query := "SELECT \\* FROM permissions WHERE id = \\$1"
// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	_, err = permissionRepo.GetByID(context.TODO(), 0)
// 	assert.NotNil(t, err)
// }

// func TestStore(t *testing.T) {
// 	now := time.Now()
// 	permission := &domain.Permission{
// 		ID:          12,
// 		Name:        "create articles",
// 		Description: "Permission to create articles",
// 		CreatedAt:   now,
// 		UpdatedAt:   now,
// 	}

// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := `
// 	  INSERT INTO permissions (
// 		name,
// 		description,
// 		created_at,
// 		updated_at
// 		)
// 		VALUES ($1, $2, $3, $4)
// 		`

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
// 		permission.Name,
// 		permission.Description,
// 		permission.CreatedAt,
// 		permission.UpdatedAt,
// 	).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	err = permissionRepo.Store(context.TODO(), permission)
// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(12), permission.ID)
// }

// func TestStoreFailure(t *testing.T) {
// 	permission := &domain.Permission{}

// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := `
// 	  INSERT INTO permissions (
// 		name,
// 		description,
// 		created_at,
// 		updated_at
// 		)
// 		VALUES ($1, $2, $3, $4)
// 		`

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
// 		"", "", "", "",
// 	).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	err = permissionRepo.Store(context.TODO(), permission)
// 	assert.NotNil(t, err)
// }

// func TestUpdate(t *testing.T) {
// 	now := time.Now()
// 	permission := &domain.Permission{
// 		ID:          12,
// 		Name:        "edit articles",
// 		Description: "Permission to edit articles",
// 		UpdatedAt:   now,
// 	}

// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := `
// 		UPDATE permissions
// 		SET
// 		name = $1,
// 		description = $2,
// 		updated_at = $3
// 		WHERE id = $4
// 	`

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
// 		permission.Name,
// 		permission.Description,
// 		permission.UpdatedAt,
// 		permission.ID,
// 	).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	err = permissionRepo.Update(context.TODO(), permission)
// 	assert.NoError(t, err)
// }

// func TestUpdateFailure(t *testing.T) {
// 	permission := &domain.Permission{}

// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := `
// 		UPDATE permissions
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
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	err = permissionRepo.Update(context.TODO(), permission)
// 	assert.NotNil(t, err)
// }

// func TestUpdateRowsAffected(t *testing.T) {
// 	now := time.Now()
// 	permission := &domain.Permission{
// 		ID:          12,
// 		Name:        "edit articles",
// 		Description: "Permission to edit articles",
// 		UpdatedAt:   now,
// 	}

// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := `
// 		UPDATE permissions
// 		SET
// 		name = $1,
// 		description = $2,
// 		updated_at = $3
// 		WHERE id = $4
// 	`

// 	mock.ExpectBegin()

// 	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
// 		permission.Name,
// 		permission.Description,
// 		permission.UpdatedAt,
// 		permission.ID,
// 	).WillReturnResult(sqlxmock.NewResult(12, 0))

// 	mock.ExpectCommit()
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	err = permissionRepo.Update(context.TODO(), permission)
// 	assert.NotNil(t, err)
// }

// func TestDelete(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := "DELETE FROM permissions WHERE id = \\$1"
// 	mock.ExpectBegin()
// 	mock.ExpectExec(query).WithArgs(12).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	err = permissionRepo.Delete(context.TODO(), int64(12))
// 	assert.NoError(t, err)
// }

// func TestDeleteFailure(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := "DELETE FROM permissions WHERE id = \\$1"
// 	mock.ExpectBegin()
// 	mock.ExpectExec(query).WithArgs(0).WillReturnResult(sqlxmock.NewResult(12, 1))
// 	mock.ExpectCommit()
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	err = permissionRepo.Delete(context.TODO(), int64(12))
// 	assert.NotNil(t, err)
// }

// func TestDeleteRowsAffected(t *testing.T) {
// 	db, mock, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer db.Close()

// 	query := "DELETE FROM permissions WHERE id = \\$1"
// 	mock.ExpectBegin()
// 	mock.ExpectExec(query).WithArgs(12).WillReturnResult(sqlxmock.NewResult(12, 0))
// 	mock.ExpectCommit()
// 	permissionRepo := postgreRepository.NewPostgrePermissionRepository(db)
// 	err = permissionRepo.Delete(context.TODO(), int64(12))
// 	assert.NotNil(t, err)
// }
