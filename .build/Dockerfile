FROM golang:1.13-alpine3.11
USER nobody:nobody
COPY mailAssistant.amd64_linux /usr/local/bin/mailAssistant
RUN chmod 0555 /usr/local/bin/mailAssistant
VOLUME ["/opt/mailAssistant/resources"]
WORKDIR /opt/mailAssistant
ENTRYPOINT /usr/local/bin/mailAssistant