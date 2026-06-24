-- name: CreateAttachment :one
INSERT INTO attachments (
  id, file_id, mime_type, patient_id, note_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $6
)
RETURNING *;

-- name: GetAttachmentByID :one
SELECT att.*
FROM attachments att
JOIN patients p ON p.id = att.patient_id
WHERE att.id = $1
  AND p.psychologist_id = $2
  AND att.deleted_at IS NULL
  AND p.deleted_at IS NULL;

-- name: ListAttachments :many
SELECT att.*
FROM attachments att
JOIN patients p ON p.id = att.patient_id
WHERE p.psychologist_id = $1
  AND att.patient_id = $2
  AND (sqlc.narg(note_id)::uuid IS NULL OR att.note_id = sqlc.narg(note_id)::uuid)
  AND att.deleted_at IS NULL
  AND p.deleted_at IS NULL
  AND (
    sqlc.narg(after_id)::uuid IS NULL
    OR att.id < sqlc.narg(after_id)::uuid
  )
  AND (
    sqlc.narg(before_id)::uuid IS NULL
    OR att.id > sqlc.narg(before_id)::uuid
  )
ORDER BY
  CASE WHEN sqlc.arg(is_backward)::bool THEN att.id END ASC,
  CASE WHEN NOT sqlc.arg(is_backward)::bool THEN att.id END DESC
LIMIT sqlc.arg(limit_count)::int;

-- name: DeleteAttachment :execrows
UPDATE attachments att
SET deleted_at = $3, updated_at = $3
FROM patients p
WHERE p.id = att.patient_id
  AND att.id = $1
  AND p.psychologist_id = $2
  AND att.deleted_at IS NULL
  AND p.deleted_at IS NULL;
