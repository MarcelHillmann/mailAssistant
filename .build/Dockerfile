FROM scratch
LABEL maintainer="Marcel Hillmann"
COPY --chown=nobody:nobody mailAssistant.amd64_linux /usr/bin/mailAssistant
VOLUME ["/opt/mailAssistant/resources"]
WORKDIR /opt/mailAssistant
ENTRYPOINT /usr/local/bin/mailAssistant