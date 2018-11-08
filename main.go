package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const version = "v0.12"

// Main function
func main() {

	initLogger("MicroHTTP-")

	args := os.Args
	if len(args) == 1 {
		showHelp()
	}

	if _, err := os.Stat(args[1]); err == nil {
		var mCfg microConfig
		loadConfigFromFile(args[1], &mCfg)
		if valid, err := validateConfig(args[1], &mCfg); valid && err == nil {
			startServer(&mCfg)
		} else {
			logAction(logERROR, err)
			os.Exit(1)
		}

	} else {
		showHelp()
	}
}

// Function to start Server
func startServer(mCfg *microConfig) {

	m := micro{
		config: *mCfg,
		vhosts: make(map[string]microConfig),
	}

	if m.config.Serve.VirtualHosting {
		for k, v := range m.config.Serve.VirtualHosts {
			var cfg microConfig
			loadConfigFromFile(v, &cfg)
			if valid, err := validateConfigVhost(v, &cfg); !valid || err != nil {
				logAction(logERROR, err)
				os.Exit(1)
			}
			m.vhosts[k] = cfg
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", m.httpServe)

	if m.config.TLS && httpCheckTLS(&m.config) {
		logAction(logNONE, fmt.Errorf("MicroHTTP is listening on port %s with TLS", mCfg.Port))
		tlsc := httpCreateTLSConfig()
		ms := http.Server{
			Addr:      mCfg.Address + ":" + mCfg.Port,
			Handler:   mux,
			TLSConfig: tlsc,
		}

		done := make(chan bool)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		go func() {
			<-quit
			logAction(logNONE, fmt.Errorf("Server is shutting down..."))

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			ms.SetKeepAlivesEnabled(false)
			if err := ms.Shutdown(ctx); err != nil {
				logAction(logNONE, fmt.Errorf("Could not gracefully shutdown the server: %v\n", err))
			}
			close(done)
		}()

		err := ms.ListenAndServeTLS(mCfg.TLSCert, mCfg.TLSKey)
		if err != nil && err != http.ErrServerClosed {
			logAction(logERROR, fmt.Errorf("Starting server failed: %s", err))
			return
		}

		<-done
		logAction(logNONE, fmt.Errorf("MicroHTTP stopped"))

	} else {
		logAction(logDEBUG, fmt.Errorf("MicroHTTP is listening on port %s", mCfg.Port))
		http.ListenAndServe(mCfg.Address+":"+mCfg.Port, mux)
	}
}

// Function to show help
func showHelp() {
	fmt.Printf("MicroHTTP version %s\n\nUsage: microhttp </path/to/config.json>\n\n", version)
	os.Exit(1)
}