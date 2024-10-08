# == DO NEVER EDIT THE NEXT LINES for Raspberry Pi Pico VS Code Extension to work ==
if(WIN32)
    set(USERHOME $ENV{USERPROFILE})
else()
    set(USERHOME $ENV{HOME})
endif()
set(sdkVersion 2.0.0)
set(toolchainVersion 13_2_Rel1)
set(picotoolVersion 2.0.0)
set(picoVscode ${USERHOME}/.pico-sdk/cmake/pico-vscode.cmake)
if (EXISTS ${picoVscode})
    include(${picoVscode})
endif()
# ====================================================================================
# == DO NEVER EDIT THE NEXT LINES for Raspberry Pi Pico VS Code Extension to work ==
if(WIN32)
    set(USERHOME $ENV{USERPROFILE})
else()
    set(USERHOME $ENV{HOME})
endif()
set(sdkVersion 2.0.0)
set(toolchainVersion 13_2_Rel1)
set(picotoolVersion 2.0.0)
set(picoVscode ${USERHOME}/.pico-sdk/cmake/pico-vscode.cmake)
if (EXISTS ${picoVscode})
    include(${picoVscode})
endif()
# ====================================================================================
cmake_minimum_required(VERSION 3.13)

include(pico_sdk_import.cmake)

project(packetti)

set(CMAKE_C_STANDARD 11)
set(CMAKE_CXX_STANDARD 23)

add_compile_definitions(
  USBD_PID=0x5052
  USBD_MANUFACTURER="@UTBDK"
  USBD_PRODUCT="Packetti"
)

pico_sdk_init()

add_executable(packetti
  main.cc
  slip.cc
  usb.cc
)

pico_generate_pio_header(packetti ${CMAKE_CURRENT_LIST_DIR}/usb_sniff.pio)

target_link_libraries(packetti pico_stdlib pico_multicore pico_stdio_usb hardware_pio hardware_dma hardware_irq)

pico_add_extra_outputs(packetti)
