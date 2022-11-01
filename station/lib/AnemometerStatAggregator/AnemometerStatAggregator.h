#include "Anemometer.h"

#ifndef ANEMOMETER_STAT_AGGREGATOR_H
#define ANEMOMETER_STAT_AGGREGATOR_H

#define ANEMOMETER_BUFFER_SIZE 5000

typedef struct AnemometerStatsSet
{
  AnemometerStats min;
  AnemometerStats max;
  AnemometerStats average;
} AnemometerStatsSet;

class AnemometerStatAggregator
{
private:
  AnemometerStats buffer[ANEMOMETER_BUFFER_SIZE];
  int end = 0;
public:
  void append(AnemometerStats stat);
  AnemometerStatsSet getAndReset();
};

#endif
