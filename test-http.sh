#!/bin/bash

# Test script for CloudGenie MCP Server HTTP Transport

set -e

SERVER_URL="${SERVER_URL:-http://localhost:8080}"
echo "Testing MCP Server at: $SERVER_URL"
echo "====================================="
echo ""

# Test 1: Initialize
echo "Test 1: Initialize Connection"
curl -s -X POST "$SERVER_URL" \
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
  }' | jq '.'
echo ""

# Test 2: List Tools
echo "Test 2: List Available Tools"
curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  }' | jq '.'
echo ""

# Test 3: List Prompts
echo "Test 3: List Available Prompts"
curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "prompts/list",
    "params": {}
  }' | jq '.'
echo ""

# Test 4: List Resources
echo "Test 4: List Resource Templates"
curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "resources/templates/list",
    "params": {}
  }' | jq '.'
echo ""

echo "====================================="
echo "All tests completed!"
