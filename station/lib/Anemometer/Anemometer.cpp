#include <Arduino.h>
#include "Anemometer.h"

#ifndef ANEMOMETER_IMPL
#define ANEMOMETER_IMPL

Anemometer::Anemometer(int _pin, int _debounceWait, double _circumfrence)
{
  pin = _pin;
  pinMode(pin, INPUT);
  debounceWait = _debounceWait;
  circumfrence = _circumfrence;
  lastReadTime = millis();
  reset();
}

void Anemometer::reset()
{
  rotations = 0;
}

void Anemometer::takeReading()
{
  int val = digitalRead(pin);
  unsigned long now = millis();
  bool state = val == 0;
  if (state && !lastState && now - lastReadTime > debounceWait) {
    rotations++;
    lastReadTime = now;
    lastState = true;
  } else if (!state) {
    lastState = false;
  }
}

double Anemometer::getSpeed()
{
  unsigned long now = millis();
  double secondsElapsed = double(now - startTime);
  double meters = circumfrence * double(rotations);
  return meters / secondsElapsed;
}

#endif
