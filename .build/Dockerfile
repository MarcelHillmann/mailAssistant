FROM alpine:3.11
ADD --chown=nobody:nobody ./upload/mailAssistant.amd64_linux /usr/local/mailAssistant
RUN chmod o-rw,o+x,g-rw,g+x,o-rw,o+x /usr/local
USER nobody:nobody

ENTRYPOINT /usr/local/mailAssistant