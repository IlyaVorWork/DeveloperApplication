-- name: CreateDeveloperApplication :one
INSERT INTO developer_applications (developer_id, apk_filename, version)
VALUES ($1, $2, $3)
RETURNING id, developer_id, apk_filename, version, verification_process_id, verification_status, published, created_at, updated_at;

-- name: GetDeveloperApplication :one
SELECT id, developer_id, apk_filename, version, verification_process_id, verification_status, published, created_at, updated_at
FROM developer_applications
WHERE id = $1;

-- name: ListDeveloperApplications :many
SELECT id, developer_id, apk_filename, version, verification_process_id, verification_status, published, created_at, updated_at
FROM developer_applications
WHERE developer_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SetVerificationProcess :one
UPDATE developer_applications
SET verification_process_id = $2,
    verification_status     = 'started',
    updated_at              = NOW()
WHERE id = $1
RETURNING id, developer_id, apk_filename, version, verification_process_id, verification_status, published, created_at, updated_at;

-- name: UpdateVerificationStatus :exec
UPDATE developer_applications
SET verification_status = $2,
    updated_at          = NOW()
WHERE verification_process_id = $1;

-- name: MarkAsPublished :exec
UPDATE developer_applications
SET published  = TRUE,
    updated_at = NOW()
WHERE id = $1;
