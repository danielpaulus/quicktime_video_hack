# Dockerfiles
Here you can see how to get everything running easily. I included a Dockerfile for debugging on Debian
Build the container and run something like:
`docker run --privileged -v /dev/bus/usb:/dev/bus/usb -it 83e1cdeb059a bash
` to try. Then you can use the `/bin/linux/qvh` binary.
