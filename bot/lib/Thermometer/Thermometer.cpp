#ifndef THERMOMETER_I
#define THERMOMETER_I
#include "Thermometer.h"
#include <Arduino.h>

// the value of the 'other' resistor
#define SERIESRESISTOR 100000    

// resistance at 25 degrees C
#define THERMISTORNOMINAL 100000      
// temp. for nominal resistance (almost always 25 C)
#define TEMPERATURENOMINAL 25  
// The beta coefficient of the thermistor (usually 3000-4000)
#define BCOEFFICIENT 3950
// Samples to take per reading
#define NUMSAMPLES 10

Thermometer::Thermometer(int pin)
{
  this->pin = pin;
}

void Thermometer::setCalibrationFactor(double calibrationFactor)
{
  this->calibrationFactor = calibrationFactor;
}

double takeReading(int pin)
{
  int samples[NUMSAMPLES];
  int total;
  for (int i=0; i<NUMSAMPLES; i++) {
    total += analogRead(pin);
    delay(10);
  }
  return ((float)total/(float)NUMSAMPLES);
}

double calculateTempCelsius(double reading)
{
  reading = (4095 / reading)  - 1;     // (1023/ADC - 1) 
  reading = SERIESRESISTOR / reading;  // 100K / (1023/ADC - 1)
  double steinhart;
  steinhart = reading / THERMISTORNOMINAL;     // (R/Ro)
  steinhart = log(steinhart);                  // ln(R/Ro)
  steinhart /= BCOEFFICIENT;                   // 1/B * ln(R/Ro)
  steinhart += 1.0 / (TEMPERATURENOMINAL + 273.15); // + (1/To)
  steinhart = 1.0 / steinhart;                 // Invert
  steinhart -= 273.15; 
  return steinhart;
}

double convertToFahrenheit(double temp)
{
  return 1.8 * temp + 32;
}

double Thermometer::readTemperature()
{
  double r = takeReading(this->pin);
  double c = calculateTempCelsius(r);
  double f1 = convertToFahrenheit(c);
  return f1 + calibrationFactor;
}

#endif
