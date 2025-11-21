#include <Arduino.h>
#include "Anemometer.h"

bool Anemometer::begin() 
{
  pinMode(ANEMOMETER_PIN, INPUT);
  reset();
  return true;
}

void Anemometer::step()
{
  int val = digitalRead(ANEMOMETER_PIN);
  bool state = val == 0;
  unsigned long now = millis();
  unsigned long elapsed = now - lastReadTime;
  if (state && !lastState && elapsed > ANEMOMETER_DEBOUNCE) {
    readings[currentReadingIndex % ANEMOMETER_BUFFER_SIZE] = elapsed;
    currentReadingIndex++;
    lastReadTime = now;
    lastState = true;
  } else if (!state) {
    lastState = false;
  }
  if (now - firstReadTime > SPEED_CALC_WAIT) {
    this->calculateSpeed();
  }
}

void Anemometer::recordWeather(WeatherReport* report)
{
  report->windSpeed = &currentSpeed;
}





void Anemometer::reset()
{
  currentReadingIndex = 0;
  lastReadTime = millis();
  firstReadTime = millis();
  currentSpeed = 0;
  lastState = false;
}

void Anemometer::calculateSpeed()
{
  unsigned long min = ULONG_MAX;
  for (int i = 0; i < currentReadingIndex; i++) {
    if (readings[currentReadingIndex] < min) {
      min = readings[currentReadingIndex];
    }
  }
  if (min == ULONG_MAX) {
    return;
  }
  double nextSpeed = ANEMOMETER_CIRCUMFERENCE / (double)min;
  if (nextSpeed > currentSpeed) {
    currentSpeed = nextSpeed;
  }
}
