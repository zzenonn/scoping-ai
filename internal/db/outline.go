package db

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
	outline "github.com/zzenonn/trainocate-tna/internal/outline"
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

const COURSE_OUTLINE_COLLECTION_NAME = "course_outlines"

type CourseOutlineRepository struct {
	client *firestore.Client
}

func convertOutlineToMap(courseOutline outline.CourseOutline) map[string]interface{} {
	return map[string]interface{}{
		"technology_name": courseOutline.TechnologyName,
		"course_code":     courseOutline.CourseCode,
		"course_name":     courseOutline.CourseName,
		"outline":         courseOutline.Outline,
	}
}

func NewCourseOutlineRepository(client *firestore.Client) CourseOutlineRepository {
	return CourseOutlineRepository{
		client: client,
	}
}

func (repo *CourseOutlineRepository) PostCourseOutline(ctx context.Context, cOutline outline.CourseOutline) (outline.CourseOutline, error) {
	cOutlineMap := convertOutlineToMap(cOutline)

	_, err := repo.client.Collection(COURSE_OUTLINE_COLLECTION_NAME).Doc(cOutline.Id).Set(ctx, cOutlineMap)
	if err != nil {
		return outline.CourseOutline{}, err
	}

	return cOutline, nil
}

func (repo *CourseOutlineRepository) GetCourseOutline(ctx context.Context, docID string) (outline.CourseOutline, error) {
	doc, err := repo.client.Collection(COURSE_OUTLINE_COLLECTION_NAME).Doc(docID).Get(ctx)
	if err != nil {
		return outline.CourseOutline{}, err
	}

	var cOutline outline.CourseOutline

	// Print contents of doc
	log.Debug(doc.Data())

	doc.DataTo(&cOutline)

	cOutline.Id = docID

	return cOutline, nil
}

func (repo *CourseOutlineRepository) GetCourseOutlinesByFilter(
	ctx context.Context, page int, pageSize int,
	filterName string, filterValue string,
) ([]outline.CourseOutline, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Query broken down for readability
	query := repo.client.Collection(COURSE_OUTLINE_COLLECTION_NAME)
	filteredQuery := query.Where(filterName, "==", filterValue)
	orderedQuery := filteredQuery.OrderBy(filterName, firestore.Asc)
	paginatedQuery := orderedQuery.Offset(offset).Limit(pageSize)

	iter := paginatedQuery.Documents(ctx)

	var cOutlines []outline.CourseOutline

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var outline outline.CourseOutline
		err = doc.DataTo(&outline)
		if err != nil {
			return nil, err
		}

		cOutlines = append(cOutlines, outline)
	}

	return cOutlines, nil
}

func (repo *CourseOutlineRepository) GetAllCourseOutlines(ctx context.Context, page int, pageSize int) ([]outline.CourseOutline, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	iter := repo.client.Collection(COURSE_OUTLINE_COLLECTION_NAME).OrderBy("course_code", firestore.Asc).Offset(offset).Limit(pageSize).Documents(ctx)
	var cOutlines []outline.CourseOutline

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var outline outline.CourseOutline
		err = doc.DataTo(&outline)
		if err != nil {
			return nil, err
		}

		cOutlines = append(cOutlines, outline)
	}

	return cOutlines, nil
}

func (repo *CourseOutlineRepository) UpdateCourseOutline(ctx context.Context, cOutline outline.CourseOutline) (outline.CourseOutline, error) {
	cOutlineMap := convertOutlineToMap(cOutline)

	_, err := repo.client.Collection(COURSE_OUTLINE_COLLECTION_NAME).Doc(cOutline.Id).Set(ctx, cOutlineMap)
	if err != nil {
		return outline.CourseOutline{}, err
	}

	return cOutline, nil
}

func (repo *CourseOutlineRepository) DeleteCourseOutline(ctx context.Context, docID string) error {
	_, err := repo.client.Collection(COURSE_OUTLINE_COLLECTION_NAME).Doc(docID).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
