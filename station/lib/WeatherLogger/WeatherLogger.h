#include <AnemometerStatAggregator.h>

#ifndef WEATHER_LOGGER_H
#define WEATHER_LOGGER_H

typedef struct WeatherLoggerData
{
  AnemometerStatsSet anemometerStatsSet;
  float temperature;
  int pressure;
  float humidity;
  int gas;
  double vaneDirection;
} WeatherLoggerData;

class WeatherLogger
{
private:
  char* url;
public:
  WeatherLogger(char* url);
  bool postWeather(WeatherLoggerData data);
};

#endif
