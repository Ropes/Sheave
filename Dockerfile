
FROM ubuntu:14.04
MAINTAINER @joshroppo

#Insert the sheavebot binary as root
ADD sheave /usr/local/bin/
ADD bot/words /usr/share/words

#Dictionary dependencies for Anagrams
RUN apt-get update 
RUN apt-get install debconf-utils
RUN debconf-get-selections > olddeb.conf
RUN sed 's/Dialog/Noninteractive/g' olddeb.conf > newdeb.conf
RUN cat newdeb.conf | grep -C 2 frontend
RUN debconf-set-selections newdeb.conf
RUN rm olddeb.conf newdeb.conf
RUN apt-get install -y dictionaries-common

#Create sheave user
RUN useradd -m -s /bin/bash -u 9534 sheave
USER sheave

#Add config file to the sheave user's .config folder 
RUN mkdir -p /home/sheave/.config
ADD sheave.conf /home/sheave/.config/
RUN ls -la /home/sheave/.config

WORKDIR /home/sheave/

EXPOSE 6667
ENTRYPOINT ["sheave", "--confPath", "/home/sheave/.config/sheave.conf", "--dictPath", "/usr/share/words"]
#ENTRYPOINT /bin/sh

