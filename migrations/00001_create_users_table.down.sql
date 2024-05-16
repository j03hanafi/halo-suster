DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "pg_trgm";

DROP INDEX IF EXISTS idx_users_created_at_asc;
DROP INDEX IF EXISTS idx_users_created_at_desc;
DROP INDEX IF EXISTS idx_users_name;
DROP INDEX IF EXISTS idx_users_nip;
DROP INDEX IF EXISTS idx_users_is_it;