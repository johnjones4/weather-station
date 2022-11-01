#include "WeatherLogger.h"
#include <ArduinoJson.h>
#include <HTTPClient.h>

#ifndef WEATHER_LOGGER_IMPL
#define WEATHER_LOGGER_IMPL

WeatherLogger::WeatherLogger(char* _url)
{
  url = _url;
}

bool WeatherLogger::postWeather(WeatherLoggerData data)
{
  DynamicJsonDocument doc(1024);
  doc["temperature"] = data.temperature;
  doc["gas"] = data.gas;
  doc["pressure"] = data.pressure;
  doc["relativeHumidity"] = data.humidity;
  doc["anemometerMin"]["rotationsPerSecond"] = data.anemometerStatsSet.min.rotationsPerSecond;
  doc["anemometerMin"]["metersPerSecond"] = data.anemometerStatsSet.min.metersPerSecond;
  doc["anemometerMax"]["rotationsPerSecond"] = data.anemometerStatsSet.max.rotationsPerSecond;
  doc["anemometerMax"]["metersPerSecond"] = data.anemometerStatsSet.max.metersPerSecond;
  doc["anemometerAverage"]["rotationsPerSecond"] = data.anemometerStatsSet.average.rotationsPerSecond;
  doc["anemometerAverage"]["metersPerSecond"] = data.anemometerStatsSet.average.metersPerSecond;

  String postData;
  serializeJson(doc, postData);

  HTTPClient http;
  http.begin(url);
  http.addHeader("Content-Type", "application/json");
  int httpResponseCode = http.POST(postData);
  Serial.println(postData);
  Serial.println(httpResponseCode);
  return httpResponseCode == 200;
}

#endif
