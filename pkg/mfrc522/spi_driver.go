package mfrc522

type Driver interface {
	WriteRegister(reg Register, tx []byte) error
	ReadRegister(reg Register, rx []byte) error
}

func NewSpiDriver(spi SPI) Driver {
	return &SpiDriver{spi: spi}
}

// SpiDriver implementing the interface documented in datasheet section 8.1.2
type SpiDriver struct {
	spi SPI
}

func (d *SpiDriver) WriteRegister(reg Register, tx []byte) error {

	cmd := byte(reg << 1) // address byte, MSB 0 to indicate write
	d.spi.Begin()
	err := d.spi.Tx([]byte{cmd}, nil)
	if err != nil {
		return err
	}

	err = d.spi.Tx(tx, nil)
	if err != nil {
		return err
	}
	d.spi.Commit()
	return nil
}

func (d *SpiDriver) ReadRegister(reg Register, rx []byte) error {

	cmd := 0b1000_0000 | byte(reg<<1) // address byte + MSB 1 to indicate read
	d.spi.Begin()
	err := d.spi.Tx([]byte{cmd}, nil)
	if err != nil {
		return err
	}

	err = d.spi.Tx(nil, rx)
	if err != nil {
		return err
	}
	d.spi.Commit()
	return nil
}
