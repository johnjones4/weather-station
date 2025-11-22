CREATE TABLE IF NOT EXISTS weather (
  timestamp TIMESTAMP NOT NULL,
  source varchar(128) NOT NULL,
  wind_speed REAL,
  vane_direction REAL,
  temperature REAL,
  pressure REAL,
  humidity REAL,
  gas REAL,
  rainfall REAL,
  PRIMARY KEY(timestamp, source)
);
