# Build both dns-proxy-api (HTTP API) and dns-proxy-cli (CLI)

all: dns-proxy-api dns-proxy-cli

dns-proxy-api:
	go build -o dns-proxy-api ./cmd/dns-proxy-api

dns-proxy-cli:
	go build -o dns-proxy-cli ./cmd/dns-proxy-cli

install: all
	cp dns-proxy-api /usr/local/bin/
	cp dns-proxy-cli /usr/local/bin/

clean:
	rm -f dns-proxy-api dns-proxy-cli
