CREATE TABLE IF NOT EXISTS users
(
    id         bytea PRIMARY KEY,
    nip        varchar(13)  NOT NULL UNIQUE,
    name       varchar(50)  NOT NULL,
    password   varchar(255) NOT NULL,
    is_it      BOOLEAN      NOT NULL,
    img_url    text         NOT NULL,
    created_at timestamp    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_nip ON users (nip);
CREATE INDEX IF NOT EXISTS idx_users_is_it ON users USING HASH (is_it);