package mfrc522

import (
	"errors"
	"strconv"
	"trelligo/pkg/debug"
)

type UID []byte

func (d *Device) PiccReadCardSerial() (UID, error) {

	// that's equal to selecting it :)
	return d.PiccSelect()
}

type PiccSelectCommand [9]byte

func (p *PiccSelectCommand) setCommand(cmd PiccCommand) {
	p[0] = byte(cmd)
}

func (p *PiccSelectCommand) setNumberOfValidBits(bits int) {
	fullBytes := bits / 8
	additionalBits := bits % 8
	p[1] = byte(fullBytes<<4 | additionalBits)
}

func (p *PiccSelectCommand) setUuidData(bytes *[4]byte) {
	p[2] = bytes[0]
	p[3] = bytes[1]
	p[4] = bytes[2]
	p[5] = bytes[3]
}

func (p *PiccSelectCommand) uuidData() []byte {
	return p[2:6]
}

func (p *PiccSelectCommand) updateBlockCheckCharacter() {
	p[6] = p[2] ^ p[3] ^ p[4] ^ p[5]
}

func (p *PiccSelectCommand) blockCheckCharacter() byte {
	return p[2] ^ p[3] ^ p[4] ^ p[5]
}

func (p *PiccSelectCommand) updateCrc(calculator func(data []byte, crc []byte) error) error {
	return calculator(p[:7], p[7:9])
}

func (p *PiccSelectCommand) slice() []byte {
	return p[:]
}

// https://github.com/OSSLibraries/Arduino_MFRC522v2/blob/95ac65b2d5f6a0bcb155f30d51add600f7d73000/src/MFRC522v2.cpp#L536
func (d *Device) PiccSelect() (UID, error) {

	// ValuesAfterColl=1 => Bits received after collision are cleared.
	if err := d.clearRegisterBitMask(CollReg, 0x80); err != nil {
		return nil, err
	}

	cascadeLevel := 1
	uidData := make([]byte, 3*4) // we can have 4 bytes per level and at most 3 cascade levels
	for {
		debug.Log("cascading at level " + strconv.Itoa(cascadeLevel))
		var selectCommand PiccCommand
		switch cascadeLevel {
		case 1:
			selectCommand = PiccCommandSelCl1
		case 2:
			selectCommand = PiccCommandSelCl2
		case 3:
			selectCommand = PiccCommandSelCl3
		default:
			return nil, errors.New("bad cascade level")
		}

		//start with cascade level 1
		uidChunk, err := d.doAnticollisionLoop(selectCommand)
		if err != nil {
			return nil, err
		}

		debug.Log("uuid chunk: " + debug.FmtSliceToHex(uidChunk))

		done, err := d.doSelect(selectCommand, uidChunk)
		if err != nil {
			return nil, err
		}
		copy(uidData[((cascadeLevel-1)*4):], uidChunk)
		if done {
			//we don't need to cascade further, device selected
			break
		}
		cascadeLevel++
	}

	return parseUid(uidData)
}

func parseUid(uidData []byte) (UID, error) {
	// 'single' 4 byte UID
	if uidData[0] != byte(PiccCommandCt) {
		uid := make([]byte, 4)
		copy(uid, uidData)
		return uid, nil
	}

	// double UID tag, 7 bytes
	if uidData[4] != byte(PiccCommandCt) {
		uid := make([]byte, 7)
		copy(uid, uidData[1:])
		return uid, nil
	}

	// we have a triple UID tag, 10 bytes
	uid := make([]byte, 10)
	copy(uid, uidData[1:4])
	copy(uid[3:], uidData[5:])
	return uid, nil
}

