package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/labstack/echo"
	"github.com/spf13/viper"

	_costumerHttpDelivery "github.com/bxcodec/go-clean-arch/costumer/delivery/http"
	_costumerHttpDeliveryMiddleware "github.com/bxcodec/go-clean-arch/costumer/delivery/http/middleware"
	_costumerRepo "github.com/bxcodec/go-clean-arch/costumer/repository/postgres"
	_costumerUcase "github.com/bxcodec/go-clean-arch/costumer/usecase"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

type config struct {
	port int
	env  string
	db   struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  string
	}
	limiter struct {
			enabled bool
			rps     float64
			burst   int
	}
	smtp struct {
			host     string
			port     int
			username string
			password string
			sender   string
	}
}

func main() {
	// dbHost := "localhost"
	// dbPort := "5432"
	// dbUser := "postgres"
	// dbPass := "admin"
	// dbName := "go-challange"
	// connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	// // connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	// val := url.Values{}
	// val.Add("parseTime", "1")
	// val.Add("loc", "Asia/Jakarta")
	// dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	// dbConn, err := sql.Open(`postgres`, dsn)

	var cfg config
	flag.IntVar(&cfg.port, "port", 9090, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:admin@localhost/go-challange?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.Parse()
	// logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	dbConn, err := openDB(cfg)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _costumerHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	ar := _costumerRepo.NewPostgresCostumerRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _costumerUcase.NewCostumerUsecase(ar, timeoutContext)
	_costumerHttpDelivery.NewCostumerHandler(e, au)

	log.Fatal(e.Start(viper.GetString("server.address"))) //nolint
}


func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
			return nil, err
	}
	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
			return nil, err
	}
	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
			return nil, err
	}
	return db, nil
}