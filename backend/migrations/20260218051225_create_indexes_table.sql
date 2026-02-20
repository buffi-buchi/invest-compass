-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS indexes
(
    ticker      VARCHAR,
    name        VARCHAR     NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (ticker)
);

CREATE INDEX IF NOT EXISTS indexes_code_idx ON indexes (ticker);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS indexes;

-- +goose StatementEnd