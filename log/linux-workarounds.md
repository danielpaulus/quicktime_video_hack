
THe usbmuxd will be killed via udev when the last device is removed. That is annoying in our case
because we want to use the Listen command. Remove the following line in the udev conf:
 
https://bbs.archlinux.org/viewtopic.php?id=229475
```
I just find a wordaround. The udev rule of usbmuxd (in /lib/udev/rules.d/39-usbmuxd.rules) is as following:

============================================================
# usbmuxd (Apple Mobile Device Muxer listening on /var/run/usbmuxd)
 
# Initialize iOS devices into "deactivated" USB configuration state and activate usbmuxd
ACTION=="add", SUBSYSTEM=="usb", ATTR{idVendor}=="05ac", ATTR{idProduct}=="12[9a][0-9a-f]", ENV{USBMUX_SUPPORTED}="1", ATTR{bConfigurationValue}="0", OWNER="usbmux", TAG+="systemd", ENV{SYSTEMD_WANTS}="usbmuxd.service"

# Exit usbmuxd when the last device is removed
ACTION=="remove", SUBSYSTEM=="usb", ENV{PRODUCT}=="5ac/12[9a][0-9a-f]/*", ENV{INTERFACE}=="255/*", RUN+="/usr/bin/usbmuxd -x"
============================================================


It seems that when the last device is removed, the usbmuxd will exit automatically. So I comment out the last line of this file, then my phone can be recognized at everytime it is pluged in.
```