FROM alpine:3.11
ADD --chown=nobody:nobody mailAssistant.amd64_linux /usr/local/bin/mailAssistant
RUN chmod 0500 /usr/local/bin/mailAssistant
USER nobody:nobody
VOLUME /opt/mailAssistant
WORKDIR /opt/mailAssistant
ENTRYPOINT mailAssistant