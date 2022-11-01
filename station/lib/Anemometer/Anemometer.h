
#ifndef ANEMOMETER_H
#define ANEMOMETER_H

typedef struct AnemometerStats
{
  double rotationsPerSecond;
  double metersPerSecond;
} AnemometerStats;

class Anemometer
{
private:
// Configurations
  int pin;
  double circumfrence;
  unsigned long debounceWait;

// State
  int rotations;
  unsigned long lastReadTime;
  unsigned long startTime;
public:
  Anemometer(int _pin, int _debounceWait, double _circumfrence);
  void reset();
  void takeReading();
  AnemometerStats getStats();
};

#endif
