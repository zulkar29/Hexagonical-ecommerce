#!/bin/bash

# Production Deployment Script for E-commerce SaaS Platform
# This script helps deploy the application to production with proper checks
# Supports both local and remote server deployment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Deployment configuration
DEPLOY_MODE="local"  # local or remote
REMOTE_HOST=""
REMOTE_USER=""
REMOTE_PASSWORD=""
REMOTE_PATH="/opt/esass"
BACKUP_BEFORE_DEPLOY="true"
ROLLBACK_MODE="false"
BACKUP_NAME=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --remote)
            DEPLOY_MODE="remote"
            shift
            ;;
        --host)
            REMOTE_HOST="$2"
            shift 2
            ;;
        --user)
            REMOTE_USER="$2"
            shift 2
            ;;
        --password)
            REMOTE_PASSWORD="$2"
            shift 2
            ;;
        --path)
            REMOTE_PATH="$2"
            shift 2
            ;;
        --rollback)
            ROLLBACK_MODE="true"
            BACKUP_NAME="$2"
            shift 2
            ;;
        --no-backup)
            BACKUP_BEFORE_DEPLOY="false"
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --remote              Deploy to remote server"
            echo "  --host HOST          Remote server IP/hostname"
            echo "  --user USER          Remote server username"
            echo "  --password PASS      Remote server password"
            echo "  --path PATH          Remote deployment path (default: /opt/esass)"
            echo "  --rollback BACKUP    Rollback to specific backup"
            echo "  --no-backup          Skip pre-deployment backup"
            echo "  --help               Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                                    # Local deployment with backup"
            echo "  $0 --remote --host 1.2.3.4 --user root --password mypass"
            echo "  $0 --rollback backup_20241201_123456 # Rollback to specific backup"
            echo "  $0 --no-backup                       # Deploy without backup"
            echo ""
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🚀 E-commerce SaaS Production Deployment${NC}"
echo "=============================================="
echo -e "${BLUE}📍 Deployment Mode: ${DEPLOY_MODE}${NC}"

if [ "$DEPLOY_MODE" = "remote" ]; then
    echo -e "${BLUE}🌐 Remote Server: ${REMOTE_HOST}${NC}"
    echo -e "${BLUE}👤 Remote User: ${REMOTE_USER}${NC}"
    echo -e "${BLUE}📁 Remote Path: ${REMOTE_PATH}${NC}"
fi

if [ "$ROLLBACK_MODE" = "true" ]; then
    echo -e "${YELLOW}🔄 Rollback Mode: ${BACKUP_NAME}${NC}"
fi

# Remote deployment functions
run_remote_command() {
    local command="$1"
    if command -v sshpass >/dev/null 2>&1; then
        sshpass -p "$REMOTE_PASSWORD" ssh -o StrictHostKeyChecking=no "$REMOTE_USER@$REMOTE_HOST" "$command"
    else
        echo -e "${YELLOW}⚠️  sshpass not found. Please install it or use SSH keys for authentication${NC}"
        echo "Ubuntu/Debian: sudo apt-get install sshpass"
        echo "macOS: brew install sshpass"
        exit 1
    fi
}

copy_to_remote() {
    local local_path="$1"
    local remote_path="$2"
    if command -v sshpass >/dev/null 2>&1; then
        sshpass -p "$REMOTE_PASSWORD" scp -r -o StrictHostKeyChecking=no "$local_path" "$REMOTE_USER@$REMOTE_HOST:$remote_path"
    else
        echo -e "${YELLOW}⚠️  sshpass not found. Please install it or use SSH keys for authentication${NC}"
        echo "Ubuntu/Debian: sudo apt-get install sshpass"
        echo "macOS: brew install sshpass"
        exit 1
    fi
}

# Backup and rollback functions
create_backup() {
    local backup_dir="/opt/backups"
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_name="backup_${timestamp}"

    echo -e "${BLUE}📦 Creating pre-deployment backup: ${backup_name}${NC}"

    if [ "$DEPLOY_MODE" = "remote" ]; then
        run_remote_command "mkdir -p ${backup_dir}"
        run_remote_command "cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml exec -T postgres pg_dump -U \${POSTGRES_USER} \${POSTGRES_DB} | gzip > ${backup_dir}/${backup_name}_db.sql.gz"
        run_remote_command "cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml exec -T redis redis-cli -a \${REDIS_PASSWORD} --rdb - | gzip > ${backup_dir}/${backup_name}_redis.rdb.gz"
        run_remote_command "cd $REMOTE_PATH && tar -czf ${backup_dir}/${backup_name}_code.tar.gz --exclude=node_modules --exclude=.git ."
        if run_remote_command "[ -d '$REMOTE_PATH/uploads' ]"; then
            run_remote_command "cd $REMOTE_PATH && tar -czf ${backup_dir}/${backup_name}_uploads.tar.gz uploads/"
        fi
    else
        mkdir -p ${backup_dir}
        docker-compose -f docker-compose.prod.yml exec -T postgres pg_dump -U ${POSTGRES_USER} ${POSTGRES_DB} | gzip > ${backup_dir}/${backup_name}_db.sql.gz
        docker-compose -f docker-compose.prod.yml exec -T redis redis-cli -a ${REDIS_PASSWORD} --rdb - | gzip > ${backup_dir}/${backup_name}_redis.rdb.gz
        tar -czf ${backup_dir}/${backup_name}_code.tar.gz --exclude=node_modules --exclude=.git .
        if [ -d "uploads" ]; then
            tar -czf ${backup_dir}/${backup_name}_uploads.tar.gz uploads/
        fi
    fi

    echo -e "${GREEN}✅ Backup created: ${backup_name}${NC}"
    echo "$backup_name"
}

