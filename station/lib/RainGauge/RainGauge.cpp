#include <RainGauge.h>
#include <Arduino.h>

bool RainGauge::begin()
{
    pinMode(RAIN_GAUGE_PIN, INPUT_PULLUP);
    lastFlip = millis();
}

void RainGauge::recordWeather(WeatherReport* report)
{
    volume = (double)this->flips * (double)ML_PER_FLIP;
    report->rainfall = &volume;
    this->flips = 0;
}

void RainGauge::step()
{
    unsigned long now = millis();
    if (digitalRead(RAIN_GAUGE_PIN) == LOW && now - lastFlip > RAINFALL_DEBOUNCE) {
        this->flips++;
        lastFlip = now;
    }
}
