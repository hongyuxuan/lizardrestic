FROM alpine:3.19.1
COPY bin /opt/bin
COPY configuration /opt/configuration
COPY cache /opt/cache
COPY etc /opt/etc
WORKDIR /opt
RUN chmod +x /opt/bin/lizardrestic-agent
EXPOSE 7038
ENTRYPOINT [ "./bin/lizardrestic-agent" ]