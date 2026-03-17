-- Schema for Beer Management
CREATE TABLE raw_materials (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock DECIMAL(10,2) NOT NULL,
    unit VARCHAR(50) NOT NULL
);

CREATE TABLE batches (
    id UUID PRIMARY KEY,
    batch_number VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    details TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
