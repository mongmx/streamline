
-- +migrate Up
CREATE TABLE streamline.topics (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    user_id BIGINT NOT NULL DEFAULT 0,
    title VARCHAR NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    FOREIGN KEY (user_id) REFERENCES streamline.users (id)
);

-- +migrate Down
DROP TABLE streamline.topics;
