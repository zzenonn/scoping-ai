package db

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
	tnauser "github.com/zzenonn/trainocate-tna/internal/user"
	"google.golang.org/api/iterator"
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

const USER_COLLECTION_NAME = "users"

type UserRepository struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) UserRepository {
	return UserRepository{
		client: client,
	}
}

func convertUserToMap(user tnauser.User) (map[string]interface{}, error) {
	if user.Name == nil || user.EmailAddress == nil {
		return nil, ErrMissingRequiredFields
	}

	userMap := map[string]interface{}{
		"id":            user.ID,
		"name":          *user.Name,
		"email_address": *user.EmailAddress,
		"corporate":     user.Corporate,
	}

	if user.Company != nil {
		userMap["company"] = *user.Company
	}

	return userMap, nil
}

func (repo *UserRepository) GetUser(ctx context.Context, id string) (tnauser.User, error) {
	doc, err := repo.client.Collection(USER_COLLECTION_NAME).Doc(id).Get(ctx)
	if err != nil {
		return tnauser.User{}, err
	}

	var user tnauser.User
	err = doc.DataTo(&user)
	if err != nil {
		return tnauser.User{}, err
	}

	user.ID = id

	return user, nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context, page int, pageSize int) ([]tnauser.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	iter := repo.client.Collection(USER_COLLECTION_NAME).OrderBy("email_address", firestore.Asc).Offset(offset).Limit(pageSize).Documents(ctx)
	var users []tnauser.User

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var u tnauser.User
		err = doc.DataTo(&u)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, user tnauser.User) (tnauser.User, error) {
	userMap, err := convertUserToMap(user)
	if err != nil {
		return tnauser.User{}, err
	}

	_, err = repo.client.Collection(USER_COLLECTION_NAME).Doc(user.ID).Set(ctx, userMap)
	if err != nil {
		return tnauser.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, user tnauser.User) (tnauser.User, error) {
	userMap, err := convertUserToMap(user)
	if err != nil {
		return tnauser.User{}, err
	}

	_, err = repo.client.Collection(USER_COLLECTION_NAME).Doc(user.ID).Set(ctx, userMap)
	if err != nil {
		return tnauser.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) DeleteUser(ctx context.Context, id string) error {
	_, err := repo.client.Collection(USER_COLLECTION_NAME).Doc(id).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
