#ifndef RainGauge_h
#define RainGauge_h

#include <Provider.h>

#define RAIN_GAUGE_PIN 14
#define ML_PER_FLIP 3.125
#define RAINFALL_DEBOUNCE 200

class RainGauge : public Provider
{
private:
  int flips = 0;
  unsigned long lastFlip = 0;
  double volume = 0;
public:
  bool begin();
  void step();
  void recordWeather(WeatherReport* report);
};

#endif
