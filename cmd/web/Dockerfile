FROM alpine
RUN apk update && apk add tzdata

WORKDIR /app
ADD web.tar.gz .

CMD ["./run.sh"]