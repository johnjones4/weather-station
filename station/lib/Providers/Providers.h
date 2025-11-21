#ifndef Providers_h
#define Providers_h

#include <Provider.h>

#define PROVIDERS_BUFFER 5

class Providers
{
private:
    Provider* list[PROVIDERS_BUFFER];
    int index = 0;
public:
    void add(Provider* p);
    int size();
    Provider* at(int index);
};

#endif