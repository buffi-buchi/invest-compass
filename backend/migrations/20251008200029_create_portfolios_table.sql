-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS portfolios
(
    id          UUID        NOT NULL,
    user_id     UUID        NOT NULL,
    name        VARCHAR     NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS portfolios;

-- +goose StatementEnd
