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
  currentReadingIndex = 0;
  reset();
}

void Anemometer::reset()
{
  currentReadingIndex = 0;
  lastReadTime = millis();
}

void Anemometer::takeReading()
{
  int val = digitalRead(pin);
  bool state = val == 0;
  unsigned long now = millis();
  unsigned long elapsed = now - lastReadTime;
  if (state && !lastState && elapsed > debounceWait) {
    readings[currentReadingIndex % ANEMOMETER_BUFFER_SIZE] = elapsed;
    currentReadingIndex++;
    lastReadTime = now;
    lastState = true;
  } else if (!state) {
    lastState = false;
  }
}

double Anemometer::getSpeed()
{
  unsigned long min = ULONG_MAX;
  for (int i = 0; i < currentReadingIndex; i++) {
    Serial.printf("Reading at %d: %lu\n", i, readings[currentReadingIndex]);
    if (readings[currentReadingIndex] < min) {
      min = readings[currentReadingIndex];
    }
  }
  if (min == ULONG_MAX) {
    return 0;
  }
  return circumfrence / (double)min;
}

#endif
