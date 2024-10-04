#include <stdio.h>
#include <string.h>

#include "hardware/clocks.h"
#include "pico/stdlib.h"

#include "usb.h"
#include "slip.h"


namespace
{
  constexpr uint32_t kDpPin = 11;
  constexpr size_t kUSBMaxPacketLen = 1028;

  Packetti::Slip::Stdout slipOut_;

  void initializeStdio()
  {
    stdio_usb_init();
    // Disable CR/LF conversion built into stdio of Pico SDK.
    stdio_set_translate_crlf(&stdio_usb, false);
  }
}

int main()
{
  // Change system clock to 120 MHz (10 times the frequency of USB Full Speed)
  set_sys_clock_khz(120000, true);

  initializeStdio();
  Packetti::USB::Initialize(kDpPin);

  std::basic_string<uint8_t> usbPacket = {};
  usbPacket.reserve(kUSBMaxPacketLen);

  while (true) {
    if (Packetti::USB::NextPacket(usbPacket)) {
      slipOut_.WritePacket(usbPacket);
    }
  }
}
