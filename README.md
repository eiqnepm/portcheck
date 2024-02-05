# portcheck

<details>
  <summary>Docker Compose example</summary>

```yaml
version: "3"

services:
  portcheck:
    container_name: portcheck
    environment:
      - PORT=6881
    image: eiqnepm/portcheck:restart
    network_mode: service:gluetun
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

  deluge:
    container_name: deluge
    environment:
      - DELUGE_LOGLEVEL=error
      - PGID=1000
      - PUID=1000
      - TZ=Etc/UTC
    image: lscr.io/linuxserver/deluge
    labels:
      # Restart this container when the port is inaccessible
      - io.github.eiqnepm.portcheck.enable=true
    network_mode: service:gluetun
    restart: unless-stopped
    volumes:
      - ./deluge/config:/config
      - ./deluge/downloads:/downloads

  gluetun:
    cap_add:
      - NET_ADMIN
    container_name: gluetun
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
```

</details>

## Environment variables

| Variable       | Default | Description                                           |
| -------------- | ------- | ----------------------------------------------------- |
| `PORT`         | `6881`  | TCP port to check                                     |
| `TIMEOUT`      | `300`   | Seconds between each port check                       |
| `DIAL_TIMEOUT` | `5`     | Seconds before the port check is considered a failure |
| `IP_VERSION`   | `4`     | IP version (4 or 6)                                   |
