-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS profiles
(
    user_id UUID    NOT NULL,
    ticker  VARCHAR NOT NULL,

    PRIMARY KEY (user_id, ticker),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS profiles;

-- +goose StatementEnd
