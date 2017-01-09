package ads111x

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"golang.org/x/exp/io/i2c"
)

// Resolution is the resolution of the ADC.
const Resolution = 1 << 16

// I2CAddress represents one of the possible I2C bus addresses.
type I2CAddress uint8

const (
	Addr48 I2CAddress = 0x48
	Addr49            = 0x49
	Addr4A            = 0x4A
	Addr4B            = 0x4B
)

const (
	// ConversionReg is the address of the conversion register.
	ConversionReg byte = iota
	// ConfigReg is the address of the config register.
	ConfigReg
	// LoThresReg is the address of the Lo_thresh register.
	LoThreshReg
	// HiThreshReg is the address of the Hi_thresh register.
	HiThreshReg
)

type Status uint16

const (
	Status_LSB  uint8  = 15
	Status_Mask uint16 = ^uint16(1 << Status_LSB)
)

const (
	// Busy means a conversion is currently being performed.
	Busy = iota << Status_LSB
	// Idle means a conversion is not currently being performed.
	Idle
)

type AIN uint16

const (
	AIN_LSB  uint8  = 12
	AIN_Mask uint16 = uint16(7 << AIN_LSB)
)

const (
	// AIN_0_1 is used to select AIN0 (pos) & AIN1 (neg) inputs (default).
	AIN_0_1 AIN = iota << AIN_LSB
	// AIN_0_3 is used to select AIN0 (pos) & AIN3 (neg) inputs.
	AIN_0_3
	// AIN_1_3 is used to select AIN1 (pos) & AIN3 (neg) inputs.
	AIN_1_3
	// AIN_2_3 is used to select AIN2 (pos) & AIN3 (neg) inputs.
	AIN_2_3
	// AIN_0_GND is used to select AIN0 (pos) & GND (neg) inputs.
	AIN_0_GND
	// AIN_1_GND is used to select AIN1 (pos) & GND (neg) inputs.
	AIN_1_GND
	// AIN_2_GND is used to select AIN2 (pos) & GND (neg) inputs.
	AIN_2_GND
	// AIN_3_GND is used to select AIN3 (pos) & GND (neg) inputs.
	AIN_3_GND
)

type FS uint16

const (
	FS_LSB  uint8  = 9
	FS_Mask uint16 = uint16(7 << FS_LSB)
)

const (
	// FS_6_144V is used to set full scale range to +/- 6.144V.
	FS_6_144V FS = iota << FS_LSB
	// FS_4_096V is used to set full scale range to +/- 4.096V.
	FS_4_096V
	// FS_2_048V is used to set full scale range to +/- 2.048V (default).
	FS_2_048V
	// FS_1_024V is used to set full scale range to +/- 1.024V.
	FS_1_024V
	// FS_0_512V is used to set full scale range to +/- 0.512V.
	FS_0_512V
	// FS_0_256V is used to set full scale range to +/- 0.256V.
	FS_0_256V
)

// FSMinMax returns the min and max voltages for the given full scale value.
func FSMinMax(fs FS) (min, max float64) {
	switch fs {
	case FS_6_144V:
		return -6.144, 6.144
	case FS_4_096V:
		return -4.096, 4.096
	case FS_2_048V:
		return -2.048, 2.048
	case FS_1_024V:
		return -1.024, 1.024
	case FS_0_512V:
		return -0.512, 0.512
	case FS_0_256V:
		return -0.256, 0.256
	default:
		panic("invalid fs value")
	}
}

// FSRange returns the difference between max and min for the full scale value.
func FSRange(fs FS) float64 {
	min, max := FSMinMax(fs)
	return max - min
}

type Mode uint16

const (
	Mode_LSB  uint8  = 8
	Mode_Mask uint16 = uint16(1 << Mode_LSB)
)

const (
	// Continuous is used to set continuous conversion mode.
	Continuous Mode = iota << Mode_LSB
	// Single is used to set Power-down single-shot mode (default).
	Single
)

type DataRate uint16

const (
	DataRate_LSB  uint8  = 5
	DataRate_Mask uint16 = uint16(7 << DataRate_LSB)
)

const (
	// DR_8SPS is used to set the data rate to 8 samples per second.
	DR_8SPS DataRate = iota << DataRate_LSB
	// DR_16SPS is used to set the data rate to 16 samples per second.
	DR_16SPS
	// DR_32SPS is used to set the data rate to 32 samples per second.
	DR_32SPS
	// DR_64SPS is used to set the data rate to 64 samples per second.
	DR_64SPS
	// DR_128SPS is used to set the data rate to 128 samples per second (default).
	DR_128SPS
	// DR_250SPS is used to set the data rate to 250 samples per second.
	DR_250SPS
	// DR_475SPS is used to set the data rate to 475 samples per second.
	DR_475SPS
	// DR_860SPS is used to set the data rate to 860 samples per second.
	DR_860SPS
)

type ComparatorMode uint16

const (
	ComparatorMode_LSB  uint8  = 4
	ComparatorMode_Mask uint16 = uint16(1 << ComparatorMode_LSB)
)

const (
	// Traditional is used to set traditional comparator with histeresis (default).
	Traditional ComparatorMode = iota << ComparatorMode_LSB
	// Window is used to set window comparator mode.
	Window
)

type ComparatorPolarity uint16

const (
	ComparatorPolarity_LSB  uint8  = 3
	ComparatorPolarity_Mask uint16 = uint16(1 << ComparatorPolarity_LSB)
)

const (
	// ActiveLow is used to set polarity of ALERT/RDY pin to active low (default).
	ActiveLow ComparatorPolarity = iota << ComparatorPolarity_LSB
	// ActiveHigh is used to set polarity of ALERT/RDY pin to active high.
	ActiveHigh
)

