#include "Logger.h"
#include <ArduinoJson.h>
#include <HTTPClient.h>

#ifndef LOGGER_IMPL
#define LOGGER_IMPL

Logger::Logger(String _url)
{
  url = _url;
}

bool Logger::post(LoggerData data)
{
  DynamicJsonDocument doc(1024);
  doc["temperature"] = data.temperature;
  doc["gas"] = data.gas;
  doc["pressure"] = data.pressure;
  doc["relativeHumidity"] = data.humidity;
  doc["anemometerMin"] = data.anemometerStatsSet.min;
  doc["anemometerMax"] = data.anemometerStatsSet.max;
  doc["anemometerAverage"] = data.anemometerStatsSet.average;

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
