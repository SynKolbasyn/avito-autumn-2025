BEGIN;

DROP INDEX idx_reviewers_user_id;
DROP INDEX idx_reviewers_pr_id;

DROP INDEX idx_pr_author_status;
DROP INDEX idx_pr_created_at;
DROP INDEX idx_pr_author;
DROP INDEX idx_pr_status;

DROP INDEX idx_user_teams_team_user;
DROP INDEX idx_user_teams_user_id;
DROP INDEX idx_user_teams_team_id;

DROP INDEX idx_teams_name;

DROP INDEX idx_users_name;
DROP INDEX idx_users_is_active;

COMMIT;
