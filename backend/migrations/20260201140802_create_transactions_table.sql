-- +goose Up
-- +goose StatementBegin

CREATE TYPE transaction_type AS ENUM ('buy', 'sell');

CREATE TABLE IF NOT EXISTS transactions
(
    id           UUID       NOT NULL,
    portfolio_id UUID       NOT NULL,
    security_id INT         NOT NULL,
    amount      INTEGER,
    price       NUMERIC(18, 4) NOT NULL,          -- цена за единицу
    trade_date  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    type        transaction_type NOT NULL,
    note        TEXT,

    PRIMARY KEY (id),
    FOREIGN KEY (portfolio_id) REFERENCES portfolios (id) ON DELETE CASCADE,
    FOREIGN KEY (security_id) REFERENCES securities (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS transactions_portfolio_id_idx ON transactions (portfolio_id);
CREATE INDEX IF NOT EXISTS transactions_security_id_idx ON transactions (security_id);
CREATE INDEX IF NOT EXISTS transactions_trade_date_idx ON transactions (trade_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS transactions;
DROP TYPE IF EXISTS transaction_type;

-- +goose StatementEnd
