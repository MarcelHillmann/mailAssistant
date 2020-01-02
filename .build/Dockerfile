FROM golang:1.13-alpine3.11
LABEL maintainer="Marcel Hillmann"
COPY --chown=nobody:nobody mailAssistant.amd64_linux /usr/local/bin/
RUN chmod 0555 /usr/local/bin/mailAssistant
USER nobody:nobody
VOLUME ["/opt/mailAssistant/resources"]
WORKDIR /opt/mailAssistant
ENTRYPOINT /usr/local/bin/mailAssistant