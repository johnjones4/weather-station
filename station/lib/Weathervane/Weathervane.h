#include <Adafruit_MCP23X17.h>

#ifndef WEATHERVANE_H
#define WEATHERVANE_H

class Weathervane
{
private:
  Adafruit_MCP23X17 mcp;
public:
  double direction = 0;
  bool begin();
  bool performReading();
};

#endif
