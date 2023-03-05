#include <Adafruit_SharpMem.h>

class Display
{
private:
  Adafruit_SharpMem *display;
  void printLine(char* buff, int max);
public:
  Display();
  void begin();
  void reportData(double t1, double c1, double t2, double c2);
};
