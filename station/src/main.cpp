#include <WiFi.h>
#include <Wire.h>
#include <SPI.h>
#include <Adafruit_Sensor.h>
#include <esp_task_wdt.h>
#include <Adafruit_BME680.h>
#include <Arduino.h>
#include <Anemometer.h>
#include <AnemometerStatAggregator.h>
#include <Period.h>
#include <secrets.h>
#include <WeatherLogger.h>
#include <Weathervane.h>

#define MINOR_PERIOD_SECONDS 5
#define MAJOR_PERIOD_SECONDS 60

#define ANEMOMETER_PIN 33
#define ANEMOMETER_DEBOUNCE 100
#define ANEMOMETER_CIRCUMFERENCE 0.50265482457

Adafruit_BME680 bme;

Weathervane weatherVane;

Anemometer anemometer(ANEMOMETER_PIN, ANEMOMETER_DEBOUNCE, ANEMOMETER_CIRCUMFERENCE);
AnemometerStatAggregator anemometerStatAggregator;

Period minorPeriod(MINOR_PERIOD_SECONDS * 1000);
Period majorPeriod(MAJOR_PERIOD_SECONDS * 1000);

WeatherLogger logger(POST_URL);

void fatal() 
{
  delay(1000);
  ESP.restart();
}

bool initWifi() 
{
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  const unsigned long start = millis();
  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.println("Establishing connection to WiFi..");
    if (millis() - start > 60000) {
      return false;
    }
  }
  return true;
}

void setup() 
{
  esp_task_wdt_init(MAJOR_PERIOD_SECONDS * 5, true);
  esp_task_wdt_add(NULL);
  
  Serial.begin(115200);

  Serial.println("Starting up");

  if (!initWifi()) 
  {
    Serial.println("WiFi timeout ... rebooting");
    fatal();
  }

  if (!bme.begin()) 
  {
    Serial.println("Could not find a valid BME680 sensor, check wiring!");
    fatal();
  }
  bme.setTemperatureOversampling(BME680_OS_8X);
  bme.setHumidityOversampling(BME680_OS_2X);
  bme.setPressureOversampling(BME680_OS_4X);
  bme.setIIRFilterSize(BME680_FILTER_SIZE_3);
  bme.setGasHeater(320, 150); // 320*C for 150 ms 
  bme.beginReading();

  if (!weatherVane.begin()) {
    Serial.println("Could not find a valid MCP23017 sensor, check wiring!");
    fatal();
  }
  
  Serial.println("Ready");
}

void loop() 
{
  // Take readings
  anemometer.takeReading();

  // Do our minor period things
  if (minorPeriod.isComplete())
  {
    AnemometerStats stats = anemometer.getStats();
    anemometer.reset();

    anemometerStatAggregator.append(stats);

    minorPeriod.reset();
  }

  // Do our major period things
  if (majorPeriod.isComplete())
  {
    if (!bme.endReading())
    {
      Serial.println("Bad BME680 reading");
      fatal();
    }
    bme.beginReading();

    if (!weatherVane.performReading()) {
      Serial.println("Bad weather vane reading");
      fatal();
    }

    WeatherLoggerData data;
    data.anemometerStatsSet = anemometerStatAggregator.getAndReset();
    data.gas = bme.gas_resistance;
    data.humidity = bme.humidity;
    data.pressure = bme.pressure;
    data.temperature = bme.temperature;
    data.vaneDirection = weatherVane.direction;

    if (!logger.postWeather(data)) {
      Serial.println("Bad HTTP response");
      fatal();
    }

    esp_task_wdt_reset();

    majorPeriod.reset();
  }
}
