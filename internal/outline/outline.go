package CourseOutlines

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"

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

const COLLECTION_NAME = "course_outlines"

var (
	ErrFetchingOutline = errors.New("failed to fetch TNA questions by technology")
	ErrNotImplemented  = errors.New("this function is not yet implemented")
)

type CourseOutline struct {
	Id             string `json:"id,omitempty"`
	TechnologyName string `json:"technology_name,omitempty"`
	CourseCode     string `json:"course_code,omitempty"`
	CourseName     string `json:"course_name,omitempty"`
	Outline        string `json:"outline,omitempty"`
}

type CourseOutlineRepository interface {
	PostCourseOutline(ctx context.Context, courseOutline CourseOutline) (CourseOutline, error)
	GetCourseOutline(ctx context.Context, id string) (CourseOutline, error)
	GetAllCourseOutlines(ctx context.Context, page int, pageSize int) ([]CourseOutline, error)
	GetCourseOutlinesByFilter(ctx context.Context, page int, pageSize int, filterName string, filterValue string) ([]CourseOutline, error)
	UpdateCourseOutline(ctx context.Context, courseOutline CourseOutline) (CourseOutline, error)
	DeleteCourseOutline(ctx context.Context, id string) error
}

type CourseOutlineService struct {
	courseOutlineRepository CourseOutlineRepository
}

func NewCourseOutlineService(courseOutlineRepository CourseOutlineRepository) *CourseOutlineService {
	return &CourseOutlineService{courseOutlineRepository: courseOutlineRepository}
}

func (service *CourseOutlineService) PostCourseOutline(ctx context.Context, courseOutline CourseOutline) (CourseOutline, error) {
	log.Debug("Posting course outline . . .")

	courseOutline.Id = uuid.New().String()

	courseOutline, err := service.courseOutlineRepository.PostCourseOutline(ctx, courseOutline)

	if err != nil {
		log.Error("Failed to post course outline")
		return CourseOutline{}, err
	}

	return courseOutline, nil
}

func (service *CourseOutlineService) GetCourseOutline(ctx context.Context, id string) (CourseOutline, error) {
	log.Debug("Retreiving course outline by id . . .")

	courseOutline, err := service.courseOutlineRepository.GetCourseOutline(ctx, id)

	if err != nil {
		log.Error("Failed to retrieve course outline by id")
		return CourseOutline{}, err
	}

	return courseOutline, nil
}

// Generic function to get course outline by filter. Most common filter is technology name.
// e.g. GetCourseOutlineByFilter(ctx, "technology_name", "Java")
// Can also be used to get course outline by course code.
// e.g. GetCourseOutlineByFilter(ctx, "course_code", "JAV101")
func (service *CourseOutlineService) GetCourseOutlinesByFilter(ctx context.Context, page int, pageSize int, filterName string, filterValue string) ([]CourseOutline, error) {
	log.Debug("Retreiving all course outlines by filter . . .")

	courseOutlines, err := service.courseOutlineRepository.GetCourseOutlinesByFilter(ctx, page, pageSize, filterName, filterValue)

	if err != nil {
		log.Error("Failed to retrieve course outlines by filter")
		return nil, err
	}

	return courseOutlines, nil
}

func (service *CourseOutlineService) GetAllCourseOutlines(ctx context.Context, page int, pageSize int) ([]CourseOutline, error) {
	log.Debug("Retreiving all course outlines . . .")

	courseOutlines, err := service.courseOutlineRepository.GetAllCourseOutlines(ctx, page, pageSize)

	if err != nil {
		log.Error("Failed to retrieve all course outlines")
		return nil, err
	}

	return courseOutlines, nil
}

func (service *CourseOutlineService) UpdateCourseOutline(ctx context.Context, courseOutline CourseOutline) (CourseOutline, error) {
	log.Debug("Updating course outline . . .")

	courseOutline, err := service.courseOutlineRepository.UpdateCourseOutline(ctx, courseOutline)

	if err != nil {
		log.Error("Failed to update course outline")
		return CourseOutline{}, err
	}

	return courseOutline, nil
}

func (service *CourseOutlineService) DeleteCourseOutline(ctx context.Context, id string) error {
	log.Debug("Deleting course outline . . .")

	err := service.courseOutlineRepository.DeleteCourseOutline(ctx, id)

	if err != nil {
		log.Error("Failed to delete course outline")
		return err
	}

	return nil
}
