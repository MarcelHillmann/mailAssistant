FROM alpine:3.11
LABEL maintainer="Marcel Hillmann"
COPY mailAssistant /usr/bin/mailAssistant
VOLUME ["/opt/mailAssistant/resources"]
WORKDIR /opt/mailAssistant
USER nobody:nogroup
ENTRYPOINT ["/usr/bin/mailAssistant"]