package seesaw

type ModuleBaseAddress byte

/** Module Base Addreses
 *  The module base addresses for different seesaw modules.
 */
const (
	SEESAW_STATUS_BASE  ModuleBaseAddress = 0x00
	SEESAW_GPIO_BASE    ModuleBaseAddress = 0x01
	SEESAW_SERCOM0_BASE ModuleBaseAddress = 0x02

	SEESAW_TIMER_BASE     ModuleBaseAddress = 0x08
	SEESAW_ADC_BASE       ModuleBaseAddress = 0x09
	SEESAW_DAC_BASE       ModuleBaseAddress = 0x0A
	SEESAW_INTERRUPT_BASE ModuleBaseAddress = 0x0B
	SEESAW_DAP_BASE       ModuleBaseAddress = 0x0C
	SEESAW_EEPROM_BASE    ModuleBaseAddress = 0x0D
	SEESAW_NEOPIXEL_BASE  ModuleBaseAddress = 0x0E
	SEESAW_TOUCH_BASE     ModuleBaseAddress = 0x0F
	SEESAW_KEYPAD_BASE    ModuleBaseAddress = 0x10
	SEESAW_ENCODER_BASE   ModuleBaseAddress = 0x11
	SEESAW_SPECTRUM_BASE  ModuleBaseAddress = 0x12
)

type FunctionAddress byte

/** GPIO module function address registers
 */
const (
	SEESAW_GPIO_DIRSET_BULK FunctionAddress = 0x02
	SEESAW_GPIO_DIRCLR_BULK FunctionAddress = 0x03
	SEESAW_GPIO_BULK        FunctionAddress = 0x04
	SEESAW_GPIO_BULK_SET    FunctionAddress = 0x05
	SEESAW_GPIO_BULK_CLR    FunctionAddress = 0x06
	SEESAW_GPIO_BULK_TOGGLE FunctionAddress = 0x07
	SEESAW_GPIO_INTENSET    FunctionAddress = 0x08
	SEESAW_GPIO_INTENCLR    FunctionAddress = 0x09
	SEESAW_GPIO_INTFLAG     FunctionAddress = 0x0A
	SEESAW_GPIO_PULLENSET   FunctionAddress = 0x0B
	SEESAW_GPIO_PULLENCLR   FunctionAddress = 0x0C
)

/** status module function address registers
 */
const (
	SEESAW_STATUS_HW_ID   FunctionAddress = 0x01
	SEESAW_STATUS_VERSION FunctionAddress = 0x02
	SEESAW_STATUS_OPTIONS FunctionAddress = 0x03
	SEESAW_STATUS_TEMP    FunctionAddress = 0x04
	SEESAW_STATUS_SWRST   FunctionAddress = 0x7F
)

/** timer module function address registers
 */
const (
	SEESAW_TIMER_STATUS FunctionAddress = 0x00
	SEESAW_TIMER_PWM    FunctionAddress = 0x01
	SEESAW_TIMER_FREQ   FunctionAddress = 0x02
)

/** ADC module function address registers
 */
const (
	SEESAW_ADC_STATUS         FunctionAddress = 0x00
	SEESAW_ADC_INTEN          FunctionAddress = 0x02
	SEESAW_ADC_INTENCLR       FunctionAddress = 0x03
	SEESAW_ADC_WINMODE        FunctionAddress = 0x04
	SEESAW_ADC_WINTHRESH      FunctionAddress = 0x05
	SEESAW_ADC_CHANNEL_OFFSET FunctionAddress = 0x07
)

/** Sercom module function address registers
 */
const (
	SEESAW_SERCOM_STATUS   FunctionAddress = 0x00
	SEESAW_SERCOM_INTEN    FunctionAddress = 0x02
	SEESAW_SERCOM_INTENCLR FunctionAddress = 0x03
	SEESAW_SERCOM_BAUD     FunctionAddress = 0x04
	SEESAW_SERCOM_DATA     FunctionAddress = 0x05
)

/** neopixel module function address registers
 */
const (
	SEESAW_NEOPIXEL_STATUS     FunctionAddress = 0x00
	SEESAW_NEOPIXEL_PIN        FunctionAddress = 0x01
	SEESAW_NEOPIXEL_SPEED      FunctionAddress = 0x02
	SEESAW_NEOPIXEL_BUF_LENGTH FunctionAddress = 0x03
	SEESAW_NEOPIXEL_BUF        FunctionAddress = 0x04
	SEESAW_NEOPIXEL_SHOW       FunctionAddress = 0x05
)

/** touch module function address registers
 */
const (
	SEESAW_TOUCH_CHANNEL_OFFSET FunctionAddress = 0x10
)

/** keypad module function address registers
 */
const (
	SEESAW_KEYPAD_STATUS   FunctionAddress = 0x00
	SEESAW_KEYPAD_EVENT    FunctionAddress = 0x01
	SEESAW_KEYPAD_INTENSET FunctionAddress = 0x02
	SEESAW_KEYPAD_INTENCLR FunctionAddress = 0x03
	SEESAW_KEYPAD_COUNT    FunctionAddress = 0x04
	SEESAW_KEYPAD_FIFO     FunctionAddress = 0x10
)

/** keypad module edge definitions
 */
const (
	SEESAW_KEYPAD_EDGE_HIGH = 0
	SEESAW_KEYPAD_EDGE_LOW
	SEESAW_KEYPAD_EDGE_FALLING
	SEESAW_KEYPAD_EDGE_RISING
)

/** encoder module edge definitions
 */
const (
	SEESAW_ENCODER_STATUS   = 0x00
	SEESAW_ENCODER_INTENSET = 0x10
	SEESAW_ENCODER_INTENCLR = 0x20
	SEESAW_ENCODER_POSITION = 0x30
	SEESAW_ENCODER_DELTA    = 0x40
)

/** Audio spectrum module function address registers
 */
const (
	SEESAW_SPECTRUM_RESULTS_LOWER FunctionAddress = 0x00 // Audio spectrum bins 0-31
	SEESAW_SPECTRUM_RESULTS_UPPER FunctionAddress = 0x01 // Audio spectrum bins 32-63
	// If some future device supports a larger spectrum can add additional
	// "bins" working upward from here. Configurable setting registers then
	// work downward from the top to avoid collision between spectrum bins
	// and configurables.
	SEESAW_SPECTRUM_CHANNEL FunctionAddress = 0xFD
	SEESAW_SPECTRUM_RATE    FunctionAddress = 0xFE
	SEESAW_SPECTRUM_STATUS  FunctionAddress = 0xFF
)
