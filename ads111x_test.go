package ads111x

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"testing"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/exp/io/i2c/driver"
)

func Test_ReadRegUint16(t *testing.T) {
	adc := newTestADC()
	defer mustClose(adc)
	if cfg, err := adc.ReadRegUint16(ConfigReg); err != nil {
		t.Fatal(err)
	} else if cfg != DefaultConfig {
		t.Fatalf("exp = 0x%x, got 0x%x", DefaultConfig, cfg)
	}
}

func Test_WriteReg(t *testing.T) {
	if err := newTestADC().WriteReg(ConfigReg, DefaultConfig); err != nil {
		t.Fatal(err)
	}
}

func Test_WriteConfig(t *testing.T) {
	if err := newTestADC().WriteConfig(DefaultConfig); err != nil {
		t.Fatal(err)
	}
}

func Test_Open(t *testing.T) {
	var expDev *i2c.Device = &i2c.Device{}
	var expErr error

	i2cOpen = func(o driver.Opener, addr int) (*i2c.Device, error) {
		return expDev, expErr
	}

	if d, err := Open("test", 0x48); err != nil {
		t.Fatal(err)
	} else if d == nil {
		t.Fatal("unexpected nil pointer")
	}

	expErr = errors.New("failed")
	if _, err := Open("test", 0x48); err != expErr {
		t.Fatal("expected error")
	}
}

func Test_Status(t *testing.T) {
	adc := newTestADC()
	defer mustClose(adc)

	if val, err := adc.Status(); err != nil {
		t.Fatal(err)
	} else if val != Idle {
		t.Fatalf("exp = %v, got = %v", Idle, val)
	}
}

func Test_ReadVolts(t *testing.T) {
	adc := newTestADC()
	defer mustClose(adc)
	i2c := adc.i2c.(*mockI2C)
	i2c.WriteRegFn = func(reg byte, b []byte) error {
		if reg != i2c.expReg {
			return fmt.Errorf("exp = %d, got = %d", i2c.expReg, reg)
		}

		if reg == ConfigReg {
			copy(i2c.cfg, b)
		}
		return nil
	}
	copy(i2c.expDat, []byte{0x40, 0x00})
	exp := 1.024
	if got, err := adc.ReadVolts(AIN_0_3); err != nil {
		t.Fatal(err)
	} else if got != exp {
		t.Fatalf("exp = %f, got = %f", exp, got)
	}
}

func Test_Mode(t *testing.T) {
	adc := newTestADC()
	defer mustClose(adc)

	// Read mode.
	if val, err := adc.Mode(); err != nil {
		t.Fatal(err)
	} else if val != Single {
		t.Fatalf("exp = %v, got = %v", Single, val)
	}
	// Write mode.
	i2c := adc.i2c.(*mockI2C)
	copy(i2c.expDat, []byte{0x84, 0x83})

	if err := adc.SetMode(Continuous); err != nil {
		t.Fatal(err)
	}
	if val, err := adc.Mode(); err != nil {
		t.Fatal(err)
	} else if val != Continuous {
		t.Fatalf("exp = %v, got = %v", Continuous, val)
	}
}

func Test_Scale(t *testing.T) {
	adc := newTestADC()
	defer mustClose(adc)

	// Read full range scale.
	if val, err := adc.Scale(); err != nil {
		t.Fatal(err)
	} else if val != Scale_2_048V {
		t.Fatalf("exp = %v, got = %v", Scale_2_048V, val)
	}
	// Write full range scale.
	i2c := adc.i2c.(*mockI2C)
	copy(i2c.expDat, []byte{0x81, 0x83})

	if err := adc.SetScale(Scale_6_144V); err != nil {
		t.Fatal(err)
	}
	if val, err := adc.Scale(); err != nil {
		t.Fatal(err)
	} else if val != Scale_6_144V {
		t.Fatalf("exp = %v, got = %v", Scale_6_144V, val)
	}
}

func Test_DataRate(t *testing.T) {
	adc := newTestADC()
	defer mustClose(adc)

	// Read data rate.
	if val, err := adc.DataRate(); err != nil {
		t.Fatal(err)
	} else if val != DR_128SPS {
		t.Fatalf("exp = %v, got = %v", DR_128SPS, val)
	}
	// Write data rate.
	i2c := adc.i2c.(*mockI2C)
	copy(i2c.expDat, []byte{0x85, 0xe3})

	if err := adc.SetDataRate(DR_860SPS); err != nil {
		t.Fatal(err)
	}
	if val, err := adc.DataRate(); err != nil {
		t.Fatal(err)
	} else if val != DR_860SPS {
		t.Fatalf("exp = %v, got = %v", DR_860SPS, val)
	}
}

func Test_ComparatorMode(t *testing.T) {
	adc := newTestADC()
	defer mustClose(adc)

	// Read comparator mode.
	if val, err := adc.ComparatorMode(); err != nil {
		t.Fatal(err)
	} else if val != Traditional {
		t.Fatalf("exp = %v, got = %v", Traditional, val)
	}

	// Write comparator mode.
	i2c := adc.i2c.(*mockI2C)
	copy(i2c.expDat, []byte{0x85, 0x93})

	if err := adc.SetComparatorMode(Window); err != nil {
		t.Fatal(err)
	}
	if val, err := adc.ComparatorMode(); err != nil {
		t.Fatal(err)
	} else if val != Window {
		t.Fatalf("exp = %v, got = %v", Window, val)
	}
}

