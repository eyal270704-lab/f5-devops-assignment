# F5 DevOps Assignment

This repository demonstrates a containerized Nginx deployment with automated testing and CI/CD, built for the F5 Networks DevOps internship assignment.

## Overview

The project consists of two Docker containers orchestrated with Docker Compose:

1. **Nginx Container** - Custom Nginx server with:
   - HTTP/HTTPS server on ports 80/443 serving custom HTML
   - Error server on port 8080 returning 403 Forbidden
   - Rate limiting (5 requests/second)
   - Self-signed SSL certificate for HTTPS

2. **Test Container** - Go-based test suite that validates:
   - HTTP server functionality (200 OK)
   - HTTPS server functionality (200 OK with self-signed cert)
   - Error server (403 Forbidden)
   - Rate limiting (503 after threshold)

## Quick Start

### Prerequisites
- Docker
- Docker Compose

### Build and Run

```bash
# Clone the repository
git clone https://github.com/yourusername/f5-devops-assignment.git
cd f5-devops-assignment

# Build and run all containers
docker compose up --build

# The test container will automatically run and exit with status 0 if all tests pass
```

### Manual Testing

```bash
# Test HTTP server (port 80)
curl http://localhost:80

# Test HTTPS server (port 443)
curl -k https://localhost:443

# Test error server (port 8080)
curl http://localhost:8080

# Test rate limiting (send 20 rapid requests)
for i in {1..20}; do curl -s -o /dev/null -w "%{http_code}\n" http://localhost; done
```

## Architecture

### Project Structure

```
f5-devops-assignment/
├── nginx/
│   ├── Dockerfile              # Nginx container build
│   ├── nginx.conf              # Main Nginx configuration
│   ├── default-http.conf       # HTTP/HTTPS server block
│   ├── default-error.conf      # Error server block (port 8080)
│   ├── index.html              # Custom HTML page
│   └── generate_cert.sh        # Self-signed SSL cert generator
├── test/
│   ├── Dockerfile              # Test container build (multi-stage)
│   ├── main.go                 # Go test suite
│   ├── go.mod                  # Go module definition
│   └── go.sum                  # Go dependencies (empty - stdlib only)
├── docker-compose.yml          # Container orchestration
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions CI/CD workflow
└── README.md
```

### Nginx Configuration

#### Main Config (`nginx.conf`)
- Worker processes set to `auto` for optimal performance
- Rate limiting zone defined: `limit_req_zone $binary_remote_addr zone=global:2m rate=5r/s;`
- Includes separate server block configuration files

#### HTTP/HTTPS Server (`default-http.conf`)
- Listens on port 80 (HTTP) and 443 (HTTPS)
- SSL certificate: `/etc/nginx/ssl/cert.pem`
- SSL key: `/etc/nginx/ssl/key.pem`
- Rate limiting applied: `limit_req zone=global burst=5 nodelay;`
- Serves HTML from `/usr/share/nginx/html`

#### Error Server (`default-error.conf`)
- Listens on port 8080
- Returns 403 Forbidden for all requests
- Rate limiting applied: `limit_req zone=global burst=5 nodelay;`

### HTTPS Implementation

The Nginx container generates a self-signed SSL certificate at build time:

```bash
# generate_cert.sh creates:
# - /etc/nginx/ssl/cert.pem (certificate)
# - /etc/nginx/ssl/key.pem (private key)
# Valid for 365 days
```

**Note**: This is a self-signed certificate suitable for development/testing only. In production, use a certificate from a trusted CA.

### Rate Limiting

Rate limiting is configured globally and applied to all server blocks:

- **Rate**: 5 requests per second per IP address
- **Zone**: 2MB memory (stores ~32,000 IP addresses)
- **Burst**: 5 extra requests allowed temporarily
- **Behavior**: Returns 503 Service Temporarily Unavailable when exceeded

#### Changing Rate Limiting Threshold

To modify the rate limit, edit `nginx/nginx.conf`:

```nginx
# Change the rate (currently 5r/s):
limit_req_zone $binary_remote_addr zone=global:2m rate=10r/s;  # 10 requests/sec

# Change the burst (currently 5):
# Edit in server block configs (default-http.conf, default-error.conf):
limit_req zone=global burst=10 nodelay;  # Allow 10 burst requests
```

## CI/CD with GitHub Actions

The repository includes a GitHub Actions workflow (`.github/workflows/ci.yml`) that:

1. Triggers on push and pull requests to `main`/`master`
2. Builds both containers using Docker Compose
3. Runs the test suite automatically
4. Creates and uploads a test result artifact
5. Fails the workflow if tests don't pass

### Viewing Test Results

After each workflow run:
1. Go to the "Actions" tab in your GitHub repository
2. Click on the latest workflow run
3. Download the `test-result` artifact
4. Check for `succeeded` or `failed` file

## Test Suite

The Go-based test suite (`test/main.go`) performs the following tests:

1. **HTTP Server Test**: Verifies port 80 returns 200 with HTML content
2. **HTTPS Server Test**: Verifies port 443 returns 200 (accepts self-signed cert)
3. **Error Server Test**: Verifies port 8080 returns 403
4. **Rate Limiting Test**: Sends 20 rapid requests and verifies some return 503

All tests must pass for the container to exit with status 0.

## Development

### Building Individual Containers

```bash
# Build only nginx
docker build -t f5-nginx nginx/

# Build only test
docker build -t f5-test test/

# Run nginx standalone
docker run -d -p 80:80 -p 443:443 -p 8080:8080 --name nginx-test f5-nginx

# Run test against running nginx
docker run --network container:nginx-test f5-test
```

### Viewing Logs

```bash
# View nginx logs
docker logs <nginx-container-id>

# Follow nginx logs in real-time
docker logs -f <nginx-container-id>

# View Docker Compose logs
docker compose logs
docker compose logs -f  # follow mode
```

### Stopping Containers

```bash
# Stop all containers
docker compose down

# Stop and remove all containers, networks, and volumes
docker compose down -v
```