#include <Providers.h>

void Providers::add(Provider* p) {
    this->list[this->index] = p;
    this->index++;
}

int Providers::size() {
    return this->index;
}

Provider* Providers::at(int index) {
    return this->list[index];
}