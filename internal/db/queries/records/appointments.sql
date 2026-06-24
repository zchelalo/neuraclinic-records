-- name: CreateAppointment :one
INSERT INTO appointments (
  id, start_time, end_time, reason, status, patient_id, cancelled_by_user_id,
  rescheduled_from_appointment_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $9
)
RETURNING *;

-- name: GetAppointmentByID :one
SELECT a.*
FROM appointments a
JOIN patients p ON p.id = a.patient_id
WHERE a.id = $1
  AND p.psychologist_id = $2
  AND p.deleted_at IS NULL;

-- name: ListAppointments :many
SELECT a.*
FROM appointments a
JOIN patients p ON p.id = a.patient_id
WHERE p.psychologist_id = $1
  AND p.deleted_at IS NULL
  AND (sqlc.narg(patient_id)::uuid IS NULL OR a.patient_id = sqlc.narg(patient_id)::uuid)
  AND (sqlc.narg(start_date)::timestamptz IS NULL OR a.start_time >= sqlc.narg(start_date)::timestamptz)
  AND (sqlc.narg(end_date)::timestamptz IS NULL OR a.start_time <= sqlc.narg(end_date)::timestamptz)
  AND (
    cardinality(sqlc.arg(statuses)::text[]) = 0
    OR a.status = ANY(sqlc.arg(statuses)::text[])
  )
  AND (
    sqlc.narg(after_id)::uuid IS NULL
    OR a.id < sqlc.narg(after_id)::uuid
  )
  AND (
    sqlc.narg(before_id)::uuid IS NULL
    OR a.id > sqlc.narg(before_id)::uuid
  )
ORDER BY
  CASE WHEN sqlc.arg(is_backward)::bool THEN a.id END ASC,
  CASE WHEN NOT sqlc.arg(is_backward)::bool THEN a.id END DESC
LIMIT sqlc.arg(limit_count)::int;

-- name: UpdateAppointmentStatus :one
UPDATE appointments a
SET
  status = $3,
  cancelled_by_user_id = $4,
  updated_at = $5
FROM patients p
WHERE p.id = a.patient_id
  AND a.id = $1
  AND p.psychologist_id = $2
  AND p.deleted_at IS NULL
RETURNING a.*;

-- name: AppointmentBelongsToPatient :one
SELECT EXISTS (
  SELECT 1
  FROM appointments a
  JOIN patients p ON p.id = a.patient_id
  WHERE a.id = $1
    AND a.patient_id = $2
    AND p.psychologist_id = $3
    AND p.deleted_at IS NULL
) AS exists;
