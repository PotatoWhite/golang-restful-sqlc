package authors

import (
	"context"
	"fmt"
	"github.com/potatowhite/restfulapi/pkg/database"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

type AuthorService interface {
	Create(ctx context.Context, cmd database.CreateAuthorParams) (*Author, error)
	Get(ctx context.Context, id int64) (*Author, error)
	Put(ctx context.Context, cmd database.UpdateAuthorParams) (*Author, error)
	Patch(ctx context.Context, cmd database.PartialUpdateAuthorParams) (*Author, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*Author, error)
}

type authorService struct {
	queries *database.Queries
}

func (a *authorService) Create(ctx context.Context, cmd database.CreateAuthorParams) (*Author, error) {
	author, err := a.queries.CreateAuthor(ctx, cmd)
	if err != nil {
		return nil, logging(fmt.Errorf("error creating author: %w", err))
	}

	return fromDB(author), nil
}

func logging(err error) error {
	logger.Printf(err.Error())
	return err
}

func (a *authorService) Patch(ctx context.Context, cmd database.PartialUpdateAuthorParams) (*Author, error) {
	author, err := a.queries.PartialUpdateAuthor(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("error updating author: %w", err)
	}
	return fromDB(author), nil
}

func (a *authorService) Get(ctx context.Context, id int64) (*Author, error) {
	author, err := a.queries.GetAuthor(ctx, id)
	if err != nil {
		return nil, err
	}

	return fromDB(author), nil
}

func (a *authorService) Put(ctx context.Context, cmd database.UpdateAuthorParams) (*Author, error) {
	author, err := a.queries.UpdateAuthor(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("error updating author: %w", err)
	}
	return fromDB(author), nil
}

func (a *authorService) Delete(ctx context.Context, id int64) error {
	err := a.queries.DeleteAuthor(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting author: %w", err)
	}
	return nil
}

func (a *authorService) List(ctx context.Context) ([]*Author, error) {
	authorList, err := a.queries.ListAuthors(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing authors: %w", err)
	}

	var apiAuthors []*Author
	for _, author := range authorList {
		apiAuthors = append(apiAuthors, fromDB(author))
	}

	return apiAuthors, nil
}

func fromDB(dbAuthor database.Author) *Author {
	return &Author{
		ID:   dbAuthor.ID,
		Name: dbAuthor.Name,
		Bio:  dbAuthor.Bio,
	}
}

func NewAuthorService(dbQueries *database.Queries) AuthorService {
	return &authorService{queries: dbQueries}
}
