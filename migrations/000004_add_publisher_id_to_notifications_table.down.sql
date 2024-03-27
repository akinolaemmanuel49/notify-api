-- 000004_add_publisher_id_to_notifications_table.down.sql
ALTER TABLE notifications
DROP COLUMN publisher_id;