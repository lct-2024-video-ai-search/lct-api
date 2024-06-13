package lct_backend

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"lct-backend/api"
	"net"
	"net/http"
	"time"
)

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		fmt.Printf("%+v", err)
		log.Fatal().Err(err).Msg("cannot load config:")
	}
	err = runMigration(config.MigrationURL, config.DBSource)
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("Migrations successful, starting...")

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := http.Client{Timeout: 10 * time.Second, Transport: transport}

	ch := getClickhouseClient(config.ClickHouseHost)
	server, err := api.NewServer(ch, client)

	if err != nil {
		log.Fatal().Err(err).Msg("error init server:")
	}

	server.Run(config.HTTPServerAddress)
}

func getClickhouseClient(addr string) *sql.DB {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                true,
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

func runMigration(migrationUrl string, dbSource string) error {
	m, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
