# Dockerfiles
Here you can see how to get everything running easily. I included a Dockerfile for debugging on Debian
Build the container and run something like:
- `docker build . -t qvh:v0.2` to build a docker image using the Dockerfile and tag it as qvh:v0.2
- `docker run --privileged -v /dev/bus/usb:/dev/bus/usb -it qvh:v0.2 bash` to get a bash session and mount all usb devices inside the container
- use the `/bin/linux/qvh` binary.
