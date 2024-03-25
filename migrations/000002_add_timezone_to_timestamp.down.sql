-- 000002_add_timezone_to_timestamp.down.sql
ALTER TABLE notifications
ALTER COLUMN created_at SET DATA TYPE TIMESTAMP;

ALTER TABLE notifications
ALTER COLUMN updated_at SET DATA TYPE TIMESTAMP;