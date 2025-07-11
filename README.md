# dns-proxy

dns-proxy is a lightweight HTTP server written in Go that provides an API endpoint to set DNS TXT records via cPanel's API. It is designed to automate DNS-01 challenges (such as for Let's Encrypt) or other scenarios where programmatic TXT record management is required.

## Features

- Exposes a `/set_txt` HTTP endpoint for setting TXT records
- Authenticates requests using a Bearer API key
- Reads configuration from a simple config file
- Interacts with cPanel's ZoneEdit API to add TXT records

## Configuration

Create a config file (default: `/etc/dns-proxy.conf`) with the following format:

```ini
API_KEY=your_api_key_here
cpanel_url=https://your-cpanel-domain:2083
cpanel_user=cpanel_username
cpanel_apikey=cpanel_api_token
```

- `API_KEY`: The Bearer token required for API requests
- `cpanel_url`: The base URL of your cPanel instance
- `cpanel_user`: The cPanel username
- `cpanel_apikey`: The cPanel API token

## Installation

1. **Build the binary:**

   ```sh
   go build -o dns-proxy main.go
   ```

1. **Install the binary:**

   ```sh
   cp dns-proxy /usr/local/bin/
   ```

   This will make `dns-proxy` available system-wide.

## Running as a Service (SystemD)

To run `dns-proxy` as a systemd service on Linux:

1. Create a systemd service file `/etc/systemd/system/dns-proxy.service` with the following content:

   ```ini
   [Unit]
   Description=DNS Proxy API Service
   After=network.target

   [Service]
   Type=simple
   ExecStart=/usr/local/bin/dns-proxy
   Restart=on-failure
   User=nobody
   Group=nogroup

   [Install]
   WantedBy=multi-user.target
   ```

   Adjust `User` and `Group` as needed for your environment.

1. Reload systemd and start the service:

   ```sh
   systemctl daemon-reload
   systemctl enable dns-proxy
   systemctl start dns-proxy
   ```

1. Check the status:

   ```sh
   systemctl status dns-proxy
   ```

### OpenRC (Alpine Linux)

For Alpine Linux (OpenRC), create `/etc/init.d/dns-proxy` with:

```sh
#!/sbin/openrc-run
command="/usr/local/bin/dns-proxy"
command_background="yes"
description="DNS Proxy API Service"

pidfile="/var/run/dns-proxy.pid"

start_pre() {
    checkpath --directory /var/run
}
```

Make it executable and enable/start the service:

```sh
chmod +x /etc/init.d/dns-proxy
rc-update add dns-proxy
dns-proxy start
```

## Usage

1. **Build and run the server:**

   ```sh
   go build -o dns-proxy main.go
   ./dns-proxy
   ```

   The server will listen on port `5000` by default.

1. **Send a request to set a TXT record:**

   - Endpoint: `POST /set_txt`
   - Headers:
     - `Authorization: Bearer <API_KEY>`
     - `Content-Type: application/json`
   - Body:

     ```json
     {
       "domain": "example.com",
       "key": "_acme-challenge",
       "value": "your_txt_value"
     }
     ```

   Example using `curl`:

   ```sh
   curl -X POST http://localhost:5000/set_txt \
     -H "Authorization: Bearer your_api_key_here" \
     -H "Content-Type: application/json" \
     -d '{"domain":"example.com","key":"_acme-challenge","value":"txt_value_here"}'
   ```

## Notes

- The server must have access to the cPanel API endpoint.
- The config file path can be changed in the source code if needed.
- Only TXT records are supported by this proxy.

## License

MIT License
