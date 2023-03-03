package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/potatowhite/restfulapi/pkg/database"
	"github.com/potatowhite/restfulapi/pkg/handler/authors"
	"net/http"
)

func NewService(queries *database.Queries) *Service {
	return &Service{queries: queries}
}

type Service struct {
	queries *database.Queries
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.POST("/authors", s.Create)
	router.GET("/authors/:id", s.Get)
	router.PUT("/authors/:id", s.FullUpdate)
	router.PATCH("/authors/:id", s.PartialUpdate)
	router.DELETE("/authors/:id", s.Delete)
	router.GET("/authors", s.List)
}

func (s *Service) Create(c *gin.Context) {
	// Parse request
	var request authors.ApiAuthor
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create author
	params := database.CreateAuthorParams{
		Name: request.Name,
		Bio:  request.Bio,
	}
	author, err := s.queries.CreateAuthor(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromDB(author)
	c.IndentedJSON(http.StatusCreated, response)
}

func (s *Service) Get(c *gin.Context) {
	// Parse request
	var pathParams authors.PathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get author
	author, err := s.queries.GetAuthor(context.Background(), pathParams.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromDB(author)
	c.IndentedJSON(http.StatusOK, response)
}

func (s *Service) FullUpdate(c *gin.Context) {
	// Parse request
	var pathParams authors.PathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request authors.ApiAuthor
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update author
	params := database.UpdateAuthorParams{
		ID:   pathParams.ID,
		Name: request.Name,
		Bio:  request.Bio,
	}
	fmt.Println(params)
	author, err := s.queries.UpdateAuthor(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromDB(author)
	c.IndentedJSON(http.StatusOK, response)
}

func (s *Service) PartialUpdate(c *gin.Context) {
	// Parse request
	var pathParams authors.PathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request authors.ApiAuthorPartialUpdate
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update author
	params := database.PartialUpdateAuthorParams{ID: pathParams.ID}
	if request.Name != nil {
		params.UpdateName = true
		params.Name = *request.Name
	}
	if request.Bio != nil {
		params.UpdateBio = true
		params.Bio = *request.Bio
	}
	author, err := s.queries.PartialUpdateAuthor(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := fromDB(author)
	c.IndentedJSON(http.StatusOK, response)
}

func (s *Service) Delete(c *gin.Context) {
	// Parse request
	var pathParams authors.PathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete author
	if err := s.queries.DeleteAuthor(context.Background(), pathParams.ID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	c.Status(http.StatusOK)
}

func (s *Service) List(c *gin.Context) {
	// List authors
	authorList, err := s.queries.ListAuthors(context.Background())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	if len(authorList) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Build response
	var response []*authors.ApiAuthor
	for _, author := range authorList {
		response = append(response, fromDB(author))
	}
	c.IndentedJSON(http.StatusOK, authorList)
}

func fromDB(author database.Author) *authors.ApiAuthor {
	return &authors.ApiAuthor{
		ID:   author.ID,
		Name: author.Name,
		Bio:  author.Bio,
	}
}
