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

// DefaultConfig is the configuration the device boots with.
const DefaultConfig = uint16(0x8583)

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

type Scale uint16

const (
	Scale_LSB  uint8  = 9
	Scale_Mask uint16 = uint16(7 << Scale_LSB)
)

const (
	// Scale_6_144V is used to set full scale range to +/- 6.144V.
	Scale_6_144V Scale = iota << Scale_LSB
	// Scale_4_096V is used to set full scale range to +/- 4.096V.
	Scale_4_096V
	// Scale_2_048V is used to set full scale range to +/- 2.048V (default).
	Scale_2_048V
	// Scale_1_024V is used to set full scale range to +/- 1.024V.
	Scale_1_024V
	// Scale_0_512V is used to set full scale range to +/- 0.512V.
	Scale_0_512V
	// Scale_0_256V is used to set full scale range to +/- 0.256V.
	Scale_0_256V
)

// ScaleMinMax returns the min and max voltages for the given full scale value.
func ScaleMinMax(fs Scale) (min, max float64) {
	switch fs {
	case Scale_6_144V:
		return -6.144, 6.144
	case Scale_4_096V:
		return -4.096, 4.096
	case Scale_2_048V:
		return -2.048, 2.048
	case Scale_1_024V:
		return -1.024, 1.024
	case Scale_0_512V:
		return -0.512, 0.512
	case Scale_0_256V:
		return -0.256, 0.256
	default:
		panic("invalid fs value")
	}
}

// ScaleRange returns the difference between max and min for the full scale value.
func ScaleRange(fs Scale) float64 {
	min, max := ScaleMinMax(fs)
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

type i2cdevice interface {
	Close() error
	Read([]byte) error
	ReadReg(reg byte, buf []byte) error
	Write(buf []byte) (err error)
	WriteReg(reg byte, buf []byte) (err error)
}

// ADC represents an ADS1113, ADS1114, or ADS1115 analog to digital converter.
type ADC struct {
	i2c i2cdevice
}

// Open returns a new ADC initialized and ready for use.
// dev is the I2C bus device, e.g., /dev/i2c-1
func Open(dev string, addr I2CAddress) (*ADC, error) {
	d, err := i2c.Open(&i2c.Devfs{Dev: dev}, int(addr))
	if err != nil {
		return nil, err
	}

	return &ADC{
		i2c: d,
	}, nil
}

// Close closes the ADC connection.
func (adc *ADC) Close() error {
	return adc.i2c.Close()
}

// Mode returns the mode config setting.
func (adc *ADC) Mode() (Mode, error) {
	cfg, err := adc.Config()
	if err != nil {
		return Continuous, err
	}
	return Mode(cfg & Mode_Mask), nil
}

// SetMode sets the mode of operation (continuous or single).
func (adc *ADC) SetMode(m Mode) error {
	cfg, err := adc.Config()
	if err != nil {
		return err
	}
	cfg &= ^Mode_Mask
	cfg |= uint16(m)
	return adc.WriteConfig(cfg)
}

// Scale returns the full scale config setting.
func (adc *ADC) Scale() (Scale, error) {
	cfg, err := adc.Config()
	if err != nil {
		return Scale_0_256V, err
	}
	return Scale(cfg & Scale_Mask), nil
}

// SetScale sets the full scale range.
func (adc *ADC) SetScale(fs Scale) error {
	cfg, err := adc.Config()
	if err != nil {
		return err
	}
	cfg &= ^Scale_Mask
	cfg |= uint16(fs)
	return adc.WriteConfig(cfg)
}

// SetDataRate sets the number of samples per second.
func (adc *ADC) SetDataRate(dr DataRate) error {
	cfg, err := adc.Config()
	if err != nil {
		return err
	}
	cfg &= ^DataRate_Mask
	cfg |= uint16(dr)
	return adc.WriteConfig(cfg)
}

// Config returns the device config.
func (adc *ADC) Config() (uint16, error) {
	buf := make([]byte, 2)
	if err := adc.ReadReg(ConfigReg, buf); err != nil {
		return 0, err
	}

	var cfg uint16
	if err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &cfg); err != nil {
		return 0, err
	}

	return cfg, nil
}

// WriteConfig writes a new config to the device.
func (adc *ADC) WriteConfig(cfg uint16) error {
	return adc.WriteReg(ConfigReg, cfg)
}

// ReadVolts reads the voltage from the specified input.
func (adc *ADC) ReadVolts(input AIN) (float64, error) {
	cfg, err := adc.Config()
	if err != nil {
		return 0.0, err
	}
	cnt, err := adc.ReadAIN(input)
	if err != nil {
		return 0, err
	}

	fsrange := ScaleRange(Scale(cfg & Scale_Mask))
	voltsPerCnt := fsrange / Resolution

	return float64(cnt) * voltsPerCnt, nil
}

// ReadAIN reads the value from the specified input.
func (adc *ADC) ReadAIN(input AIN) (uint16, error) {
	cfg, err := adc.Config()
	if err != nil {
		return 0, err
	}
	// If the input isn't currently selected, select it.
	currentInput := AIN(cfg & AIN_Mask)
	if input != currentInput {
		// Clear input select bits.
		newConfig := cfg & ^AIN_Mask
		// Set new input select bits.
		newConfig |= uint16(input)
		// Write new config.
		if err := adc.WriteConfig(newConfig); err != nil {
			return 0, err
		}
	}

	// Read value from the conversion register.
	buf := make([]byte, 2)
	if err := adc.ReadReg(ConversionReg, buf); err != nil {
		return 0, err
	}

	var n uint16
	if err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &n); err != nil {
		return 0, err
	}

	return n, nil
}

// Read reads from the device.
func (adc *ADC) Read(buf []byte) error {
	return adc.i2c.Read(buf)
}

// ReadReg reads a register.
func (adc *ADC) ReadReg(reg byte, buf []byte) error {
	if err := adc.i2c.ReadReg(reg, buf); err != nil {
		return err
	}
	//fmt.Printf("ReadReg(0x%x) = {0x%x, 0x%x}\n", reg, buf[0], buf[1])
	return nil
}

// Write writes bytes to the device.
func (adc *ADC) Write(buf []byte) error {
	return adc.i2c.Write(buf)
}

// WriteReg writes a value to a register on the device.
func (adc *ADC) WriteReg(reg byte, data uint16) error {
	b, err := toBytes(data)
	if err != nil {
		return err
	}
	//fmt.Printf("WriteReg(0x%x, {0x%x, 0x%x})\n", reg, b[0], b[1])
	//println(hex.Dump(b))
	return adc.i2c.WriteReg(reg, b)
}

// toBytes converts data to a []byte suitable to send to the device.
func toBytes(data uint16) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, data)
	return buf.Bytes(), err
}
