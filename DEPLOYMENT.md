# CloudGenie MCP Server - Kubernetes Deployment Guide

This guide covers deploying the CloudGenie MCP Server to Kubernetes with HTTP transport.

## Prerequisites

- Kubernetes cluster (1.20+)
- kubectl configured
- Docker registry access
- CloudGenie backend service running

## Architecture

```
┌─────────────────┐
│   MCP Clients   │
│  (AI Models)    │
└────────┬────────┘
         │ HTTP/JSON-RPC
         ▼
┌─────────────────────┐
│   Load Balancer     │
│   (k8s Service)     │
└────────┬────────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
┌────────┐ ┌────────┐
│  Pod 1 │ │  Pod 2 │
│  MCP   │ │  MCP   │
│ Server │ │ Server │
└───┬────┘ └───┬────┘
    │          │
    └────┬─────┘
         │ HTTP
         ▼
┌─────────────────┐
│   CloudGenie    │
│    Backend      │
└─────────────────┘
```

## Quick Start

### 1. Build and Push Docker Image

```bash
# Build the image
docker build -t your-registry/cloudgenie-mcp-server:latest .

# Push to registry
docker push your-registry/cloudgenie-mcp-server:latest
```

### 2. Update Kubernetes Manifest

Edit `k8s-deployment.yaml` and update:

- `image: your-registry/cloudgenie-mcp-server:latest` (line 40)
- `CLOUDGENIE_BACKEND_URL` in ConfigMap (line 12)
- Ingress host if using external access (line 93)

### 3. Deploy to Kubernetes

```bash
# Apply the manifests
kubectl apply -f k8s-deployment.yaml

# Check deployment status
kubectl get pods -n cloudgenie
kubectl get svc -n cloudgenie

# Check logs
kubectl logs -n cloudgenie -l app=cloudgenie-mcp-server --tail=100 -f
```

## Configuration

### Environment Variables

| Variable                 | Required | Default                  | Description                |
| ------------------------ | -------- | ------------------------ | -------------------------- |
| `CLOUDGENIE_BACKEND_URL` | Yes      | `http://localhost:50051` | CloudGenie backend API URL |
| `MCP_HTTP_PORT`          | No       | `8080`                   | HTTP server port           |

### ConfigMap

Update the ConfigMap in `k8s-deployment.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cloudgenie-mcp-config
  namespace: cloudgenie
data:
  CLOUDGENIE_BACKEND_URL: "http://cloudgenie-backend:50051"
  MCP_HTTP_PORT: "8080"
```

## Accessing the Service

### Internal Access (within cluster)

```bash
# From another pod in the same namespace
curl http://cloudgenie-mcp-server:8080

# From another namespace
curl http://cloudgenie-mcp-server.cloudgenie.svc.cluster.local:8080
```

### External Access

If you deployed the Ingress:

```bash
# Test the endpoint
curl -X POST https://cloudgenie-mcp.example.com \
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

### Port Forward (for testing)

```bash
# Forward local port 8080 to service
kubectl port-forward -n cloudgenie svc/cloudgenie-mcp-server 8080:8080

# Test locally
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
```

## MCP Client Integration

Configure your MCP client to connect via HTTP:

### Python Client Example

```python
import requests

MCP_SERVER_URL = "http://cloudgenie-mcp-server.cloudgenie.svc.cluster.local:8080"

def call_mcp_tool(tool_name, arguments):
    payload = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "tools/call",
        "params": {
            "name": tool_name,
            "arguments": arguments
        }
    }

    response = requests.post(MCP_SERVER_URL, json=payload)
    return response.json()

# Example: Create a resource
result = call_mcp_tool("create_resource", {
    "name": "my-app",
    "type": "application",
    "blueprint_id": "webapp-blueprint",
    "properties": {
        "framework": "react"
    }
})
```

### TypeScript/JavaScript Client Example

```typescript
const MCP_SERVER_URL =
  "http://cloudgenie-mcp-server.cloudgenie.svc.cluster.local:8080";

async function callMcpTool(toolName: string, args: any) {
  const response = await fetch(MCP_SERVER_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      jsonrpc: "2.0",
      id: 1,
      method: "tools/call",
      params: {
        name: toolName,
        arguments: args,
      },
    }),
  });

  return response.json();
}
```

## Scaling

### Horizontal Pod Autoscaler

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: cloudgenie-mcp-hpa
  namespace: cloudgenie
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: cloudgenie-mcp-server
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

Apply with:

```bash
kubectl apply -f hpa.yaml
```

## Monitoring

### Check Health

```bash
# Check pod health
kubectl get pods -n cloudgenie -w

# View logs
kubectl logs -n cloudgenie -l app=cloudgenie-mcp-server --tail=100 -f

# Describe pod (for troubleshooting)
kubectl describe pod -n cloudgenie <pod-name>
```

### Prometheus Metrics (Optional)

Add metrics endpoint to main.go and configure Prometheus scraping:

```yaml
annotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "8080"
  prometheus.io/path: "/metrics"
```

## Troubleshooting

### Pod Not Starting

```bash
# Check pod status
kubectl get pods -n cloudgenie

# Check events
kubectl get events -n cloudgenie --sort-by='.lastTimestamp'

# Check logs
kubectl logs -n cloudgenie <pod-name>
```

### Connection Issues

```bash
# Test service connectivity
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -n cloudgenie -- \
  curl -v http://cloudgenie-mcp-server:8080

# Check service endpoints
kubectl get endpoints -n cloudgenie cloudgenie-mcp-server
```

### Backend Connection Issues

```bash
# Test backend connectivity from pod
kubectl exec -it -n cloudgenie <pod-name> -- wget -O- http://cloudgenie-backend:50051/health
```

## Security Considerations

1. **Network Policies**: Restrict traffic to authorized clients
2. **Authentication**: Add API key or OAuth for production
3. **TLS/SSL**: Use HTTPS with valid certificates
4. **Resource Limits**: Set appropriate CPU/memory limits
5. **RBAC**: Use Kubernetes RBAC for access control

## Updates and Rollbacks

### Rolling Update

```bash
# Update image
kubectl set image deployment/cloudgenie-mcp-server \
  mcp-server=your-registry/cloudgenie-mcp-server:v2.0.0 \
  -n cloudgenie

# Check rollout status
kubectl rollout status deployment/cloudgenie-mcp-server -n cloudgenie
```

### Rollback

```bash
# Rollback to previous version
kubectl rollout undo deployment/cloudgenie-mcp-server -n cloudgenie

# Rollback to specific revision
kubectl rollout undo deployment/cloudgenie-mcp-server --to-revision=2 -n cloudgenie
```

## Production Checklist

- [ ] Configure resource limits
- [ ] Set up monitoring and alerting
- [ ] Configure HPA for autoscaling
- [ ] Enable TLS/SSL
- [ ] Configure network policies
- [ ] Set up backup/disaster recovery
- [ ] Document runbooks
- [ ] Configure log aggregation
- [ ] Set up health checks
- [ ] Test failover scenarios

## Support

For issues and questions:

- GitHub: https://github.com/deepakvbansode/idp-cloudgenie-mcp-server
- Documentation: See README_IDP.md for MCP protocol details
