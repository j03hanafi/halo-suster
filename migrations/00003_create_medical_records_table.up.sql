CREATE TABLE IF NOT EXISTS medical_records
(
    id                   bytea         NOT NULL PRIMARY KEY,
    patient_id           VARCHAR(16)   NOT NULL,
    patient_phone_number VARCHAR(15)   NOT NULL,
    patient_name         VARCHAR(30)   NOT NULL,
    patient_birth_date   DATE          NOT NULL,
    patient_is_male      BOOLEAN       NOT NULL,
    patient_img_url      TEXT          NOT NULL,
    symptoms             VARCHAR(2000) NOT NULL,
    medications          VARCHAR(2000) NOT NULL,
    staff_id             bytea         NOT NULL,
    staff_nip            varchar(13)   NOT NULL,
    staff_name           VARCHAR(30)   NOT NULL,
    created_at           timestamp     NOT NULL
);

-- patient_id, staff_id, staff_nip, created_at
CREATE INDEX IF NOT EXISTS idx_medical_records_patient_id ON medical_records USING hash (patient_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_staff_id ON medical_records USING hash (staff_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_staff_nip ON medical_records USING hash (staff_nip);
CREATE INDEX IF NOT EXISTS idx_medical_records_created_at_asc ON medical_records (created_at ASC);
CREATE INDEX IF NOT EXISTS idx_medical_records_created_at_desc ON medical_records (created_at DESC);