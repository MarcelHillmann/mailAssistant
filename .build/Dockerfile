FROM alpine:3.11
ADD --chown=nobody:nobody mailAssistant.amd64_linux /usr/local/mailAssistant
RUN chmod o-rw,o+x,g-rw,g+x,o-rw,o+x /usr/local/mailAssistant
USER nobody:nobody
VOLUME /opt/mailAssistant
WORKDIR /opt/mailAssistant
ENTRYPOINT /usr/local/mailAssistant