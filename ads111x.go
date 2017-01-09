package ads111x

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/exp/io/i2c"
)

// Resolution is the resolution of the ADC.
const Resolution = 1 << 16

// I2CAddress represents one of the possible I2C bus addresses.
type I2CAddress uint8

const (
	I2CA_48 = 0x48
	I2CA_49 = 0x49
	I2CA_4A = 0x4A
	I2CA_4B = 0x4B
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
	// S_Busy means a conversion is currently being performed.
	S_Busy = iota << Status_LSB
	// S_Idle means a conversion is not currently being performed.
	S_Idle
)

type AIN uint16

const (
	AIN_LSB  uint8  = 12
	AIN_Mask uint16 = ^uint16(7 << AIN_LSB)
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
	FS_Mask uint16 = ^uint16(7 << FS_LSB)
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
	Mode_Mask uint16 = ^uint16(1 << Mode_LSB)
)

const (
	// M_Continuous is used to set continuous conversion mode.
	M_Continuous Mode = iota << Mode_LSB
	// M_Single is used to set Power-down single-shot mode (default).
	M_Single
)

type DataRate uint16

const (
	DataRate_LSB  uint8  = 5
	DataRate_Mask uint16 = ^uint16(7 << DataRate_LSB)
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
	ComparatorMode_Mask uint16 = ^uint16(1 << ComparatorMode_LSB)
)

const (
	// CM_Traditional is used to set traditional comparator with histeresis (default).
	CM_Traditional ComparatorMode = iota << ComparatorMode_LSB
	// CM_Window is used to set window comparator mode.
	CM_Window
)

type ComparatorPolarity uint16

const (
	ComparatorPolarity_LSB  uint8  = 3
	ComparatorPolarity_Mask uint16 = ^uint16(1 << ComparatorPolarity_LSB)
)

const (
	// CP_ActiveLow is used to set polarity of ALERT/RDY pin to active low (default).
	CP_ActiveLow ComparatorPolarity = iota << ComparatorPolarity_LSB
	// CP_ActiveHigh is used to set polarity of ALERT/RDY pin to active high.
	CP_ActiveHigh
)

type ComparatorLatching uint16

const (
	ComparatorLatching_LSB  uint8  = 2
	ComparatorLatching_Mask uint16 = ^uint16(1 << ComparatorLatching_LSB)
)

const (
	// CL_Off is used to set the comparator to non-latching (default).
	CL_Off ComparatorLatching = iota << ComparatorLatching_LSB
	//CL_On is used to set the comparator to latching.
	CL_On
)

type ComparatorQueue uint16

const (
	ComparatorQueue_LSB  uint8  = 0
	ComparatorQueue_Mask uint16 = ^uint16(1 << ComparatorQueue_LSB)
)

const (
	// CQ_AfterOne is used to set number of successive conversions exceeding upper or lower
	// thresholds before asserting ALERT/RDY pin.
	CQ_AfterOne ComparatorQueue = iota << ComparatorQueue_LSB
	CQ_AfterTwo
	CQ_AfterFour
	CQ_Disable
)

//const unknownConfig = uint32(1 << 31)

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

// Config returns the device config.
func (adc *ADC) Config() (uint16, error) {
	if adc.open {
		return adc.config, nil
	}

	var buf [2]byte
	if err := adc.ReadReg(ConfigReg, buf[:]); err != nil {
		return 0, err
	}

	config, err := toUint16(buf[:])
	if err != nil {
		return 0, err
	}

	adc.config = config

	return config, nil
}

// WriteConfig writes a new config to the device.
func (adc *ADC) WriteConfig(cfg uint16) error {
	return adc.WriteReg(ConfigReg, cfg)
}

// ReadVolts reads the voltage from the specified input.
func (adc *ADC) ReadVolts(input AIN) (float64, error) {
	cnt, err := adc.Read(input)
	if err != nil {
		return 0, err
	}

	fsrange := FSRange(FS(adc.config & FS_Mask))
	voltsPerCnt := fsrange / Resolution

	return float64(cnt) * voltsPerCnt, nil
}

// Read reads the value from the specified input.
func (adc *ADC) Read(input AIN) (uint16, error) {
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

// ReadReg reads a register.
func (adc *ADC) ReadReg(reg byte, buf []byte) error {
	return adc.i2c.ReadReg(reg, buf)
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
