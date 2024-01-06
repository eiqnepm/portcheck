# portcheck

<details>
  <summary>Docker Compose example</summary>

```yaml
version: "3"

services:
  portcheck:
    depends_on:
      - "qbittorrent"
    environment:
      CLIENT_PORT: "6881"
      CLIENT_WEBUI_PORT: "8080"
      CLIENT_USERNAME: "admin"
      CLIENT_PASSWORD: "adminadmin"
    image: "eiqnepm/portcheck:latest"
    network_mode: "service:gluetun"
    restart: "unless-stopped"

  gluetun:
    cap_add:
      - "NET_ADMIN"
    devices:
      - "/dev/net/tun:/dev/net/tun"
    environment:
      VPN_SERVICE_PROVIDER: "mullvad"
      VPN_TYPE: "wireguard"
      WIREGUARD_PRIVATE_KEY: "👀"
      WIREGUARD_ADDRESSES: "👀"
      SERVER_CITIES: "Amsterdam"
      OWNED_ONLY: "yes"
      FIREWALL_VPN_INPUT_PORTS: "6881"
    image: "qmcgaw/gluetun:latest"
    ports:
      # qBittorrent
      - "8080:8080"
    restart: "unless-stopped"
    volumes:
      - "./gluetun:/gluetun"

  qbittorrent:
    environment:
      PUID: "1000"
      PGID: "1000"
      TZ: "Etc/UTC"
      WEBUI_PORT: "8080"
    image: "lscr.io/linuxserver/qbittorrent:latest"
    network_mode: "service:gluetun"
    restart: "unless-stopped"
    volumes:
      - "./qbittorrent:/config"
      - "./torrents:/downloads"
```

</details>

## Environment variables

| Variable              | Default       | Description                                           |
| --------------------- | ------------- | ----------------------------------------------------- |
| `CLIENT`              | `qBittorrent` | Either `qBittorrent` or `Deluge`                      |
| `CLIENT_PORT`         | `6881`        | Client incoming connections port                      |
| `CLIENT_WEBUI_SCHEME` | `http`        | Client WebUI scheme                                   |
| `CLIENT_WEBUI_HOST`   | `localhost`   | Client WebUI host                                     |
| `CLIENT_WEBUI_PORT`   | `8080`        | Client WebUI port                                     |
| `CLIENT_USERNAME`     | `admin`       | Client WebUI username (not required for Deluge)                                 |
| `CLIENT_PASSWORD`     | `adminadmin`  | Client WebUI password                                 |
| `TIMEOUT`             | `300`         | Seconds between each port check                       |
| `DIAL_TIMEOUT`        | `5`           | Seconds before the port check is considered a failure |
