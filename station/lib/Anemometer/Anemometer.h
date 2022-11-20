
#ifndef ANEMOMETER_H
#define ANEMOMETER_H

class Anemometer
{
private:
// Configurations
  int pin;
  double circumfrence;
  unsigned long debounceWait;

// State
  int rotations;
  bool lastState = false;
  unsigned long lastReadTime;
  unsigned long startTime;
public:
  Anemometer(int _pin, int _debounceWait, double _circumfrence);
  void reset();
  void takeReading();
  double getSpeed();
};

#endif
