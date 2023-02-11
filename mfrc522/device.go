package mfrc522

import (
	"bytes"
	"errors"
	"time"
)

var (
	ErrCommunication   = errors.New("communication error")
	ErrCollision       = errors.New("collision error")
	ErrTimeout         = errors.New("timeout error")
	ErrNoRoom          = errors.New("buffer has not enough room error")
	ErrInternal        = errors.New("internal error")
	ErrIllegalArgument = errors.New("illegal argument error")
	ErrBadCrc          = errors.New("bad crc error")
)

type Device struct {
	driver Driver
}

func NewDevice(d Driver) *Device {
	return &Device{driver: d}
}
func (d *Device) Init() error {
	// driver init

	d.SoftReset()

	// Reset baud rates
	if err := d.driver.WriteRegister(TxModeReg, []byte{0}); err != nil {
		return err
	}
	if err := d.driver.WriteRegister(RxModeReg, []byte{0}); err != nil {
		return err
	}

	// Reset ModWidthReg
	if err := d.driver.WriteRegister(ModWidthReg, []byte{0x26}); err != nil {
		return err
	}

	if err := d.initTimeout(); err != nil {
		return err
	}

	// Default 0x00. Force a 100 % ASK modulation independent of the ModGsPReg register setting
	if err := d.driver.WriteRegister(TxASKReg, []byte{0x40}); err != nil {
		return err
	}

	// Default 0x3F. Set the preset value for the CRC coprocessor for the CalcCRC command to 0x6363 (ISO 14443-3 part 6.2.4)
	if err := d.driver.WriteRegister(ModeReg, []byte{0x3D}); err != nil {
		return err
	}

	if err := d.AntennaOn(); err != nil {
		return err
	}

	// give it some time to start
	time.Sleep(4 * time.Millisecond)

	version, err := d.GetVersion()
	if err != nil {
		return err
	}

	if version == VersionUnknown {
		return errors.New("invalid version")
	}

	return nil
}
func (d *Device) initTimeout() error {
	// When communicating with a PICC we need a timeout if something goes wrong.
	// f_timer = 13.56 MHz / (2*TPreScaler+1) where TPreScaler = [TPrescaler_Hi:TPrescaler_Lo].
	// TPrescaler_Hi are the four low bits in TModeReg. TPrescaler_Lo is TPrescalerReg.

	// TAuto=1; timer starts automatically at the end of the transmission in all communication modes at all speeds
	if err := d.driver.WriteRegister(TModeReg, []byte{0x80}); err != nil {
		return err
	}
	// TPreScaler = TModeReg[3..0]:TPrescalerReg, ie 0x0A9 = 169 => f_timer=40kHz, ie a timer period of 25μs.
	if err := d.driver.WriteRegister(TPrescalerReg, []byte{0xA9}); err != nil {
		return err
	}

	// Reload timer with 0x3E8 = 1000, ie 25ms before timeout.
	if err := d.driver.WriteRegister(TReloadRegH, []byte{0x03}); err != nil {
		return err
	}
	if err := d.driver.WriteRegister(TReloadRegL, []byte{0xE8}); err != nil {
		return err
	}

	return nil
}

func (d *Device) GetVersion() (byte, error) {
	b, err := d.readSingleRegister(VersionReg)
	if err != nil {
		return VersionUnknown, err
	}
	i := bytes.IndexByte(knownVersions, b)
	if i < 0 {
		return VersionUnknown, errors.New("unknown error")
	}
	return b, nil
}

// AntennaOn turns the antenna on by enabling pins TX1 and TX2.
// After a reset these pins are disabled.
func (d *Device) AntennaOn() error {
	b, err := d.readSingleRegister(TxControlReg)
	if err != nil {
		return err
	}
	if (b & 0x03) == 0 {
		return d.driver.WriteRegister(TxControlReg, []byte{b | 0x03})
	}
	return nil
}

func (d *Device) IsNewCardPresent() bool {

	//reset baud rates
	d.driver.WriteRegister(TxModeReg, []byte{0x00})
	d.driver.WriteRegister(RxModeReg, []byte{0x00})

	// Reset ModWidthReg
	d.driver.WriteRegister(ModWidthReg, []byte{0x26})

	buf := make([]byte, 2)
	err := d.reqestA(buf)
	return err == nil || errors.Is(err, ErrCollision)
}

func (d *Device) SoftReset() error {
	if err := d.driver.WriteRegister(CommandReg, CommandSoftReset.ToSlice()); err != nil {
		return err
	}

	// The datasheet does not mention how long the SoftRest command takes to complete.
	// But the MFRC522 might have been in soft power-down mode (triggered by bit 4 of CommandReg) .
	// Section 8.8.2 in the datasheet says the oscillator start-up time is the start up time of the crystal + 37,74μs. Let us be generous: 50ms.
	retries := 3

	// Wait for the PowerDown bit in CommandReg to be cleared
	time.Sleep(50 * time.Millisecond)
	for {
		pd, err := d.readPowerDownBit()
		retries--
		if err != nil && pd == false {
			return nil
		}
		if retries == 0 {
			return errors.New("timeout doing soft reset")
		}
	}
}

func (d *Device) readPowerDownBit() (bool, error) {
	b, err := d.readSingleRegister(CommandReg)
	return (b & (byte(1) << 4)) != 0, err
}

func (d *Device) writeSingleRegister(reg Register, b byte) error {
	return d.driver.WriteRegister(reg, []byte{b})
}
func (d *Device) readSingleRegister(reg Register) (byte, error) {
	buf := make([]byte, 1)
	if err := d.driver.ReadRegister(reg, buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (d *Device) clearRegisterBitMask(reg Register, mask byte) error {

	b, err := d.readSingleRegister(reg)
	if err != nil {
		return err
	}
	return d.writeSingleRegister(reg, b&(^mask))
}

func (d *Device) setRegisterBitMask(reg Register, mask byte) error {

	b, err := d.readSingleRegister(reg)
	if err != nil {
		return err
	}
	return d.writeSingleRegister(reg, b|mask)
}

func (d *Device) waitForCommandCompletion(waitIRqBits byte) error {
	// Wait for the command to complete.
	// In PCD_Init() we set the TAuto flag in TModeReg. This means the timer automatically starts when the PCD stops transmitting.
	// Each iteration of the do-while-loop takes 17.86μs.
	// TODO check/modify for other architectures than Arduino Uno 16bit
	i := 0
	for i = 30; i > 0; i-- {
		r, err := d.readSingleRegister(ComIrqReg) // ComIrqReg[7..0] bits are: Set1 TxIRq RxIRq IdleIRq HiAlertIRq LoAlertIRq ErrIRq TimerIRq
		if err != nil {
			return err
		}
		if r&waitIRqBits != 0 { // One of the interrupts that signal success has been set.
			break
		}
		if r&0x01 != 0 { // Timer interrupt - nothing received in 25ms
			return errors.New("timeout: timer interrupt")
		}
		time.Sleep(1 * time.Millisecond)
	}
	// ~30ms nothing happened. Communication with the MFRC522 might be down.
	if i == 0 {
		return errors.New("timeout: no reply")
	}

	return nil
}
