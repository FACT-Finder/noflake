FROM scratch
WORKDIR /opt/noflake
ADD noflake /opt/noflake
EXPOSE 8000
ENTRYPOINT ["./noflake"]
