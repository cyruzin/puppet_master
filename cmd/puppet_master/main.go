package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`../../config.json`)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("running in DEBUG mode")
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

	log.Println(dbConnection.DriverName())
}

func databaseConnection(
	ctx context.Context,
	driverName string,
	dataSourceName string,
) *sqlx.DB {
	db, err := sqlx.ConnectContext(ctx, driverName, dataSourceName)
	if err != nil {
		log.Fatalln("could not connect to the database. error: ", err)
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatalln("could not ping the database. error: ", err)
	}

	return db
}
