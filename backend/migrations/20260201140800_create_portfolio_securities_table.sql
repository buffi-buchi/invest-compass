-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS portfolio_securities
(
    id              UUID            NOT NULL,
    portfolio_id    UUID            NOT NULL,
    security_id     UUID            NOT NULL,
    amount                          INTEGER ,          -- количество бумаг
    create_time     TIMESTAMPTZ     NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (portfolio_id) REFERENCES portfolios (id) ON DELETE CASCADE,
    FOREIGN KEY (security_id) REFERENCES securities (id) ON DELETE CASCADE,
    UNIQUE (portfolio_id, security_id)            -- одна бумага не может быть дважды в портфеле
);

CREATE INDEX IF NOT EXISTS portfolio_securities_portfolio_id_idx ON portfolio_securities (portfolio_id);
CREATE INDEX IF NOT EXISTS portfolio_securities_security_id_idx ON portfolio_securities (security_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS portfolio_securities;

-- +goose StatementEnd
