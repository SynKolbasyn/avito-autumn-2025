BEGIN;

CREATE INDEX idx_users_is_active ON users (is_active);
CREATE INDEX idx_users_name ON users (name);

CREATE INDEX idx_teams_name ON teams (name);

CREATE INDEX idx_user_teams_team_id ON user_teams (team_id);
CREATE INDEX idx_user_teams_user_id ON user_teams (user_id);
CREATE INDEX idx_user_teams_team_user ON user_teams (team_id, user_id);

CREATE INDEX idx_pr_status ON pull_requests (status);
CREATE INDEX idx_pr_author ON pull_requests (author_id);
CREATE INDEX idx_pr_created_at ON pull_requests (created_at);
CREATE INDEX idx_pr_author_status ON pull_requests (author_id, status);

CREATE INDEX idx_reviewers_pr_id ON reviewers (pr_id);
CREATE INDEX idx_reviewers_user_id ON reviewers (user_id);

COMMIT;
