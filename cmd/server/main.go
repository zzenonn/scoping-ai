package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"strings"

	"github.com/zzenonn/trainocate-tna/internal/db"
	tna "github.com/zzenonn/trainocate-tna/internal/tna"
	transportHttp "github.com/zzenonn/trainocate-tna/internal/transport/http"
)

func init() {

	// Set log level based on environment variables
	switch logLevel := strings.ToLower(os.Getenv("LOG_LEVEL")); logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}

}

// Instantiate and startup go app
func Run(projectName string) error {
	log.Println("starting up the application")

	firestoreDb, err := db.NewDatabase(projectName)

	if err != nil {
		log.Error("Failed to connect to the database")
		return err
	}

	// if err := psqlDb.MigrateDb(); err != nil {
	// 	log.Error("Failed to migrate the database")
	// 	return err
	// }

	qSetRepository := db.NewQuestionSetRepository(firestoreDb.Client)
	qSetService := tna.NewQuestionService(&qSetRepository)
	qSetHandler := transportHttp.NewQuestionSetHandler(qSetService)

	httpHandler := transportHttp.NewMainHandler()

	httpHandler.AddHandler(qSetHandler)

	httpHandler.MapRoutes()

	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	projectName := flag.String("project-name", "", "The name of the project (required)")
	flag.Parse()

	if *projectName == "" {
		log.Debug("The 'project-name' flag is required")
		flag.Usage()
		os.Exit(1)
	}

	log.Infof("the server is up with project: %s", *projectName)

	if err := Run(*projectName); err != nil {
		log.Error(err)
	}
}
