# Dockerfiles

Here you can see how to get everything running easily.
I have created separate Dockerfiles for building and for running.

### For building qvh use the Dockerfile in /build like so:

#### 1. Build Image:

- `docker build -f build/Dockerfile.debian -t "qvhbuild:$(git branch --show-current)" --build-arg GIT_BRANCH=$(git branch --show-current) .`

#### 2. Get shell in container

- `docker run -it qvhbuild:$(git branch --show-current) bash`

### For running qvh use the Dockerfile like so:

#### 1. Build Image for Running:

- `docker build -f Dockerfile.debian -t "qvhrun:$(git branch --show-current)" .`

#### 2. Get shell in container

- mount your host usb devices into the container and get a shell with the following command:
- `docker run --privileged -v /dev/bus/usb:/dev/bus/usb -it qvhrun:$(git branch --show-current) bash`
- use the `/bin/linux/qvh` binary to execute qvh
