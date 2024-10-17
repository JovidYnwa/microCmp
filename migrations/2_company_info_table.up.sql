CREATE TABLE IF NOT EXISTS COMPANY_INFO (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    cmp_filter JSONB,
    sms_data JSONB,
    action_data JSONB,
    FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE
);
