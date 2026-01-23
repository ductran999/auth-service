#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Directories
PKG_GEN_PATH="gen/openapi"

echo -e "${GREEN}🚀 Starting API generation...${NC}"

# Create directories if they don't exist
mkdir -p ${PKG_GEN_PATH}

# OpenAPI spec file
OPENAPI_SPEC="./api/openapi.yaml"

# Check if OpenAPI spec exists
if [ ! -f "$OPENAPI_SPEC" ]; then
    echo -e "${RED}❌ OpenAPI spec not found at ${OPENAPI_SPEC}${NC}"
    echo -e "${YELLOW}Please create the OpenAPI specification file first${NC}"
    exit 1
fi

# Clean up existing generated files
echo -e "${YELLOW}🧹 Cleaning up existing generated files...${NC}"
rm -f ${PKG_GEN_PATH}/*

# Generate types
echo -e "${YELLOW}🔧 Generating Go types...${NC}"
oapi-codegen -generate "types" -package openapi -o ${PKG_GEN_PATH}/types.gen.go ${OPENAPI_SPEC}

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Go types generated${NC}"
else
    echo -e "${RED}❌ Failed to generate Go types${NC}"
    exit 1
fi

# Generate server
echo -e "${YELLOW}🔧 Generating Go server code...${NC}"
oapi-codegen -generate "gin" -package openapi -o ${PKG_GEN_PATH}/server.gen.go ${OPENAPI_SPEC}

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Go server code generated${NC}"
else
    echo -e "${RED}❌ Failed to generate Go server code${NC}"
    exit 1
fi

# Update go.mod
echo -e "${YELLOW}📦 Updating go.mod...${NC}"
go mod tidy

echo -e "${GREEN}🎉 API generation completed successfully!${NC}"
echo -e "${GREEN}Generated files:${NC}"
echo -e "  📄 OpenAPI spec: ${OPENAPI_SPEC}"
echo -e "  🔧 Go types: ${PKG_GEN_PATH}/types.gen.go"
echo -e "  🔧 Go server: ${PKG_GEN_PATH}/server.gen.go"
