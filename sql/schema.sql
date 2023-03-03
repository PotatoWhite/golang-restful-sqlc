CREATE TABLE authors
(
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    bio  TEXT        NOT NULL
);

CREATE TABLE books
(
    id        BIGSERIAL PRIMARY KEY,
    title     VARCHAR(2000) NOT NULL,
    author_id BIGINT        NOT NULL,
    FOREIGN KEY (author_id) REFERENCES authors (id)
);