use gmq;
drop table if exists queue;
create table queue (
	id serial PRIMARY KEY NOT NULL AUTO_INCREMENT,
	message MEDIUMTEXT,
	processed BOOL DEFAULT FALSE,
	UNIQUE INDEX proc_idx USING BTREE (id, processed)
);