rollback_deployment() {
    local backup_name="$1"
    local backup_dir="/opt/backups"

    echo -e "${YELLOW}🔄 Rolling back to backup: ${backup_name}${NC}"

    if [ "$DEPLOY_MODE" = "remote" ]; then
        # Check if backup exists
        if ! run_remote_command "[ -f '${backup_dir}/${backup_name}_db.sql.gz' ]"; then
            echo -e "${RED}❌ Backup not found: ${backup_name}${NC}"
            exit 1
        fi

        # Stop services
        run_remote_command "cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml down"

        # Restore database
        run_remote_command "cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml up -d postgres redis"
        sleep 10
        run_remote_command "cd $REMOTE_PATH && zcat ${backup_dir}/${backup_name}_db.sql.gz | docker-compose -f docker-compose.prod.yml exec -T postgres psql -U \${POSTGRES_USER} -d \${POSTGRES_DB}"

        # Restore code
        run_remote_command "cd $REMOTE_PATH && tar -xzf ${backup_dir}/${backup_name}_code.tar.gz"

        # Restore uploads if they exist
        if run_remote_command "[ -f '${backup_dir}/${backup_name}_uploads.tar.gz' ]"; then
            run_remote_command "cd $REMOTE_PATH && tar -xzf ${backup_dir}/${backup_name}_uploads.tar.gz"
        fi

        # Start services
        run_remote_command "cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml up -d"

    else
        # Check if backup exists
        if [ ! -f "${backup_dir}/${backup_name}_db.sql.gz" ]; then
            echo -e "${RED}❌ Backup not found: ${backup_name}${NC}"
            exit 1
        fi

        # Stop services
        docker-compose -f docker-compose.prod.yml down

        # Restore database
        docker-compose -f docker-compose.prod.yml up -d postgres redis
        sleep 10
        zcat ${backup_dir}/${backup_name}_db.sql.gz | docker-compose -f docker-compose.prod.yml exec -T postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}

        # Restore code
        tar -xzf ${backup_dir}/${backup_name}_code.tar.gz

        # Restore uploads if they exist
        if [ -f "${backup_dir}/${backup_name}_uploads.tar.gz" ]; then
            tar -xzf ${backup_dir}/${backup_name}_uploads.tar.gz
        fi

        # Start services
        docker-compose -f docker-compose.prod.yml up -d
    fi

    echo -e "${GREEN}✅ Rollback completed successfully${NC}"
}

