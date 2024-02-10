-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets
(
    id      varchar(36) PRIMARY KEY,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 100.0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wallets;
-- +goose StatementEnd
