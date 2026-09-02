package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/purpose-robot/planet-express/assets"
	"github.com/purpose-robot/planet-express/auth"
	"github.com/purpose-robot/planet-express/internal/config"
	"github.com/purpose-robot/planet-express/internal/httputil"
	"github.com/purpose-robot/planet-express/sqlite"

	httpx "github.com/purpose-robot/planet-express/http"
	sqliteauth "github.com/purpose-robot/planet-express/sqlite/auth"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("load config from env: %v", err)
	}

	conn, err := sqlite.Open(cfg.DB.Name)
	if err != nil {
		log.Fatalf("open sqlite database: %v", err)
	}

	renderer, err := httputil.NewRenderer(
		assets.HTML,
		"html/base.tmpl",
		"html/partials/*.tmpl",
	)
	if err != nil {
		log.Fatalf("initiate template renderer: %v", err)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      httpx.NewServer(auth.NewService(sqliteauth.NewTransactor(conn.Writer)), renderer, scs.New()).Handle(),
		IdleTimeout:  cfg.HTTP.IdleTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	log.Fatal(server.ListenAndServe())
}
