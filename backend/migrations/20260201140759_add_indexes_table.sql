-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS indexes
(
    id          UUID        NOT NULL,
    index_code  VARCHAR     NOT NULL UNIQUE,  -- например "MOEXBC"
    name        VARCHAR     NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT now()

    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS indexes_code_idx ON indexes (index_code);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS indexes;

-- +goose StatementEnd
