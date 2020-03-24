user  nginx;
worker_processes  1;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;

events {
    worker_connections  1024;
}

http {
  server {
    listen $PORT default_server;

    location / {
      proxy_pass $TARGET;
    }
  }
}