func (d *Device) doSelect(selectionCommand PiccCommand, uidData []byte) (bool, error) {

	cmd := PiccSelectCommand{}
	cmd.setCommand(selectionCommand)

	numberOfValidBits := 7 * 8 // SEL + NVB + 4*UID data + BCC
	cmd.setNumberOfValidBits(numberOfValidBits)
	cmd.setUuidData((*[4]byte)(uidData))
	cmd.updateBlockCheckCharacter()

	if err := cmd.updateCrc(d.calculateCrc); err != nil {
		return false, err
	}

	selectAcknowledge := make([]byte, 3) // also known as SAK
	n, bits, err := d.transceiveData(cmd[:], selectAcknowledge, 0, 0, false)
	if err != nil {
		return false, errors.New("bad SAK: " + err.Error())
	}
	if n != 3 || bits != 0 {
		return false, errors.New("bad SAK: expected 24 bits")
	}

	crc := make([]byte, 2)
	if err := d.calculateCrc(selectAcknowledge[:1], crc); err != nil {
		return false, err
	}
	if crc[0] != selectAcknowledge[1] || crc[1] != selectAcknowledge[2] {
		return false, errors.New("bad CRC in SAK")
	}

	cascadingDone := selectAcknowledge[0]&0x04 == 0
	return cascadingDone, nil
}
func (d *Device) doAnticollisionLoop(selectionCommand PiccCommand) ([]byte, error) {

	cmd := PiccSelectCommand{}

	// This NVB (number of valid bits) value defines that the PCD will transmit no part of UID CLn. Consequently this
	//command forces all PICCs in the field to respond with their complete UID CLn
	numberOfValidBits := 2 * 8 // 2 bytes, SEL + NVB
	cmd.setCommand(selectionCommand)

	// receive buffer, 4*UID bytes + BCC
	for {

		debug.Log("ANTICOLLISION - NVB: " + strconv.Itoa(numberOfValidBits))
		cmd.setNumberOfValidBits(numberOfValidBits)

		numberOfValidBytes := numberOfValidBits / 8
		numberOfRemainingBits := byte(numberOfValidBits % 8)

		txNumberOfBytes := numberOfValidBytes
		if numberOfRemainingBits > 0 {
			txNumberOfBytes += 1
		}

		debug.Log("cmd: " + debug.FmtSliceToHex(cmd[:]))

		debug.Log("tx: " + debug.FmtSliceToHex(cmd[:txNumberOfBytes]) + " + " + strconv.Itoa(int(numberOfRemainingBits)) + "bits")

		rx := cmd[numberOfValidBytes:]
		debug.Log("rx_a: " + debug.FmtSliceToHex(rx))
		nbytes, rxLastBits, err := d.transceiveData(cmd[:txNumberOfBytes], rx, numberOfRemainingBits, numberOfRemainingBits, false)
		debug.Log("rx_b: len=" + strconv.Itoa(nbytes*8) + " last_bits=" + strconv.Itoa(rxLastBits) + " " + debug.FmtSliceToHex(rx))

		if err != nil {
			debug.Log("anticollision result: " + err.Error())
		}
		if err != nil && (!errors.Is(err, ErrCollision)) {
			return nil, err
		} else if errors.Is(err, ErrCollision) {
			numberOfValidUidBits, err := d.resolveSelectUidCollision(&cmd)
			if err != nil {
				return nil, err
			}
			nextNumberOfValidBits := numberOfValidUidBits + 2*8
			if nextNumberOfValidBits <= numberOfValidBits {
				// we should have learned at least one bit, something is off
				return nil, ErrInternal
			}
			numberOfValidBits = nextNumberOfValidBits
		} else {
			// no error at all, we know now 32 bits

			//TODO check BCC
			//bcc := cmd.blockCheckCharacter()

			return cmd.uuidData(), nil
		}
	}
}

// resolveSelectUidCollision finds the colliding bit and updates the command resolving to the bit set to 1
// it then returns the number of now valid bits within the UID bytes, which is between 1 and 32 bits
func (d *Device) resolveSelectUidCollision(cmd *PiccSelectCommand) (int, error) {

	// CollReg[7..0] bits are: ValuesAfterColl reserved CollPosNotValid CollPos[4:0]
	r, err := d.readSingleRegister(CollReg)
	if err != nil {
		return 0, err
	}
	if r&0x20 != 0 { // CollPosNotValid
		return 0, ErrCollision // Without a valid collision position we cannot continue
	}
	collisionPos := int(r & 0x1F) // Values 0-31, 0 means bit 32.
	if collisionPos == 0 {
		collisionPos = 32
	}

	bitIndex := (collisionPos - 1) % 8
	bufferIndex := 1 + (collisionPos / 8) // includes SEL + NVB, we always know at least one bit
	if collisionPos%8 != 0 {              // byte not full
		bufferIndex += 1
	}
	cmd[bufferIndex] |= 1 << bitIndex
	return collisionPos, nil
}

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

	if err := d.sendIdleCommand(); err != nil { // Stop any active command.
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

		r, err := d.readSingleRegister(ControlReg)
		rxLastBits = int(r & 0b111) // RxLastBits[2:0] indicates the number of valid bits in the last received byte. If this value is 000b, the whole byte is valid.

		if err := d.driver.ReadRegister(FIFODataReg, rx); err != nil {
			return bytesRead, rxLastBits, err
		}
	}

	//check for collision
	if errorRegError != nil {
		return bytesRead, rxLastBits, errorRegError
	}

	if o.CheckCRC {
		//TODO: handle CRC
	}

	return bytesRead, rxLastBits, nil
}

func (d *Device) sendIdleCommand() error {
	return d.writeSingleRegister(CommandReg, byte(CommandIdle))
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
