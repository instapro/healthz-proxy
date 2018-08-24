# healthz-proxy

Proxy requests to a target URI.

Especially useful to publicly serve targets that are behind a firewall, or are insecure to be open.

Usage:
```shell
docker pull renang/healthz-proxy:latest
docker run --publish 8080:8080 renang/healthz-proxy:latest http://numbersapi.com/42

curl http://localhost:8080
```
