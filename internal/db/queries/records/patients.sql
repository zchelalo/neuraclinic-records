-- name: CreateAddress :exec
INSERT INTO addresses (
  id, country, province, city, postal_code, neighborhood, street, street_number,
  unit_number, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10
);

-- name: CreatePatient :exec
INSERT INTO patients (
  id, first_name, middle_name, first_last_name, second_last_name, birth_date,
  birth_country, birth_state, birth_city, sex, marital_status, occupation,
  religion, phone, email, address_id, psychologist_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $18
);

-- name: GetPatientByID :one
SELECT
  p.id,
  p.first_name,
  p.middle_name,
  p.first_last_name,
  p.second_last_name,
  p.birth_date,
  p.birth_country,
  p.birth_state,
  p.birth_city,
  p.sex,
  p.marital_status,
  p.occupation,
  p.religion,
  p.phone,
  p.email,
  p.psychologist_id,
  p.created_at,
  p.updated_at,
  p.deleted_at,
  a.id AS address_id,
  a.country AS address_country,
  a.province AS address_province,
  a.city AS address_city,
  a.postal_code AS address_postal_code,
  a.neighborhood AS address_neighborhood,
  a.street AS address_street,
  a.street_number AS address_street_number,
  a.unit_number AS address_unit_number,
  a.created_at AS address_created_at,
  a.updated_at AS address_updated_at,
  a.deleted_at AS address_deleted_at
FROM patients p
JOIN addresses a ON a.id = p.address_id
WHERE p.id = $1
  AND p.psychologist_id = $2
  AND p.deleted_at IS NULL;

-- name: ListPatients :many
SELECT
  p.id,
  p.first_name,
  p.middle_name,
  p.first_last_name,
  p.second_last_name,
  p.birth_date,
  p.email,
  p.phone
FROM patients p
WHERE p.psychologist_id = $1
  AND p.deleted_at IS NULL
  AND (
    sqlc.arg(search_query)::text = ''
    OR p.first_name ILIKE '%' || sqlc.arg(search_query)::text || '%'
    OR COALESCE(p.middle_name, '') ILIKE '%' || sqlc.arg(search_query)::text || '%'
    OR p.first_last_name ILIKE '%' || sqlc.arg(search_query)::text || '%'
    OR COALESCE(p.second_last_name, '') ILIKE '%' || sqlc.arg(search_query)::text || '%'
    OR p.email ILIKE '%' || sqlc.arg(search_query)::text || '%'
    OR p.phone ILIKE '%' || sqlc.arg(search_query)::text || '%'
  )
  AND (
    NOT sqlc.arg(with_pending_appointments)::bool
    OR EXISTS (
      SELECT 1
      FROM appointments ap
      WHERE ap.patient_id = p.id
        AND ap.status = 'APPOINTMENT_STATUS_SCHEDULED'
        AND ap.start_time >= now()
    )
  )
  AND (
    NOT sqlc.arg(with_no_appointments)::bool
    OR NOT EXISTS (
      SELECT 1
      FROM appointments ap
      WHERE ap.patient_id = p.id
    )
  )
  AND (
    NOT sqlc.arg(ever_had_appointments)::bool
    OR EXISTS (
      SELECT 1
      FROM appointments ap
      WHERE ap.patient_id = p.id
    )
  )
  AND (
    sqlc.narg(after_id)::uuid IS NULL
    OR p.id < sqlc.narg(after_id)::uuid
  )
  AND (
    sqlc.narg(before_id)::uuid IS NULL
    OR p.id > sqlc.narg(before_id)::uuid
  )
ORDER BY
  CASE WHEN sqlc.arg(is_backward)::bool THEN p.id END ASC,
  CASE WHEN NOT sqlc.arg(is_backward)::bool THEN p.id END DESC
LIMIT sqlc.arg(limit_count)::int;

-- name: UpdatePatientIdentification :execrows
UPDATE patients
SET
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  middle_name = COALESCE(sqlc.narg(middle_name), middle_name),
  first_last_name = COALESCE(sqlc.narg(first_last_name), first_last_name),
  second_last_name = COALESCE(sqlc.narg(second_last_name), second_last_name),
  birth_date = COALESCE(sqlc.narg(birth_date), birth_date),
  sex = COALESCE(sqlc.narg(sex), sex),
  birth_country = COALESCE(sqlc.narg(birth_country), birth_country),
  birth_state = COALESCE(sqlc.narg(birth_state), birth_state),
  birth_city = COALESCE(sqlc.narg(birth_city), birth_city),
  occupation = COALESCE(sqlc.narg(occupation), occupation),
  marital_status = COALESCE(sqlc.narg(marital_status), marital_status),
  religion = COALESCE(sqlc.narg(religion), religion),
  updated_at = sqlc.arg(updated_at)
WHERE id = sqlc.arg(id)
  AND psychologist_id = sqlc.arg(psychologist_id)
  AND deleted_at IS NULL;

-- name: UpdatePatientContact :execrows
UPDATE patients
SET
  phone = COALESCE(sqlc.narg(phone), phone),
  email = COALESCE(sqlc.narg(email), email),
  updated_at = sqlc.arg(updated_at)
WHERE id = sqlc.arg(id)
  AND psychologist_id = sqlc.arg(psychologist_id)
  AND deleted_at IS NULL;

-- name: UpdateAddressByPatientID :execrows
UPDATE addresses a
SET
  country = COALESCE(sqlc.narg(country), a.country),
  province = COALESCE(sqlc.narg(province), a.province),
  city = COALESCE(sqlc.narg(city), a.city),
  postal_code = COALESCE(sqlc.narg(postal_code), a.postal_code),
  neighborhood = COALESCE(sqlc.narg(neighborhood), a.neighborhood),
  street = COALESCE(sqlc.narg(street), a.street),
  street_number = COALESCE(sqlc.narg(street_number), a.street_number),
  unit_number = COALESCE(sqlc.narg(unit_number), a.unit_number),
  updated_at = sqlc.arg(updated_at)
FROM patients p
WHERE p.address_id = a.id
  AND p.id = sqlc.arg(patient_id)
  AND p.psychologist_id = sqlc.arg(psychologist_id)
  AND p.deleted_at IS NULL;

-- name: PatientExists :one
SELECT EXISTS (
  SELECT 1
  FROM patients
  WHERE id = $1
    AND psychologist_id = $2
    AND deleted_at IS NULL
) AS exists;
