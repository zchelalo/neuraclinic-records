CREATE TABLE addresses (
  id uuid PRIMARY KEY,
  country varchar(100) NOT NULL,
  province varchar(100) NOT NULL,
  city varchar(100) NOT NULL,
  postal_code varchar(12) NOT NULL,
  neighborhood varchar(100) NOT NULL,
  street varchar(150) NOT NULL,
  street_number varchar(50) NOT NULL,
  unit_number varchar(50),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE patients (
  id uuid PRIMARY KEY,
  first_name varchar(100) NOT NULL,
  middle_name varchar(100),
  first_last_name varchar(100) NOT NULL,
  second_last_name varchar(100),
  birth_date date NOT NULL,
  birth_country varchar(100) NOT NULL,
  birth_state varchar(100) NOT NULL,
  birth_city varchar(100) NOT NULL,
  sex varchar(50) NOT NULL,
  marital_status varchar(50) NOT NULL,
  occupation varchar(100),
  religion varchar(50),
  phone varchar(15) NOT NULL,
  email varchar(254) NOT NULL,
  address_id uuid NOT NULL REFERENCES addresses(id) ON DELETE RESTRICT,
  psychologist_id uuid NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE familiograms (
  id uuid PRIMARY KEY,
  data jsonb NOT NULL,
  patient_id uuid UNIQUE NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE appointments (
  id uuid PRIMARY KEY,
  start_time timestamptz NOT NULL,
  end_time timestamptz NOT NULL,
  reason text NOT NULL,
  status varchar(50) NOT NULL,
  patient_id uuid NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
  cancelled_by_user_id uuid,
  rescheduled_from_appointment_id uuid REFERENCES appointments(id) ON DELETE SET NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE notes (
  id uuid PRIMARY KEY,
  patient_id uuid NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
  appointment_id uuid REFERENCES appointments(id) ON DELETE SET NULL,
  title varchar(255),
  content_html text NOT NULL,
  content_text text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE attachments (
  id uuid PRIMARY KEY,
  file_id uuid NOT NULL,
  mime_type varchar(100) NOT NULL,
  patient_id uuid NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
  note_id uuid REFERENCES notes(id) ON DELETE SET NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE INDEX idx_patients_psychologist_active
  ON patients (psychologist_id, created_at DESC, id DESC)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_patients_email_active
  ON patients (email)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_appointments_patient_start
  ON appointments (patient_id, start_time DESC);

CREATE INDEX idx_notes_patient_created_active
  ON notes (patient_id, created_at DESC, id DESC)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_attachments_patient_created_active
  ON attachments (patient_id, created_at DESC, id DESC)
  WHERE deleted_at IS NULL;

