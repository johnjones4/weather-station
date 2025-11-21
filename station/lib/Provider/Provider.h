#ifndef Provider_h
#define Provider_h

#include <Arduino.h>

typedef struct
{
  double min;
  double max;
  double average;
} AnemometerStatsSet;

typedef struct {
    double* windSpeed;
    double* vaneDirection;
    float* temperature;
    uint32_t* pressure;
    float* humidity;
    uint32_t* gas;
    double *rainfall;
} WeatherReport;

class Provider {
public:
    virtual void recordWeather(WeatherReport* report) = 0;
    virtual void step() = 0;
    virtual bool begin() = 0;
};

#endif