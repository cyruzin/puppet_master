package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	permissionRepo "github.com/cyruzin/puppet_master/modules/permission/repository/postgres"
	permissionUseCase "github.com/cyruzin/puppet_master/modules/permission/usecase"
	roleRepo "github.com/cyruzin/puppet_master/modules/role/repository/postgres"
	roleUseCase "github.com/cyruzin/puppet_master/modules/role/usecase"
	gql "github.com/cyruzin/puppet_master/modules/shared/delivery/graphql"
	userRepo "github.com/cyruzin/puppet_master/modules/user/repository/postgres"
	userUseCase "github.com/cyruzin/puppet_master/modules/user/usecase"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
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

	dataSourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName,
	)

	dbConnection := databaseConnection(ctx, dbDriver, dataSourceName)

	defer dbConnection.Close()

	pRepo := permissionRepo.NewPostgrePermissionRepository(dbConnection)
	pCase := permissionUseCase.NewPermissionUsecase(pRepo)

	rRepo := roleRepo.NewPostgreRoleRepository(dbConnection)
	rCase := roleUseCase.NewRoleUsecase(rRepo)

	uRepo := userRepo.NewPostgreUserRepository(dbConnection)
	uCase := userUseCase.NewUserUsecase(uRepo)

	root := gql.NewRoot(pCase, rCase, uCase)

	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    root.Query,
		Mutation: root.Mutation,
	})

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	log.Info().Msg("the server is running on port: 8000")

	http.Handle("/graphql", h)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal().Stack().Err(err).Msg("")
	}
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
