package main

import (
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/tylerb/graceful"
	"goji.io"
	"goji.io/pat"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	mux := goji.NewMux()
	//	mux.UseC(relicWrapper.HandleHTTPC)
	mux.HandleFunc(pat.Get("/"), Root)

	user := goji.NewMux()
	//	user.Use(middleware.JSON)
	user.HandleFuncC(pat.Post("/tenants/:tenant_uuid/users"), CreateUser)
	user.HandleFuncC(pat.Get("/tenants/:tenant_uuid/users/:user_id"), GetUser)
	mux.HandleC(pat.New("/tenants/:tenant_uuid/*"), user)

	srv := &graceful.Server{
		Timeout: 3 * time.Second,
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
