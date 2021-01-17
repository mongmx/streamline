
-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SCHEMA IF NOT EXISTS streamline;

-- +migrate Down
