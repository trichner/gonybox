package mfrc522

import "machine"

type SPI interface {
	Begin()
	Commit()
	Tx(w []byte, r []byte) error
}

type SpiImpl struct {
	spi machine.SPI
	cs  machine.Pin
}

func NewSpi(spi machine.SPI, chipSelect machine.Pin) SPI {
	return &SpiImpl{
		spi: spi,
		cs:  chipSelect,
	}
}

func (s *SpiImpl) Begin() {
	s.cs.Low()
}

func (s *SpiImpl) Commit() {
	s.cs.High()
}

func (s *SpiImpl) Tx(w []byte, r []byte) error {
	err := s.spi.Tx(w, r)
	if err != nil {
		return err
	}
	return nil
}
