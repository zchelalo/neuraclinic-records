-- name: CreateNote :one
INSERT INTO notes (
  id, patient_id, appointment_id, title, content_html, content_text, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $7
)
RETURNING *;

-- name: GetNoteByID :one
SELECT n.*
FROM notes n
JOIN patients p ON p.id = n.patient_id
WHERE n.id = $1
  AND p.psychologist_id = $2
  AND n.deleted_at IS NULL
  AND p.deleted_at IS NULL;

-- name: ListNotes :many
SELECT
  n.id,
  n.patient_id,
  n.appointment_id,
  n.title,
  n.created_at,
  n.updated_at,
  n.deleted_at
FROM notes n
JOIN patients p ON p.id = n.patient_id
WHERE p.psychologist_id = $1
  AND n.patient_id = $2
  AND n.deleted_at IS NULL
  AND p.deleted_at IS NULL
  AND (sqlc.narg(start_date)::timestamptz IS NULL OR n.created_at >= sqlc.narg(start_date)::timestamptz)
  AND (sqlc.narg(end_date)::timestamptz IS NULL OR n.created_at <= sqlc.narg(end_date)::timestamptz)
  AND (
    NOT sqlc.arg(with_appointment_associated)::bool
    OR n.appointment_id IS NOT NULL
  )
  AND (
    NOT sqlc.arg(with_files_associated)::bool
    OR EXISTS (
      SELECT 1
      FROM attachments att
      WHERE att.note_id = n.id
        AND att.deleted_at IS NULL
    )
  )
  AND (
    sqlc.arg(search_query)::text = ''
    OR COALESCE(n.title, '') ILIKE '%' || sqlc.arg(search_query)::text || '%'
    OR n.content_text ILIKE '%' || sqlc.arg(search_query)::text || '%'
  )
  AND (
    sqlc.narg(after_id)::uuid IS NULL
    OR n.id < sqlc.narg(after_id)::uuid
  )
  AND (
    sqlc.narg(before_id)::uuid IS NULL
    OR n.id > sqlc.narg(before_id)::uuid
  )
ORDER BY
  CASE WHEN sqlc.arg(is_backward)::bool THEN n.id END ASC,
  CASE WHEN NOT sqlc.arg(is_backward)::bool THEN n.id END DESC
LIMIT sqlc.arg(limit_count)::int;

-- name: UpdateNote :execrows
UPDATE notes n
SET
  appointment_id = COALESCE(sqlc.narg(appointment_id), n.appointment_id),
  title = COALESCE(sqlc.narg(title), n.title),
  content_html = COALESCE(sqlc.narg(content_html), n.content_html),
  content_text = COALESCE(sqlc.narg(content_text), n.content_text),
  updated_at = sqlc.arg(updated_at)
FROM patients p
WHERE p.id = n.patient_id
  AND n.id = sqlc.arg(id)
  AND p.psychologist_id = sqlc.arg(psychologist_id)
  AND n.deleted_at IS NULL
  AND p.deleted_at IS NULL;

-- name: DeleteNote :execrows
UPDATE notes n
SET deleted_at = $3, updated_at = $3
FROM patients p
WHERE p.id = n.patient_id
  AND n.id = $1
  AND p.psychologist_id = $2
  AND n.deleted_at IS NULL
  AND p.deleted_at IS NULL;

-- name: NoteBelongsToPatient :one
SELECT EXISTS (
  SELECT 1
  FROM notes n
  JOIN patients p ON p.id = n.patient_id
  WHERE n.id = $1
    AND n.patient_id = $2
    AND p.psychologist_id = $3
    AND n.deleted_at IS NULL
    AND p.deleted_at IS NULL
) AS exists;
