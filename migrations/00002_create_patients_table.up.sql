CREATE TABLE IF NOT EXISTS patients
(
    id           VARCHAR(16) NOT NULL PRIMARY KEY,
    phone_number VARCHAR(15) NOT NULL,
    name         VARCHAR(30) NOT NULL,
    birth_date   DATE        NOT NULL,
    is_male      BOOLEAN     NOT NULL,
    img_url      TEXT        NOT NULL,
    created_at   TIMESTAMP   NOT NULL
);

CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE INDEX IF NOT EXISTS idx_patients_name ON patients USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_patients_phone_number ON patients (phone_number varchar_pattern_ops);
CREATE INDEX IF NOT EXISTS idx_patients_created_at_asc ON patients (created_at ASC);
CREATE INDEX IF NOT EXISTS idx_patients_created_at_desc ON patients (created_at DESC);