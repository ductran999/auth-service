#!/usr/bin/env bash

docker compose -f docker-compose.yml --env-file .env up -d db

docker compose -f docker-compose.yml --env-file .env up -d redis

docker compose -f docker-compose.yml --env-file .env run --rm migrate
