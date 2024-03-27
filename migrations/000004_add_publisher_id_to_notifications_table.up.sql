-- 000004_add_publisher_id_to_notifications_table.up.sql
ALTER TABLE notifications
ADD COLUMN publisher_id INTEGER;

ALTER TABLE notifications
ADD CONSTRAINT fk_publisher
FOREIGN KEY (publisher_id)
REFERENCES users(id);