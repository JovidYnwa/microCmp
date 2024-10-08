CREATE TABLE IF NOT EXISTS company (
    id SERIAL PRIMARY KEY,          -- Auto-incrementing ID
    name VARCHAR(255) NOT NULL,     -- Company name, required field
    user VARCHAR(255) NOT NULL,     -- who created
    query_id VARCHAR(255) NOT NULL, --request id from billing
    start_time TIMESTAMP,     
    duration INTEGER,
    repition INTEGER,
    company_launched INTEGER,        -- Country ID as a foreign key reference (optional)
    subscriber_count INTEGER CHECK (subscriber_count >= 0),  -- Non-negative subscribers count
    efficiency NUMERIC(5, 2)        -- Precision for efficiency (example: 99.99)   
);
