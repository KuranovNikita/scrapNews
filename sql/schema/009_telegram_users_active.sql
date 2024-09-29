-- +goose Up

ALTER TABLE telegramUsers ADD COLUMN active BOOLEAN NOT NULL;

-- +goose Down

ALTER TABLE telegramUsers DROP COLUMN active;