#include "Display.h"
#include <Adafruit_GFX.h>

#define TEXT_BUFFER_SIZE 64

Display::Display()
{
  display = new Adafruit_SharpMem(SCK, MOSI, SS, 96, 96);
};

void Display::begin()
{
  display->begin();

  display->clearDisplay();
  display->setTextSize(1);
  display->setTextColor(0);
  display->setCursor(0,0);
  display->cp437(true);
  // char* buff = "Starting Up";
  // printLine(buff, 16);
  display->refresh();
}

void zerobuff(char* buff, int size)
{
  for (int i = 0; i < size; i++)
  {
    buff[i] = '\0';
  }
}

void Display::printLine(char* buff, int max)
{
  for (int i = 0; i < max; i++)
  {
    if (i == '\0')
    {
      break;
    }
    display->write(buff[i]);
  }
}

void Display::reportData(double t1, double c1, double t2, double c2)
{
  display->clearDisplay();;
  display->setCursor(0,0);

  char buff[TEXT_BUFFER_SIZE];
  
  zerobuff(buff, TEXT_BUFFER_SIZE);
  sprintf(buff, "%0.2f", t1);
  printLine(buff, TEXT_BUFFER_SIZE);

  display->setCursor(0,1);

  zerobuff(buff, TEXT_BUFFER_SIZE);
  sprintf(buff, "%0.2f", t2);
  printLine(buff, TEXT_BUFFER_SIZE);

  display->refresh();
}
