FROM ubuntu

COPY explorer /root/explorer

RUN chmod +x /root/explorer

VOLUME /data
WORKDIR /data

EXPOSE 8080

CMD ["/root/explorer"]