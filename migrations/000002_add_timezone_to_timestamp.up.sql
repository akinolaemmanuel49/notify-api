-- 000002_add_timezone_to_timestamp.up.sql
ALTER TABLE notifications
ALTER COLUMN created_at SET DATA TYPE TIMESTAMP WITH TIME ZONE;

ALTER TABLE notifications
ALTER COLUMN updated_at SET DATA TYPE TIMESTAMP WITH TIME ZONE;
