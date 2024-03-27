-- 000005_add_emails_are_unique_constraint.down.sql
ALTER TABLE users
DROP CONSTRAINT uc_email