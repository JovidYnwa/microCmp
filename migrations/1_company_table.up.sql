CREATE TABLE IF NOT EXISTS company (
    id SERIAL PRIMARY KEY,            
    cmp_name VARCHAR(255) NOT NULL,
    cmp_description TEXT,
    navi_user VARCHAR(255) NOT NULL,  
    query_id VARCHAR(255) NOT NULL,
    start_time TIMESTAMP,             
    duration INTEGER,                 
    repetition INTEGER
);