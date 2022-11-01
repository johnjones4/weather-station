#include "AnemometerStatAggregator.h"
#include <float.h>

#ifndef ANEMOMETER_STAT_AGGREGATOR_IMPL
#define ANEMOMETER_STAT_AGGREGATOR_IMPL

void AnemometerStatAggregator::append(AnemometerStats stat)
{
  buffer[end % ANEMOMETER_BUFFER_SIZE] = stat;
  end++;
}

AnemometerStatsSet AnemometerStatAggregator::getAndReset()
{
  AnemometerStatsSet statsSet;
  statsSet.min.rotationsPerSecond = DBL_MAX;
  statsSet.max.rotationsPerSecond = DBL_MIN;

  int start = end < ANEMOMETER_BUFFER_SIZE ? 0 : end + 1; 
  int length = end < ANEMOMETER_BUFFER_SIZE ? end : ANEMOMETER_BUFFER_SIZE;

  for (int i = start; i < start + length; i++)
  {
    AnemometerStats stats = buffer[i % ANEMOMETER_BUFFER_SIZE];

    statsSet.average.metersPerSecond += stats.metersPerSecond;
    statsSet.average.rotationsPerSecond += stats.rotationsPerSecond;

    if (stats.rotationsPerSecond > statsSet.max.rotationsPerSecond) {
      statsSet.max = stats;
    }

    if (stats.rotationsPerSecond < statsSet.max.rotationsPerSecond) {
      statsSet.min = stats;
    }
  }

  statsSet.average.metersPerSecond = statsSet.average.metersPerSecond / double(length);
  statsSet.average.rotationsPerSecond = statsSet.average.rotationsPerSecond / double(length);

  end = 0;  

  return statsSet;
}

#endif
