package ads111x

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

//func msb(n uint16) uint8 { return uint8((n & 0xFF00) >> 8) }
//func lsb(n uint16) uint8 { return uint8(n & 0x00FF) }

//func Test_toBytes(t *testing.T) {
//	b, err := toBytes(uint16(0x8583))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if b[1] != 0x85 || b[0] != 0x83 {
//		fmt.Printf("0x%x\n", bytesToUint16(b))
//		t.Fatalf("exp = 0x%x, got = 0x%x%x\n", DefaultConfig, b[1], b[0])
//	}
//}
//
//func Test_toUint16(t *testing.T) {
//	b := []byte{0x83, 0x85}
//	v, err := toUint16(b)
//	if err != nil {
//		t.Fatal(err)
//	} else if v != DefaultConfig {
//		t.Fatalf("exp = 0x%x, got = 0x%x", DefaultConfig, v)
//	}
//}

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

func Test_SetMode(t *testing.T) {
	adc := newTestADC()
	i2c := adc.i2c.(*mockI2C)
	i2c.expReg = ConfigReg
	copy(i2c.expDat, []byte{0x84, 0x83})
	if err := adc.SetMode(Continuous); err != nil {
		t.Fatal(err)
	}
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
	panic("not implemented")
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
	copy(buf, m.cfg)
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

func bytesToUint16(b []byte) uint16 {
	var v uint16
	if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &v); err != nil {
		panic(err)
	}
	return v
}
