CREATE TABLE IF NOT EXISTS company_repetion (
    company_id INTEGER PRIMARY KEY,
    FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE,
    efficiency DOUBLE PRECISION,
    sub_amount NUMERIC(7),
    start_date TIMESTAMP,
    end_date TIMESTAMP
);
