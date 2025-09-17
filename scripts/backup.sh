#!/bin/bash

# Production Backup Script for E-commerce SaaS Platform
# This script creates backups of database, Redis, and uploads

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📦 Starting backup process...${NC}"

# Configuration
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)
COMPOSE_FILE="docker-compose.prod.yml"
RETENTION_DAYS=7

# Load environment variables if .env.production exists
if [ -f ".env.production" ]; then
    export $(grep -v '^#' .env.production | xargs)
fi

# Create backup directory
mkdir -p $BACKUP_DIR

echo -e "${YELLOW}🗄️  Backing up PostgreSQL database...${NC}"
# Database backup
docker-compose -f $COMPOSE_FILE exec -T postgres pg_dump \
    -U ${POSTGRES_USER:-postgres} \
    ${POSTGRES_DB:-ecommerce_saas_prod} | \
    gzip > $BACKUP_DIR/db_$DATE.sql.gz

echo -e "${YELLOW}🔴 Backing up Redis data...${NC}"
# Redis backup
docker-compose -f $COMPOSE_FILE exec -T redis redis-cli \
    -a ${REDIS_PASSWORD} --rdb - | \
    gzip > $BACKUP_DIR/redis_$DATE.rdb.gz

echo -e "${YELLOW}📁 Backing up uploads...${NC}"
# Application files backup
if [ -d "uploads" ]; then
    tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz uploads/
else
    echo -e "${YELLOW}⚠️  No uploads directory found, skipping...${NC}"
fi

echo -e "${YELLOW}🧹 Cleaning old backups (keeping last $RETENTION_DAYS days)...${NC}"
# Keep only last 7 days
find $BACKUP_DIR -name "*.gz" -mtime +$RETENTION_DAYS -delete

# Backup summary
echo -e "${GREEN}✅ Backup completed successfully!${NC}"
echo -e "${BLUE}📊 Backup summary:${NC}"
echo "  📅 Date: $DATE"
echo "  📂 Location: $BACKUP_DIR"
echo "  📋 Files created:"
echo "    - db_$DATE.sql.gz"
echo "    - redis_$DATE.rdb.gz"
if [ -d "uploads" ]; then
    echo "    - uploads_$DATE.tar.gz"
fi

# Show backup directory size
BACKUP_SIZE=$(du -sh $BACKUP_DIR | cut -f1)
echo "  💾 Total backup size: $BACKUP_SIZE"

echo -e "${GREEN}🎉 Backup process completed!${NC}"