FROM debian:latest

RUN apt-get update && apt install -y wget
RUN wget https://github.com/danielpaulus/quicktime_video_hack/releases/download/v0.2-beta/bin.zip
RUN apt-get update && apt install -y zip
RUN unzip bin.zip
RUN apt-get update && apt install libusb-1.0
RUN apt-get update && apt-get install -y libgstreamer1.0-0 gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-doc gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl gstreamer1.0-gtk3 gstreamer1.0-qt5 gstreamer1.0-pulseaudio
RUN chmod +x /bin/linux/qvh
