package sqlh

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Attach receives a map of int64/[]int and attach the IDs on the given pivot table.
func Attach(
	ctx context.Context,
	tx *sqlx.Tx,
	query string,
	s map[int64][]int,
) error {

	for index, ids := range s {
		for _, item := range ids {
			_, err := tx.ExecContext(ctx, query, index, item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Detach receives a map of int64/[]int and Detach the IDs on the given pivot table.
func Detach(
	ctx context.Context,
	tx *sqlx.Tx,
	query string,
	s map[int64][]int,
) error {
	for index := range s {
		_, err := tx.ExecContext(ctx, query, index)
		if err != nil {
			return err
		}
	}

	return nil
}

// Sync receives a map of int64/[]int and sync the IDs on the given pivot table.
func Sync(ctx context.Context, tx *sqlx.Tx, query string, s map[int64][]int) error {
	empty := IsEmpty(s)

	if !empty {
		err := Detach(ctx, tx, query, s)
		if err != nil {
			return err
		}

		err = Attach(ctx, tx, query, s)
		if err != nil {
			return err
		}
	} else {
		err := Detach(ctx, tx, query, s)
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
