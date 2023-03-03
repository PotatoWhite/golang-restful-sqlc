package authors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/potatowhite/restfulapi/cmd/config"
	"github.com/potatowhite/restfulapi/pkg/database"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type apiError struct {
	Error string
}

type ServiceTestSuite struct {
	suite.Suite
	router  *gin.Engine
	queries *database.Queries
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupSuite() {
	cfg, err := config.Read()
	s.Require().NoError(err)

	postgres, err := database.NewPostgres(cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Dbname)
	s.Require().NoError(err)

	s.queries = database.New(postgres.DB)
	service := NewAuthorService(s.queries)
	handler := NewAuthorHandler(service)

	s.router = gin.Default()
	handler.RegisterHandlers(s.router)
}

func (s *ServiceTestSuite) SetupTest() {
	s.queries.TruncateAuthor(context.Background())
}

func (s *ServiceTestSuite) TestCreateAuthor() {
	// Arrange
	author := Author{
		Name: "test name",
		Bio:  "test bio",
	}

	var buffer bytes.Buffer
	s.Require().NoError(json.NewEncoder(&buffer).Encode(author))

	// Act
	request, err := http.NewRequest(http.MethodPost, "/authors", &buffer)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusCreated, rec.Result().StatusCode)

	// Assert Response Body
	var created Author
	if err := json.NewDecoder(rec.Result().Body).Decode(&created); err != nil {
		log.Printf("Error decoding rec body: %v", err)
	}

	s.Require().Equal(author.Name, created.Name)
	s.Require().Equal(author.Bio, created.Bio)
}

func (s *ServiceTestSuite) TestCreateAuthor_InvalidRequest() {
	// Arrange
	author := Author{
		Name: "test name",
	}

	var buffer bytes.Buffer
	s.Require().NoError(json.NewEncoder(&buffer).Encode(author))

	// Act
	request, err := http.NewRequest(http.MethodPost, "/authors", &buffer)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusBadRequest, rec.Result().StatusCode)

	// Assert Response Body
	var apiErr apiError
	if err := json.NewDecoder(rec.Result().Body).Decode(&apiErr); err != nil {
		log.Printf("Error decoding rec body: %v", err)
	}

	s.Require().Equal("Key: 'Author.Bio' Error:Field validation for 'Bio' failed on the 'required' tag", apiErr.Error)
}

func (s *ServiceTestSuite) TestGetAuthor() {
	// Arrange
	author := Author{
		Name: "test name",
		Bio:  "test bio",
	}

	created, err := s.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: author.Name,
		Bio:  author.Bio,
	})
	s.Require().NoError(err)

	// Act
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/authors/%d", created.ID), nil)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusOK, rec.Result().StatusCode)

	// Assert Response Body
	var got Author
	if err := json.NewDecoder(rec.Result().Body).Decode(&got); err != nil {
		log.Printf("Error decoding rec body: %v", err)
	}

	s.Require().Equal(author.Name, got.Name)
	s.Require().Equal(author.Bio, got.Bio)
}

func (s *ServiceTestSuite) TestGetAuthor_NotFound() {
	// Act
	request, err := http.NewRequest(http.MethodGet, "/authors/1", nil)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusNoContent, rec.Result().StatusCode)

	// Assert Response Body
	var apiErr apiError
	if err := json.NewDecoder(rec.Result().Body).Decode(&apiErr); err != nil {
		log.Printf("Error decoding rec body: %v", err)
	}
}

func (s *ServiceTestSuite) TestListAuthors() {
	// Arrange
	author1 := Author{
		Name: "test name 1",
		Bio:  "test bio 1",
	}

	author2 := Author{
		Name: "test name 2",
		Bio:  "test bio 2",
	}

	_, err := s.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: author1.Name,
		Bio:  author1.Bio,
	})
	s.Require().NoError(err)

	_, err = s.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: author2.Name,
		Bio:  author2.Bio,
	})
	s.Require().NoError(err)

	// Act
	request, err := http.NewRequest(http.MethodGet, "/authors", nil)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusOK, rec.Result().StatusCode)

	// Assert Response Body
	var got []Author
	if err := json.NewDecoder(rec.Result().Body).Decode(&got); err != nil {
		log.Printf("Error decoding rec body: %v", err)
	}

	s.Require().Len(got, 2)
	s.Require().Equal(author1.Name, got[0].Name)
	s.Require().Equal(author1.Bio, got[0].Bio)
	s.Require().Equal(author2.Name, got[1].Name)
	s.Require().Equal(author2.Bio, got[1].Bio)
}

