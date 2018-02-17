package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"gitlab.com/salismt/microservice-pattern-user-service/config"
	"gitlab.com/salismt/microservice-pattern-user-service/cache"
	"gitlab.com/salismt/microservice-pattern-user-service/app"
)

func main() {
	config := config.BaseConfig{}
	config.Load()

	fmt.Printf("Port is %s", config.GetValue("port"))
	fmt.Printf("Redis is %s", config.GetValue("APP_RD_ADDRESS"))

	cache := cache.Cache{Enable: true}

	// flag to receive information directly from the command line
	flag.StringVar(
		&cache.Address,
		"redis_address",
		config.GetValue("APP_RD_ADDRESS"),
		"Redis Address",
	)

	flag.StringVar(
		&cache.Auth,
		"redis_auth",
		config.GetValue("APP_RD_AUTH"),
		"Redis Auth",
	)

	flag.StringVar(
		&cache.DB,
		"redis_db_name",
		config.GetValue("APP_RB_DBNAME"),
		"Redis DB Name",
	)

	flag.IntVar(
		&cache.MaxActive,
		"redis_max_active",
		60,
		"Redis Max Active",
	)

	flag.IntVar(
		&cache.IdleTimeoutSecs,
		"redis_timeout",
		60,
		"Redis timeout in seconds",
	)

	flag.Parse()
	cache.Pool = cache.NewCachePool()

	connectionString := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		config.GetValue("APP_DB_USERNAME"),
		config.GetValue("APP_DB_PASSWORD"),
		config.GetValue("APP_DB_NAME"),
	)

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a := app.App{}
	a.Initialize(cache, db)
	a.Run(":8080")
}
