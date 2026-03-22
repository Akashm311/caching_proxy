# Caching Proxy

https://github.com/Akashm311/caching_proxy
A CLI tool that starts a caching proxy server. It forwards requests to the origin server and caches the responses. If the same request is made again, it returns the cached response instead of forwarding to the origin server.

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd caching_proxy

# Build the binary
go build -o caching-proxy .
```

## Usage

### Start the Proxy Server

```bash
caching-proxy --port <number> --origin <url>
```

**Options:**
| Flag | Description | Default |
|------|-------------|---------|
| `--port` | Port on which the proxy server will run | 8080 |
| `--origin` | URL of the server to forward requests to | http://localhost:3000 |
| `--clear-cache` | Clear the cache (exits after) | false |

### Examples

Start proxy on port 3000, forwarding to dummyjson.com:

```bash
./caching-proxy --port 3000 --origin http://dummyjson.com
```

Make a request through the proxy:

```bash
curl http://localhost:3000/products
```

### Response Headers

The proxy adds an `X-Cache` header to indicate cache status:

| Header | Meaning |
|--------|---------|
| `X-Cache: MISS` | Response fetched from origin server |
| `X-Cache: HIT` | Response served from cache |

### Clear Cache

**Option 1:** While server is running, send a DELETE request:

```bash
curl -X DELETE http://localhost:3000/clear-cache
```

**Option 2:** Restart the server (cache is in-memory)

## How It Works

```
Client Request → Proxy Server → Check Cache
                                    ↓
                    ┌───────────────┴───────────────┐
                    ↓                               ↓
               Cache HIT                       Cache MISS
                    ↓                               ↓
            Return cached                   Forward to origin
            response with                         ↓
            X-Cache: HIT                   Cache response
                                                  ↓
                                           Return response
                                           with X-Cache: MISS
```

## Project Structure

```
caching_proxy/
├── main.go       # Main application code
├── go.mod        # Go module definition
└── README.md     # This file
```

## License

MIT
