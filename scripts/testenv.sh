#!/usr/bin/env bash

if [ -f .test.env ]; then
    set -a
    source .test.env
    set +a
fi

docker compose -f docker-compose.yml up -d db

docker compose -f docker-compose.yml up -d redis

docker compose -f docker-compose.yml run --rm migrate
