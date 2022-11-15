#ifndef MESSAGING_MANAGER_I
#define MESSAGING_MANAGER_I
#include "MessagingManager.h"
#include <BLEDevice.h>
#include <BLEUtils.h>
#include <BLEScan.h>
#include <BLEAdvertisedDevice.h>

#define DOUBLE_LENGTH 8
#define MESSAGE_LENGTH 16

MessagingManager::MessagingManager(std::string _deviceName, std::string _serviceUuid, std::string _characteristicUuid)
{
  deviceName = _deviceName;
  serviceUuid = _serviceUuid;
  characteristicUuid = _characteristicUuid;
}

void MessagingManager::begin()
{
  BLEDevice::init(deviceName);
  pServer = BLEDevice::createServer();
  pService = pServer->createService(serviceUuid);
  pCharacteristic = pService->createCharacteristic(
                                         characteristicUuid,
                                         BLECharacteristic::PROPERTY_READ | 
                                         BLECharacteristic::PROPERTY_NOTIFY
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

void MessagingManager::reportTemperatures(double temp1, double temp2)
{
  byte buffer[MESSAGE_LENGTH];
  generateMessage(temp1, temp1, buffer);
  pCharacteristic->setValue(buffer, MESSAGE_LENGTH);
  pCharacteristic->notify(true);
}

#endif
