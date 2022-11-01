CREATE TABLE IF NOT EXISTS weather (
  tstamp TIMESTAMP NOT NULL PRIMARY KEY,
  anemometer_rotations INTEGER NOT NULL,
  wind_speed REAL NOT NULL,
  vein_direction INTEGER NOT NULL,
  temperature REAL NOT NULL,
  temperature_calibration_factor REAL NOT NULL,
  gas REAL NOT NULL,
  relative_humidity INTEGER NOT NULL
);
