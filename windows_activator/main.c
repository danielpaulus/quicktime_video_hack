#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <windows.h>
#include "./libusb-win32-bin-1.2.6.0/include/lusb0_usb.h"

int openDev(char *serial);

#if INTPTR_MAX == INT64_MAX
#pragma comment(lib, "../libusb-win32-bin-1.2.6.0/lib/msvc_x64/libusb.lib")
#elif INTPTR_MAX == INT32_MAX
#pragma comment(lib,"libusb-win32-bin-1.2.6.0/lib/msvc/libusb.lib")
#else
#error Unknown pointer size or missing size macros!
#endif

#define activate_device = "s"

int main(int argc, char **argv) {
    if (argc < 2) {
        printf("no arguments passed, need device serial");
        return 1;
    }
    if (argc > 2) {
        printf("invalid arguments passed, need device serial");
        return 1;
    }
    usb_init();
    usb_set_debug(255);
    return openDev(argv[1]);
}

int openDev(char *deviceSerial) {
    printf("openDev, search: %s", deviceSerial);
    struct usb_bus *bus = NULL;
    struct usb_device *dev = NULL;
    struct usb_dev_handle *udh = NULL;
    int code;
    code = usb_find_busses();
    if (code != 0) {
        printf("error usb_find_busses: %d", code);
    }
    code = usb_find_devices();
    if (code != 0) {
        printf("error usb_find_busses: %d", code);
    }
    for (bus = usb_get_busses(); bus; bus = bus->next) {
        for (dev = bus->devices; dev; dev = dev->next) {
            printf("checking device pid: %x vid: %x ", dev->descriptor.idProduct, dev->descriptor.idVendor);
            udh = usb_open(dev);
            if (udh == 0) {
                printf("error usb_open: pid: %x vid: %x ", dev->descriptor.idProduct, dev->descriptor.idVendor);
                continue;
            }
            char szSerialNumber[128] = {0};
            code = usb_get_string_simple(udh, dev->descriptor.iSerialNumber, szSerialNumber, sizeof(szSerialNumber));
            if (code != 0) {
                printf("error usb_get_string_simple: %d", code);
                continue;
            }
            printf("cmp: %s with %s", szSerialNumber, deviceSerial);
            if (strcmp(szSerialNumber, deviceSerial) == 0) {
                #ifdef activate_device
                printf("match, send control message to activate");
                code = usb_control_msg(udh, 0x40, 0x52, 0x00, 0x02, 0, 0, 1000);
                #else
                printf("match, send control message to de-activate");
                code = usb_control_msg(udh, 0x40, 0x52, 0x00, 0x00, 0, 0, 1000);
                #endif
                if (code != 0) {
                    printf("error usb_control_msg: %d", code);
                }
                printf("sleep");
                Sleep(1000);
                printf("done");
                return code;
            }
            code = usb_close(udh);
            if (code != 0) {
                printf("error usb_close: pid: %x vid: %x ", dev->descriptor.idProduct, dev->descriptor.idVendor);
            }
        }
    }
    return 1;
}
