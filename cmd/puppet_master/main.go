package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`../../config.json`)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	env := viper.GetBool(`debug`)

	if env {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("running in DEVELOPMENT mode")
	} else {
		log.Debug().Msg("running in PRODUCTION mode")
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}
}

func main() {
	ctx := context.Background()

	dbDriver := viper.GetString(`database.driver`)
	dbHost := viper.GetString(`database.host`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	// Postgre connection string
	dataSourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName,
	)

	dbConnection := databaseConnection(ctx, dbDriver, dataSourceName)

	defer dbConnection.Close()

	// TODO: Implement connection
	log.Info().Str("driver", dbConnection.DriverName())
}

func databaseConnection(
	ctx context.Context,
	driverName string,
	dataSourceName string,
) *sqlx.DB {
	db, err := sqlx.ConnectContext(ctx, driverName, dataSourceName)
	if err != nil {
		log.Fatal().
			Err(err).
			Stack().
			Str("database", db.DriverName()).
			Msg("could not connect to the database")
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatal().
			Err(err).
			Stack().
			Str("database", db.DriverName()).
			Msg("could not ping the database")
	}

	return db
}
