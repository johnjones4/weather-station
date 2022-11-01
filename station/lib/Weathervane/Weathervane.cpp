#include "Weathervane.h"

#ifndef WEATHERVANE_IMPL
#define WEATHERVANE_IMPL

#define ANGLE 45.0
#define N_PINS 8
int pins[N_PINS] = {0,1,2,3,4,5,6,7};

bool Weathervane::begin()
{
  if (!mcp.begin_I2C()) {
    return false;
  }

  for (int p = 0; p < N_PINS; p++) {
    mcp.pinMode(pins[p], INPUT);
  }

  return true;
}

bool Weathervane::performReading()
{
  int firstPin = -1;
  int lastPin = -1;
  for (int p = 0; p < N_PINS; p++) {
    if (mcp.digitalRead(pins[p]) == LOW) {
      if (firstPin == -1) {
        firstPin = p;
      }
      lastPin = p;
    } else if (firstPin >= 0) {
      break;
    }
  }
  if (firstPin < 0) {
    return false;
  }
  if (firstPin == lastPin) {
    direction = double(firstPin) * ANGLE;
    return true;
  }
  double mid = double(firstPin) + (double(lastPin - firstPin) / 2.0);
  direction = mid * ANGLE;
  return true;
}

#endif WEATHERVANE_IMPL
