FROM scratch
LABEL maintainer="Marcel Hillmann"
COPY mailAssistant.alpine /usr/bin/mailAssistant
VOLUME ["/opt/mailAssistant/resources"]
WORKDIR /opt/mailAssistant
ENTRYPOINT /usr/local/bin/mailAssistant