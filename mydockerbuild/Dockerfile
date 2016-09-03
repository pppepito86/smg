FROM ubuntu:14.04

RUN ["/bin/bash", "-c", "apt-get update"]
RUN ["/bin/bash", "-c", "apt-get install -y software-properties-common"]

RUN ["/bin/bash", "-c", "add-apt-repository ppa:webupd8team/java && apt-get update"]


RUN echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 select true" | debconf-set-selections
RUN echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 seen true" | debconf-set-selections

RUN apt-get install -y oracle-java8-installer
