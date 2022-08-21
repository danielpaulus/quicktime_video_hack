cmake_minimum_required(VERSION 3.23)
set(CMAKE_VERBOSE_MAKEFILE ON)
project(windows_activator C)

set(CMAKE_C_STANDARD 99)
set(CMAKE_WINDOWS_EXPORT_ALL_SYMBOLS TRUE)
set(BUILD_SHARED_LIBS TRUE)

add_library(bar STATIC IMPORTED GLOBAL)
set_property(TARGET bar PROPERTY IMPORTED_LOCATION ${CMAKE_CURRENT_SOURCE_DIR}/libusb-win32-bin-1.2.6.0/lib/msvc_x64/libusb.lib)

add_executable(windows_activator main.c)

target_link_libraries(windows_activator bar)
