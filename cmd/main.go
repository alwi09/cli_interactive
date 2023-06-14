package main

import (
	"cli_interactive/internal/config"
	"cli_interactive/internal/database"
	"cli_interactive/internal/service"
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {

	// Set log level to display only INFO level
	logrus.SetLevel(logrus.InfoLevel)

	// Set log formatter to text format
	logrus.SetFormatter(&logrus.TextFormatter{})

	var cfg config.Config

	// initialize config
	envConfigErr := envconfig.Process("", &cfg)
	if envConfigErr != nil {
		logrus.Fatal(errors.New("initialize config is error"))
	}

	// initialize database connection
	db, err := database.ConnectDB(&cfg)
	if err != nil {
		return
	}

	// initialize migration
	err = database.Migrate(db)
	if err != nil {
		logrus.Fatalf("error running schema migration %v", err)
	}

	// initialize repository wikis
	wikisRepository := database.NewWikisRepository(db)

	// initialize handler wikis
	wikisHandler := service.NewHandlerImpl(wikisRepository)

	rootCmd := &cobra.Command{
		Use:   "add-wikis",
		Short: "CLI to store topics in the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to the Wiki CLI!")
			wikisHandler.StartInteractiveCLI(cmd, args)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
