FROM ubuntu:14.04

RUN ["/bin/bash", "-c", "apt-get update"]
RUN ["/bin/bash", "-c", "apt-get install -y software-properties-common"]

RUN ["/bin/bash", "-c", "add-apt-repository ppa:webupd8team/java && apt-get update"]


RUN echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 select true" | debconf-set-selections
RUN echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 seen true" | debconf-set-selections

RUN apt-get update 

RUN apt-get install -y oracle-java8-installer

RUN apt-get install -y g++

RUN apt-get install -y curl
RUN curl -sL https://deb.nodesource.com/setup_6.x | sudo -E bash -
RUN apt-get install -y nodejs

