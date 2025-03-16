-- +goose Up
ALTER TABLE users
ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN verification_token VARCHAR(255);

CREATE INDEX idx_verification_token ON users(verification_token);

-- +goose Down
ALTER TABLE users
DROP COLUMN email_verified,
DROP COLUMN verification_token;

DROP INDEX IF EXISTS idx_verification_token;