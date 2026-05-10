CREATE TABLE developer_applications
(
    id                      UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id            UUID      NOT NULL,
    apk_filename            VARCHAR   NOT NULL,
    version                 VARCHAR   NOT NULL,
    verification_process_id UUID,
    verification_status     VARCHAR,
    published               BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMP NOT NULL DEFAULT NOW()
);
