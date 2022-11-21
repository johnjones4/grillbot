#ifndef MESSAGING_MANAGER_H
#define MESSAGING_MANAGER_H
#include <Arduino.h>
#include <BLEDevice.h>
#include <BLEUtils.h>
#include <BLEScan.h>
#include <BLEAdvertisedDevice.h>

typedef struct {
  double calibration0;
  double calibration1;
} calibration;

class MessagingManager
{
private:
  BLEServer *pServer;
  BLEService *pService;
  BLECharacteristic *tempCharacteristic;
  BLECharacteristic *calibCharacteristic;
  std::string deviceName;
  std::string serviceUuid;
  std::string tempCharacteristicUuid;
  std::string calibCharacteristicUuid;
public:
  MessagingManager(std::string _deviceName, std::string _serviceUuid, std::string _tempCharacteristicUuid, std::string _calibCharacteristicUuid);
  void begin();
  void reportTemperatures(double temp1, double temp2);
  void setCalibrations(double temp1, double temp2);
  calibration getCalibrations();
};

#endif
