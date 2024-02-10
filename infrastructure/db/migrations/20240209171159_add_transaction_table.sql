-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions
(
    time           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    from_wallet_id varchar(36) REFERENCES wallets (id),
    to_wallet_id   varchar(36) REFERENCES wallets (id),
    amount         DECIMAL(10, 2) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd
