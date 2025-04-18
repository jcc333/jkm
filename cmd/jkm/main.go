package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jcc333/jkm/internal/configure"
	"github.com/jcc333/jkm/internal/log"
	"github.com/jcc333/jkm/internal/router"
)

func main() {
	isLogging := os.Getenv("JKM_LOGGING")
	log.Init(isLogging != "false" && isLogging != "")

	cfg, err := configure.Load()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		fmt.Println("\nPlease create a .env file based on the example.env template")
		fmt.Println("and set your email server settings and credentials.")
		os.Exit(1)
	}

	app, err := router.New(cfg)
	if err != nil {
		msg := fmt.Sprintf("creating router: %v", err)
		log.Info(msg)
		fmt.Fprintf(os.Stderr, msg)
		os.Exit(1)
	}
	defer app.Disconnect()
	defer log.Info("thank you for using jkm!")
	defer log.Close()
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running UI: %v\n", err)
		os.Exit(1)
	}
}
