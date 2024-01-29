ALTER TABLE spaces ADD CONSTRAINT spaces_name_key UNIQUE (name);
ALTER TABLE folders ADD CONSTRAINT folders_space_id_name_key UNIQUE (space_id, name);

---- create above / drop below ----

ALTER TABLE spaces DROP CONSTRAINT spaces_name_key;
ALTER TABLE folders DROP CONSTRAINT folders_space_id_name_key;
