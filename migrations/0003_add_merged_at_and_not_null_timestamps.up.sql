BEGIN;

UPDATE users
SET created_at = COALESCE(created_at, CURRENT_TIMESTAMP),
    updated_at = COALESCE(updated_at, CURRENT_TIMESTAMP)
WHERE created_at IS NULL OR updated_at IS NULL;

UPDATE teams
SET created_at = COALESCE(created_at, CURRENT_TIMESTAMP),
    updated_at = COALESCE(updated_at, CURRENT_TIMESTAMP)
WHERE created_at IS NULL OR updated_at IS NULL;

UPDATE user_teams
SET created_at = COALESCE(created_at, CURRENT_TIMESTAMP)
WHERE created_at IS NULL;

UPDATE pull_requests
SET created_at = COALESCE(created_at, CURRENT_TIMESTAMP),
    updated_at = COALESCE(updated_at, CURRENT_TIMESTAMP)
WHERE created_at IS NULL OR updated_at IS NULL;

UPDATE reviewers
SET created_at = COALESCE(created_at, CURRENT_TIMESTAMP)
WHERE created_at IS NULL;

ALTER TABLE users
ALTER COLUMN created_at SET NOT NULL,
ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE teams
ALTER COLUMN created_at SET NOT NULL,
ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE user_teams
ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE pull_requests
ALTER COLUMN created_at SET NOT NULL,
ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE reviewers
ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE pull_requests
ADD COLUMN merged_at TIMESTAMPTZ DEFAULT NULL;

UPDATE pull_requests
SET merged_at = updated_at
WHERE status = 'MERGED' AND merged_at IS NULL;

CREATE INDEX idx_pr_merged_at ON pull_requests (merged_at);
CREATE INDEX idx_pr_status_merged_at ON pull_requests (status, merged_at);

COMMIT;
