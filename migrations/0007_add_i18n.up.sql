-- profile: localized text fields + hero fields
ALTER TABLE profile
  ADD COLUMN IF NOT EXISTS greeting JSONB,
  ADD COLUMN IF NOT EXISTS pitch JSONB;

ALTER TABLE profile
  ALTER COLUMN full_name DROP DEFAULT,
  ALTER COLUMN title DROP DEFAULT;

ALTER TABLE profile
  ALTER COLUMN full_name TYPE JSONB USING jsonb_build_object('ru', full_name, 'en', ''),
  ALTER COLUMN title TYPE JSONB USING jsonb_build_object('ru', title, 'en', ''),
  ALTER COLUMN about TYPE JSONB USING CASE
    WHEN about IS NULL THEN NULL
    ELSE jsonb_build_object('ru', about, 'en', '')
  END,
  ALTER COLUMN note TYPE JSONB USING CASE
    WHEN note IS NULL THEN NULL
    ELSE jsonb_build_object('ru', note, 'en', '')
  END,
  ALTER COLUMN hobbies TYPE JSONB USING CASE
    WHEN hobbies IS NULL THEN NULL
    ELSE jsonb_build_object('ru', hobbies, 'en', '')
  END;

UPDATE profile
SET
  greeting = jsonb_build_object('ru', 'Здравствуйте!', 'en', 'Hello!'),
  pitch = jsonb_build_object(
    'ru', 'Эта страница рассказывает о моих знаниях и навыках в области web-разработки.',
    'en', 'This page showcases my knowledge and skills in web development.'
  )
WHERE greeting IS NULL OR pitch IS NULL;

-- work_history: localized fields
ALTER TABLE work_history
  ALTER COLUMN name TYPE JSONB USING jsonb_build_object('ru', name, 'en', ''),
  ALTER COLUMN about TYPE JSONB USING jsonb_build_object('ru', about, 'en', ''),
  ALTER COLUMN job_title TYPE JSONB USING CASE
    WHEN job_title IS NULL THEN NULL
    ELSE jsonb_build_object('ru', job_title, 'en', '')
  END,
  ALTER COLUMN what_i_did TYPE JSONB USING jsonb_build_object(
    'ru', COALESCE(to_jsonb(what_i_did), '[]'::jsonb),
    'en', '[]'::jsonb
  ),
  ALTER COLUMN projects TYPE JSONB USING jsonb_build_object(
    'ru', COALESCE(to_jsonb(projects), '[]'::jsonb),
    'en', '[]'::jsonb
  );

-- technology: localized description
ALTER TABLE technology
  ALTER COLUMN description TYPE JSONB USING CASE
    WHEN description IS NULL THEN NULL
    ELSE jsonb_build_object('ru', description, 'en', '')
  END;

-- education: localized fields
ALTER TABLE education
  ALTER COLUMN name TYPE JSONB USING CASE
    WHEN name IS NULL THEN NULL
    ELSE jsonb_build_object('ru', name, 'en', '')
  END,
  ALTER COLUMN course TYPE JSONB USING jsonb_build_object('ru', course, 'en', ''),
  ALTER COLUMN organization TYPE JSONB USING jsonb_build_object('ru', organization, 'en', '');

-- tag: localized name
ALTER TABLE tag DROP CONSTRAINT IF EXISTS tag_name_key;
ALTER TABLE tag
  ALTER COLUMN name TYPE JSONB USING jsonb_build_object('ru', name, 'en', '');

CREATE UNIQUE INDEX IF NOT EXISTS tag_name_ru_unique ON tag ((name->>'ru'));
