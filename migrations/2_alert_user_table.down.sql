-- Migration Down: Remove email column from TEST_TABLE
ALTER TABLE TEST_TABLE
DROP COLUMN email;