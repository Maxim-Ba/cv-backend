DROP TABLE IF EXISTS profile_technology;

ALTER TABLE profile
  DROP COLUMN IF EXISTS note,
  DROP COLUMN IF EXISTS hobbies;
