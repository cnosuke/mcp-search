package main

import (
	"fmt"
	"os"

	"github.com/cnosuke/mcp-search/config"
	"github.com/cnosuke/mcp-search/logger"
	"github.com/cnosuke/mcp-search/server"
	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"
)

var (
	// Version and Revision are replaced when building.
	// To set specific version, edit Makefile.
	Version  = "0.0.1"
	Revision = "xxx"

	Name  = "mcp-search"
	Usage = "A MCP server implementation for web search using Brave Search API"
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s (%s)", Version, Revision)
	app.Name = Name
	app.Usage = Usage

	app.Commands = []*cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Start the MCP search server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "config",
					Aliases: []string{"c"},
					Value:   "config.yml",
					Usage:   "path to the configuration file",
				},
			},
			Action: func(c *cli.Context) error {
				configPath := c.String("config")

				// Read the configuration file
				cfg, err := config.LoadConfig(configPath)
				if err != nil {
					return errors.Wrap(err, "failed to load configuration file")
				}

				// Initialize logger
				if err := logger.InitLogger(cfg.Debug, cfg.Log); err != nil {
					return errors.Wrap(err, "failed to initialize logger")
				}
				defer logger.Sync()

				// Start the server
				return server.Run(cfg)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
