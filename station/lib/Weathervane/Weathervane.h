#include <Adafruit_MCP23X17.h>
#include <Provider.h>

#ifndef WEATHERVANE_H
#define WEATHERVANE_H

class Weathervane : public Provider
{
private:
  Adafruit_MCP23X17 *mcp;
  double direction = 0;
public:
  bool begin();
  void step();
  void recordWeather(WeatherReport* report);
};

#endif
