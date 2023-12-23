FROM ubuntu

COPY explore-me /root/explore-me

RUN chmod +x /root/explore-me

VOLUME /data
WORKDIR /data

EXPOSE 8080

CMD ["/root/explore-me"]