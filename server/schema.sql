CREATE TABLE IF NOT EXISTS weather (
  timestamp TIMESTAMP NOT NULL PRIMARY KEY,
  source varchar(128),
  wind_speed REAL NOT NULL,
  vane_direction REAL NOT NULL,
  temperature REAL NOT NULL,
  pressure REAL NOT NULL,
  humidity REAL NOT NULL,
  gas REAL NOT NULL,
  rainfall READ NOT NULL
);
