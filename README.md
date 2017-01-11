[![GoDoc](https://godoc.org/github.com/dgnorton/ads111x?status.svg)](https://godoc.org/github.com/dgnorton/ads111x)
# ads111x
Go package for interfacing with ADS111x 16-bit analog to digital converters from Adafruit.

**NOTE:** Never apply more than *VDD* volts to the analog inputs. See [the datasheet](https://cdn-shop.adafruit.com/datasheets/ads1115.pdf) for full details on device operation.

## Installation

```
go get github.com/dgnorton/ads111x
```
## Getting started
Here's an example program that opens the device, sets a couple config settings, and reads voltages from the device's analog inputs.
```golang
func main() {
	// Open the device.
	adc, err := ads111x.Open("/dev/i2c-1", ads111x.Addr48)
	if err != nil {
		log.Fatal(err)
	}
	// Set the mode to continuous acquisition.
	if err := adc.SetMode(ads111x.Continuous); err != nil {
		log.Fatal(err)
	}
	// Set the scale to +/- 6.144 volts. (again, see note above & datasheet)
	if err := adc.SetScale(ads111x.Scale_6_144V); err != nil {
		log.Fatal(err)
	}

	for {
		// Read the voltage across analog inputs 0 & 1.
		v, err := adc.ReadVolts(ads111x.AIN_0_1)
		if err != nil {
			log.Fatal(err)
		}
		// Print the voltage and take a short nap.
		fmt.Printf("%f v\n", v)
		time.Sleep(time.Second)
	}
}
```
## Compiling
To build for an RPi 2:
```
GOARM=7 GOARCH=arm GOOS=linux go build
```
## Compatibility
This library has only been tested on an RPi 2 running Linux.
