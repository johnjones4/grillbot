CREATE TABLE readings (
  received INT NOT NULL PRIMARY KEY,
  Temperatures[0] NUMERIC NOT NULL,
  Temperatures[1] NUMERIC NOT NULL
);

CREATE TABLE metadata (
  key TEXT NOT NULL PRIMARY KEY,
  value TEXT
);
