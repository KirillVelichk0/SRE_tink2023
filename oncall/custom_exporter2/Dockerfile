FROM golang:alpine
WORKDIR /home/custom_export/source/build
COPY src/ /home/custom_export/source/src
COPY configs/ /home/custom_export/source/configs

RUN cd /home/custom_export/source/src && go build -o /home/custom_export/source/build/main
EXPOSE 9211
CMD ["/home/custom_export/source/build/main"]
