
#ifndef ANEMOMETER_H
#define ANEMOMETER_H

#define ANEMOMETER_BUFFER_SIZE 1000

class Anemometer
{
private:
// Configurations
  int pin;
  double circumfrence;
  unsigned long debounceWait;

// State
  bool lastState = false;
  unsigned long lastReadTime;
  unsigned long readings[ANEMOMETER_BUFFER_SIZE];
  int currentReadingIndex = 0;
public:
  Anemometer(int _pin, int _debounceWait, double _circumfrence);
  void reset();
  void takeReading();
  double getSpeed();
};

#endif
