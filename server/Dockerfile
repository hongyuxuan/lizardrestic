FROM alpine:3.19.1

RUN apk add sqlite
        
COPY bin /opt/bin
COPY configuration /opt/configuration
COPY cache /opt/cache
COPY etc /opt/etc
COPY docs /opt/docs
WORKDIR /opt
RUN chmod +x /opt/bin/lizardrestic-server
EXPOSE 7138
ENTRYPOINT [ "./bin/lizardrestic-server" ]