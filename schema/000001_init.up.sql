CREATE TABLE teams (
    team_name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE users (
    user_id     VARCHAR(255) PRIMARY KEY,
    username    VARCHAR(255) NOT NULL,
    team_name   VARCHAR(255) NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE pull_requests (
    pull_request_id   VARCHAR(255) PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id         VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    status            VARCHAR(16) NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at         TIMESTAMPTZ NULL
);

CREATE TABLE pr_reviewers (
    pull_request_id VARCHAR(255) NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id     VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,

    PRIMARY KEY (pull_request_id, reviewer_id)
);
