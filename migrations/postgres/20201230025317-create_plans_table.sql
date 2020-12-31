
-- +migrate Up
CREATE TABLE streamline.plans (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    price NUMERIC NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);
INSERT INTO streamline.plans (title, description)
VALUES ('free', '1 topic'), ('silver', '10 topic'), ('gold', 'unlimited topic');

-- +migrate Down
DROP TABLE streamline.plans;
