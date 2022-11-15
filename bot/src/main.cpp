#include <Arduino.h>
#include <MessagingManager.h>
#include <Thermometer.h>

#define PERIPHERAL_NAME     "GrillBot"
#define SERVICE_UUID        "1f814c33-4191-45a8-948e-6fcc7f9c10e5"
#define CHARACTERISTIC_UUID "a3612fbb-7c00-4ab2-b925-425c4ef2a002"

MessagingManager mm(PERIPHERAL_NAME, SERVICE_UUID, CHARACTERISTIC_UUID);
Thermometer t0(A0);
Thermometer t1(A3);

void setup() {
  Serial.begin(9600);
  
  Serial.println("starting ble ...");
  mm.begin();
  Serial.println("started");
}

void loop() {
  double t0f = t0.readTemperature();
  double t1f = t1.readTemperature();

  Serial.println("Temps are: ");
  Serial.println(t0f);
  Serial.println(t1f);

  mm.reportTemperatures(t0f, t1f);
  
  delay(1000);
}
