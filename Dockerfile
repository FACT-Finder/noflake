FROM alpine:3.17
WORKDIR /var/lib/noflake
RUN mkdir /opt/noflake
ADD noflake /opt/noflake
EXPOSE 8000
ENTRYPOINT ["/opt/noflake/noflake"]
