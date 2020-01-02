FROM golang:1.13-alpine3.11
COPY mailAssistant.amd64_linux /usr/local/bin/mailAssistant
RUN chmod 0500 /usr/local/bin/mailAssistant
VOLUME ["/opt/mailAssistant/resources"]

USER nobody:nobody
WORKDIR /opt/mailAssistant
ENTRYPOINT /usr/local/bin/mailAssistant