# portcheck

<details>
  <summary>Docker Compose example</summary>

```yaml
version: "3"

services:
  portcheck:
    depends_on:
      - deluge
    environment:
      - CLIENT=Deluge
    image: eiqnepm/portcheck:dev
    network_mode: service:gluetun
    restart: unless-stopped

  gluetun:
    cap_add:
      - NET_ADMIN
    devices:
      - /dev/net/tun:/dev/net/tun
    environment:
      - VPN_SERVICE_PROVIDER=mullvad
      - VPN_TYPE=wireguard
      - WIREGUARD_PRIVATE_KEY=👀
      - WIREGUARD_ADDRESSES=👀
      - SERVER_CITIES=Amsterdam
      - OWNED_ONLY=yes
      - FIREWALL_VPN_INPUT_PORTS=6881
    image: qmcgaw/gluetun
    ports:
      - 8112:8112 # Deluge
    restart: unless-stopped
    volumes:
      - ./gluetun:/gluetun

  deluge:
    environment:
      - DELUGE_LOGLEVEL=error
      - PGID=1000
      - PUID=1000
      - TZ=Etc/UTC
    image: lscr.io/linuxserver/deluge
    network_mode: service:gluetun
    restart: unless-stopped
    volumes:
      - ./deluge/config:/config
      - ./deluge/downloads:/downloads
```

</details>

## Environment variables

| Variable              | Default                                      | Description                                           |
| --------------------- | -------------------------------------------- | ----------------------------------------------------- |
| `CLIENT`              | `qBittorrent`                                | BitTorrent client (`qBittorrent` or `Deluge`)         |
| `CLIENT_PORT`         | `6881`                                       | Client incoming connections port                      |
| `CLIENT_WEBUI_SCHEME` | `http`                                       | Client WebUI scheme                                   |
| `CLIENT_WEBUI_HOST`   | `localhost`                                  | Client WebUI host                                     |
| `CLIENT_WEBUI_PORT`   | `8080` (qBittorrent) `8112` (Deluge)         | Client WebUI port                                     |
| `CLIENT_USERNAME`     | `admin`                                      | Client WebUI username (not required for Deluge)       |
| `CLIENT_PASSWORD`     | `adminadmin` (qBittorrent) `deluge` (Deluge) | Client WebUI password                                 |
| `TIMEOUT`             | `300`                                        | Seconds between each port check                       |
| `DIAL_TIMEOUT`        | `5`                                          | Seconds before the port check is considered a failure |
