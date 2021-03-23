package postgre

import (
	"context"
	"database/sql"

	"github.com/cyruzin/puppet_master/domain"

	"github.com/jmoiron/sqlx"
)

type postgreRepository struct {
	Conn *sqlx.DB
}

// NewPostgreAuthRepository will create an object that represent
// the auth.Repository interface.
func NewPostgreAuthRepository(Conn *sqlx.DB) domain.AuthRepository {
	return &postgreRepository{Conn}
}

func (p *postgreRepository) Authenticate(ctx context.Context, email string) (string, error) {
	var password string

	query := "SELECT password from users WHERE email = $1"

	err := p.Conn.GetContext(ctx, &password, query, email)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return password, nil
}
