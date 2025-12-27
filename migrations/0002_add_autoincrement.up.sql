CREATE SEQUENCE tag_id_seq;
ALTER TABLE tag ALTER COLUMN id SET DEFAULT nextval('tag_id_seq');
ALTER SEQUENCE tag_id_seq OWNED BY tag.id;

CREATE SEQUENCE education_id_seq;
ALTER TABLE education ALTER COLUMN id SET DEFAULT nextval('education_id_seq');
ALTER SEQUENCE education_id_seq OWNED BY education.id;


CREATE SEQUENCE technology_id_seq;
ALTER TABLE technology ALTER COLUMN id SET DEFAULT nextval('technology_id_seq');
ALTER SEQUENCE technology_id_seq OWNED BY technology.id;


CREATE SEQUENCE work_history_id_seq;
ALTER TABLE work_history ALTER COLUMN id SET DEFAULT nextval('work_history_id_seq');
ALTER SEQUENCE work_history_id_seq OWNED BY work_history.id;



