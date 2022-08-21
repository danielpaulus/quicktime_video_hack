cmake_minimum_required(VERSION 3.23)
project(windows_activator C)

set(CMAKE_C_STANDARD 99)

add_library(bar STATIC IMPORTED)
set_property(TARGET bar PROPERTY IMPORTED_IMPLIB ${CMAKE_CURRENT_SOURCE_DIR}/libusb-win32-bin-1.2.6.0/lib/msvc_x64/libusb.lib)

add_executable(windows_activator main.c)

target_link_libraries(windows_activator bar)
