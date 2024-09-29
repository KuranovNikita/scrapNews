-- +goose Up

ALTER TABLE users
ADD CONSTRAINT unique_name UNIQUE (name);


-- +goose Down
ALTER TABLE users
DROP CONSTRAINT IF EXISTS unique_name;


