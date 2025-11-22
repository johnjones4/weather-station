#include "Logger.h"
#include <ArduinoJson.h>
#include <HTTPClient.h>
#include <../../include/secrets.h>

#ifndef LOGGER_IMPL
#define LOGGER_IMPL

Logger::Logger(String _url)
{
  url = _url;
}

bool Logger::post(WeatherReport *data)
{
  Serial.println("Sending weather");
  DynamicJsonDocument doc(1024);
  doc["source"] = SOURCE_NAME;
  if (data->temperature != NULL) {
    doc["temperature"] = *data->temperature;
  }
  if (data->gas != NULL) {
    doc["gas"] = *data->gas;
  }
  if (data->pressure != NULL) {
    doc["pressure"] = *data->pressure;
  }
  if (data->humidity != NULL) {
    doc["humidity"] = *data->humidity;
  }
  if (data->windSpeed != NULL) {
    doc["windSpeed"] = *data->windSpeed;
  }
  if (data->vaneDirection != NULL) {
    doc["vaneDirection"] = *data->vaneDirection;
  }
  if (data->rainfall != NULL) {
    doc["rainfall"] = *data->rainfall;
  }

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
