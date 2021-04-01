package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	authRepository "github.com/cyruzin/puppet_master/modules/auth/repository/postgres"
	authCacheRepository "github.com/cyruzin/puppet_master/modules/auth/repository/redis"
	authUseCase "github.com/cyruzin/puppet_master/modules/auth/usecase"
	permissionRepository "github.com/cyruzin/puppet_master/modules/permission/repository/postgres"
	permissionUseCase "github.com/cyruzin/puppet_master/modules/permission/usecase"
	roleRepository "github.com/cyruzin/puppet_master/modules/role/repository/postgres"
	roleUseCase "github.com/cyruzin/puppet_master/modules/role/usecase"
	gql "github.com/cyruzin/puppet_master/modules/shared/delivery/graphql"
	"github.com/cyruzin/puppet_master/modules/shared/delivery/graphql/middleware"
	userRepository "github.com/cyruzin/puppet_master/modules/user/repository/postgres"
	userUseCase "github.com/cyruzin/puppet_master/modules/user/usecase"
	"github.com/cyruzin/puppet_master/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/go-redis/redis/v8"

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
	viper.SetConfigFile(util.PathBuilder("./config.json"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if viper.GetBool(`debug`) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("running in development mode")
	} else {
		log.Info().Msg("running in production mode")
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataSourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString(`database.host`),
		viper.GetString(`database.user`),
		viper.GetString(`database.pass`),
		viper.GetString(`database.name`),
	)

	postgreDB := postgreConnection(
		ctx,
		viper.GetString(`database.driver`),
		dataSourceName,
	)

	defer postgreDB.Close()

	redisClient := redisConnection(ctx)

	defer redisClient.Close()

	authCacheRepository := authCacheRepository.NewRedisCacheRepository(redisClient)

	permissionRepository := permissionRepository.NewPostgrePermissionRepository(postgreDB)
	permissionUseCase := permissionUseCase.NewPermissionUsecase(permissionRepository)

	roleRepository := roleRepository.NewPostgreRoleRepository(postgreDB, permissionRepository)
	roleUseCase := roleUseCase.NewRoleUsecase(permissionRepository, roleRepository)

	userRepository := userRepository.NewPostgreUserRepository(postgreDB, permissionRepository, roleRepository)
	userUseCase := userUseCase.NewUserUsecase(permissionRepository, roleRepository, userRepository)

	authRepository := authRepository.NewPostgreAuthRepository(postgreDB)
	authUseCase := authUseCase.NewAuthUsecase(
		authRepository,
		authCacheRepository,
		permissionRepository,
		roleRepository,
		userRepository,
	)

	root := gql.NewRoot(authUseCase, permissionUseCase, roleUseCase, userUseCase)

	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    root.Query,
		Mutation: root.Mutation,
	})

	graphqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	chiHandler := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	chiHandler.Use(
		cors.Handler,
		render.SetContentType(render.ContentTypeJSON),
		middleware.LoggerMiddleware,
		middleware.AuthMiddleware,
	)

	chiHandler.Handle("/graphql", graphqlHandler)

	srv := &http.Server{
		Addr:              ":" + viper.GetString(`server.port`),
		ReadTimeout:       viper.GetDuration(`server.read_timeout`),
		ReadHeaderTimeout: viper.GetDuration(`server.read_header_timeout`),
		WriteTimeout:      viper.GetDuration(`server.write_timeout`),
		IdleTimeout:       viper.GetDuration(`server.idle_timeout`),
		Handler:           chiHandler,
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

	log.Info().Msgf("the server is running on port: %s", viper.GetString(`server.port`))
	log.Info().Msg("you're good to go! :)")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Stack().Err(err).Msg("server failed to start")
	}

	<-idleConnsClosed
}

func postgreConnection(
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

	log.Info().Msg("postgre connected")

	return db
}

// Redis connection.
func redisConnection(ctx context.Context) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString(`redis.address`),
		Password: viper.GetString(`redis.password`),
		DB:       viper.GetInt(`redis.db`),
	})

	_, err := client.WithContext(ctx).Ping(ctx).Result()
	if err != nil {
		log.Fatal().
			Err(err).
			Stack().
			Msg("could not connect to redis")
	}

	log.Info().Msg("redis connected")

	return client
}
