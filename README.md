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
      QBITTORRENT_PORT: "6881"
      QBITTORRENT_WEBUI_PORT: "8080"
      QBITTORRENT_USERNAME: "admin"
      QBITTORRENT_PASSWORD: "adminadmin"
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

| Variable                   | Default      | Description                                           |
| -------------------------- | ------------ | ----------------------------------------------------- |
| `QBITTORRENT_PORT`         | `6881`       | qBittorrent incoming connections port                 |
| `QBITTORRENT_WEBUI_SCHEME` | `http`       | qBittorrent WebUI scheme                              |
| `QBITTORRENT_WEBUI_HOST`   | `localhost`  | qBittorrent WebUI host                                |
| `QBITTORRENT_WEBUI_PORT`   | `8080`       | qBittorrent WebUI port                                |
| `QBITTORRENT_USERNAME`     | `admin`      | qBittorrent WebUI username                            |
| `QBITTORRENT_PASSWORD`     | `adminadmin` | qBittorrent WebUI password                            |
| `TIMEOUT`                  | `300`        | Seconds between each port check                       |
| `DIAL_TIMEOUT`             | `5`          | Seconds before the port check is considered a failure |
