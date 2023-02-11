package mfrc522

import (
	"errors"
	"strconv"
)

type UID []byte

func (d *Device) PiccReadCardSerial() (UID, error) {

	//MFRC522::StatusCode result = PICC_Select(&uid);
	err := nil
	return nil, err
}

// https://github.com/OSSLibraries/Arduino_MFRC522v2/blob/95ac65b2d5f6a0bcb155f30d51add600f7d73000/src/MFRC522v2.cpp#L536
func (d *Device) PiccSelect() (UID, error) {

	// ValuesAfterColl=1 => Bits received after collision are cleared.
	if err := d.clearRegisterBitMask(CollReg, 0x80); err != nil {
		return nil, err
	}

	cascadeLevel := 1
	uuidComplete := false

	// The SELECT/ANTICOLLISION commands uses a 7 byte standard frame + 2 bytes CRC_A
	buffer := make([]byte, 9)
	// Description of buffer structure:
	//        Byte 0: SEL                 Indicates the Cascade Level: PICC_Command::PICC_CMD_SEL_CL1, PICC_Command::PICC_CMD_SEL_CL2 or PICC_Command::PICC_CMD_SEL_CL3
	//        Byte 1: NVB                 Number of Valid Bits (in complete command, not just the UID): High nibble: complete bytes, Low nibble: Extra bits.
	//        Byte 2: UID-data or CT      See explanation below. CT means Cascade Tag.
	//        Byte 3: UID-data
	//        Byte 4: UID-data
	//        Byte 5: UID-data
	//        Byte 6: BCC                 Block Check Character - XOR of bytes 2-5
	//        Byte 7: CRC_A
	//        Byte 8: CRC_A
	// The BCC and CRC_A are only transmitted if we know all the UID bits of the current Cascade Level.
	//
	// Description of bytes 2-5: (Section 6.5.4 of the ISO/IEC 14443-3 draft: UID contents and cascade levels)
	//        UID size    Cascade level   Byte2   Byte3   Byte4   Byte5
	//        ========    =============   =====   =====   =====   =====
	//         4 bytes        1           uid0    uid1    uid2    uid3
	//         7 bytes        1           CT      uid0    uid1    uid2
	//                        2           uid3    uid4    uid5    uid6
	//        10 bytes        1           CT      uid0    uid1    uid2
	//                        2           CT      uid3    uid4    uid5
	//                        3           uid6    uid7    uid8    uid9

	uidSize := 0
	validBits := 0
	useCascadeTag := false
	for !uuidComplete {

		switch cascadeLevel {
		case 1:
			buffer[0] = byte(PiccCommandSelCl1)
			useCascadeTag = validBits > 0 && uidSize > 4
		case 2:
		case 3:
		default:
			return nil, ErrInternal
		}

		//TODO
		_ = useCascadeTag
	}

	return nil, nil
}

/*
bool MFRC522::PICC_ReadCardSerial() {
  MFRC522::StatusCode result = PICC_Select(&uid);
  return (result == StatusCode::STATUS_OK);
} // End PICC_ReadCardSerial()
*/

func (d *Device) reqestA(rx []byte) error {
	return d.sendReqAOrWupa(PiccCommandReqA, rx)
}

func (d *Device) sendReqAOrWupa(cmd PiccCommand, res []byte) error {

	if len(res) != 2 {
		return errors.New("inalid response buffer length 2 != " + strconv.Itoa(len(res)))
	}

	// ValuesAfterColl=1 => Bits received after collision are cleared.
	if err := d.clearRegisterBitMask(CollReg, 0x80); err != nil {
		return err
	}

	txBits := uint8(7) // For REQA and WUPA we need the short frame format - transmit only 7 bits of the last (and only) byte. TxLastBits = BitFramingReg[2..0]

	n, rxBits, err := d.transceiveData(cmd.ToSlice(), res, txBits, 0, false)
	if err != nil {
		return err
	}
	if n != 2 || rxBits != 0 {
		return errors.New("invalid command response")
	}

	return nil
}

