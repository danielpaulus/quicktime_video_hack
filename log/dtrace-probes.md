I had a look at a few USB API Calls:


sudo dtrace -n '*:*:*IOUSBDevice*:entry { stack(); }'

[iOS8 "com.apple.mobile.screenshotr" is replaced with the "com.apple.cmio.iOSScreenCaptureAssistant" service · Issue #122 · libimobiledevice/libimobiledevice · GitHub](https://github.com/libimobiledevice/libimobiledevice/issues/122)

iOSScreenCaptureAssistant
https://www.youtube.com/watch?v=A9gqzn3XcDM&feature=youtu.be

/System/Library/Frameworks/CoreMediaIO.framework/Versions/A/Resources/iOSScreenCapture.plugin/Contents/Resources/iOSScreenCaptureAssistant




sudo dtrace -n '*:*:*IOUSBDevice*:entry/execname=="iOSScreenCapture"/ { stack(); }'



sudo dtrace -n '*:*:*IOUSBDevice*:entry/execname=="iOSScreenCapture"/ { printf("--%s--", execname); }'



sudo dtrace -n '*:*IOUSBFamily*:*:entry/execname=="iOSScreenCapture"/ { printf("--%s--", execname); }'

sudo dtrace -n '*:*:*ControlRequest*:entry/execname=="iOSScreenCapture"/ { printf("--%s--", probefunc); }'






that is how you write data to usb: [objective c - USB device send/receive data - Stack Overflow](https://stackoverflow.com/questions/41038150/usb-device-send-receive-data)
```IOReturn WriteToDevice(IOUSBDeviceInterface **dev, UInt16 deviceAddress,
                        UInt16 length, UInt8 writeBuffer[])
{

    IOUSBDevRequest     request;
    request.bmRequestType = USBmakebmRequestType(kUSBOut, kUSBVendor,
                                                kUSBDevice);
    request.bRequest = 0xa0;
    request.wValue = deviceAddress;
    request.wIndex = 0;
    request.wLength = length;
    request.pData = writeBuffer;

    return (*dev)->DeviceRequest(dev, &request);
}
```
so let's check out all the deviceRequests then


sudo dtrace -n '*:*:*DeviceRequest*:entry/execname=="iOSScreenCapture"/ { tracemem(arg1, 8); }'
sudo dtrace -n '*:*:*DeviceRequest*:entry { tracemem(arg1, 8); }'
sudo dtrace -n '*:*:*DeviceRequest*:entry/execname=="iOSScreenCapture"/ { printf("devpointer:%#010x -- struct_ptr: %#010x", arg0, arg1);  }'

die methoden hier:
[IOUSBFamily/IOUSBInterfaceUserClient.h at master · opensource-apple/IOUSBFamily · GitHub](https://github.com/opensource-apple/IOUSBFamily/blob/master/IOUSBUserClient/Headers/IOUSBInterfaceUserClient.h)


who is setting the config? 
sudo dtrace -n '*:*:*SetConfiguration*:entry { printf("c:%s-b:%d", execname, arg1); ustack(); }'
```
 0  86572 _ZN21IOUSBDeviceUserClient16SetConfigurationEh:entry c:usbmuxd-b:6
              libsystem_kernel.dylib`mach_msg_trap+0xa
              IOKit`io_connect_method+0x176
              IOKit`IOConnectCallScalarMethod+0x4c
              IOUSBLib`IOUSBDeviceClass::SetConfiguration(unsigned char)+0x51
              usbmuxd`0x0000000104d6a02d+0x339
              IOKit`IODispatchCalloutFromCFMessage+0x164
              CoreFoundation`__CFMachPortPerform+0x11a
              CoreFoundation`__CFRUNLOOP_IS_CALLING_OUT_TO_A_SOURCE1_PERFORM_FUNCTION__+0x29
              CoreFoundation`__CFRunLoopDoSource1+0x20f
              CoreFoundation`__CFRunLoopRun+0x9dc
              CoreFoundation`CFRunLoopRunSpecific+0x1c7
              CoreFoundation`CFRunLoopRun+0x28
              usbmuxd`0x0000000104d5d7b3+0x59e
              libdyld.dylib`start+0x1
              usbmuxd`0x2
```

seems like usbmuxd is doing it, so let's check it out:/System/Library/PrivateFrameworks/MobileDevice.framework/Versions/A/Resources/usbmuxd
/Library/Preferences/com.apple.usbmuxd.plis
interesting strings: AddNewInterface, FoundNewInterfaces
sub_10000e02d


ups:
first
![a0b55237.png](:storage/49993860-1195-4bd7-9568-ea254440d571/a0b55237.png)
and then
![1b82ea57.png](:storage/49993860-1195-4bd7-9568-ea254440d571/1b82ea57.png)
