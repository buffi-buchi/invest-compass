- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS securities
(
    ticker      VARCHAR, -- Тикер
    short_name  VARCHAR     NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (ticker)
    );


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS securities;

-- +goose StatementEnd
