#include <Arduino.h>
#include "Anemometer.h"

#ifndef ANEMOMETER_IMPL
#define ANEMOMETER_IMPL

Anemometer::Anemometer(int _pin, int _debounceWait, double _circumfrence)
{
  pin = _pin;
  debounceWait = _debounceWait;
  circumfrence = _circumfrence;
  reset();
}

void Anemometer::reset()
{
  rotations = 0;
  lastReadTime = startTime = millis();
}

void Anemometer::takeReading()
{
  int val = analogRead(pin);
  unsigned long now = millis();
  if (val == 0 && now - lastReadTime > debounceWait) {
    rotations++;
    lastReadTime = now;
  }
}

AnemometerStats Anemometer::getStats()
{
  AnemometerStats stats;
  unsigned long now = millis();
  double secondsElapsed = double(now - startTime);
  double rotationsD = double(rotations);
  stats.rotationsPerSecond = rotationsD / secondsElapsed;
  double meters = circumfrence * rotationsD;
  stats.metersPerSecond = meters / secondsElapsed;
  return stats;
}

#endif