#!/bin/bash
# Initialize Go Microservice Template
# Usage: ./scripts/init.sh <module-name> <service-name>
# Example: ./scripts/init.sh github.com/myorg/myservice myservice

set -e

CURRENT_MODULE="github.com/minisource/template_go"
NEW_MODULE="${1:-}"
SERVICE_NAME="${2:-}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_usage() {
    echo "Usage: $0 <module-name> <service-name>"
    echo ""
    echo "Arguments:"
    echo "  module-name   Go module path (e.g., github.com/your-org/your-service)"
    echo "  service-name  Service name for configs and docker (e.g., your-service)"
    echo ""
    echo "Example:"
    echo "  $0 github.com/mycompany/payment-service payment-service"
}

# Validate arguments
if [ -z "$NEW_MODULE" ] || [ -z "$SERVICE_NAME" ]; then
    echo -e "${RED}Error: Missing required arguments${NC}"
    echo ""
    print_usage
    exit 1
fi

if [ "$NEW_MODULE" = "github.com/your-org/your-service" ]; then
    echo -e "${RED}Error: Please provide your actual module name${NC}"
    print_usage
    exit 1
fi

echo -e "${GREEN}Initializing Go Microservice Template${NC}"
echo "  Module: $NEW_MODULE"
echo "  Service: $SERVICE_NAME"
echo ""

# Get script directory (works even if called from another location)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

echo -e "${YELLOW}[1/5]${NC} Replacing module name in go.mod..."
sed -i "s|$CURRENT_MODULE|$NEW_MODULE|g" src/go.mod

echo -e "${YELLOW}[2/5]${NC} Replacing imports in all Go files..."
find src -name "*.go" -type f -exec sed -i "s|$CURRENT_MODULE|$NEW_MODULE|g" {} \;

echo -e "${YELLOW}[3/5]${NC} Updating configuration files..."
find src/config -name "*.yml" -type f -exec sed -i "s|DiviPay|$SERVICE_NAME|g" {} \;

echo -e "${YELLOW}[4/5]${NC} Updating Docker files..."
sed -i "s|backend|$SERVICE_NAME|g" src/Dockerfile
sed -i "s|backend|$SERVICE_NAME|g" docker/docker-compose.yml 2>/dev/null || true

echo -e "${YELLOW}[5/5]${NC} Running go mod tidy..."
cd src && go mod tidy

echo ""
echo -e "${GREEN}✅ Project initialized successfully!${NC}"
echo ""
echo "Next steps:"
echo "  1. Update src/config/config-development.yml with your database settings"
echo "  2. Update src/docs/docs.go with your API info"
echo "  3. Run 'make run' or 'cd src && go run ./cmd/main.go' to start the server"
echo ""
echo "Useful commands:"
echo "  make run        - Run the application"
echo "  make test       - Run tests"
echo "  make swagger    - Generate API docs"
echo "  make docker-run - Run with Docker"
