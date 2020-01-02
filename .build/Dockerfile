FROM alpine:3.11
LABEL maintainer="Marcel Hillmann"
COPY mailAssistant /usr/bin/mailAssistant
RUN apk add --no-cache bash libc-dev inotify-tools
VOLUME ["/opt/mailAssistant/resources"]
COPY entrypoint.sh /opt/mailAssistant
WORKDIR /opt/mailAssistant
USER nobody:nogroup
ENTRYPOINT /usr/bin/mailAssistant run