ALTER TABLE tag ALTER COLUMN id DROP DEFAULT;
SELECT pg_get_serial_sequence('tag', 'id') AS sequence_name;
ALTER SEQUENCE tag_id_seq OWNED BY NONE;
DROP SEQUENCE IF EXISTS tag_id_seq;

ALTER TABLE education ALTER COLUMN id DROP DEFAULT;
SELECT pg_get_serial_sequence('education', 'id') AS sequence_name;
ALTER SEQUENCE education_id_seq OWNED BY NONE;
DROP SEQUENCE IF EXISTS education_id_seq;

ALTER TABLE technology ALTER COLUMN id DROP DEFAULT;
SELECT pg_get_serial_sequence('technology', 'id') AS sequence_name;
ALTER SEQUENCE technology_id_seq OWNED BY NONE;
DROP SEQUENCE IF EXISTS technology_id_seq;

ALTER TABLE work_history ALTER COLUMN id DROP DEFAULT;
SELECT pg_get_serial_sequence('work_history', 'id') AS sequence_name;
ALTER SEQUENCE work_history_id_seq OWNED BY NONE;
DROP SEQUENCE IF EXISTS work_history_id_seq;
