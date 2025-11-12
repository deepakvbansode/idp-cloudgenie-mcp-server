#!/bin/bash

# Test script for CloudGenie MCP Server
# This script sends test messages to the MCP server via stdin

echo "==================================="
echo "CloudGenie MCP Server Test Suite"
echo "==================================="
echo ""

# Start the server in the background
(
    # Wait a moment for the server to start
    sleep 0.5

    echo "1. Sending initialize request..."
    echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
    sleep 0.3

    echo ""
    echo "2. Sending initialized notification..."
    echo '{"jsonrpc":"2.0","method":"initialized","params":{}}'
    sleep 0.3

    echo ""
    echo "3. Listing available tools..."
    echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
    sleep 0.3

    echo ""
    echo "4. Testing health check..."
    echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"cloudgenie_health_check","arguments":{}}}'
    sleep 0.3

    echo ""
    echo "5. Testing get blueprints..."
    echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"cloudgenie_get_blueprints","arguments":{}}}'
    sleep 0.3

    echo ""
    echo "6. Testing get resources..."
    echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"cloudgenie_get_resources","arguments":{}}}'
    sleep 0.3

    echo ""
    echo "7. Listing available resources..."
    echo '{"jsonrpc":"2.0","id":6,"method":"resources/list","params":{}}'
    sleep 0.3

    echo ""
    echo "8. Listing available prompts..."
    echo '{"jsonrpc":"2.0","id":7,"method":"prompts/list","params":{}}'
    sleep 0.3

    echo ""
    echo "9. Testing ping..."
    echo '{"jsonrpc":"2.0","id":8,"method":"ping","params":{}}'
    sleep 0.3

    # Give time to see the responses
    sleep 1

) | ./cloudgenie-mcp-server

echo ""
echo "==================================="
echo "Test completed!"
echo "==================================="
echo ""
echo "Note: Some tool calls may fail if the CloudGenie backend is not running."
echo "Set CLOUDGENIE_BACKEND_URL environment variable to point to your backend."

