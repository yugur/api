CREATE TABLE languages (
	id			bigserial	PRIMARY KEY,
	name		text		NOT NULL,
	code		text		NOT NULL UNIQUE
);

CREATE TABLE wordtypes (
	id			bigserial	PRIMARY KEY,
	name		text		NOT NULL UNIQUE
);

CREATE TABLE entries (
	id			bigserial 	PRIMARY KEY,
	headword	text 		NOT NULL,
	type		bigint		REFERENCES wordtypes(id) ON DELETE CASCADE,
	definition	text		,
	tags		text[]		,
	language	bigint		REFERENCES languages(id) ON DELETE CASCADE
);

CREATE TABLE users (
	id			bigserial	PRIMARY KEY,
	username	text		NOT NULL UNIQUE,
	hash		text		NOT NULL,
	email		text		NOT NULL UNIQUE,
	dob			date		,
	gender		text		,
	joindate	timestamp	NOT NULL,
	language	bigint		REFERENCES languages(id) ON DELETE CASCADE,
	fluency		int[]
);

INSERT INTO languages (name, code) VALUES
	('English (AU)', 'en-AU'),
	('Western Yugur', 'yge'),
	('Eastern Yugur', 'yuy'),
	('Chinese', 'zh'),
	('한국어', 'ko-KR');

INSERT INTO wordtypes (name) VALUES
	('noun'),
	('verb'),
	('adjective');

INSERT INTO entries (headword, type, definition, tags, language) VALUES
	('fire', (SELECT id FROM languages WHERE name='noun'), 'Burning fuel or other material: a cooking fire; a forest fire.', '{"flame", "blaze"}', (SELECT id FROM languages WHERE code='en-AU')),
	('fire', (SELECT id FROM languages WHERE name='noun'), 'Burning intensity of feeling; ardor.', '{"passion", "devotion"}', (SELECT id FROM languages WHERE code='en-AU'));