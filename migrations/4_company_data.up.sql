CREATE TABLE IF NOT EXISTS company_data (
    msisdn VARCHAR(255) NOT NULL,
    company_id INTEGER NOT NULL,
    action_commited BOOLEAN DEFAULT FALSE,-- did committed action
    notification_amount NUMERIC(10, 2),   -- reconsider precision for this column
    last_update TIMESTAMP,                -- fixed typo TIMESTAP -> TIMESTAMP
    FOREIGN KEY(company_id) REFERENCES company(id) ON DELETE CASCADE
);