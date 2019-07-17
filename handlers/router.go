package handlers

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"github.com/spf13/viper"

	"github.com/kenjones-cisco/dapperdox/config"
	"github.com/kenjones-cisco/dapperdox/handlers/guides"
	"github.com/kenjones-cisco/dapperdox/handlers/home"
	"github.com/kenjones-cisco/dapperdox/handlers/reference"
	"github.com/kenjones-cisco/dapperdox/handlers/specs"
	"github.com/kenjones-cisco/dapperdox/handlers/static"
	"github.com/kenjones-cisco/dapperdox/handlers/timeout"
	log "github.com/kenjones-cisco/dapperdox/logger"
	"github.com/kenjones-cisco/dapperdox/network"
	"github.com/kenjones-cisco/dapperdox/proxy"
	"github.com/kenjones-cisco/dapperdox/render"
	"github.com/kenjones-cisco/dapperdox/spec"
	"github.com/kenjones-cisco/dapperdox/version"
)

// NewRouterChain creates a router with a chain of middlewares that acts as an http.Handler
func NewRouterChain() http.Handler {
	router := pat.New()
	chain := alice.New(withLogger, timeoutHandler, withCsrf, injectHeaders).Then(router)

	listener, err := network.NewListener()
	if err != nil {
		log.Logger().Fatalf("%s", err)
	}

	var wg sync.WaitGroup
	var sg sync.WaitGroup
	sg.Add(1)

	go func() {
		log.Logger().Trace("Listen for and serve swagger spec requests for start up")
		wg.Add(1)
		sg.Done()
		_ = http.Serve(listener, chain)
		log.Logger().Trace("Finished service swagger specs for start up")
		wg.Done()
	}()

	sg.Wait()

	// Register the spec routes (Listener and server must be up and running by now)
	specs.Register(router)

	if err = spec.LoadSpecifications(); err != nil {
		log.Logger().Fatalf("Load specification error: %s", err)
	}

	render.Register()

	reference.Register(router)
	guides.Register(router)
	static.Register(router)

	home.Register(router)
	proxy.Register(router)

	_ = listener.Close() // Stop serving specs
	wg.Wait()            // wait for go routine serving specs to terminate

	return chain
}

func withLogger(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, h)
}

func withCsrf(h http.Handler) http.Handler {
	csrfHandler := nosurf.New(h)
	csrfHandler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rsn := nosurf.Reason(req).Error()
		log.Logger().Warnf("failed csrf validation: %s", rsn)
		render.HTML(w, http.StatusBadRequest, "error", map[string]interface{}{"error": rsn})
	}))
	return csrfHandler
}

func timeoutHandler(h http.Handler) http.Handler {
	return timeout.Handler(h, 1*time.Second, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Logger().Warn("request timed out")
		render.HTML(w, http.StatusRequestTimeout, "error", map[string]interface{}{"error": "Request timed out"})
	}))
}

// Handle additional headers such as strict transport security for TLS, and
// giving the Server name.
func injectHeaders(h http.Handler) http.Handler {
	tlsEnabled := viper.GetString(config.TLSCert) != "" && viper.GetString(config.TLSKey) != ""

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", fmt.Sprintf("%s %s", version.ProductName, version.Version))

		if tlsEnabled {
			w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}

		h.ServeHTTP(w, r)
	})
}
