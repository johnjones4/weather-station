#include <WiFi.h>
#include <Wire.h>
#include <SPI.h>
#include <Adafruit_Sensor.h>
#include <esp_task_wdt.h>
#include <Arduino.h>
#include <secrets.h>
#include <BME680.h>
#include <Weathervane.h>
#include <Provider.h>
#include <Logger.h>
#include <Anemometer.h>
#include <features.h>
#include <Providers.h>
#include <RainGauge.h>

#define WATCHDOG_WAIT 60000
#define TRANSMIT_WAIT 300000

Providers providers;
unsigned long lastTransmit;
Logger logger(SERVER_URL);

void fatal() 
{
  Serial.println("Fatal error ... rebooting");
  delay(1000);
  ESP.restart();
}

bool initWifi() 
{
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  const unsigned long start = millis();
  while (WiFi.status() != WL_CONNECTED)
  {
    delay(1000);
    Serial.println("Establishing connection to WiFi..");
    if (millis() - start > 60000)
    {
      return false;
    }
  }
  return true;
}

void setup() 
{
  esp_task_wdt_init(WATCHDOG_WAIT, true);
  esp_task_wdt_add(NULL);
  
  Serial.begin(9600);

  Serial.println("Starting up");

  if (!initWifi()) 
  {
    Serial.println("WiFi timeout");
    fatal();
  }

#ifdef F_ANEMOMETER
  providers.add(new Anemometer());
#endif
#ifdef F_BME680
  providers.add(new BME680());
#endif
#ifdef F_WEATHERVANE
  providers.add(new Weathervane());
#endif
#ifdef F_RAINGAUGE
  providers.add(new RainGauge());
#endif

  for (int i = 0; i < providers.size(); i++) {
    if (!providers.at(i)->begin()) {
      Serial.printf("Bad provider init at %d\n", i);
      fatal();
    }
  }
  
  Serial.println("Ready");
}

void loop() 
{
  esp_task_wdt_reset();

  for (int i = 0; i < providers.size(); i++) {
    providers.at(i)->step();
  }

  unsigned long now = millis();
  if (lastTransmit == 0 || now - lastTransmit > TRANSMIT_WAIT) {
    lastTransmit = now;
    WeatherReport report = {
      windSpeed: NULL,
      vaneDirection: NULL,
      temperature: NULL,
      pressure: NULL,
      humidity: NULL,
      gas: NULL,
      rainfall: NULL,
    };
    for (int i = 0; i < providers.size(); i++) {
      providers.at(i)->recordWeather(&report);
    }
    logger.post(&report);
  }
}
