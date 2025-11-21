#include <BME680.h>

bool BME680::begin()
{
    this->bme = new Adafruit_BME680();
    if (!bme->begin()) 
    {
        Serial.println("Could not find a valid BME680 sensor, check wiring!");
        return false;
    }
    bme->setTemperatureOversampling(BME680_OS_8X);
    bme->setHumidityOversampling(BME680_OS_2X);
    bme->setPressureOversampling(BME680_OS_4X);
    bme->setIIRFilterSize(BME680_FILTER_SIZE_3);
    bme->setGasHeater(320, 150); // 320*C for 150 ms 
    bme->beginReading();
    return true;
}


void BME680::step()
{

}

void BME680::recordWeather(WeatherReport* report)
{
    bme->endReading();
    report->gas = &(bme->gas_resistance);
    report->humidity = &bme->humidity;
    report->pressure = &bme->pressure;
    report->temperature = &bme->temperature;
    bme->beginReading();
}
