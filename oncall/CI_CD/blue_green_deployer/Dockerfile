FROM golang:alpine
WORKDIR /home/custom_export/source/build
COPY cmd/ /home/custom_export/source/cmd
COPY configs/ /home/custom_export/source/configs
COPY internal/ /home/custom_export/source/internal
COPY go.mod /home/custom_export/source/go.mod
COPY go.sum /home/custom_export/source/go.sum
RUN cd /home/custom_export/source/cmd && go build -o /home/custom_export/source/build/main
EXPOSE 8085
CMD ["/home/custom_export/source/build/main"]