func (s *ServiceTestSuite) TestListAuthors_Empty() {
	// Act
	request, err := http.NewRequest(http.MethodGet, "/authors", nil)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusNoContent, rec.Result().StatusCode)

	// Assert Response Body
	var got []Author
	if err := json.NewDecoder(rec.Result().Body).Decode(&got); err != nil {
		log.Printf("Error decoding rec body: %v", err)
	}

	s.Require().Len(got, 0)
}

func (s *ServiceTestSuite) TestUpdateAuthor() {
	// Arrange
	author := Author{
		Name: "test name",
		Bio:  "test bio",
	}

	created, err := s.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: author.Name,
		Bio:  author.Bio,
	})
	s.Require().NoError(err)

	author.Name = "updated name"
	author.Bio = "updated bio"

	var buffer bytes.Buffer
	s.Require().NoError(json.NewEncoder(&buffer).Encode(author))

	// Act
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/authors/%d", created.ID), &buffer)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusOK, rec.Result().StatusCode)

	// Assert Response Body
	var updated Author
	if err := json.NewDecoder(rec.Result().Body).Decode(&updated); err != nil {
		log.Printf("Error decoding rec body: %v", err)
	}

	s.Require().Equal(author.Name, updated.Name)
	s.Require().Equal(author.Bio, updated.Bio)
}

func (s *ServiceTestSuite) TestUpdateAuthor_NotFound() {
	// Arrange
	author := Author{
		Name: "test name",
		Bio:  "test bio",
	}

	var buffer bytes.Buffer
	s.Require().NoError(json.NewEncoder(&buffer).Encode(author))

	// Act
	request, err := http.NewRequest(http.MethodPut, "/authors/1", &buffer)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusNotFound, rec.Result().StatusCode)

	// Assert Response Body
	var apiErr apiError
	if err := json.NewDecoder(rec.Result().Body).Decode(&apiErr); err != nil {
		log.Printf("Error decoding rec body: %v", err)
	}
}

func (s *ServiceTestSuite) TestPartialUpdateAuthor() {
	// Arrange
	author := Author{
		Name: "test name",
		Bio:  "test bio",
	}

	created, err := s.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: author.Name,
		Bio:  author.Bio,
	})
	s.Require().NoError(err)

	author.Name = "updated name"

	var buffer bytes.Buffer
	s.Require().NoError(json.NewEncoder(&buffer).Encode(author))

	// Act
	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/authors/%d", created.ID), &buffer)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusOK, rec.Result().StatusCode)

	// Assert Response Body
	var updated Author
	if err := json.NewDecoder(rec.Result().Body).Decode(&updated); err != nil {
		log.Printf("Error decoding rec body: %v", err)
	}

	s.Require().Equal(author.Name, updated.Name)
	s.Require().Equal(author.Bio, updated.Bio)
}

func (s *ServiceTestSuite) TestDeleteAuthor() {
	// Arrange
	author := Author{
		Name: "test name",
		Bio:  "test bio",
	}

	created, err := s.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: author.Name,
		Bio:  author.Bio,
	})
	s.Require().NoError(err)

	// Act
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/authors/%d", created.ID), nil)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusNoContent, rec.Result().StatusCode)

}

func (s *ServiceTestSuite) TestDeleteAuthor_NotFound() {
	// Act
	request, err := http.NewRequest(http.MethodDelete, "/authors/1", nil)
	s.Require().NoError(err)

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, request)

	// Assert Status Code
	s.Require().Equal(http.StatusNoContent, rec.Result().StatusCode)
}
