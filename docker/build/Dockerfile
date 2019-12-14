FROM debian:latest
RUN apt-get update && apt install -y git
RUN git clone https://github.com/danielpaulus/quicktime_video_hack
RUN wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
RUN tar -xvf go1.13.3.linux-amd64.tar.gz
RUN mv go /usr/local
RUN export GOROOT=/usr/local/go
RUN export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
RUN apt-get update && apt install -y libgstreamer-plugins-base1.0-dev

