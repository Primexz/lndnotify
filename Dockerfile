FROM alpine
COPY lndnotify /usr/bin/lndnotify
ENTRYPOINT ["/usr/bin/lndnotify"]