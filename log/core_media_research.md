
### CMTime
http://developer.apple.com/library/ios/#documentation/CoreMedia/Reference/CMTime/Reference/reference.html

https://stackoverflow.com/questions/3684883/coremedia-cmtime

The Apple documentation contains information on the CMTime struct.

As I understand it, you set the "TimeScale" to a timescale suitable for the media (e.g. 44100 = 1/44100 sec - which might be suitable for a CD). Then the "Value" represents units of that timescale. So, a value of 88200 would be 2 secs.
```
CMTime cmTime = new CMTime();
cmTime.TimeScale = 44100;
cmTime.Value = 88200;
``` 

### Xamarin C# native interface for CoreMedia

https://github.com/xamarin/xamarin-macios/tree/master/src/CoreMedia

### Core Media Headers
https://github.com/phracker/MacOSX-SDKs/tree/master/MacOSX10.9.sdk/System/Library/Frameworks/CoreMedia.framework