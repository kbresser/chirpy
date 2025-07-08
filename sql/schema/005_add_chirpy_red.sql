-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red bool DEFAULT FALSE NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;
