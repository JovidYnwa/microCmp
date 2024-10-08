CREATE TABLE IF NOT EXISTS company (
    id SERIAL PRIMARY KEY,            -- Auto-incrementing ID
    cmp_name VARCHAR(255) NOT NULL,   -- Company name, required field
    navi_user VARCHAR(255) NOT NULL,  -- Who created
    query_id VARCHAR(255) NOT NULL,   -- Request ID from billing
    start_time TIMESTAMP,             -- Start time
    duration INTEGER,                 -- Duration
    repetition INTEGER,               -- Repetition count
);
