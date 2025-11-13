# HTTP Transport Migration Summary

## Overview

The CloudGenie MCP Server has been successfully migrated from stdio transport to HTTP transport using the official MCP Go SDK's `StreamableHTTPHandler`. This enables deployment to Kubernetes and allows AI clients to connect over HTTP.

## Changes Made

### 1. Main Server Code (main.go)

**Before (stdio):**

```go
// Run blocks until the connection is closed
if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
    log.Fatalf("[MAIN] Server error: %v", err)
}
```

**After (HTTP):**

```go
// Create HTTP handler that serves the MCP server
handler := mcp.NewStreamableHTTPHandler(
    func(r *http.Request) *mcp.Server {
        return server
    },
    nil, // Use default options
)

// Set up HTTP server
httpServer := &http.Server{
    Addr:    fmt.Sprintf(":%s", port),
    Handler: handler,
}

// Start the HTTP server with graceful shutdown
if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
    log.Fatalf("[MAIN] Server error: %v", err)
}
```

### 2. Key Features Added

- **HTTP Server**: Listens on configurable port (default 8080)
- **Graceful Shutdown**: Handles SIGTERM/SIGINT signals
- **Environment Configuration**: `MCP_HTTP_PORT` environment variable
- **Health Checks**: Ready for Kubernetes liveness/readiness probes
- **Multiple Clients**: Can handle concurrent HTTP connections

### 3. Protocol Details

The server now accepts HTTP POST requests with JSON-RPC 2.0 payloads:

**Endpoint**: `POST http://<server>:<port>/`

**Request Format**:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list",
  "params": {}
}
```

**Response Format**:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [...]
  }
}
```

## Deployment Architecture

```
┌──────────────────────────────────────┐
│         Kubernetes Cluster            │
│                                       │
│  ┌────────────────────────────────┐  │
│  │      Ingress / Load Balancer   │  │
│  │   (External HTTP Endpoint)      │  │
│  └──────────────┬─────────────────┘  │
│                 │                     │
│  ┌──────────────▼─────────────────┐  │
│  │      Service (ClusterIP)        │  │
│  │    cloudgenie-mcp-server:8080  │  │
│  └──────────────┬─────────────────┘  │
│                 │                     │
│       ┌─────────┴─────────┐          │
│       │                   │          │
│  ┌────▼────┐        ┌────▼────┐     │
│  │  Pod 1  │        │  Pod 2  │     │
│  │  :8080  │        │  :8080  │     │
│  └────┬────┘        └────┬────┘     │
│       │                   │          │
│       └─────────┬─────────┘          │
│                 │                     │
│  ┌──────────────▼─────────────────┐  │
│  │    CloudGenie Backend Service   │  │
│  │         :50051                   │  │
│  └─────────────────────────────────┘  │
└──────────────────────────────────────┘
```

## Files Created/Modified

### New Files

1. **Dockerfile** - Multi-stage build for container image
2. **k8s-deployment.yaml** - Complete Kubernetes manifests
   - Namespace
   - ConfigMap
   - Deployment (2 replicas)
   - Service (ClusterIP)
   - Ingress (optional)
3. **DEPLOYMENT.md** - Comprehensive deployment guide
4. **.dockerignore** - Optimize Docker builds

### Modified Files

1. **main.go** - Migrated from stdio to HTTP transport
   - Added HTTP server setup
   - Added graceful shutdown
   - Added port configuration

## Environment Variables

| Variable                 | Required | Default                  | Purpose                    |
| ------------------------ | -------- | ------------------------ | -------------------------- |
| `CLOUDGENIE_BACKEND_URL` | Yes      | `http://localhost:50051` | CloudGenie backend API URL |
| `MCP_HTTP_PORT`          | No       | `8080`                   | HTTP server listen port    |

## Client Connection Examples

### cURL

```bash
curl -X POST http://cloudgenie-mcp-server:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  }'
```

### Python

```python
import requests

response = requests.post(
    "http://cloudgenie-mcp-server:8080",
    json={
        "jsonrpc": "2.0",
        "id": 1,
        "method": "tools/call",
        "params": {
            "name": "create_resource",
            "arguments": {
                "name": "my-app",
                "type": "application",
                "blueprint_id": "webapp-blueprint",
                "properties": {}
            }
        }
    }
)
print(response.json())
```

