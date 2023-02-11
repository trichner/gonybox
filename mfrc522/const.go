package mfrc522

type Register byte

const (

	// Page 0: Command and status
	_             Register = 0x00 // reserved for future use
	CommandReg    Register = 0x01 // starts and stops command execution
	ComIEnReg     Register = 0x02 // enable and disable interrupt request control bits
	DivIEnReg     Register = 0x03 // enable and disable interrupt request control bits
	ComIrqReg     Register = 0x04 // interrupt request bits
	DivIrqReg     Register = 0x05 // interrupt request bits
	ErrorReg      Register = 0x06 // error bits showing the error status of the last command executed
	Status1Reg    Register = 0x07 // communication status bits
	Status2Reg    Register = 0x08 // receiver and transmitter status bits
	FIFODataReg   Register = 0x09 // input and output of 64 byte FIFO buffer
	FIFOLevelReg  Register = 0x0A // number of bytes stored in the FIFO buffer
	WaterLevelReg Register = 0x0B // level for FIFO underflow and overflow warning
	ControlReg    Register = 0x0C // miscellaneous control registers
	BitFramingReg Register = 0x0D // adjustments for bit-oriented frames
	CollReg       Register = 0x0E // bit position of the first bit-collision detected on the RF interface
	_             Register = 0x0F // reserved for future use

	// Page 1: Command
	_              Register = 0x10 // reserved for future use
	ModeReg        Register = 0x11 // defines general modes for transmitting and receiving
	TxModeReg      Register = 0x12 // defines transmission data rate and framing
	RxModeReg      Register = 0x13 // defines reception data rate and framing
	TxControlReg   Register = 0x14 // controls the logical behavior of the antenna driver pins TX1 and TX2
	TxASKReg       Register = 0x15 // controls the setting of the transmission modulation
	TxSelReg       Register = 0x16 // selects the internal sources for the antenna driver
	RxSelReg       Register = 0x17 // selects internal receiver settings
	RxThresholdReg Register = 0x18 // selects thresholds for the bit decoder
	DemodReg       Register = 0x19 // defines demodulator settings
	_              Register = 0x1A // reserved for future use
	_              Register = 0x1B // reserved for future use
	MfTxReg        Register = 0x1C // controls some MIFARE communication transmit parameters
	MfRxReg        Register = 0x1D // controls some MIFARE communication receive parameters
	_              Register = 0x1E // reserved for future use
	SerialSpeedReg Register = 0x1F // selects the speed of the serial UART interface

	// Page 2: Configuration
	_                 Register = 0x20 // reserved for future use
	CRCResultRegH     Register = 0x21 // shows the MSB and LSB values of the CRC calculation
	CRCResultRegL     Register = 0x22
	_                 Register = 0x23 // reserved for future use
	ModWidthReg       Register = 0x24 // controls the ModWidth setting?
	_                 Register = 0x25 // reserved for future use
	RFCfgReg          Register = 0x26 // configures the receiver gain
	GsNReg            Register = 0x27 // selects the conductance of the antenna driver pins TX1 and TX2 for modulation
	CWGsPReg          Register = 0x28 // defines the conductance of the p-driver output during periods of no modulation
	ModGsPReg         Register = 0x29 // defines the conductance of the p-driver output during periods of modulation
	TModeReg          Register = 0x2A // defines settings for the internal timer
	TPrescalerReg     Register = 0x2B // the lower 8 bits of the TPrescaler value. The 4 high bits are in TModeReg.
	TReloadRegH       Register = 0x2C // defines the 16-bit timer reload value
	TReloadRegL       Register = 0x2D
	TCounterValueRegH Register = 0x2E // shows the 16-bit timer value
	TCounterValueRegL Register = 0x2F

	// Page 3: Test Registers
	_               Register = 0x30 // reserved for future use
	TestSel1Reg     Register = 0x31 // general test signal configuration
	TestSel2Reg     Register = 0x32 // general test signal configuration
	TestPinEnReg    Register = 0x33 // enables pin output driver on pins D1 to D7
	TestPinValueReg Register = 0x34 // defines the values for D1 to D7 when it is used as an I/O bus
	TestBusReg      Register = 0x35 // shows the status of the internal test bus
	AutoTestReg     Register = 0x36 // controls the digital self-test
	VersionReg      Register = 0x37 // shows the software version
	AnalogTestReg   Register = 0x38 // controls the pins AUX1 and AUX2
	TestDAC1Reg     Register = 0x39 // defines the test value for TestDAC1
	TestDAC2Reg     Register = 0x3A // defines the test value for TestDAC2
	TestADCReg      Register = 0x3B // shows the value of ADC I and Q channels
)

