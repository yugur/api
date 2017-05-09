CREATE TABLE languages (
	lang_id		bigserial	PRIMARY KEY,
	name		text		NOT NULL,
	code		text		NOT NULL UNIQUE
);

CREATE TABLE wordtypes (
	wordtype_id	bigserial	PRIMARY KEY,
	name		text		NOT NULL UNIQUE
);

CREATE TABLE tags (
	tag_id		bigserial		PRIMARY KEY,
	name		text		NOT NULL UNIQUE
);

CREATE TABLE entries (
	entry_id		bigserial 	PRIMARY KEY,
	headword	text 		NOT NULL,
	wordtype	bigint		REFERENCES wordtypes (wordtype_id),
	definition	text		,
	hw_lang		bigint		REFERENCES languages (lang_id) ON DELETE CASCADE,
	def_lang	bigint		REFERENCES languages (lang_id) ON DELETE CASCADE
);

CREATE TABLE entry_tags (
	tag_id		bigint,
	entry_id	bigint,
	CONSTRAINT PK_entry_tags PRIMARY KEY
    (
        tag_id,
        entry_id
    ),
    FOREIGN KEY (tag_id) REFERENCES tags (tag_id),
    FOREIGN KEY (entry_id) REFERENCES entries (entry_id)
);

CREATE TABLE users (
	user_id		bigserial	PRIMARY KEY,
	username	text		NOT NULL UNIQUE,
	hash		text		NOT NULL,
	email		text		NOT NULL UNIQUE,
	dob			date		,
	gender		text		,
	joindate	timestamp	NOT NULL,
	language	bigint		REFERENCES languages (lang_id),
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

INSERT INTO tags (name) VALUES
	('flame'),
	('blaze'),
	('passion'),
	('fervor');

INSERT INTO entries (headword, wordtype, definition, hw_lang, def_lang) VALUES
	('fire', (SELECT wordtype_id FROM wordtypes WHERE name='noun'), 'Burning fuel or other material: a cooking fire; a forest fire.', (SELECT lang_id FROM languages WHERE code='en-AU'), (SELECT lang_id FROM languages WHERE code='en-AU')),
	('fire', (SELECT wordtype_id FROM wordtypes WHERE name='noun'), 'Burning intensity of feeling; ardor.', (SELECT lang_id FROM languages WHERE code='en-AU'), (SELECT lang_id FROM languages WHERE code='en-AU'));

-- add 'flame' tag to all entries with the headword 'fire'
INSERT INTO
	entry_tags (tag_id, entry_id)
SELECT
  t.tag_id,
  e.entry_id
FROM
  (
    SELECT
      tag_id
    FROM
      tags
    WHERE
      name='flame'
  ) AS t
CROSS JOIN
  (
    SELECT
      entry_id
    FROM
      entries
    WHERE
      headword='fire'
  ) AS e
;