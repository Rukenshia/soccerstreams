FROM alpine
RUN apk update && apk add tzdata
ADD ca-certificates.cer /etc/ssl/certs/

WORKDIR /app
ADD web.tar.gz .

CMD ["./run.sh"]