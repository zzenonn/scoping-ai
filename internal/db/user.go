package db

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
	scopingUser "github.com/zzenonn/scoping-ai/internal/user"
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

type UserRepository struct {
	client         *firestore.Client
	CollectionName string
}

func NewUserRepository(client *firestore.Client, collectionName string) UserRepository {
	return UserRepository{
		client:         client,
		CollectionName: collectionName,
	}
}

func convertUserToMap(user scopingUser.User) (map[string]interface{}, error) {
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

func (repo *UserRepository) GetUser(ctx context.Context, id string) (scopingUser.User, error) {
	doc, err := repo.client.Collection(repo.CollectionName).Doc(id).Get(ctx)
	if err != nil {
		return scopingUser.User{}, err
	}

	var user scopingUser.User
	err = doc.DataTo(&user)
	if err != nil {
		return scopingUser.User{}, err
	}

	user.ID = id

	return user, nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context, page int, pageSize int) ([]scopingUser.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	iter := repo.client.Collection(repo.CollectionName).OrderBy("email_address", firestore.Asc).Offset(offset).Limit(pageSize).Documents(ctx)
	var users []scopingUser.User

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var u scopingUser.User
		err = doc.DataTo(&u)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, user scopingUser.User) (scopingUser.User, error) {
	userMap, err := convertUserToMap(user)
	if err != nil {
		return scopingUser.User{}, err
	}

	_, err = repo.client.Collection(repo.CollectionName).Doc(user.ID).Set(ctx, userMap)
	if err != nil {
		return scopingUser.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, user scopingUser.User) (scopingUser.User, error) {
	userMap, err := convertUserToMap(user)
	if err != nil {
		return scopingUser.User{}, err
	}

	_, err = repo.client.Collection(repo.CollectionName).Doc(user.ID).Set(ctx, userMap, firestore.MergeAll)
	if err != nil {
		return scopingUser.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) DeleteUser(ctx context.Context, id string) error {
	_, err := repo.client.Collection(repo.CollectionName).Doc(id).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