list_backups() {
    local backup_dir="/opt/backups"

    echo -e "${BLUE}📋 Available backups:${NC}"

    if [ "$DEPLOY_MODE" = "remote" ]; then
        run_remote_command "ls -la ${backup_dir}/*_db.sql.gz 2>/dev/null | sed 's/.*backup_\\([0-9_]*\\)_db.sql.gz/\\1/' | sort -r" || echo "No backups found"
    else
        ls -la ${backup_dir}/*_db.sql.gz 2>/dev/null | sed 's/.*backup_\([0-9_]*\)_db.sql.gz/\1/' | sort -r || echo "No backups found"
    fi
}

# Validate remote deployment parameters
if [ "$DEPLOY_MODE" = "remote" ]; then
    if [ -z "$REMOTE_HOST" ] || [ -z "$REMOTE_USER" ] || [ -z "$REMOTE_PASSWORD" ]; then
        echo -e "${RED}❌ Error: Remote deployment requires --host, --user, and --password${NC}"
        echo "Use --help for usage information"
        exit 1
    fi

    echo -e "${BLUE}🔍 Testing remote connection...${NC}"
    if run_remote_command "echo 'Connection successful'"; then
        echo -e "${GREEN}✅ Remote connection established${NC}"
    else
        echo -e "${RED}❌ Failed to connect to remote server${NC}"
        exit 1
    fi
fi

# Handle rollback mode
if [ "$ROLLBACK_MODE" = "true" ]; then
    if [ -z "$BACKUP_NAME" ]; then
        echo -e "${YELLOW}📋 Available backups:${NC}"
        list_backups
        echo ""
        echo -e "${RED}❌ Please specify a backup name to rollback to${NC}"
        echo "Usage: $0 --rollback backup_YYYYMMDD_HHMMSS"
        exit 1
    fi

    echo -e "${YELLOW}⚠️  WARNING: This will restore your system to backup: ${BACKUP_NAME}${NC}"
    echo -e "${YELLOW}⚠️  All current data and code changes will be lost!${NC}"
    echo ""
    read -p "Are you sure you want to proceed with rollback? (yes/NO): " confirm
    if [[ ! "$confirm" =~ ^[Yy][Ee][Ss]$ ]]; then
        echo -e "${BLUE}Rollback cancelled${NC}"
        exit 0
    fi

    rollback_deployment "$BACKUP_NAME"
    echo -e "${GREEN}🎉 Rollback completed successfully!${NC}"
    exit 0
fi

# Check if .env.production exists
if [ ! -f ".env.production" ]; then
    echo -e "${RED}❌ Error: .env.production file not found${NC}"
    echo -e "${YELLOW}💡 Please copy .env.production.example to .env.production and configure it${NC}"
    echo "   cp .env.production.example .env.production"
    echo "   nano .env.production"
    exit 1
fi

# Load environment variables
export $(grep -v '^#' .env.production | xargs)

# Check required environment variables
echo -e "${BLUE}🔍 Checking required environment variables...${NC}"

required_vars=(
    "DOMAIN"
    "POSTGRES_DB"
    "POSTGRES_USER"
    "POSTGRES_PASSWORD"
    "REDIS_PASSWORD"
    "JWT_SECRET"
    "SSLCOMMERZ_STORE_ID"
    "SSLCOMMERZ_STORE_PASSWORD"
)

missing_vars=()
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        missing_vars+=("$var")
    fi
done

if [ ${#missing_vars[@]} -ne 0 ]; then
    echo -e "${RED}❌ Missing required environment variables:${NC}"
    printf '%s\n' "${missing_vars[@]}"
    echo -e "${YELLOW}💡 Please configure these in .env.production${NC}"
    exit 1
fi

echo -e "${GREEN}✅ All required environment variables are set${NC}"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker is running${NC}"

# Check if ports are available
echo -e "${BLUE}🔍 Checking if required ports are available...${NC}"

check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null ; then
        echo -e "${RED}❌ Port $1 is already in use${NC}"
        echo "Process using port $1:"
        lsof -Pi :$1 -sTCP:LISTEN
        return 1
    else
        echo -e "${GREEN}✅ Port $1 is available${NC}"
        return 0
    fi
}

ports_ok=true
for port in 80 443 5432 6379; do
    if ! check_port $port; then
        ports_ok=false
    fi
done

if [ "$ports_ok" = false ]; then
    echo -e "${YELLOW}⚠️  Some ports are in use. Continue anyway? (y/N)${NC}"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Deploy application
if [ "$DEPLOY_MODE" = "remote" ]; then
    echo -e "${BLUE}🏗️  Remote deployment to ${REMOTE_HOST}...${NC}"

    # Prepare remote server
    echo -e "${BLUE}📦 Preparing remote server...${NC}"
    run_remote_command "mkdir -p $REMOTE_PATH"
    run_remote_command "cd $REMOTE_PATH && which docker || (curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh)"
    run_remote_command "cd $REMOTE_PATH && which docker-compose || curl -L \"https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-\$(uname -s)-\$(uname -m)\" -o /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose"

    # Copy project files to remote server
    echo -e "${BLUE}📁 Copying project files to remote server...${NC}"
    copy_to_remote "." "$REMOTE_PATH/"

    # Deploy on remote server
    echo -e "${BLUE}🚀 Starting deployment on remote server...${NC}"
    run_remote_command "cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml down || true"
    run_remote_command "cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml build --no-cache"
    run_remote_command "cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml up -d"

else
    echo -e "${BLUE}🏗️  Local deployment...${NC}"

    # Stop existing containers if running
    echo -e "${YELLOW}🛑 Stopping existing containers...${NC}"
    docker-compose -f docker-compose.prod.yml down || true

    # Pull latest images and build
    echo -e "${BLUE}🔨 Building containers...${NC}"
    docker-compose -f docker-compose.prod.yml build --no-cache

    # Start services
    echo -e "${BLUE}🚀 Starting services...${NC}"
    docker-compose -f docker-compose.prod.yml up -d
fi

# Wait for services to be healthy
echo -e "${BLUE}⏳ Waiting for services to be healthy...${NC}"
sleep 10

# Check service health
echo -e "${BLUE}🏥 Checking service health...${NC}"

check_service_health() {
    local service=$1
    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if [ "$DEPLOY_MODE" = "remote" ]; then
            if run_remote_command "cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml ps $service | grep -q 'healthy\\|Up'"; then
                echo -e "${GREEN}✅ $service is healthy${NC}"
                return 0
            fi
        else
            if docker-compose -f docker-compose.prod.yml ps $service | grep -q "healthy\|Up"; then
                echo -e "${GREEN}✅ $service is healthy${NC}"
                return 0
            fi
        fi
        echo -e "${YELLOW}⏳ Waiting for $service to be healthy (attempt $attempt/$max_attempts)${NC}"
        sleep 5
        ((attempt++))
    done

    echo -e "${RED}❌ $service failed to become healthy${NC}"
    return 1
}

services=("postgres" "redis" "backend" "storefront" "admin-panel" "caddy")
all_healthy=true

for service in "${services[@]}"; do
    if ! check_service_health $service; then
        all_healthy=false
    fi
done

if [ "$all_healthy" = false ]; then
    echo -e "${RED}❌ Some services are not healthy. Check logs:${NC}"
    if [ "$DEPLOY_MODE" = "remote" ]; then
        echo "ssh $REMOTE_USER@$REMOTE_HOST 'cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml logs'"
    else
        echo "docker-compose -f docker-compose.prod.yml logs"
    fi
    exit 1
fi

# Test endpoints
echo -e "${BLUE}🧪 Testing endpoints...${NC}"

test_endpoint() {
    local url=$1
    local description=$2

    if curl -sSf "$url" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ $description is accessible${NC}"
    else
        echo -e "${RED}❌ $description is not accessible${NC}"
        return 1
    fi
}

# Wait a bit more for services to be fully ready
sleep 10

# Test endpoints
if [ "$DEPLOY_MODE" = "remote" ]; then
    test_endpoint "http://$REMOTE_HOST/health" "Remote health check"
    test_endpoint "https://$DOMAIN/health" "HTTPS health check" || echo -e "${YELLOW}⚠️  HTTPS may take a few minutes for SSL certificate generation${NC}"
else
    test_endpoint "http://localhost/health" "Local health check"
    test_endpoint "https://$DOMAIN/health" "HTTPS health check" || echo -e "${YELLOW}⚠️  HTTPS may take a few minutes for SSL certificate generation${NC}"
fi

# Show deployment summary
echo ""
echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
echo "=============================================="
echo -e "${BLUE}📋 Service URLs:${NC}"
echo "🌐 Main site: https://$DOMAIN"
echo "🏪 Admin panel: https://admin.$DOMAIN"
echo "📊 API: https://$DOMAIN/api/v1"

if [ "$DEPLOY_MODE" = "remote" ]; then
    echo ""
    echo -e "${BLUE}🌐 Remote Server Details:${NC}"
    echo "🖥️  Server: $REMOTE_HOST"
    echo "👤 User: $REMOTE_USER"
    echo "📁 Path: $REMOTE_PATH"
fi

echo ""
echo -e "${BLUE}🔧 Management commands:${NC}"
if [ "$DEPLOY_MODE" = "remote" ]; then
    echo "📜 View logs: ssh $REMOTE_USER@$REMOTE_HOST 'cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml logs -f'"
    echo "📊 Check status: ssh $REMOTE_USER@$REMOTE_HOST 'cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml ps'"
    echo "🛑 Stop services: ssh $REMOTE_USER@$REMOTE_HOST 'cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml down'"
    echo "🔄 Restart: ssh $REMOTE_USER@$REMOTE_HOST 'cd $REMOTE_PATH && docker-compose -f docker-compose.prod.yml restart'"
    echo "🔌 SSH Access: ssh $REMOTE_USER@$REMOTE_HOST"
else
    echo "📜 View logs: docker-compose -f docker-compose.prod.yml logs -f"
    echo "📊 Check status: docker-compose -f docker-compose.prod.yml ps"
    echo "🛑 Stop services: docker-compose -f docker-compose.prod.yml down"
    echo "🔄 Restart: docker-compose -f docker-compose.prod.yml restart"
fi

echo ""
echo -e "${YELLOW}💡 Next steps:${NC}"
echo "1. Configure DNS records for $DOMAIN to point to $([ "$DEPLOY_MODE" = "remote" ] && echo "$REMOTE_HOST" || echo "your server")"
echo "2. Create your first tenant via admin panel"
echo "3. Test subdomain functionality"
echo "4. Configure monitoring and backups"
echo ""
echo -e "${GREEN}✨ Your multi-tenant e-commerce SaaS platform is now live!${NC}"