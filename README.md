

```bash
# Debug with openocd (Fedora openocd doesn't compute)
sudo ~/.platformio/packages/tool-openocd/bin/openocd -f /usr/share/openocd/scripts/board/st_nucleo_l4.cfg
```

https://wiki.dfrobot.com/DFPlayer_Mini_SKU_DFR0299


## Serial for DFPlayer Mini

```
playFirst := []byte{0x7e, 0xff, 0x06, 0x03, 0x00, 0x00, 0x01, 0xff, 0xe6, 0xef}
reset := []byte{0x7e, 0xff, 0x06, 0x0c, 0x01, 0x00, 0x00, 0xfe, 0xee, 0xef}

00: Start Bit 0x7E
01: Version, 0xFF
02: Length, e.g. 6
03: Command Byte, e.g. 0x0C for reset
04: Should ACK (0x01 for enabled, otherwise 0x00)
05: Query high byte
06: Query low byte
06: Checksum
07: End Byte, 0xEF
```
