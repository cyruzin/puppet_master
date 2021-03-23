package sql_helper

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// If you're using MySQL change the $1, $2 to ?, ?.
// Default: PostgreSQL
const (
	// INSERT INTO pivot_table VALUES (x, y)
	insertQuery = `INSERT INTO $1 VALUES ($2, $3)`
	// DELETE FROM pivot_table WHERE field = x
	deleteQuery = `DELETE FROM $1 WHERE $2 = $3"`
)

// Attach receives a map of int64/[]int and attach the IDs on the given pivot table.
func Attach(ctx context.Context, conn *sqlx.DB, s map[int64][]int, pivot string) error {
	for index, ids := range s {
		for _, values := range ids {
			tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = conn.ExecContext(ctx, insertQuery, pivot, index, values)
			if err != nil {
				tx.Rollback()
				return err
			}

			tx.Commit()
		}
	}

	return nil
}

// Detach receives a map of int64/[]int and Detach the IDs on the given pivot table.
func Detach(ctx context.Context, conn *sqlx.DB, s map[int64][]int, pivot, field string) error {
	for index := range s {
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = conn.ExecContext(ctx, deleteQuery, pivot, field, index)
		if err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
	}

	return nil
}

// Sync receives a map of int64/[]int and sync the IDs on the given pivot table.
func Sync(ctx context.Context, conn *sqlx.DB, s map[int64][]int, pivot, field string) error {
	empty := IsEmpty(s)

	if !empty {
		err := Detach(ctx, conn, s, pivot, field)
		if err != nil {
			return err
		}

		err = Attach(ctx, conn, s, pivot)
		if err != nil {
			return err
		}
	} else {
		err := Detach(ctx, conn, s, pivot, field)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsEmpty checks if a given map of int64/[]int is empty.
func IsEmpty(s map[int64][]int) bool {
	empty := true

	for _, ids := range s {
		if len(ids) > 0 {
			empty = false
		}
	}

	return empty
}
