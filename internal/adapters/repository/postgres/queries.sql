-- name: CreateDeveloperApplication :one
INSERT INTO developer_applications (
    developer_id, code_name, category_id, android_package_name, default_locale,
    web_video, inapp_video, web_background_image, inapp_background_image,
    name, short_title, description, goals, tasks, results, challenges, location,
    video_cover, safety, apk_filename, version
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
RETURNING id, developer_id, code_name, category_id, android_package_name, default_locale, web_video, inapp_video, web_background_image, inapp_background_image, name, short_title, description, goals, tasks, results, challenges, location, video_cover, safety, apk_filename, version, verification_process_id, verification_status, verification_failed_step, published, created_at, updated_at;

-- name: GetDeveloperApplication :one
SELECT id, developer_id, code_name, category_id, android_package_name, default_locale, web_video, inapp_video, web_background_image, inapp_background_image, name, short_title, description, goals, tasks, results, challenges, location, video_cover, safety, apk_filename, version, verification_process_id, verification_status, verification_failed_step, published, created_at, updated_at
FROM developer_applications
WHERE id = $1;

-- name: ListDeveloperApplications :many
SELECT id, developer_id, code_name, category_id, android_package_name, default_locale, web_video, inapp_video, web_background_image, inapp_background_image, name, short_title, description, goals, tasks, results, challenges, location, video_cover, safety, apk_filename, version, verification_process_id, verification_status, verification_failed_step, published, created_at, updated_at
FROM developer_applications
WHERE developer_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SetVerificationProcess :one
UPDATE developer_applications
SET verification_process_id = $2,
    verification_status     = 'pending',
    updated_at              = NOW()
WHERE id = $1
RETURNING id, developer_id, code_name, category_id, android_package_name, default_locale, web_video, inapp_video, web_background_image, inapp_background_image, name, short_title, description, goals, tasks, results, challenges, location, video_cover, safety, apk_filename, version, verification_process_id, verification_status, verification_failed_step, published, created_at, updated_at;

-- name: UpdateVerificationStatus :exec
UPDATE developer_applications
SET verification_status = $2,
    verification_failed_step = $3,
    updated_at          = NOW()
WHERE verification_process_id = $1;

-- name: MarkAsPublished :exec
UPDATE developer_applications
SET published  = TRUE,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateDeveloperApplication :one
UPDATE developer_applications
SET code_name              = $2,
    category_id            = $3,
    android_package_name   = $4,
    default_locale         = $5,
    web_video              = $6,
    inapp_video            = $7,
    web_background_image   = $8,
    inapp_background_image = $9,
    name                   = $10,
    short_title            = $11,
    description            = $12,
    goals                  = $13,
    tasks                  = $14,
    results                = $15,
    challenges             = $16,
    location               = $17,
    video_cover            = $18,
    safety                 = $19,
    version                = $20,
    updated_at             = NOW()
WHERE id = $1
RETURNING id, developer_id, code_name, category_id, android_package_name, default_locale, web_video, inapp_video, web_background_image, inapp_background_image, name, short_title, description, goals, tasks, results, challenges, location, video_cover, safety, apk_filename, version, verification_process_id, verification_status, verification_failed_step, published, created_at, updated_at;
