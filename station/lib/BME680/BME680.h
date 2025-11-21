#ifndef BME680_h
#define BME680_h

#include <Provider.h>
#include <Adafruit_BME680.h>

class BME680 : public Provider
{
private:
  Adafruit_BME680 *bme;
public:
  bool begin();
  void step();
  void recordWeather(WeatherReport* report);
};

#endif