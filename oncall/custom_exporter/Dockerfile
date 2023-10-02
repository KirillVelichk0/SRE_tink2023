FROM ubuntu:22.04


RUN DEBIAN_FRONTEND=noninteractive apt-get update\
&& DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::=--force-confold install apt-utils wget \
 && DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::=--force-confold install sudo libssl-dev systemctl \
    && DEBIAN_FRONTEND=noninteractive rm -rf /var/cache/apt/archives/*
    
RUN DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::=--force-confold upgrade && sudo DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::=--force-confold update&&\
DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::=--force-confold install git

RUN sudo DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::=--force-confold install golang


RUN useradd -m -s /bin/bash custom_export
COPY src /home/custom_export/source/src
COPY installer /home/custom_export/installer
RUN rm -f /home/custom_export/source/src/go.sum && rm -f /home/custom_export/source/src/go.mod
COPY README.md /home/custom_export/source/README.md
WORKDIR /home/custom_export/build


RUN chown -R custom_export:custom_export /home/custom_export/source
RUN cd /home/custom_export/source/src && bash /home/custom_export/installer && sudo go build -o /home/custom_export/build
RUN chown -R custom_export:custom_export /home/custom_export/build

COPY configs /home/custom_export/configs

EXPOSE 9210

CMD ["sudo", "-EHu", "custom_export", "/home/custom_export/build/main"] 
