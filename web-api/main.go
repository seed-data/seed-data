package main

import (
	"net/http"
	"os"
	"sync"
	"time"

	"./adapters"
	"./models"
	"./routes"
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	cli "gopkg.in/urfave/cli.v1"
)

// Opts for the command line arguments we accept
type Opts struct {
	verbose              bool
	autoMigrate          bool
	backfillFrom         string
	importWorkerCount    number
	importWorkerInterval time.Duration
}

func checkError(err error) {
	if nil != err {
		log.WithError(err).Fatal("Fatal error")
		panic(err)
	}
}

func dbConnection() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost user=docker dbname=docker password=docker sslmode=disable")
	checkError(err)
	return db
}

func main() {
	opts := Opts{
		verbose:              false,
		autoMigrate:          false,
		importWorkerCount:    0,
		importWorkerInterval: time.Duration(15 * time.Minute),
	}

	app := cli.NewApp()
	app.Name = "@seed-data/web-api"
	app.Version = "2016.12.28"
	app.Usage = "Powers your API needs"

	// Explicitly setting so that the short code for version is "V" and verbose can be "v"
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "Output the app version",
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "verbose, v",
			Usage:       "Verbose logging",
			Destination: &opts.verbose,
		},
		cli.StringFlag{
			Name:        "backfill-from",
			Usage:       "Seed the data with this dataset.  Duplicates are automatically ignored",
			Destination: &opts.backfillFrom,
		},
		cli.IntFlag{
			Name:        "import-workers",
			Usage:       "The number of workers to spawn in the background.  These workers listen for new intra-day and end-of-day quotes.",
			Value:       1,
			Destination: &opts.importWorkerCount,
		},
		cli.DurationFlag{
			Name:        "import-interval",
			Usage:       "How often should we kick off a new round of workers to import intra-day data.",
			Value:       time.Duration(15 * time.Minute),
			Destination: &opts.importWorkerInterval,
		},
	}

	app.Commands = []cli.Command{
	// Add additional commands we want to run here
	}

	// Create a closure & pass the opts down with the cli context
	app.Action = func(context *cli.Context) error {
		return runMain(opts, context)
	}
	app.Run(os.Args)
}

func importWorkerTicker(ticker *time.Ticker, done chan interface{}, workerCount int) {
	counter := int64(0)
	for {
		select {
		case <-done:
			log.Info("Exiting importWorerTicker")
			return
		case tick := <-ticker:
			counter++
			log.Infof("Running tick #%v at %v\n", counter, tick)
			checkError(adapters.ImportDailyData(dbConnection, workerCount))
			log.Infof("Completed tick #%v at %v\n", counter, tick)
		}
	}
}

func runMain(opts Opts, context *cli.Context) error {
	if opts.verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("Opts = %+v\n", opts)

	// Execute the migrations
	if opts.autoMigrate {
		log.Error("Running migrations")
		models.AutoMigrate(dbConnection)
		log.Error("Completed migrations")
	}

	// Execute the back-fill operation
	if opts.backfillFrom {
		log.Error("Running backfill")
		models.Backfill(dbConnection, opts.backfillFrom)
		log.Error("Completed backfill")
	}

	var wg sync.WaitGroup
	doneChannel := make(chan interface{})

	// Kick off the background interval worker
	if opts.importWorkerCount > 0 && opts.importWorkerInterval > time.Duration(0) {
		ticker := time.NewTicker(opts.importWorkerInterval)
		tickHandler := func() {
			defer wg.Done()
			importWorkerTicker(ticker, doneChannel, opts.importWorkerCount)
		}
		wg.Add(1)
		go tickHandler()
	}

	// Wait for the background workers to complete
	defer func() {
		log.Info("Closing the done-channel.")
		close(doneChannel)
		log.Info("Waiting for background workers to finish.")
		wg.Wait()
		log.Info("All background workers are done.")
		log.Info("Completing shutdown")
	}()

	// Create the router
	router := mux.NewRouter()
	// Helper method to memoize the handler functions and pass in a db instance
	addGet := func(path string, handler func(*gorm.DB, http.ResponseWriter, *http.Request)) {
		f := func(rw http.ResponseWriter, req *http.Request) {
			db := dbConnection()
			defer db.Close()
			handler(db, rw, req)
		}
		router.HandleFunc(path, f).Methods("GET")
	}
	addGet("/", routes.HelloWorldHandler)
	addGet("/health-check.json", routes.HealthCheckHandler)
	addGet("/status.json", routes.StatuHandler)
	addGet("/symbols.json", routes.GetSymbolsHandler)
	addGet("/symbols/{id}.json", routes.GetSymbolHandler)
	addGet("/contracts.json", routes.GetContractsHandler)
	addGet("/contracts/{id}.json", routes.GetContractHandler)

	// This will serve files under http://localhost:80/static/<filename>
	router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// Create the new server & attach the mux router
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:80",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		// Limit the max-header size
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	checkError(err)
}