### JavaScript/TypeScript

```typescript
const response = await fetch("http://cloudgenie-mcp-server:8080", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    jsonrpc: "2.0",
    id: 1,
    method: "tools/list",
    params: {},
  }),
});
const data = await response.json();
```

## Deployment Steps

### 1. Build Docker Image

```bash
docker build -t your-registry/cloudgenie-mcp-server:latest .
docker push your-registry/cloudgenie-mcp-server:latest
```

### 2. Deploy to Kubernetes

```bash
# Update image in k8s-deployment.yaml
vim k8s-deployment.yaml

# Apply manifests
kubectl apply -f k8s-deployment.yaml

# Verify deployment
kubectl get pods -n cloudgenie
kubectl logs -n cloudgenie -l app=cloudgenie-mcp-server -f
```

### 3. Test the Deployment

```bash
# Port forward for testing
kubectl port-forward -n cloudgenie svc/cloudgenie-mcp-server 8080:8080

# Test locally
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
```

## Advantages of HTTP Transport

### For Kubernetes Deployment

1. **Native HTTP Support**: Standard Kubernetes Services and Ingress
2. **Load Balancing**: Built-in with K8s Service
3. **Health Checks**: HTTP-based liveness/readiness probes
4. **Horizontal Scaling**: Multiple pods can handle requests
5. **Service Discovery**: DNS-based service discovery
6. **TLS Termination**: At Ingress level

### For Clients

1. **Standard Protocol**: HTTP/JSON-RPC is widely supported
2. **Firewall Friendly**: Port 80/443 usually open
3. **Language Agnostic**: Any HTTP client works
4. **Debugging**: Easy to test with curl/Postman
5. **Monitoring**: Standard HTTP metrics

## Backward Compatibility

The server no longer supports stdio transport. To use stdio (e.g., for local Claude Desktop):

```go
// In main.go, replace HTTP setup with:
if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
    log.Fatalf("[MAIN] Server error: %v", err)
}
```

## Performance Considerations

- **Concurrent Connections**: HTTP handler supports multiple simultaneous clients
- **Resource Limits**: Configure in k8s-deployment.yaml
  - Default: 100m CPU request, 500m CPU limit
  - Default: 128Mi memory request, 512Mi memory limit
- **Scaling**: HPA can scale based on CPU/memory metrics
- **Timeout**: Default 30s for CloudGenie API calls

## Security

1. **Network Policies**: Restrict access to authorized namespaces
2. **TLS**: Use Ingress with valid certificates
3. **Authentication**: Add API key middleware if needed
4. **RBAC**: Kubernetes role-based access control
5. **Secrets**: Store sensitive data in K8s Secrets

## Monitoring

### Logs

```bash
kubectl logs -n cloudgenie -l app=cloudgenie-mcp-server --tail=100 -f
```

### Metrics (Optional)

Add Prometheus metrics endpoint:

```go
// In main.go
http.Handle("/metrics", promhttp.Handler())
```

### Tracing (Optional)

Integrate OpenTelemetry for distributed tracing.

## Troubleshooting

### Common Issues

1. **Pod CrashLoopBackOff**

   - Check backend URL is correct
   - Verify backend service is accessible
   - Check resource limits

2. **Connection Refused**

   - Verify service/pod is running
   - Check service ports match container ports
   - Test with kubectl port-forward

3. **Slow Response**
   - Check backend API latency
   - Increase resource limits
   - Enable connection pooling

## Next Steps

1. ✅ HTTP transport implemented
2. ✅ Kubernetes manifests created
3. ✅ Deployment guide written
4. ⏳ Update main documentation (README_IDP.md)
5. ⏳ Add authentication middleware
6. ⏳ Implement metrics endpoint
7. ⏳ Add distributed tracing

## References

- MCP SDK: https://github.com/modelcontextprotocol/go-sdk
- StreamableHTTPHandler: https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk/mcp#StreamableHTTPHandler
- Kubernetes Documentation: https://kubernetes.io/docs/
