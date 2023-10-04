package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"

	"github.com/zzenonn/trainocate-tna/internal/db"
	tnamessage "github.com/zzenonn/trainocate-tna/internal/message"
	outline "github.com/zzenonn/trainocate-tna/internal/outline"
	tna "github.com/zzenonn/trainocate-tna/internal/tna"
	transportHttp "github.com/zzenonn/trainocate-tna/internal/transport/http"
	tnauser "github.com/zzenonn/trainocate-tna/internal/user"
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

func getSecret(secretName string) (string, error) {
	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	secretRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	result, err := client.AccessSecretVersion(ctx, secretRequest)
	if err != nil {
		return "", err
	}

	log.Debug(result.Payload.Data)

	return string(result.Payload.Data), nil
}

// Instantiate and startup go app
func Run(projectName string) error {
	log.Println("starting up the application")

	firestoreDb, err := db.NewDatabase(projectName)

	if err != nil {
		log.Error("Failed to connect to the database")
		return err
	}

	openAIKeySecretName := fmt.Sprintf("projects/%s/secrets/OpenAIAPIKey/versions/latest", projectName)

	openAPIKey, err := getSecret(openAIKeySecretName)

	if err != nil {
		log.Error("Failed to get the OpenAI API key")
	}

	qSetRepository := db.NewQuestionSetRepository(firestoreDb.Client, "question_sets")
	qSetService := tna.NewQuestionService(&qSetRepository)
	qSetHandler := transportHttp.NewQuestionSetHandler(qSetService)

	cOutlineRepository := db.NewCourseOutlineRepository(firestoreDb.Client, "course_outlines")
	cOutlineService := outline.NewCourseOutlineService(&cOutlineRepository)
	cOutlineHandler := transportHttp.NewCourseOutlineHandler(cOutlineService)

	userRepository := db.NewUserRepository(firestoreDb.Client, "users")
	userService := tnauser.NewUserService(&userRepository)
	userHandler := transportHttp.NewUserHandler(userService)

	openAiRepository := db.NewOpenAiRepository(openAPIKey, "https://api.openai.com/v1/chat/completions", "gpt-4", 1.0)
	messageRepository := db.NewMessageRepository(firestoreDb.Client, "messages")
	messageService := tnamessage.NewMessageService(&messageRepository, &openAiRepository)
	messageHandler := transportHttp.NewMessageHandler(messageService)

	httpHandler := transportHttp.NewMainHandler()

	httpHandler.AddHandler(qSetHandler)
	httpHandler.AddHandler(cOutlineHandler)
	httpHandler.AddHandler(userHandler)
	httpHandler.AddHandler(messageHandler)

	httpHandler.MapRoutes()

	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	projectId := flag.String("project-id", "", "The id of the project (required)")
	flag.Parse()

	if *projectId == "" {
		log.Debug("The 'project-id' flag is required")
		flag.Usage()
		os.Exit(1)
	}

	log.Infof("the server is up with project: %s", *projectId)

	if err := Run(*projectId); err != nil {
		log.Error(err)
	}
}
