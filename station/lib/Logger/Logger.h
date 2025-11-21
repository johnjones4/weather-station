#include <Provider.h>

#ifndef LOGGER_H
#define LOGGER_H

class Logger
{
private:
  String url;
public:
  Logger(String url);
  bool post(WeatherReport *data);
};

#endif
