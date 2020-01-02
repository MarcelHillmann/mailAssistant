FROM centos:8
LABEL maintainer="Marcel Hillmann"
COPY --chown=nobody:nobody mailAssistant.amd64_linux /usr/bin/mailAssistant
RUN apk --no-cache add ca-certificates bash; \
    chmod 0555 /usr/bin/mailAssistant
USER nobody:nobody
VOLUME ["/opt/mailAssistant/resources"]
WORKDIR /opt/mailAssistant
ENTRYPOINT /usr/local/bin/mailAssistant