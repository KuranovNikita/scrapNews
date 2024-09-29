-- +goose Up

CREATE TABLE siteParse (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE,
    url_site TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL UNIQUE
);


-- +goose Down
DROP TABLE siteParse;