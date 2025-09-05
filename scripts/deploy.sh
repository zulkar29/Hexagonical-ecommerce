#!/bin/bash

# Production Deployment Script for Hexagonal E-commerce SaaS
# Usage: ./scripts/deploy.sh [environment]
# Environment: production (default) | staging

set -e

ENVIRONMENT=${1:-production}
COMPOSE_FILE="docker-compose.prod.yml"
ENV_FILE=".env.production"

echo "🚀 Starting deployment for $ENVIRONMENT environment..."

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   echo "❌ This script should not be run as root for security reasons"
   exit 1
fi

# Check prerequisites
echo "📋 Checking prerequisites..."

if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed"
    exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
    echo "❌ Environment file $ENV_FILE not found"
    echo "   Copy .env.production.example to $ENV_FILE and configure it"
    exit 1
fi

# Load environment variables
set -a
source $ENV_FILE
set +a

echo "✅ Prerequisites checked"

# Create required directories
echo "📁 Creating required directories..."
mkdir -p uploads
mkdir -p backups
mkdir -p logs
echo "✅ Directories created"

# Pull latest code (if in production)
if [ "$ENVIRONMENT" = "production" ]; then
    echo "📥 Pulling latest code..."
    git pull origin main
    echo "✅ Code updated"
fi

# Build and deploy services
echo "🔨 Building services..."
docker-compose -f $COMPOSE_FILE build --no-cache

echo "🏃 Starting services..."
docker-compose -f $COMPOSE_FILE up -d

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
timeout 60s bash -c 'while ! docker-compose -f '$COMPOSE_FILE' exec -T postgres pg_isready -U '$POSTGRES_USER' -d '$POSTGRES_DB'; do sleep 2; done'

# Run database migrations
echo "🗄️ Running database migrations..."
docker-compose -f $COMPOSE_FILE exec -T backend ./main migrate

# Check service health
echo "🩺 Checking service health..."
sleep 10

# Check backend health
if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "❌ Backend health check failed"
    echo "📋 Backend logs:"
    docker-compose -f $COMPOSE_FILE logs backend --tail=20
    exit 1
fi

# Check database connection
if ! docker-compose -f $COMPOSE_FILE exec -T postgres pg_isready -U $POSTGRES_USER -d $POSTGRES_DB > /dev/null 2>&1; then
    echo "❌ Database health check failed"
    exit 1
fi

echo "✅ All services are healthy"

# Display service status
echo "📊 Service status:"
docker-compose -f $COMPOSE_FILE ps

# Display access URLs
echo ""
echo "🌐 Access URLs:"
echo "   Main Site: https://$DOMAIN"
echo "   Admin Panel: https://$ADMIN_SUBDOMAIN.$DOMAIN"
echo "   API Health: https://$DOMAIN/api/health"

# Display next steps
echo ""
echo "🎉 Deployment completed successfully!"
echo ""
echo "📝 Next steps:"
echo "   1. Configure DNS records for your domain"
echo "   2. Set up monitoring and alerts"
echo "   3. Configure backup schedule"
echo "   4. Review security settings"
echo ""
echo "📚 For more information, see DEPLOYMENT.md"

# Setup log rotation (optional)
if [ "$ENVIRONMENT" = "production" ]; then
    echo "⚙️  Setting up log rotation..."
    cat > /tmp/docker-compose-logs << 'EOF'
/var/lib/docker/containers/*/*-json.log {
    rotate 7
    daily
    compress
    size=1M
    missingok
    delaycompress
    copytruncate
}
EOF

    if [ -d /etc/logrotate.d ]; then
        sudo cp /tmp/docker-compose-logs /etc/logrotate.d/docker-compose-logs
        echo "✅ Log rotation configured"
    fi
fi

echo "✨ Deployment script completed!"