package seesaw

import (
	"errors"
	"time"
	"trelligo/pkg/debug"
	"trelligo/pkg/shims/ufmt"
)

const DefaultSeesawAddress = 0x49

const defaultDelay = 250 * time.Microsecond

const (
	seesawHwIdCodeSAMD09  = 0x55 // HW ID code for SAMD09
	seesawHwIdCodeTINY8x7 = 0x87 // HW ID code for ATtiny817
)

type Device struct {
	bus  I2C
	addr uint16
	hwid byte
}

func New(addr uint16, bus I2C) *Device {
	return &Device{
		bus:  bus,
		addr: addr,
	}
}

func (d *Device) Begin() error {

	debug.Log("soft reset")
	err := d.SoftReset()
	if err != nil {
		return err
	}

	debug.Log("wait for hwid")
	var lastErr error
	for i := 0; i < 10; i++ {
		hwid, err := d.readHwId()
		if err != nil {
			d.hwid = hwid
			lastErr = nil
			break
		}
		debug.Log("fail: " + err.Error())
		lastErr = err
		time.Sleep(10 * time.Millisecond)
	}

	if lastErr != nil {
		return lastErr
	}

	return nil
}

func (d *Device) readHwId() (byte, error) {
	hwid, err := d.ReadRegister(SEESAW_STATUS_BASE, SEESAW_STATUS_HW_ID)
	if err != nil {
		return 0, err
	}

	if hwid == seesawHwIdCodeSAMD09 || hwid == seesawHwIdCodeTINY8x7 {
		return hwid, nil
	}

	return 0, errors.New("unknown hardware ID: " + ufmt.ByteToHexString(hwid))
}

func (d *Device) SoftReset() error {
	return d.WriteRegister(SEESAW_STATUS_BASE, SEESAW_STATUS_SWRST, 0xFF)
}

func (d *Device) WriteRegister(module ModuleBaseAddress, function FunctionAddress, value byte) error {
	buf := []byte{byte(module), byte(function), value}
	return d.bus.Tx(d.addr, buf, nil)
}

// ReadRegister reads a single register from seesaw
func (d *Device) ReadRegister(module ModuleBaseAddress, function FunctionAddress) (byte, error) {
	buf := make([]byte, 1)
	err := d.Read(module, function, buf, defaultDelay)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

// Read reads a number of bytes from the device after sending the read command and waiting 'delay'. The delays depend
// on the module and function and are documented in the seesaw datasheet
func (d *Device) Read(module ModuleBaseAddress, function FunctionAddress, buf []byte, delay time.Duration) error {
	prefix := []byte{byte(module), byte(function)}
	err := d.bus.Tx(d.addr, prefix, nil)
	if err != nil {
		return err
	}

	//see seesaw datasheet
	time.Sleep(delay)

	return d.bus.Tx(d.addr, nil, buf)
}

func (d *Device) Write(module ModuleBaseAddress, function FunctionAddress, buf []byte) error {
	prefix := []byte{byte(module), byte(function)}
	data := append(prefix, buf...)
	return d.bus.Tx(d.addr, data, nil)
}
