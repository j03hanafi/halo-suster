CREATE TABLE IF NOT EXISTS users
(
    id         bytea        NOT NULL PRIMARY KEY,
    nip        varchar(13)  NOT NULL UNIQUE,
    name       varchar(50)  NOT NULL,
    password   varchar(255) NOT NULL,
    is_it      BOOLEAN      NOT NULL,
    img_url    text         NOT NULL,
    created_at timestamp    NOT NULL
);

CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE INDEX IF NOT EXISTS idx_users_name ON users USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_users_nip ON users (nip varchar_pattern_ops);
CREATE INDEX IF NOT EXISTS idx_users_is_it ON users USING hash (is_it);
CREATE INDEX IF NOT EXISTS idx_users_created_at_asc ON users (created_at ASC);
CREATE INDEX IF NOT EXISTS idx_users_created_at_desc ON users (created_at DESC);
