
CREATE TABLE IF NOT EXISTS company_info (
    company_id INTEGER PRIMARY KEY,
    FOREIGN KEY (company_id) REFERENCES company(id),
    trpl_type VARCHAR(255),
    trpl_type INTEGER,
    trpl_name VARCHAR(255),
    trpl_id INTEGER,
    balance_begin NUMERIC(5,2),
    balance_end NUMERIC(5,2),
    subs_status_id INTEGER,
    subs_status_name VARCHAR(255),
    subs_device_id INTEGER,
    subs_device_name VARCHAR(255),
    region VARCHAR(255), 
    sms_tj VARCHAR(255),
    sms_ru VARCHAR(255),
    sms_eng VARCHAR(255),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);