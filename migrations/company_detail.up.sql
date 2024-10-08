CREATE TABLE IF NOT EXISTS company_detail (
    company_id INTEGER PRIMARY KEY,
    efficiency DOUBLE PRECISION,
    sub_amount NUMERIC(7),
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE  -- Cascade delete
);