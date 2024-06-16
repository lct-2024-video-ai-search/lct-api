package main

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"lct-backend/api"
	"lct-backend/config"
	"lct-backend/db"
	"net"
	"net/http"
	"time"
)

// @title           Zvezdolet Search API
// @version         1.0
// @description     API к сервису индексации и поиска видео

// @host      api-zvezdolet.ddns.net
// @BasePath  /

func main() {
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("%+v", err)
		log.Fatal().Err(err).Msg("cannot load appConfig:")
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   60 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}
	client := http.Client{Timeout: 600 * time.Second, Transport: transport}
	ch := getClickhouseClient(appConfig.ClickHouseHost)
	store := db.NewSQLVideoStore(ch)
	server, err := api.NewServer(&store, appConfig.VideoProcessingServiceAddress, appConfig.VideoIndexingServiceAddress, client)

	if err != nil {
		log.Fatal().Err(err).Msg("error init server:")
	}

	server.Run(appConfig.HTTPServerAddress)
}

func getClickhouseClient(addr string) *sql.DB {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                false,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})
	conn.SetMaxIdleConns(500)
	conn.SetMaxOpenConns(1000)
	conn.SetConnMaxLifetime(time.Hour)
	if err := conn.Ping(); err != nil {
		log.Fatal().Err(err).Msg("cannot connect to clickhouse")
	}
	return conn
}
