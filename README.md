# Minimal Reverse Proxy via NGINX as Docker Image

## Configuration via ENV

- __PORT__: port the reverse proy serves on, default: `8080`
- __TARGET__: target uri to proxy (incl. schema & path), default: `https://localhost`

## Run

`docker run --rm -it --name nginx-revprox -p 8999:8999  -e PORT=8999 -e TARGET=https://postman-echo.com thlcodes/nginx-revprox

## Licence

MIT. Do whatever you want with this.