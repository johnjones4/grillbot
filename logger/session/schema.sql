CREATE TABLE readings (
  received INT NOT NULL PRIMARY KEY,
  temp1 NUMERIC NOT NULL,
  temp2 NUMERIC NOT NULL
);

CREATE TABLE metadata (
  key TEXT NOT NULL PRIMARY KEY,
  value TEXT
);
