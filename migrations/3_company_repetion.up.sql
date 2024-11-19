CREATE TABLE IF NOT EXISTS company_repetion (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    efficiency DOUBLE PRECISION,
    sub_amount NUMERIC(7),
    start_date TIMESTAMP,
    --end_date TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE  -- Changed to reference company_type
);