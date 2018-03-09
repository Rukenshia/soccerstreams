FROM scratch
ADD ca-certificates.pem /etc/ssl/certs/
ADD sweeper /
RUN chmod +x /sweeper
CMD ["/sweeper"]