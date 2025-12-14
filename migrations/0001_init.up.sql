CREATE TABLE
  IF NOT EXISTS tag (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    hex_color TEXT NOT NULL UNIQUE
  );

CREATE TABLE
  IF NOT EXISTS education (
    id BIGINT PRIMARY KEY,
    name TEXT,
    year INT NOT NULL,
    course TEXT NOT NULL,
    organization TEXT NOT NULL
  );

CREATE TABLE
  IF NOT EXISTS technology (
    id BIGINT PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    logo_url TEXT -- URL к файлу в S3
  );

CREATE TABLE
  IF NOT EXISTS technologies_tag (
    tag_id BIGINT NOT NULL REFERENCES tag (id) ON DELETE CASCADE,
    technology_id BIGINT NOT NULL REFERENCES technology (id) ON DELETE CASCADE,
    PRIMARY KEY (tag_id, technology_id)
  );

CREATE TABLE
  IF NOT EXISTS work_history (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    about TEXT NOT NULL,
    logo_url JSONB, -- URL к файлу в S3
    period_start DATE NOT NULL,
    period_end DATE,
    what_i_did TEXT[],
    projects TEXT[]
  );

CREATE TABLE
  IF NOT EXISTS work_history_technology (
    work_history_id BIGINT NOT NULL REFERENCES work_history (id) ON DELETE CASCADE,
    technology_id BIGINT NOT NULL REFERENCES technology (id) ON DELETE CASCADE,
    PRIMARY KEY (work_history_id, technology_id)
  );

