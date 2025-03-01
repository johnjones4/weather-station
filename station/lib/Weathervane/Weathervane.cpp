#include "Weathervane.h"

#ifndef WEATHERVANE_IMPL
#define WEATHERVANE_IMPL

#define ANGLE 45.0
#define N_PINS 8
#define SENSOR_OFFSET 6
int pins[N_PINS] = {0,1,2,3,4,5,6,7};

bool Weathervane::begin()
{
  if (!mcp.begin_I2C())
  {
    return false;
  }

  for (int p = 0; p < N_PINS; p++)
  {
    mcp.pinMode(pins[p], INPUT);
  }

  return true;
}

bool Weathervane::performReading()
{
  int firstPin = -1;
  int lastPin = -1;
  for (int p = 0; p < N_PINS; p++)
  {
    int pp = (SENSOR_OFFSET + p) % N_PINS;
    if (mcp.digitalRead(pins[p]) == LOW)
    {
      if (firstPin == -1)
      {
        firstPin = pp;
      }
      lastPin = pp;
    } 
    else if (firstPin >= 0)
    {
      break;
    }
  }
  if (firstPin < 0) 
  {
    return false;
  }
  if (firstPin == lastPin)
  {
    direction = double(firstPin) * ANGLE;
  } else {
    double mid = double(firstPin) + (double(lastPin - firstPin) / 2.0);
    direction = (mid * ANGLE);
  }
  return true;
}

#endif
