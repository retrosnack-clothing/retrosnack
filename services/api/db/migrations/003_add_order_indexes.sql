-- +goose Up
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_checkout_session_id ON orders(checkout_session_id);
CREATE INDEX idx_orders_status ON orders(status);

-- +goose Down
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_checkout_session_id;
DROP INDEX IF EXISTS idx_orders_user_id;