// Command has the constants from the datasheet section '10. MFRC522 command set'
type Command byte

const (
	CommandIdle             Command = 0x00 // no action, cancels current command execution
	CommandMem              Command = 0x01 // stores 25 bytes into the internal buffer
	CommandGenerateRandomID Command = 0x02 // generates a 10-byte random ID number
	CommandCalcCRC          Command = 0x03 // activates the CRC coprocessor or performs a self-test
	CommandTransmit         Command = 0x04 // transmits data from the FIFO buffer
	CommandNoCmdChange      Command = 0x07 // no command change, can be used to modify the CommandReg register bits without affecting the command, for example, the PowerDown bit
	CommandReceive          Command = 0x08 // activates the receiver circuits
	CommandTransceive       Command = 0x0C // transmits data from FIFO buffer to antenna and automatically activates the receiver after transmission
	CommandMFAuthent        Command = 0x0E // performs the MIFARE standard authentication as a reader
	CommandSoftReset        Command = 0x0F // resets the MFRC522
)

func (c Command) ToSlice() []byte {
	return []byte{byte(c)}
}

const (
	BIT0 byte = 1 << iota
	BIT1
	BIT2
	BIT3
	BIT4
	BIT5
	BIT6
	BIT7
)

type BitFraming byte

func (b *BitFraming) SetStartSend(on bool) {
	if on {
		*b |= BitFraming(BIT7)
	}
	*b &= ^BitFraming(BIT7)
}

func (b *BitFraming) SetRxAlign(align uint8) {
	mask := byte(0b111 << 4)
	t := byte(*b)
	t &= ^mask
	*b = BitFraming(t | (mask & (align << 4)))
}

func (b *BitFraming) SetTxLastBits(n uint8) {
	mask := byte(0b111)
	t := byte(*b)
	t &= ^mask
	*b = BitFraming(t | (mask & n))
}

// Datasheet 9.3.1.7 - ErrorReg register
const (
	//ErrorRegWrErr data is written into the FIFO buffer by the host during the MFAuthent
	//command or if data is written into the FIFO buffer by the host during the
	//time between sending the last bit on the RF interface and receiving the
	//last bit on the RF interface
	ErrorRegWrErr = BIT7

	//ErrorRegTempErr internal temperature sensor detects overheating, in which case the
	//antenna drivers are automatically switched off
	ErrorRegTempErr = BIT6

	//ErrorRegBufferOvfl the host or a MFRC522’s internal state machine (e.g. receiver) tries to
	//write data to the FIFO buffer even though it is already full
	ErrorRegBufferOvfl = BIT4

	//ErrorRegCollErr a bit-collision is detected
	//- cleared automatically at receiver start-up phase
	//- only valid during the bitwise anticollision at 106 kBd
	//- always set to logic 0 during communication protocols at 212 kBd, 424 kBd and 848 kBd
	ErrorRegCollErr = BIT3

	// ErrorRegCRCErr the RxModeReg register’s RxCRCEn bit is set and the CRC calculation fails
	//automatically cleared to logic 0 during receiver start-up phase
	ErrorRegCRCErr = BIT2
	//ErrorRegParityErr parity check failed
	//automatically cleared during receiver start-up phase
	//only valid for ISO/IEC 14443 A/MIFARE communication at 106 kBd
	ErrorRegParityErr = BIT1

	//ErrorRegProtocolErr set to logic 1 if the SOF is incorrect
	//automatically cleared during receiver start-up phase
	//bit is only valid for 106 kBd
	//during the MFAuthent command, the ProtocolErr bit is set to logic 1 if the
	//number of bytes received in one data stream is incorrect
	ErrorRegProtocolErr = BIT0
)

// Version of chip / firmware.
const (
	VersionCounterfeit = 0x12
	VersionFM17522     = 0x88
	VersionFM17522_1   = 0xb2
	VersionFM17522E    = 0x89
	Version0_0         = 0x90
	Version1_0         = 0x91
	Version2_0         = 0x92
	VersionUnknown     = 0xff
)

var knownVersions = []byte{VersionCounterfeit, VersionFM17522, VersionFM17522_1, VersionFM17522E, Version0_0, Version1_0, Version2_0}