func Test_ComparatorPolarity(t *testing.T) {
	adc := newTestADC()
	defer mustClose(adc)

	// Read comparator polarity.
	if val, err := adc.ComparatorPolarity(); err != nil {
		t.Fatal(err)
	} else if val != ActiveLow {
		t.Fatalf("exp = %v, got = %v", ActiveLow, val)
	}

	// Write comparator polarity.
	i2c := adc.i2c.(*mockI2C)
	copy(i2c.expDat, []byte{0x85, 0x8b})

	if err := adc.SetComparatorPolarity(ActiveHigh); err != nil {
		t.Fatal(err)
	}
	if val, err := adc.ComparatorPolarity(); err != nil {
		t.Fatal(err)
	} else if val != ActiveHigh {
		t.Fatalf("exp = %v, got = %v", ActiveHigh, val)
	}
}

func Test_ComparatorLatching(t *testing.T) {
	adc := newTestADC()
	defer mustClose(adc)

	// Read comparator latching.
	if val, err := adc.ComparatorLatching(); err != nil {
		t.Fatal(err)
	} else if val != Off {
		t.Fatalf("exp = %v, got = %v", Off, val)
	}

	// Write comparator latching.
	i2c := adc.i2c.(*mockI2C)
	copy(i2c.expDat, []byte{0x85, 0x87})

	if err := adc.SetComparatorLatching(On); err != nil {
		t.Fatal(err)
	}
	if val, err := adc.ComparatorLatching(); err != nil {
		t.Fatal(err)
	} else if val != On {
		t.Fatalf("exp = %v, got = %v", On, val)
	}
}

func Test_ComparatorQueue(t *testing.T) {
	adc := newTestADC()
	defer mustClose(adc)

	// Read comparator queuing mode.
	if val, err := adc.ComparatorQueue(); err != nil {
		t.Fatal(err)
	} else if val != Disable {
		t.Fatalf("exp = %v, got = %v", Disable, val)
	}

	// Write comparator queuing mode.
	i2c := adc.i2c.(*mockI2C)
	copy(i2c.expDat, []byte{0x85, 0x80})

	if err := adc.SetComparatorQueue(AfterOne); err != nil {
		t.Fatal(err)
	}
	if val, err := adc.ComparatorQueue(); err != nil {
		t.Fatal(err)
	} else if val != AfterOne {
		t.Fatalf("exp = %v, got = %v", AfterOne, val)
	}
}

func Test_ScaleMinMax(t *testing.T) {
	test := func(s Scale, expMin, expMax float64) {
		min, max := ScaleMinMax(s)
		if min != expMin || max != expMax {
			t.Fatalf("exp = %f and %f, got = %f and %f", expMin, expMax, min, max)
		}
	}

	test(Scale_6_144V, -6.144, 6.144)
	test(Scale_4_096V, -4.096, 4.096)
	test(Scale_2_048V, -2.048, 2.048)
	test(Scale_1_024V, -1.024, 1.024)
	test(Scale_0_512V, -0.512, 0.512)
	test(Scale_0_256V, -0.256, 0.256)
	test(Scale_0_256V, -0.256, 0.256)
}

func newTestADC() *ADC {
	return &ADC{
		i2c: &mockI2C{
			expReg: ConfigReg,
			expDat: []byte{0x85, 0x83},
			cfg:    []byte{0x85, 0x83},
		},
	}
}

type mockI2C struct {
	CloseFn    func() error
	ReadFn     func(buf []byte) error
	ReadRegFn  func(reg byte, buf []byte) error
	WriteFn    func(buf []byte) error
	WriteRegFn func(reg byte, buf []byte) error
	expReg     byte
	expDat     []byte
	cfg        []byte
}

func (m *mockI2C) Close() error {
	if m.CloseFn != nil {
		return m.CloseFn()
	}
	return nil
}

func (m *mockI2C) Read(buf []byte) error {
	if m.ReadFn != nil {
		return m.ReadFn(buf)
	}
	panic("not implemented")
}

func (m *mockI2C) ReadReg(reg byte, buf []byte) error {
	if m.ReadRegFn != nil {
		return m.ReadRegFn(reg, buf)
	}

	if reg == ConfigReg {
		copy(buf, m.cfg)
	} else if reg == ConversionReg {
		copy(buf, m.expDat)
	}

	return nil
}

func (m *mockI2C) Write(buf []byte) (err error) {
	if m.WriteFn != nil {
		return m.WriteFn(buf)
	}
	panic("not implemented")
}

func (m *mockI2C) WriteReg(reg byte, buf []byte) error {
	if m.WriteRegFn != nil {
		return m.WriteRegFn(reg, buf)
	}

	if reg != m.expReg {
		return fmt.Errorf("exp = %d, got = %d", m.expReg, reg)
	}

	if buf[0] != m.expDat[0] || buf[1] != m.expDat[1] {
		return fmt.Errorf("exp = 0x%x%x, got = 0x%x%x\n", m.expDat[0], m.expDat[1], buf[0], buf[1])
	}

	if reg == ConfigReg {
		copy(m.cfg, buf)
	}

	return nil
}

func mustClose(adc *ADC) {
	if err := adc.Close(); err != nil {
		panic(err)
	}
}

func bytesToUint16(b []byte) uint16 {
	var v uint16
	if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &v); err != nil {
		panic(err)
	}
	return v
}
