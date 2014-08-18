
FROM ubuntu:14.04
MAINTAINER @joshroppo

#Insert the binary as root
ADD sheave /usr/local/bin/

#Dictionary dependencies for Anagrams
RUN apt-get update 
RUN apt-get install debconf-utils
RUN debconf-get-selections > olddeb.conf
RUN sed 's/debconf\/frontend	select Dialog/debconf\/frontend	select Noninteractive/g' olddeb.conf > newdeb.conf
RUN debconf-set-selections newdeb.conf
RUN apt-get install -y dictionaries-common

#Create sheave user
RUN useradd -m -s /bin/bash -u 9534 sheave
USER sheave

#Add config file to the sheave user's .config folder 
RUN mkdir -p /home/sheave/.config
ADD sheave.conf /home/sheave/.config/

WORKDIR /home/sheave/

EXPOSE 6667
ENTRYPOINT ["sheave", "--confPath", "~/.config/sheave.conf"]

