CREATE TABLE developer_applications
(
    id                      UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id            UUID      NOT NULL,

    -- → application_service: applications
    code_name               VARCHAR   NOT NULL,
    category_id             BIGINT    NOT NULL DEFAULT 0,
    android_package_name    VARCHAR   NOT NULL,
    default_locale          VARCHAR   NOT NULL,
    web_video               VARCHAR   NOT NULL DEFAULT '',
    inapp_video             VARCHAR   NOT NULL DEFAULT '',
    web_background_image    VARCHAR   NOT NULL DEFAULT '',
    inapp_background_image  VARCHAR,

    -- → application_service: application_translations
    name                    VARCHAR   NOT NULL,
    short_title             VARCHAR   NOT NULL,
    description             TEXT,
    goals                   VARCHAR   NOT NULL DEFAULT '',
    tasks                   VARCHAR   NOT NULL DEFAULT '',
    results                 TEXT,
    challenges              VARCHAR,
    location                VARCHAR,

    -- developer portal-specific
    video_cover             VARCHAR,
    safety                  TEXT,

    apk_filename            VARCHAR   NOT NULL,
    version                 VARCHAR   NOT NULL,
    verification_process_id  UUID,
    verification_status      VARCHAR,
    verification_failed_step VARCHAR,
    published                BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMP NOT NULL DEFAULT NOW()
);
