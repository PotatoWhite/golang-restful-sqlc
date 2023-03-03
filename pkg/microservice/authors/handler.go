package authors

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/potatowhite/restfulapi/pkg/database"
	"net/http"
)

// logger
type AuthorHandler interface {
	RegisterHandlers(router *gin.Engine)
}

type authorHandler struct {
	service AuthorService
}

func NewAuthorHandler(service AuthorService) AuthorHandler {
	return &authorHandler{service: service}
}

func (h *authorHandler) RegisterHandlers(router *gin.Engine) {
	router.POST("/authors", h.Create)
	router.GET("/authors/:id", h.Get)
	router.PUT("/authors/:id", h.Put)
	router.PATCH("/authors/:id", h.Patch)
	router.DELETE("/authors/:id", h.Delete)
	router.GET("/authors", h.List)
}

func (h *authorHandler) Create(c *gin.Context) {
	var req Author
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if author, err := h.service.Create(c, database.CreateAuthorParams{Name: req.Name, Bio: req.Bio}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusCreated, author)
	}
}

func (h *authorHandler) Get(c *gin.Context) {
	var pathParams PathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := h.service.Get(c, pathParams.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, author)
}

func (h *authorHandler) Put(c *gin.Context) {
	var pathParams PathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req Author
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := h.service.Put(c, database.UpdateAuthorParams{ID: pathParams.ID, Name: req.Name, Bio: req.Bio})
	if err != nil {
		// no row
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	} else {
		c.JSON(http.StatusOK, author)
	}
}

func (h *authorHandler) Patch(c *gin.Context) {
	var pathParam PathParameters
	if err := c.ShouldBindUri(&pathParam); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// print the request body
	var req AuthorPartialUpdate
	if err := c.ShouldBindJSON(&req); err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params := database.PartialUpdateAuthorParams{ID: pathParam.ID}
	if req.Name != nil {
		params.Name = *req.Name
		params.UpdateName = true
	}
	if req.Bio != nil {
		params.Bio = *req.Bio
		params.UpdateBio = true
	}

	author, err := h.service.Patch(c, params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update author"})
		return
	}

	c.JSON(http.StatusOK, author)
}

func (h *authorHandler) Delete(c *gin.Context) {
	var pathParams PathParameters
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(c, pathParams.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *authorHandler) List(c *gin.Context) {
	authors, err := h.service.List(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(authors) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, authors)
}
