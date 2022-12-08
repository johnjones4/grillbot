#include <Arduino.h>
#include <MessagingManager.h>
#include <Thermometer.h>
#include <Preferences.h>

#define PERIPHERAL_NAME     "GrillBot"
#define SERVICE_UUID        "1f814c33-4191-45a8-948e-6fcc7f9c10e5"
#define TEMP_CHARACTERISTIC_UUID  "a3612fbb-7c00-4ab2-b925-425c4ef2a002"
#define CALIB_CHARACTERISTIC_UUID "09222388-fd96-4194-822b-fa052786c130"

#define CALIBRATION_0_KEY "calibration_0"
#define CALIBRATION_1_KEY "calibration_1"

Preferences preferences;

MessagingManager mm(PERIPHERAL_NAME, SERVICE_UUID, TEMP_CHARACTERISTIC_UUID, CALIB_CHARACTERISTIC_UUID);
Thermometer t0(A10);
Thermometer t1(A13);

void setup() {
  Serial.begin(9600);

  Serial.println("loading preferences ...");
  preferences.begin("grillbot", false);

  double calibration0 = preferences.getDouble(CALIBRATION_0_KEY, 0);
  t0.setCalibrationFactor(calibration0);

  double calibration1 = preferences.getDouble(CALIBRATION_1_KEY, 0);
  t1.setCalibrationFactor(calibration1);
  Serial.printf("Stored calibrations are: %f, %f\n", calibration0, calibration1);

  Serial.println("loaded");

  Serial.println("starting ble ...");
  mm.begin();
  mm.setCalibrations(calibration0, calibration1);
  Serial.println("started");
}

void loop() {
  double t0f = t0.readTemperature();
  double t1f = t1.readTemperature();
  Serial.printf("Temps are: %f, %f\n", t0f, t1f);
  mm.reportTemperatures(t0f, t1f);

  calibration c = mm.getCalibrations();
  Serial.printf("BLE calibrations are: %f, %f\n", c.calibration0, c.calibration1);
  t0.setCalibrationFactor(c.calibration0);
  t1.setCalibrationFactor(c.calibration1);
  preferences.putDouble(CALIBRATION_0_KEY, c.calibration0);
  preferences.putDouble(CALIBRATION_1_KEY, c.calibration1);
  
  delay(1000);
}