func (d *Device) transceiveData(tx []byte, rx []byte, txLastBits byte, rxAlign byte, checkCrc bool) (bytesRead int, rxLastBits int, err error) {

	waitIRq := byte(0x30) // RxIRq | IdleIRq
	return d.communicateWithPicc(CommandTransceive, tx, rx, communicateWithPiccOpts{
		WaitForIRqMask: waitIRq,
		TxLastBits:     txLastBits,
		RxAlign:        rxAlign,
		CheckCRC:       checkCrc,
	})
}

type communicateWithPiccOpts struct {
	WaitForIRqMask byte

	// 9.3.1.14 - BitFramingReg register
	RxAlign    byte
	TxLastBits byte

	CheckCRC bool
}

func (d *Device) communicateWithPicc(cmd Command, tx []byte, rx []byte, o communicateWithPiccOpts) (bytesRead int, rxLastBits int, err error) {

	bitFraming := ((0b111 & o.RxAlign) << 4) | (o.TxLastBits & 0b111)

	if err := d.writeSingleRegister(CommandReg, byte(CommandIdle)); err != nil { // Stop any active command.
		return 0, 0, err
	}
	if err := d.writeSingleRegister(ComIrqReg, 0x7F); err != nil { // Clear all seven interrupt request bits
		return 0, 0, err
	}
	if err := d.writeSingleRegister(FIFOLevelReg, 0x80); err != nil { // FlushBuffer = 1, FIFO initialization
		return 0, 0, err
	}
	if err := d.driver.WriteRegister(FIFODataReg, tx); err != nil { // Write sendData to the FIFO
		return 0, 0, err
	}
	if err := d.writeSingleRegister(BitFramingReg, bitFraming); err != nil { // Bit adjustments
		return 0, 0, err
	}
	if err := d.writeSingleRegister(CommandReg, byte(cmd)); err != nil {
		return 0, 0, err
	}

	if cmd == CommandTransceive {
		if err := d.setRegisterBitMask(BitFramingReg, 0x80); err != nil { // StartSend=1, transmission of data starts
			return 0, 0, err
		}
	}

	if err := d.waitForCommandCompletion(o.WaitForIRqMask); err != nil {
		return 0, 0, err
	}

	// if there is any error except a collision, abort
	errorRegError := d.readErrorReg()
	if errorRegError != nil && (!errors.Is(errorRegError, ErrCollision)) {
		return 0, 0, errorRegError
	}

	// read data if requested
	if rx != nil {
		rawN, err := d.readSingleRegister(FIFOLevelReg)
		if err != nil {
			return 0, 0, err
		}
		bytesRead = int(rawN)

		if bytesRead > len(rx) {
			return bytesRead, 0, ErrNoRoom
		}

		if err := d.driver.ReadRegister(FIFODataReg, rx); err != nil {
			return bytesRead, 0, err
		}
		r, err := d.readSingleRegister(ControlReg)
		rxLastBits = int(r & 0b111) // RxLastBits[2:0] indicates the number of valid bits in the last received byte. If this value is 000b, the whole byte is valid.
	}

	//check for collision
	if errorRegError != nil {
		return 0, 0, errorRegError
	}

	if o.CheckCRC {
		//TODO: handle CRC
	}

	return bytesRead, rxLastBits, nil
}

func (d *Device) readErrorReg() error {

	// Stop now if any errors except collisions were detected.
	// ErrorReg[7..0] bits are: WrErr TempErr reserved BufferOvfl CollErr CRCErr ParityErr ProtocolErr BufferOvfl ParityErr ProtocolErr
	errorValue, err := d.readSingleRegister(ErrorReg)
	if err != nil {
		return err
	}
	return mapErrorValue(errorValue)
}

func mapErrorValue(errorRegValue byte) error {

	errorRegValue = errorRegValue & (^BIT5)
	if errorRegValue == 0 {
		return nil
	}

	if errorRegValue&ErrorRegCollErr != 0 {
		return ErrCollision
	}

	if errorRegValue&ErrorRegCRCErr != 0 {
		return ErrBadCrc
	}

	return ErrInternal
}
