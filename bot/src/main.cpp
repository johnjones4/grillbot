#include <Arduino.h>
#include <BLEDevice.h>
#include <BLEUtils.h>
#include <BLEScan.h>
#include <BLEAdvertisedDevice.h>

#define PERIPHERAL_NAME     "GrillBot"
#define SERVICE_UUID        "1f814c33-4191-45a8-948e-6fcc7f9c10e5"
#define CHARACTERISTIC_UUID "a3612fbb-7c00-4ab2-b925-425c4ef2a002"

BLEServer *pServer;
BLEService *pService;
BLECharacteristic *pCharacteristic;

void setup() {
  Serial.begin(9600);

  Serial.println("starting ble ...");
  BLEDevice::init(PERIPHERAL_NAME);
  pServer = BLEDevice::createServer();
  pService = pServer->createService(SERVICE_UUID);
  pCharacteristic = pService->createCharacteristic(
                                         CHARACTERISTIC_UUID,
                                         BLECharacteristic::PROPERTY_READ | 
                                         BLECharacteristic::PROPERTY_NOTIFY
                                       );
  pService->start();
  BLEAdvertising *pAdvertising = BLEDevice::getAdvertising();
  pAdvertising->addServiceUUID(SERVICE_UUID);
  BLEDevice::startAdvertising();
  Serial.println("started");
}

void doubleToBytes(double val, byte* buffer){
  // Create union of shared memory space
  union {
    double value;
    char array[8];
  } u;
  // Overite bytes of union with float variable
  u.value = val;
  // Assign bytes to input array
  memcpy(buffer, u.array, 8);
}

void generateMessage(double valueA, double valueB, byte* buffer) {
  doubleToBytes(valueA, buffer);
  byte secondBuffer[8];
  doubleToBytes(valueB, secondBuffer);
  for (int i = 0; i < 8; i++) {
    buffer[8 + i] = secondBuffer[i];
  }
}

void loop() {
  double now = double(millis()) / 100000.0;
  double c1 = sin(now) * 100.0 + 100.0;
  double c2 = cos(now) * 100.0 + 100.0;
  byte buffer[16];
  generateMessage(c1, c2, buffer);
  pCharacteristic->setValue(buffer, 16);
  pCharacteristic->notify(true);
  delay(1000);
}