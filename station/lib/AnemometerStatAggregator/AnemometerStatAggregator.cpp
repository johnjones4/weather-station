#include "AnemometerStatAggregator.h"
#include <float.h>

#ifndef ANEMOMETER_STAT_AGGREGATOR_IMPL
#define ANEMOMETER_STAT_AGGREGATOR_IMPL

bool AnemometerStatAggregator::append(double speed)
{
  if (end >= ANEMOMETER_BUFFER_SIZE)
  {
    return false;
  }
  buffer[end] = speed;
  end++;
  return true;
}

AnemometerStatsSet AnemometerStatAggregator::getStats()
{
  AnemometerStatsSet statsSet;
  statsSet.min = DBL_MAX;
  statsSet.max = DBL_MIN;

  for (int i = 0; i < end; i++)
  {
    double speed = buffer[i % ANEMOMETER_BUFFER_SIZE];

    statsSet.average += speed;

    if (speed > statsSet.max)
    {
      statsSet.max = speed;
    }

    if (speed < statsSet.min)
    {
      statsSet.min = speed;
    }
  }

  statsSet.average = statsSet.average / double(end);

  return statsSet;
}

void AnemometerStatAggregator::reset()
{
  end = 0;
}

#endif
