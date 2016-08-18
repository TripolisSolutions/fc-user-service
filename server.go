package main

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/tylerb/graceful"
	"goji.io"
	"goji.io/pat"

	mongo "github.com/TripolisSolutions/go-helper/mgojuice"
	"github.com/TripolisSolutions/go-helper/middleware"
	"github.com/TripolisSolutions/go-helper/rabbitcage"
	"github.com/TripolisSolutions/go-helper/redis"
	"github.com/TripolisSolutions/go-helper/relic"
	"github.com/TripolisSolutions/go-helper/settings"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	EnvSettingsInit()
	settings.EnvSettingsInit()

	if err := mongo.Startup(); err != nil {
		log.WithFields(log.Fields{}).Fatal("Mongo startup failed")
		os.Exit(1)
	}

	if err := rabbitcage.SetupConn(
		settings.ProjectEnvSettings.RabbitMQHost,
		ProjectEnvSettings.Buffer,
		true); err != nil {
		log.WithFields(log.Fields{}).Fatal("RabbitMQ startup failed")
		os.Exit(1)
	}

	if err := rabbitcage.Make(ProjectEnvSettings.EventUsersCreatedQueue,
		ProjectEnvSettings.EventUsersExchange,
		ProjectEnvSettings.EventUsersCreatedQueue); err != nil {
		log.WithFields(log.Fields{}).Fatal("RabbitMQ make failed")
		os.Exit(1)
	}

	if err := rabbitcage.Make(ProjectEnvSettings.EventUsersUpdatedQueue,
		ProjectEnvSettings.EventUsersExchange,
		ProjectEnvSettings.EventUsersUpdatedQueue); err != nil {
		log.WithFields(log.Fields{}).Fatal("RabbitMQ make failed")
		os.Exit(1)
	}

	if err := rabbitcage.Make(ProjectEnvSettings.EventUsersDeletedQueue,
		ProjectEnvSettings.EventUsersExchange,
		ProjectEnvSettings.EventUsersDeletedQueue); err != nil {
		log.WithFields(log.Fields{}).Fatal("RabbitMQ make failed")
		os.Exit(1)
	}

	relicWrapper, err := relic.Startup()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("New Relic startup failed")
		os.Exit(1)
	}

	redis.Setup(redis.Config{settings.ProjectEnvSettings.RedisAddress})

	go func() {
		for {
			select {
			case <-rabbitcage.AmqpReady:
			}
		}
	}()

	mux := goji.NewMux()
	mux.UseC(relicWrapper.HandleHTTPC)
	mux.HandleFunc(pat.Get("/"), Root)

	user := goji.NewMux()
	user.Use(middleware.JSON)
	user.HandleFuncC(pat.Post("/tenants/:tenant_id/users"), CreateUser)
	user.HandleFuncC(pat.Get("/tenants/:tenant_id/users/:user_id"), GetUser)
	user.HandleFuncC(pat.Get("/tenants/:tenant_id/users"), GetUsers)
	user.HandleFuncC(pat.Put("/tenants/:tenant_id/users/:user_id"), UpdateUser)
	user.HandleFuncC(pat.Delete("/tenants/:tenant_id/users/:user_id"), DeleteUser)
	mux.HandleC(pat.New("/tenants/:tenant_id/*"), user)

	srv := &graceful.Server{
		Timeout: 3 * time.Second,
		BeforeShutdown: func() bool {
			mongo.Shutdown()
			redis.Shutdown()
			rabbitcage.Shutdown()
			return true
		},
		Server: &http.Server{
			Addr:    ":6969",
			Handler: mux,
		},
	}

	srv.ListenAndServe()
}

func Root(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "FCARE\nTripolis Solutions\n=============================\n")
}
