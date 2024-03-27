-- 000005_add_emails_are_unique_constraint.up.sql
ALTER TABLE users
ADD CONSTRAINT uc_email
UNIQUE (email)