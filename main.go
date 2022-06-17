package main

import (
	"fmt"
	"os"

	"github.com/FACT-Finder/noflake/api"
	"github.com/FACT-Finder/noflake/database"
	"github.com/FACT-Finder/noflake/logger"
	"github.com/FACT-Finder/noflake/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

//go:generate oapi-codegen -config ./openapi-gen.yml openapi.yaml
func main() {
	logger.Init(zerolog.InfoLevel)
	app := &cli.App{
		Name: "noflake",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "the address to listen on",
				EnvVars: []string{"NOFLAKE_ADDRESS"},
				Value:   ":8000",
			},
			&cli.StringFlag{
				Name:    "db",
				Usage:   "the path to the sqlite database",
				EnvVars: []string{"NOFLAKE_DB"},
				Value:   "noflake.sqlite3",
			},
			&cli.StringFlag{
				Name:     "token",
				Usage:    "token to secure the POST endpoints",
				EnvVars:  []string{"NOFLAKE_TOKEN"},
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			db := database.New(c.String("db"))
			token := c.String("token")
			webapi := api.New(db, token)

			listenAddr := c.String("address")
			log.Info().Str("address", listenAddr).Msg("HTTP")
			err := server.Start(webapi, listenAddr)

			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
