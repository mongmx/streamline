
-- +migrate Up
CREATE TABLE streamline.auths (
    user_id BIGINT NOT NULL DEFAULT 0,
    type VARCHAR NOT NULL DEFAULT '',
    secret TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    FOREIGN KEY (user_id) REFERENCES streamline.users (id)
);

-- +migrate Down
DROP TABLE streamline.auths;
