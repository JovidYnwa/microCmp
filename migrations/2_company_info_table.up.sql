CREATE TABLE IF NOT EXISTS COMPANY (
    id SERIAL PRIMARY KEY,
    company_type_id INTEGER NOT NULL,
    cmp_desc JSONB, --name e.t.c
    cmp_filter JSONB,
    sms_data JSONB,
    action_data JSONB,
    FOREIGN KEY (company_id) REFERENCES company_type(id) ON DELETE CASCADE
);

    -- query_id VARCHAR(255) NOT NULL,
