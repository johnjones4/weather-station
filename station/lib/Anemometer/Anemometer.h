
#ifndef ANEMOMETER_H
#define ANEMOMETER_H

#include <Provider.h>

#define ANEMOMETER_BUFFER_SIZE 1000
#define ANEMOMETER_PIN 34
#define ANEMOMETER_DEBOUNCE 125
#define ANEMOMETER_CIRCUMFERENCE 0.502654824574
#define SPEED_CALC_WAIT 5000

class Anemometer : public Provider
{
private:
  bool lastState = false;
  unsigned long lastReadTime;
  unsigned long firstReadTime;
  unsigned long readings[ANEMOMETER_BUFFER_SIZE];
  int currentReadingIndex = 0;
  double currentSpeed = 0;

  bool append(double speed);
  void reset();
  void calculateSpeed();
public:
  bool begin();
  void step();
  void recordWeather(WeatherReport* report);
};

#endif
