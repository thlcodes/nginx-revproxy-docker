FROM nginx:alpine

ENV PORT 8080
ENV TARGET "http://localhost"

COPY ./nginx.conf.tpl /nginx.conf.tpl

EXPOSE $PORT

CMD ["/bin/sh" , "-c" , "envsubst < /nginx.conf.tpl > /etc/nginx/nginx.conf && exec nginx -g 'daemon off;'"]