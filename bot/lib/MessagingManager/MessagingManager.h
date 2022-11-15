#ifndef MESSAGING_MANAGER_H
#define MESSAGING_MANAGER_H
#include <Arduino.h>
#include <BLEDevice.h>
#include <BLEUtils.h>
#include <BLEScan.h>
#include <BLEAdvertisedDevice.h>

class MessagingManager
{
private:
  BLEServer *pServer;
  BLEService *pService;
  BLECharacteristic *pCharacteristic;
  std::string deviceName;
  std::string serviceUuid;
  std::string characteristicUuid;
public:
  MessagingManager(std::string _deviceName, std::string _serviceUuid, std::string _characteristicUuid);
  void begin();
  void reportTemperatures(double temp1, double temp2);
};

#endif
