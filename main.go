package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	cache := Cache{Enable: true}

	// flag to receive information directly from the command line
	flag.StringVar(
		&cache.Address,
		"redis_address",
		os.Getenv("APP_RD_ADDRESS"),
		"Redis Address",
	)

	flag.StringVar(
		&cache.Auth,
		"redis_auth",
		os.Getenv("APP_RD_AUTH"),
		"Redis Auth",
	)

	flag.StringVar(
		&cache.DB,
		"redis_db_name",
		os.Getenv("APP_RB_DBNAME"),
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
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a := App{}
	a.Initialize(cache, db)
	a.Run(":8080")
}
