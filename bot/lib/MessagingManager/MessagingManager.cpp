#ifndef MESSAGING_MANAGER_I
#define MESSAGING_MANAGER_I
#include "MessagingManager.h"
#include <BLEDevice.h>
#include <BLEUtils.h>
#include <BLEScan.h>
#include <BLEAdvertisedDevice.h>

#define DOUBLE_LENGTH 8
#define MESSAGE_LENGTH 16

MessagingManager::MessagingManager(std::string _deviceName, std::string _serviceUuid, std::string _tempCharacteristicUuid, std::string _calibCharacteristicUuid)
{
  deviceName = _deviceName;
  serviceUuid = _serviceUuid;
  tempCharacteristicUuid = _tempCharacteristicUuid;
  calibCharacteristicUuid = _calibCharacteristicUuid;
}

void MessagingManager::begin()
{
  BLEDevice::init(deviceName);
  pServer = BLEDevice::createServer();
  pService = pServer->createService(serviceUuid);
  tempCharacteristic = pService->createCharacteristic(
                                         tempCharacteristicUuid,
                                         BLECharacteristic::PROPERTY_READ
                                       );
  calibCharacteristic = pService->createCharacteristic(
                                         calibCharacteristicUuid,
                                         BLECharacteristic::PROPERTY_READ | BLECharacteristic::PROPERTY_WRITE
                                       );
  pService->start();
  BLEAdvertising *pAdvertising = BLEDevice::getAdvertising();
  pAdvertising->addServiceUUID(serviceUuid);
  BLEDevice::startAdvertising();
}

void doubleToBytes(double val, byte* buffer)
{
  union {
    double value;
    char array[DOUBLE_LENGTH];
  } u;
  u.value = val;
  memcpy(buffer, u.array, DOUBLE_LENGTH);
}

void generateMessage(double temp1, double temp2, byte* buffer)
{
  doubleToBytes(temp1, buffer);
  byte secondBuffer[DOUBLE_LENGTH];
  doubleToBytes(temp2, secondBuffer);
  for (int i = 0; i < DOUBLE_LENGTH; i++) {
    buffer[8 + i] = secondBuffer[i];
  }
}

bool isNoise(double t)
{
  return t < 10 || t > 800;
}

void MessagingManager::reportTemperatures(double temp1, double temp2)
{
  if (isNoise(temp1) || isNoise(temp2)) {
    Serial.println("Readings are noisy");
    return;
  }
  byte buffer[MESSAGE_LENGTH];
  generateMessage(temp1, temp2, buffer);
  tempCharacteristic->setValue(buffer, MESSAGE_LENGTH);
  char* pHex = BLEUtils::buildHexData(nullptr, buffer, MESSAGE_LENGTH);
  Serial.printf("Temp Value: %s\n", pHex);
}

void MessagingManager::setCalibrations(double temp1, double temp2)
{
  byte buffer[MESSAGE_LENGTH];
  generateMessage(temp1, temp2, buffer);
  calibCharacteristic->setValue(buffer, MESSAGE_LENGTH);
  char* pHex = BLEUtils::buildHexData(nullptr, buffer, MESSAGE_LENGTH);
  Serial.printf("Calibration Value: %s\n", pHex);
}

calibration MessagingManager::getCalibrations()
{
  std::string data = calibCharacteristic->getValue();
  byte buffer[MESSAGE_LENGTH];
  memcpy(buffer, data.c_str(), MESSAGE_LENGTH);

  calibration c;

  for (int i = 0; i < 2; i++)
  {
    union {
      double value;
      char array[DOUBLE_LENGTH];
    } u;
    for (int j = 0; j < DOUBLE_LENGTH; j++) {
      int index = (i * DOUBLE_LENGTH) + j;
      u.array[j] = buffer[index];
    }
    if (i == 0) {
      c.calibration0 = u.value;
    } else if (i == 1) {
      c.calibration1 = u.value;
    }
  }

  Serial.printf("Calibrations: %d, %d\n", c.calibration0, c.calibration1);

  return c;
}

#endif
