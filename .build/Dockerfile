FROM alpine:3.11
LABEL maintainer="Marcel Hillmann"
COPY mailAssistant entrypoint.sh /usr/bin/
RUN apk add --no-cache bash libc-dev inotify-tools && chmod +x /usr/bin/entrypoint.sh
WORKDIR /opt/mailAssistant
VOLUME ["/opt/mailAssistant/resources","/opt/homes"]
RUN chmod +rwx /opt/homes
USER nobody:nogroup
ENTRYPOINT entrypoint.sh