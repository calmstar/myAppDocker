FROM ubuntu:latest

LABEL maintainer="cwx@email.com"

RUN apt-get update && apt-get install -y  nginx

COPY ./index.html /var/www/html/

EXPOSE 80

RUN nginx -t

CMD nginx -g "daemon off;"
