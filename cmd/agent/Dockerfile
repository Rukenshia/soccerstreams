FROM scratch
ADD ca-certificates.pem /etc/ssl/certs/
ADD agent /
RUN chmod +x /agent
CMD ["/agent"]