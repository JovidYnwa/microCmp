CREATE TABLE IF NOT EXISTS COMPANY_INFO (
    company_id INTEGER PRIMARY KEY,
    FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE,
    cmp_filter JSONB,
    sms_data JSONB,
    action_data JSONB
);
