FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD zoneinfo /usr/share/zoneinfo
ADD agent /

EXPOSE 9000

CMD ["/agent"]