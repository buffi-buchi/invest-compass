-- +goose Up
-- +goose StatementBegin

SELECT version();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

SELECT version();

-- +goose StatementEnd
