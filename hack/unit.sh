#!/bin/sh
echo "Running Unit Tests"
PROMPT_PATH=$(pwd)/pkg/knowledge/test/knowledge_prompts go test ./...