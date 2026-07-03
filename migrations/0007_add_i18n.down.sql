DROP INDEX IF EXISTS tag_name_ru_unique;

ALTER TABLE tag
  ALTER COLUMN name TYPE TEXT USING COALESCE(name->>'ru', '');

ALTER TABLE tag ADD CONSTRAINT tag_name_key UNIQUE (name);

ALTER TABLE education
  ALTER COLUMN organization TYPE TEXT USING COALESCE(organization->>'ru', ''),
  ALTER COLUMN course TYPE TEXT USING COALESCE(course->>'ru', ''),
  ALTER COLUMN name TYPE TEXT USING CASE WHEN name IS NULL THEN NULL ELSE name->>'ru' END;

ALTER TABLE technology
  ALTER COLUMN description TYPE TEXT USING CASE WHEN description IS NULL THEN NULL ELSE description->>'ru' END;

ALTER TABLE work_history
  ALTER COLUMN projects TYPE TEXT[] USING ARRAY(SELECT jsonb_array_elements_text(projects->'ru')),
  ALTER COLUMN what_i_did TYPE TEXT[] USING ARRAY(SELECT jsonb_array_elements_text(what_i_did->'ru')),
  ALTER COLUMN job_title TYPE TEXT USING CASE WHEN job_title IS NULL THEN NULL ELSE job_title->>'ru' END,
  ALTER COLUMN about TYPE TEXT USING COALESCE(about->>'ru', ''),
  ALTER COLUMN name TYPE TEXT USING COALESCE(name->>'ru', '');

ALTER TABLE profile
  DROP COLUMN IF EXISTS greeting,
  DROP COLUMN IF EXISTS pitch,
  ALTER COLUMN hobbies TYPE TEXT USING CASE WHEN hobbies IS NULL THEN NULL ELSE hobbies->>'ru' END,
  ALTER COLUMN note TYPE TEXT USING CASE WHEN note IS NULL THEN NULL ELSE note->>'ru' END,
  ALTER COLUMN about TYPE TEXT USING CASE WHEN about IS NULL THEN NULL ELSE about->>'ru' END,
  ALTER COLUMN title TYPE TEXT USING COALESCE(title->>'ru', ''),
  ALTER COLUMN full_name TYPE TEXT USING COALESCE(full_name->>'ru', '');
