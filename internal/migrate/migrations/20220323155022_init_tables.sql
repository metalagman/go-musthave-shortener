-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "urls"
(
    id           BIGSERIAL primary key,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,
    uid          UUID,
    original_url TEXT
);
CREATE UNIQUE INDEX IF NOT EXISTS urls_unique_original_url_null
    ON urls (original_url)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "urls";
-- +goose StatementEnd
