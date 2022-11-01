#include "Period.h"
#include <Arduino.h>

#ifndef PERIOD_IMPL
#define PERIOD_IMPL

Period::Period(unsigned long _wait)
{
  wait = _wait;
  reset();
}

bool Period::isComplete()
{
  unsigned long now = millis();
  return now - start > wait;
}

void Period::reset()
{
  start = millis();
}

#endif
