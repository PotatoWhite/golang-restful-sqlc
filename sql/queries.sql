-- name: CreateAuthor :one
INSERT INTO authors (name, bio)
VALUES ($1, $2)
    RETURNING *;

-- name: GetAuthor :one
SELECT *
FROM authors
WHERE id = $1
    LIMIT 1;

-- name: UpdateAuthor :one
UPDATE authors
SET name = $2,
    bio  = $3
WHERE id = $1
    RETURNING *;

-- name: PartialUpdateAuthor :one
UPDATE authors
SET name = CASE WHEN @update_name::boolean THEN @name::VARCHAR(32) ELSE name END,
    bio  = CASE WHEN @update_bio::boolean THEN @bio::TEXT ELSE bio END
WHERE id = @id
RETURNING *;

-- name: DeleteAuthor :exec
DELETE
FROM authors
WHERE id = $1;

-- name: ListAuthors :many
SELECT *
FROM authors
ORDER BY name;

-- name: TruncateAuthor :exec
TRUNCATE authors;







-- name: CreateBook :one
INSERT INTO books (title, author_id)
VALUES ($1, $2)
    RETURNING *;

-- name: GetBook :one
SELECT *
FROM books
WHERE id = $1
    LIMIT 1;

-- name: UpdateBook :one
UPDATE books
SET title = $2,
    author_id = $3
WHERE id = $1
    RETURNING *;

-- name: PartialUpdateBook :one
UPDATE books
SET title = CASE WHEN @update_title::boolean THEN @title::VARCHAR(2000) ELSE title END,
    author_id = CASE WHEN @update_author_id::boolean THEN @author_id::BIGINT ELSE author_id END
WHERE id = @id
RETURNING *;

-- name: DeleteBook :exec
DELETE
FROM books
WHERE id = $1;

-- name: ListBooks :many
SELECT *
FROM books
ORDER BY title;

-- name: TruncateBook :exec
TRUNCATE books;

-- name: ListBooksByAuthor :many
SELECT *
FROM books
WHERE author_id = $1
ORDER BY title;

-- name: ListBooksByAuthorName :many
SELECT *
FROM books
WHERE author_id = (
    SELECT id
    FROM authors
    WHERE name = $1
    LIMIT 1
)
ORDER BY title;

