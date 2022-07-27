package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/myOmikron/q-scheduler/modules/scheduler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	log2 "github.com/labstack/gommon/log"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/middleware"

	"github.com/myOmikron/q-scheduler/models"
)

var ascii = `
  ___    ___      _             _        _
 / _ \  / __| __ | |_   ___  __| | _  _ | | ___  _ _
| (_) | \__ \/ _||   \ / -_)/ _\ || || || |/ -_)| '_|
 \__\_\ |___/\__||_||_|\___|\__/_| \_._||_|\___||_|
`

func StartServer(configPath string) {
	// Config
	config := models.GetConfig(configPath)

	// Check for required files
	if _, err := os.Stat(config.HTTP.TLSKeyPath); os.IsNotExist(err) {
		color.Println(color.RED, "[File Error]")
		fmt.Printf("Private key not found: %s\n", config.HTTP.TLSKeyPath)
		os.Exit(1)
	}

	if _, err := os.Stat(config.HTTP.TLSCertPath); os.IsNotExist(err) {
		color.Println(color.RED, "[File Error]")
		fmt.Printf("Certificate not found: %s\n", config.HTTP.TLSCertPath)
		os.Exit(1)
	}

	// Initialize database
	db := initDatabase(config)

	// Initialize scheduler
	s := scheduler.New(db)
	go s.Start()

	// Initialize web server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(log2.OFF)

	// Set middleware
	e.Use(middleware.Panic())

	// Define routes
	defineRoutes(e, db, config, s)

	// Display art & start server
	color.Println(color.RED, ascii)

	fmt.Print("Start listening on ")
	color.Println(color.PURPLE, config.HTTP.ListenAddress)

	control := make(chan os.Signal, 1)
	signal.Notify(control, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		// Start server
		if err := e.StartTLS(
			config.HTTP.ListenAddress, config.HTTP.TLSCertPath, config.HTTP.TLSKeyPath,
		); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println(err.Error())
		}
	}()

	restart := false
	for {
		sig := <-control

		if sig == syscall.SIGHUP { // Reload server
			color.Println(color.PURPLE, "Server is restarting")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			e.Shutdown(ctx)
			s.Quit()
			cancel()
			restart = true
			break
		} else if sig == syscall.SIGINT { // Shutdown gracefully
			color.Println(color.PURPLE, "Server is stopping gracefully")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			e.Shutdown(ctx)
			s.Quit()
			cancel()
			break
		} else if sig == syscall.SIGTERM { // Shutdown immediately
			e.Close()
			s.Quit()
			color.Println(color.PURPLE, "Server was shut down")
			break
		} else {
			fmt.Printf("Received unknown signal: %s\n", sig.String())
		}
	}
	if restart {
		StartServer(configPath)
	}

}
