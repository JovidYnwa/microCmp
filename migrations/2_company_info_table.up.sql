CREATE TABLE IF NOT EXISTS company (
    id SERIAL PRIMARY KEY,
    company_type_id INTEGER NOT NULL,
    cmp_desc JSONB,  -- company description and other details
    cmp_filter JSONB, -- company filters like phoneType, balanceLimits
    sms_data JSONB, -- SMS related data
    action_data JSONB, -- action-related data
    FOREIGN KEY (company_type_id) REFERENCES company_type(id) ON DELETE CASCADE
);