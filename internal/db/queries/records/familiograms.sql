-- name: CreateFamiliogram :exec
INSERT INTO familiograms (
  id, data, patient_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $4
);

-- name: GetFamiliogramByPatientID :one
SELECT f.id, f.data, f.patient_id, f.created_at, f.updated_at
FROM familiograms f
JOIN patients p ON p.id = f.patient_id
WHERE f.patient_id = $1
  AND p.psychologist_id = $2
  AND p.deleted_at IS NULL;

-- name: UpdateFamiliogram :one
UPDATE familiograms f
SET data = $3, updated_at = $4
FROM patients p
WHERE p.id = f.patient_id
  AND f.id = $1
  AND p.psychologist_id = $2
  AND p.deleted_at IS NULL
RETURNING f.id, f.data, f.patient_id, f.created_at, f.updated_at;
