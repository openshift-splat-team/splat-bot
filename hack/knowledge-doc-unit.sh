#!/bin/sh
echo "Running Unit Tests Against SPLAT Bot Doc"
go clean -testcache
PROMPT_PATH=$(pwd)/../splat-bot-doc/knowledge_prompts go test ./... -v
