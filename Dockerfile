FROM alpine:3.16
WORKDIR /opt/noflake
ADD noflake /opt/noflake
EXPOSE 8000
ENTRYPOINT ["./noflake"]
