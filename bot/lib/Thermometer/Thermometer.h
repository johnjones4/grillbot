#ifndef THERMOMETER_H
#define THERMOMETER_H

class Thermometer
{
private:
  int pin;
  double calibrationFactor = 0;
public:
  Thermometer(int _pin);
  void setCalibrationFactor(double calibrationFactor);
  double readTemperature();
};

#endif
