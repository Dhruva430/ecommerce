CREATE INDEX idx_cart_user ON cart (user_id);

CREATE INDEX idx_cart_product ON cart (product_id);

CREATE INDEX idx_order_user ON orders (user_id);

CREATE INDEX idx_order_seller ON orders (seller_id);

CREATE INDEX idx_product_seller ON product (seller_id);

CREATE INDEX idx_refresh_token_user ON refresh_token (user_id);