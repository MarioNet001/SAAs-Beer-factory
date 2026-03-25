CREATE TABLE inventory (
    product_id UUID PRIMARY KEY,
    stock_level INT NOT NULL DEFAULT 0,
    minimum_threshold INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
