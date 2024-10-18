CREATE TABLE IF NOT EXISTS company_type (
    id SERIAL PRIMARY KEY,            
    cmp_name VARCHAR(255) NOT NULL,
    navi_user VARCHAR(255) NOT NULL
);