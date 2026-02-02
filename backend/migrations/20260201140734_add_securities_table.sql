-- +goose Up
-- +goose StatementBegin

CREATE TYPE security_type AS ENUM ('share', 'bond', 'fund');

CREATE TABLE IF NOT EXISTS securities
(
    id          UUID        NOT NULL,
    sec_id      VARCHAR     NOT NULL UNIQUE,
    ticker      VARCHAR     NOT NULL,
    short_name  VARCHAR     NOT NULL,
    type        security_type NOT NULL,
    extra       JSONB,                        -- для специфических полей
    create_time TIMESTAMPTZ NOT NULL DEFAULT now()

    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS securities_sec_id_idx ON securities (sec_id);
CREATE INDEX IF NOT EXISTS securities_ticker_idx ON securities (ticker);
CREATE INDEX IF NOT EXISTS securities_type_idx ON securities (type);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS securities;
DROP TYPE IF EXISTS security_type;

-- +goose StatementEnd
