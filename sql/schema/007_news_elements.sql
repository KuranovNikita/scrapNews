-- +goose Up

CREATE TABLE newsElements (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    site_parse_id UUID NOT NULL REFERENCES siteParse(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    news_date TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE
);


-- +goose Down
DROP TABLE newsElements;