#include "Anemometer.h"

#ifndef ANEMOMETER_STAT_AGGREGATOR_H
#define ANEMOMETER_STAT_AGGREGATOR_H

#define ANEMOMETER_BUFFER_SIZE 5000

typedef struct AnemometerStatsSet
{
  double min;
  double max;
  double average;
} AnemometerStatsSet;

class AnemometerStatAggregator
{
private:
  double buffer[ANEMOMETER_BUFFER_SIZE];
  int end = 0;
public:
  bool append(double speed);
  AnemometerStatsSet getStats();
  void reset();
};

#endif
