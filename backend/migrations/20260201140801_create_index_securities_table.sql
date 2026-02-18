-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS index_securities
(
    id          UUID        NOT NULL,
    index_id    UUID         NOT NULL,
    security_id UUID         NOT NULL,
    weight      NUMERIC(4, 2) NOT NULL,     -- доля бумаги в индексе (например 5.5%)
    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (index_id) REFERENCES indexes (id) ON DELETE CASCADE,
    FOREIGN KEY (security_id) REFERENCES securities (id) ON DELETE CASCADE,
    UNIQUE (index_id, security_id)                    -- одна бумага не может быть дважды в индексе
);

CREATE INDEX IF NOT EXISTS index_securities_index_id_idx ON index_securities (index_id);
CREATE INDEX IF NOT EXISTS index_securities_security_id_idx ON index_securities (security_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS index_securities;

-- +goose StatementEnd
