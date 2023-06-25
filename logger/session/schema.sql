CREATE TABLE readings (
  received INT NOT NULL PRIMARY KEY,
  temperature0 NUMERIC NOT NULL,
  temperature1 NUMERIC NOT NULL
);

CREATE TABLE metadata (
  key TEXT NOT NULL PRIMARY KEY,
  value TEXT
);
