-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_history
(
    client_name           String,
    exchange_name         String,
    label                 String,
    pair                  String,
    side                  String,
    type                  String,
    base_qty              Float64,
    price                 Float64,
    algorithm_name_placed String,
    lowest_sell_prc       Float64,
    highest_buy_prc       Float64,
    commission_quote_qty  Float64,
    time_placed           TIMESTAMP DEFAULT now()
)
    ENGINE = MergeTree
    ORDER BY time_placed;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_history;
-- +goose StatementEnd
