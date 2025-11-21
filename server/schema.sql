CREATE TABLE IF NOT EXISTS weather (
  tstamp TIMESTAMP NOT NULL PRIMARY KEY,
  anemometer_avg REAL NOT NULL,
  anemometer_min REAL NOT NULL,
  anemometer_max REAL NOT NULL,
  vein_direction REAL NOT NULL,
  temperature REAL NOT NULL,
  temperature_calibration_factor REAL NOT NULL,
  gas INTEGER NOT NULL,
  pressure INTEGER NOT NULL,
  relative_humidity REAL NOT NULL
);
