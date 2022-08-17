#include <stdio.h>
#include <stdint.h>
#include <unistd.h>
#include "libusb-win32-bin-1.2.6.0/include/lusb0_usb.h"

#if INTPTR_MAX == INT64_MAX
#pragma comment(lib,"libusb-win32-bin-1.2.6.0/lib/msvc_x64/libusb.lib")
#elif INTPTR_MAX == INT32_MAX
#pragma comment(lib,"libusb-win32-bin-1.2.6.0/lib/msvc/libusb.lib")
#else
#error Unknown pointer size or missing size macros!
#endif

int main() {
    printf("Hello, World!\n");
    return 0;
}

void openDev(){
    struct usb_bus* bus = NULL;
    struct usb_device* dev = NULL;
    struct usb_dev_handle* udh = NULL;
    usb_find_busses();
    usb_find_devices();
    for (bus = usb_get_busses(); bus; bus = bus->next)
    {
        for (dev = bus->devices; dev; dev = dev->next)
        {
            if (dev->descriptor.idVendor == 0x05AC)
            {
                udh = usb_open(dev);
            }
        }
    }
    usb_control_msg(udh, 0x40, 0x52, 0x00, 0x02, 0, 0, 1000);
    sleep(1);
}
