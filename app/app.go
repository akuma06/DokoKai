package app

import (
	"github.com/spf13/viper"
	"net/http"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/tylerb/graceful"
	"fmt"
	"github.com/akuma06/DokoKai/app/templates"
)

var Version string
var Commit string

func GetVersion() string {
	return Version
}

func GetVersionName() string {
	return viper.Get("version.name").(string)
}

// RunServer runs webapp mainloop
func RunServer() {
	http.Handle("/", CSRFRouter)
	serverPort := viper.GetInt64("http.port")
	serverHost := viper.GetString("http.host")
	certFile := viper.GetString("http.certfile")
	keyFile := viper.GetString("http.keyfile")
	drainIntervalString := viper.GetString("http.draininterval")

	drainInterval, err := time.ParseDuration(drainIntervalString)
	if err != nil {
		log.Fatal(err)
	}
	serverAddress := fmt.Sprintf("%s:%d", serverHost, serverPort)
	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{
			Addr: serverAddress,
			WriteTimeout: 30 * time.Second,
			ReadTimeout:  10 * time.Second,
		},
	}

	log.Infoln("Running server on " + serverAddress)

	if certFile != "" && keyFile != "" {
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		log.Fatal(err)
	}
}

func Serve() {
	if viper.GetString("app.version") != Version {
		viper.Set("app.version", Version)
		if viper.GetString("app.commit") != Commit {
			viper.Set("app.commit", Commit)
		}
		viper.WriteConfig()
	}
	if viper.GetString("env") == "DEVELOPMENT" {
		templates.View.SetDevelopmentMode(true)
		log.Info("Template Live Update enabled")
	}
	RunServer()
}