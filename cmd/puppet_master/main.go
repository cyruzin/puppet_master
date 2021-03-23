package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	authRepo "github.com/cyruzin/puppet_master/modules/auth/repository/postgres"
	authUseCase "github.com/cyruzin/puppet_master/modules/auth/usecase"
	permissionRepo "github.com/cyruzin/puppet_master/modules/permission/repository/postgres"
	permissionUseCase "github.com/cyruzin/puppet_master/modules/permission/usecase"
	roleRepo "github.com/cyruzin/puppet_master/modules/role/repository/postgres"
	roleUseCase "github.com/cyruzin/puppet_master/modules/role/usecase"
	gql "github.com/cyruzin/puppet_master/modules/shared/delivery/graphql"
	"github.com/cyruzin/puppet_master/modules/shared/delivery/graphql/middleware"
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

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if env {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("running in DEVELOPMENT mode")
	} else {
		log.Info().Msg("running in PRODUCTION mode")
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}
}

func main() {
	ctx := context.Background()

	dataSourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString(`database.host`),
		viper.GetString(`database.user`),
		viper.GetString(`database.pass`),
		viper.GetString(`database.name`),
	)

	dbConnection := databaseConnection(
		ctx,
		viper.GetString(`database.driver`),
		dataSourceName,
	)

	defer dbConnection.Close()

	aRepo := authRepo.NewPostgreAuthRepository(dbConnection)
	aCase := authUseCase.NewAuthUsecase(aRepo)

	pRepo := permissionRepo.NewPostgrePermissionRepository(dbConnection)
	pCase := permissionUseCase.NewPermissionUsecase(pRepo)

	rRepo := roleRepo.NewPostgreRoleRepository(dbConnection)
	rCase := roleUseCase.NewRoleUsecase(rRepo)

	uRepo := userRepo.NewPostgreUserRepository(dbConnection)
	uCase := userUseCase.NewUserUsecase(uRepo)

	root := gql.NewRoot(aCase, pCase, rCase, uCase)

	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    root.Query,
		Mutation: root.Mutation,
	})

	graphqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle(
		"/graphql",
		middleware.LoggerMiddleware(graphqlHandler),
	)

	srv := &http.Server{
		Addr:              viper.GetString(`server.port`),
		ReadTimeout:       viper.GetDuration(`server.read_timeout`),
		ReadHeaderTimeout: viper.GetDuration(`server.read_header_timeout`),
		WriteTimeout:      viper.GetDuration(`server.write_timeout`),
		IdleTimeout:       viper.GetDuration(`server.idle_timeout`),
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		gracefulStop := make(chan os.Signal, 1)
		signal.Notify(gracefulStop, os.Interrupt)
		<-gracefulStop

		log.Info().Msg("shutting down the server...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("server failed to shutdown")
		}
		close(idleConnsClosed)
	}()

	log.Info().Msgf("the server is running on port %s: ", viper.GetString(`server.port`))
	log.Info().Msg("you're good to go! :)")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Stack().Err(err).Msg("server failed to start")
	}

	<-idleConnsClosed
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
