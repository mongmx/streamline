
-- +migrate Up
CREATE TABLE streamline.users (
	id BIGSERIAL PRIMARY KEY,
	uuid UUID DEFAULT uuid_generate_v4(),
    email VARCHAR NOT NULL DEFAULT '',
	is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    plan_id BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

-- +migrate Down
DROP TABLE streamline.users;
