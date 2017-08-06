CREATE TABLE entries (
  headword  	varchar(255) NOT NULL,
  definition 	varchar(255) NOT NULL
);

INSERT INTO entries (headword, definition) VALUES
('happy', 'enjoying or showing or marked by joy or pleasure'),
('sad', 'experiencing or showing sorrow or unhappiness'),
('dog', 'a member of the genus Canis that has been domesticated by man since prehistoric times'),
('cat', 'feline mammal usually having thick soft fur and no ability to roar: domestic cats');

ALTER TABLE entries ADD PRIMARY KEY (headword);

CREATE TABLE users (
	username 	varchar(255) NOT NULL,
	hash		varchar(255) NOT NULL
);

ALTER TABLE users ADD PRIMARY KEY (username);