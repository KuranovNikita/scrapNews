-- +goose Up

ALTER TABLE telegramUsers 
ADD CONSTRAINT unique_chat_id UNIQUE (chat_id);

-- +goose Down

ALTER TABLE telegramUsers 
DROP CONSTRAINT unique_chat_id;