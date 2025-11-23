BEGIN;

DROP INDEX idx_pr_status_merged_at;
DROP INDEX idx_pr_merged_at;

ALTER TABLE pull_requests
DROP COLUMN merged_at;

ALTER TABLE reviewers
ALTER COLUMN created_at DROP NOT NULL;

ALTER TABLE pull_requests
ALTER COLUMN created_at DROP NOT NULL,
ALTER COLUMN updated_at DROP NOT NULL;

ALTER TABLE user_teams
ALTER COLUMN created_at DROP NOT NULL;

ALTER TABLE teams
ALTER COLUMN created_at DROP NOT NULL,
ALTER COLUMN updated_at DROP NOT NULL;

ALTER TABLE users
ALTER COLUMN created_at DROP NOT NULL,
ALTER COLUMN updated_at DROP NOT NULL;

COMMIT;
