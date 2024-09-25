-- Migration Up: Add email column to TEST_TABLE
ALTER TABLE TEST_TABLE
ADD COLUMN email VARCHAR(255);