package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var (
	dbUrl = os.Getenv("DB_URL")
	step  = os.Getenv("STEP")
)
var Conn *sql.DB

func main() {
	for {
		db, err := sql.Open("postgres", dbUrl)
		if err == nil {
			err = db.Ping()
			if err == nil {
				fmt.Println("Successfully connected to the database")
				Conn = db
				break
			}
		}
		fmt.Println("Error connecting to the database, retrying in 1 second...")
		time.Sleep(time.Second)
	}
	defer Conn.Close()
	var upCmd = &cobra.Command{
		Use:   "up",
		Short: "Run all available migrations",
		Run: func(cmd *cobra.Command, args []string) {
			driver, err := postgres.WithInstance(Conn, &postgres.Config{})
			if err != nil {
				log.Fatalf("Could not start migration: %v", err)
			}
			m, err := migrate.NewWithDatabaseInstance(
				"file://migrations",
				"postgres", driver)
			if err != nil {
				log.Fatalf("Could not start migration: %v", err)
			}
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("An error occurred while migrating: %v", err)
			}
		},
	}
	var downCmd = &cobra.Command{
		Use:   "down",
		Short: "Revert the latest migration",
		Run: func(cmd *cobra.Command, args []string) {
			driver, err := postgres.WithInstance(Conn, &postgres.Config{})
			if err != nil {
				log.Fatalf("Could not start migration: %v", err)
			}
			m, err := migrate.NewWithDatabaseInstance(
				"file://migrations",
				"postgres", driver)
			if err != nil {
				log.Fatalf("Could not start migration: %v", err)
			}
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("An error occurred while migrating: %v", err)
			}
		},
	}
	var stepCmd = &cobra.Command{
		Use:   "step",
		Short: "Run a specific number of migrations",
		Run: func(cmd *cobra.Command, args []string) {
			driver, err := postgres.WithInstance(Conn, &postgres.Config{})
			if err != nil {
				log.Fatalf("Could not start migration: %v", err)
			}
			m, err := migrate.NewWithDatabaseInstance(
				"file://migrations",
				"postgres", driver)
			if err != nil {
				log.Fatalf("Could not start migration: %v", err)
			}
			stepInt, err := strconv.Atoi(step)
			if err := m.Steps(stepInt); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("An error occurred while migrating: %v", err)
			}
		},
	}
	var rootCmd = &cobra.Command{Use: "migrate"}
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(stepCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Could not execute command: %v", err)
	}
}
