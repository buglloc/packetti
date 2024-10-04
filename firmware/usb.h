#pragma once

#include <stdlib.h>
#include <stdint.h>
#include <string>
#include <string_view>

#include "pico/types.h"
#include "pico/error.h"
#include "pico/util/queue.h"


namespace Packetti::USB
{
  pico_error_codes Initialize(uint32_t dpPin);
  bool NextPacket(std::basic_string<uint8_t>& out);
}
