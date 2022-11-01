#include <AnemometerStatAggregator.h>
#include <Arduino.h>

#ifndef LOGGER_H
#define LOGGER_H

typedef struct LoggerData
{
  AnemometerStatsSet anemometerStatsSet;
  float temperature;
  int pressure;
  float humidity;
  int gas;
  double vaneDirection;
} LoggerData;

class Logger
{
private:
  String url;
public:
  Logger(String url);
  bool post(LoggerData data);
};

#endif
