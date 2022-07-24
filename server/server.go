package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/myOmikron/q-scheduler/models"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	log2 "github.com/labstack/gommon/log"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/middleware"
)

var ascii = `
▄▄▄▄▄▄▄         ▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄ ▄▄   ▄▄ ▄▄▄▄▄▄▄ ▄▄▄▄▄▄  ▄▄   ▄▄ ▄▄▄     ▄▄▄▄▄▄▄ ▄▄▄▄▄▄
█       █      █       █       █  █ █  █       █      ██  █ █  █   █   █       █   ▄  █
█   ▄   █      █  ▄▄▄▄▄█       █  █▄█  █    ▄▄▄█  ▄    █  █ █  █   █   █    ▄▄▄█  █ █ █
█  █ █  █      █  █▄▄▄▄█     ▄▄█       █   █▄▄▄█ █ █   █  █▄█  █   █   █   █▄▄▄█  █▄▄█▄▄
█  █▄█  █      █▄▄▄▄▄  █    █  █   ▄   █    ▄▄▄█ █▄█   █       █   █▄▄▄█    ▄▄▄█    ▄▄  █
█      █       ▄▄▄▄▄█  █    █▄▄█  █ █  █   █▄▄▄█       █       █       █   █▄▄▄█   █  █ █
█▄▄▄▄██▄█      █▄▄▄▄▄▄▄█▄▄▄▄▄▄▄█▄▄█ █▄▄█▄▄▄▄▄▄▄█▄▄▄▄▄▄██▄▄▄▄▄▄▄█▄▄▄▄▄▄▄█▄▄▄▄▄▄▄█▄▄▄█  █▄█
`

func StartServer(configPath string) {
	// Config
	config := models.GetConfig(configPath)

	// Initialize database
	_ = initDatabase(config)

	// Initialize web server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(log2.OFF)

	// Set middleware
	e.Use(middleware.Panic())

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
			cancel()
			restart = true
			break
		} else if sig == syscall.SIGINT { // Shutdown gracefully
			color.Println(color.PURPLE, "Server is stopping gracefully")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			e.Shutdown(ctx)
			cancel()
			break
		} else if sig == syscall.SIGTERM { // Shutdown immediately
			e.Close()
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