type ComparatorLatching uint16

const (
	ComparatorLatching_LSB  uint8  = 2
	ComparatorLatching_Mask uint16 = uint16(1 << ComparatorLatching_LSB)
)

const (
	// Off is used to set the comparator to non-latching (default).
	Off ComparatorLatching = iota << ComparatorLatching_LSB
	//On is used to set the comparator to latching.
	On
)

type ComparatorQueue uint16

const (
	ComparatorQueue_LSB  uint8  = 0
	ComparatorQueue_Mask uint16 = uint16(1 << ComparatorQueue_LSB)
)

const (
	// CQ_AfterOne is used to set number of successive conversions exceeding upper or lower
	// thresholds before asserting ALERT/RDY pin.
	AfterOne ComparatorQueue = iota << ComparatorQueue_LSB
	AfterTwo
	AfterFour
	Disable
)

// ADC represents an ADS1113, ADS1114, or ADS1115 analog to digital converter.
type ADC struct {
	i2c    *i2c.Device
	open   bool
	config uint16 // current device config
}

// Open returns a new ADC initialized and ready for use.
// dev is the I2C bus device, e.g., /dev/i2c-1
func Open(dev string, addr I2CAddress) (*ADC, error) {
	d, err := i2c.Open(&i2c.Devfs{Dev: dev}, int(addr))
	if err != nil {
		return nil, err
	}

	adc := &ADC{
		i2c: d,
	}

	if _, err := adc.Config(); err != nil {
		return nil, err
	}

	adc.open = true

	return adc, nil
}

// Close closes the ADC connection.
func (adc *ADC) Close() error {
	return adc.i2c.Close()
}

// Mode returns the mode config setting.
func (adc *ADC) Mode() (Mode, error) {
	return Mode(adc.config & Mode_Mask), nil
}

// SetMode sets the mode of operation (continuous or single).
func (adc *ADC) SetMode(m Mode) error {
	cfg := adc.config & ^Mode_Mask
	cfg |= uint16(m)
	return adc.WriteConfig(cfg)
}

// FullScale returns the full scale config setting.
func (adc *ADC) FullScale() (FS, error) {
	return FS(adc.config & FS_Mask), nil
}

// SetFullScale sets the full scale range.
func (adc *ADC) SetFullScale(fs FS) error {
	cfg := adc.config & ^FS_Mask
	cfg |= uint16(fs)
	return adc.WriteConfig(cfg)
}

// Config returns the device config.
func (adc *ADC) Config() (uint16, error) {
	if adc.open {
		return adc.config, nil
	}

	if err := adc.WriteReg(ConfigReg, []byte{}); err != nil {
		return 0, err
	}

	var buf [2]byte
	if err := adc.Read(buf[:]); err != nil {
		return 0, err
	}

	println("Config() read...")
	println(hex.Dump(buf[:]))

	config, err := leUint16(buf[:])
	if err != nil {
		return 0, err
	}

	adc.config = config

	return config, nil
}

// WriteConfig writes a new config to the device.
func (adc *ADC) WriteConfig(cfg uint16) error {
	println("writing config...")
	println(cfg)
	if err := adc.WriteReg(ConfigReg, cfg); err != nil {
		return err
	}
	adc.config = cfg
	return nil
}

// ReadVolts reads the voltage from the specified input.
func (adc *ADC) ReadVolts(input AIN) (float64, error) {
	cnt, err := adc.ReadAIN(input)
	if err != nil {
		return 0, err
	}

	fmt.Printf("FS_Mask = %v\n", FS_Mask)
	fmt.Printf("adc.config = %v\n", adc.config)
	println(FS(adc.config & FS_Mask))

	fsrange := FSRange(FS(adc.config & FS_Mask))
	voltsPerCnt := fsrange / Resolution

	return float64(cnt) * voltsPerCnt, nil
}

// ReadAIN reads the value from the specified input.
func (adc *ADC) ReadAIN(input AIN) (uint16, error) {
	// If the input isn't currently selected, select it.
	currentInput := AIN(adc.config & AIN_Mask)
	if input != currentInput {
		// Clear input select bits.
		newConfig := adc.config & ^AIN_Mask
		// Set new input select bits.
		newConfig |= uint16(input)
		// Write new config.
		if err := adc.WriteConfig(newConfig); err != nil {
			return 0, err
		}
	}

	// Read value from the conversion register.
	var buf [2]byte
	if err := adc.ReadReg(ConversionReg, buf[:]); err != nil {
		return 0, err
	}

	val, err := toUint16(buf[:])
	if err != nil {
		return 0, err
	}

	return val, nil
}

// Read reads from the device.
func (adc *ADC) Read(buf []byte) error {
	return adc.i2c.Read(buf)
}

// ReadReg reads a register.
func (adc *ADC) ReadReg(reg byte, buf []byte) error {
	fmt.Printf("ReadReg(reg = %v)\n", reg)
	err := adc.i2c.ReadReg(reg, buf)
	println(hex.Dump(buf))
	return err
}

// Write writes bytes to the device.
func (adc *ADC) Write(buf []byte) error {
	return adc.i2c.Write(buf)
}

// WriteReg writes a value to a register on the device.
func (adc *ADC) WriteReg(reg byte, data interface{}) error {
	b, err := toBytes(data)
	if err != nil {
		return err
	}
	return adc.i2c.WriteReg(reg, b)
}

// toBytes converts data to a []byte suitable to send to the device.
func toBytes(data interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, data)
	return buf.Bytes(), err
}

// toUint16 converts []byte to a uint16.
func toUint16(b []byte) (uint16, error) {
	var v uint16
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &v)
	return v, err
}

func leUint16(b []byte) (uint16, error) {
	var v uint16
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &v)
	return v, err
}
