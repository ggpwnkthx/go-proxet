# go-proxet
A very simple, multi-threaded, IO proxy/relay written in GoLang.
# Usage
Each argument begins with a named network, which is immediately followed by a comma, and ends with a path or address. An example of a single argument is below.
```
tcp,localhost:443
```
A relay is set up with each argument pair. An example of an argument pair is below:
```
unix,/srv/dockerhub.sock tcp,hub.docker.com:443
```
The above would create a unix socket on the host machine and relay traffic to port 443 at hub.docker.com.

The supported named networks are as follows:
 * ```tcp```
 * ```tcp4``` (IPv4-only)
 * ```tcp6``` (IPv6-only)
 * ```udp```
 * ```udp4``` (IPv4-only)
 * ```udp6``` (IPv6-only)
 * ```ip```
 * ```ip4``` (IPv4-only)
 * ```ip6``` (IPv6-only)
 * ```unix```
 * ```unixgram```
 * ```unixpacket``` 
## Best if used via Docker
Docker run example:
```
SOCK_PATH=/path/to/unix/sockets
docker run \
    --network overlay_network_name \
    -v $SOCK_PATH:$SOCK_PATH \
    ggpwnkthx/go-proxet:latest \
        unix,$SOCK_PATH/http.sock tcp,container_alias:80 \
        unix,$SOCK_PATH/https.sock tcp,container_alias:443
```