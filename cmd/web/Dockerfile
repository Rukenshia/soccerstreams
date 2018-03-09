FROM alpine
RUN apk update && apk add tzdata ca-certificates
WORKDIR /app

ADD web.tar.gz .
CMD ["./run.sh"]