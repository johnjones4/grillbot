# Grillbot

A device for monitoring and logging thermistor readings from a grill or smoker. For more background on this, watch the video or read the blog post.

## Bot (Microcontroller Device)

The bot reads the thermistors via an ADC pin and emmits the two double values to the serial console in base64.

It's built using PlatformIO and designed to work on an ESP32 platform such as `esp32doit-devkit-v1` but any decent microcontroller platform that supports the Arduino toolchain, has two ADC pins, and serial communication should suffice.

## Logger (Desktop Software)

The logger reads the serial port connected to the bot and displays and logs the temperatue value coming from the bot. It has a few different configuations that can be passed in via command line flags. It is a Go project and also has a minimal web interface written in vanilla HTML/CSS/JS.

### Command Line Flags:

```
  -change-threshold float
    	Percent change threshold that should register before logging a new reading (default 0.05)
  -debug
    	Run with debug logging
  -file string
    	Resume a previous cook by passing in a cook file
  -food string
    	Food being prepared (ie brisket)
  -host string
    	Hostname to listen on for the web interface (default ":8080")
  -method string
    	Method of cooking (ie smoking) (default "smoked")
  -serial string
    	Serial device to use (ie /dev/ttyUSB0)
  -simulated
    	Use simulated data
  -time-threshold duration
    	Time threshold that should pass before logging a new reading (default 30s)
```
