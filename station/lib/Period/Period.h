#ifndef PERIOD_H
#define PERIOD_H

class Period
{
private:
  unsigned long start;
  unsigned long wait;
public:
  Period(unsigned long wait);
  bool isComplete();
  void reset();
};

#endif
