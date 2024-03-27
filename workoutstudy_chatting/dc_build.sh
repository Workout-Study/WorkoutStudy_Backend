#!/bin/bash

# 이미지 빌드
echo "Docker 이미지 빌드 중..."
docker build -t chatting-service:latest .