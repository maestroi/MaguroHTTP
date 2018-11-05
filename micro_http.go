package main

import (
	"fmt"
	"net/http"
)

var mCfg = microConfig{}

func main() {

	initLogger("MicroHTTP-")
	loadConfigFromFile("./config.json", &mCfg)

	mux := http.NewServeMux()

	// If ProxyMode is enabled, use proxy handler
	if mCfg.ProxyMode {
		mux.HandleFunc("/", handleProxy)
	} else {
		mux.HandleFunc("/", handleHTTP)
	}

	logAction(logDEBUG, fmt.Errorf("MicroHTTP is listening on port %s", mCfg.Port))
	http.ListenAndServe(mCfg.Address+":"+mCfg.Port, mux)

}
