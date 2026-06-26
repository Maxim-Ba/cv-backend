ALTER TABLE profile
  ADD COLUMN IF NOT EXISTS note TEXT,
  ADD COLUMN IF NOT EXISTS hobbies TEXT;

CREATE TABLE IF NOT EXISTS profile_technology (
  profile_id    BIGINT NOT NULL REFERENCES profile (id) ON DELETE CASCADE,
  technology_id BIGINT NOT NULL REFERENCES technology (id) ON DELETE CASCADE,
  PRIMARY KEY (profile_id, technology_id)
);

UPDATE profile
SET
  about = E'Занимаюсь программированием с 2015 года. Мой основной стек Typescript/JavaScript, React, Redux.\n\nПродолжительное время занимался программированием станков, это направление мне очень нравилось, но из-за отсутствия возможности дальнейшего развития в 2018 году сменил вектор с производства на IT.\nВыбрал для себя на тот момент самый интересный для себя путь - Web разработка.',
  note = 'Люблю принципы Clean Architecture и функциональное программирование.',
  hobbies = 'Мои хобби: спорт, игра на барабанах и электрогитаре.'
WHERE id = (SELECT id FROM profile LIMIT 1);

INSERT INTO profile_technology (profile_id, technology_id)
SELECT p.id, t.id
FROM profile p
CROSS JOIN technology t
WHERE t.title IN (
  'TypeScript',
  'JavaScript',
  'React',
  'Redux',
  'Python',
  'Django',
  'C++',
  'Haskell'
)
ON CONFLICT DO NOTHING;
