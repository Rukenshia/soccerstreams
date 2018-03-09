FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD zoneinfo /usr/share/zoneinfo
ADD sweeper /
CMD ["/sweeper"]