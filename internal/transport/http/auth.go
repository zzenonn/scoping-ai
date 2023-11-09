package http

import (
	"context"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	log "github.com/sirupsen/logrus"
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

var firebaseApp *firebase.App

func VerifyFirebaseToken(ctx context.Context, idToken string) bool {
	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Errorf("error getting Auth client: %v\n", err)
		return false
	}

	_, err = client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Errorf("error verifying ID token: %v\n", err)
		return false
	}

	return true
}
