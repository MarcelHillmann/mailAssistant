FROM alpine:3.11
LABEL maintainer="Marcel Hillmann"
COPY mailAssistant /usr/bin/mailAssistant
VOLUME ["/opt/mailAssistant/resources"]
WORKDIR /opt/mailAssistant
COPY entrypoint.sh /opt/mailAssistant/entrypoint.sh
RUN apk add --no-cache bash libc-dev inotify-tools && chmod +x /opt/mailAssistant/entrypoint.sh
USER nobody:nogroup
ENTRYPOINT /opt/mailAssistant/entrypoint.sh