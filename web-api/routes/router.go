package routes

import (
	"net/http"

	redis "gopkg.in/redis.v5"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Router for the command line arguments we accept
type Router struct {
	dbConnection    func() *gorm.DB
	redisConnection func() *redis.Client
	router          *mux.Router
}

func checkError(err error) {
	if nil != err {
		log.WithError(err).Fatal("Fatal error")
		panic(err)
	}
}

// NewRouter returns a new instance of *Router
func NewRouter(
	dbConnection func() *gorm.DB,
	redisConnection func() *redis.Client,
) *Router {
	// Create a new instance of router
	r := &Router{
		dbConnection:    dbConnection,
		redisConnection: redisConnection,
		router:          mux.NewRouter(),
	}
	// Bind the routes to the router
	r.router.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) { r.getHelloWorldHandler(rw, req) }).Methods("GET")
	r.router.HandleFunc("/health-check.json", func(rw http.ResponseWriter, req *http.Request) { r.getHealthCheckHandler(rw, req) }).Methods("GET")
	r.router.HandleFunc("/status.json", func(rw http.ResponseWriter, req *http.Request) { r.StatusHandler(rw, req) }).Methods("GET")
	r.router.HandleFunc("/symbols.json", func(rw http.ResponseWriter, req *http.Request) { r.getSymbolsHandler(rw, req) }).Methods("GET")
	r.router.HandleFunc("/symbols/{id}.json", func(rw http.ResponseWriter, req *http.Request) { r.getSymbolHandler(rw, req) }).Methods("GET")
	r.router.HandleFunc("/contracts.json", func(rw http.ResponseWriter, req *http.Request) { r.getContractsHandler(rw, req) }).Methods("GET")
	r.router.HandleFunc("/contracts/{id}.json", func(rw http.ResponseWriter, req *http.Request) { r.getContractHandler(rw, req) }).Methods("GET")
	return r
}

// Handler Return the http.Handler interface to be used by the http.Server
func (r *Router) Handler() http.Handler {
	return r.router
}

func (r *Router) getHelloWorldHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content Type", "application/json")
	rw.Write([]byte("{ 'status': 'hello world!' }"))
}

func (r *Router) getHealthCheckHandler(rw http.ResponseWriter, req *http.Request) {
}

func (r *Router) getStatusHandler(rw http.ResponseWriter, req *http.Request) {
}

func (r *Router) getSymbolsHandler(rw http.ResponseWriter, req *http.Request) {
}

func (r *Router) getSymbolHandler(rw http.ResponseWriter, req *http.Request) {
}

func (r *Router) getContractsHandler(rw http.ResponseWriter, req *http.Request) {
}

func (r *Router) getContractHandler(rw http.ResponseWriter, req *http.Request) {
}
