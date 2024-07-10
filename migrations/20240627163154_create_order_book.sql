-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_book
(
    id       UUID DEFAULT generateUUIDv4(),
    exchange String,
    pair     String,
    asks     Tuple(Float64, Float64),
    bids     Tuple(Float64, Float64)
)
    ENGINE = MergeTree
    ORDER BY (exchange, pair);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_book;
-- +goose StatementEnd
