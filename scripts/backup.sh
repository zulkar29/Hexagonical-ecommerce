#!/bin/bash

# Backup Script for Hexagonal E-commerce SaaS
# Usage: ./scripts/backup.sh [destination]
# Destination: local (default) | s3

set -e

BACKUP_DEST=${1:-local}
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)
COMPOSE_FILE="docker-compose.prod.yml"
ENV_FILE=".env.production"

echo "💾 Starting backup process..."

# Load environment variables
if [ -f "$ENV_FILE" ]; then
    set -a
    source $ENV_FILE
    set +a
fi

# Create backup directory
mkdir -p $BACKUP_DIR

echo "📊 Creating database backup..."
# PostgreSQL backup
docker-compose -f $COMPOSE_FILE exec -T postgres pg_dump \
    -U $POSTGRES_USER \
    -h localhost \
    --clean \
    --no-owner \
    --no-acl \
    $POSTGRES_DB | gzip > $BACKUP_DIR/database_$DATE.sql.gz

echo "🔄 Creating Redis backup..."
# Redis backup
docker-compose -f $COMPOSE_FILE exec -T redis redis-cli --rdb - | \
    gzip > $BACKUP_DIR/redis_$DATE.rdb.gz

echo "📁 Creating application files backup..."
# Application uploads backup
if [ -d "uploads" ]; then
    tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz uploads/
fi

# Environment configuration backup
if [ -f "$ENV_FILE" ]; then
    cp $ENV_FILE $BACKUP_DIR/env_$DATE.backup
fi

# Docker compose configuration backup
if [ -f "$COMPOSE_FILE" ]; then
    cp $COMPOSE_FILE $BACKUP_DIR/compose_$DATE.yml
fi

# Create backup manifest
cat > $BACKUP_DIR/manifest_$DATE.txt << EOF
Backup created: $DATE
Database: database_$DATE.sql.gz
Redis: redis_$DATE.rdb.gz
Uploads: uploads_$DATE.tar.gz
Environment: env_$DATE.backup
Compose: compose_$DATE.yml
EOF

echo "📋 Backup manifest created"

# Calculate backup size
BACKUP_SIZE=$(du -sh $BACKUP_DIR/*_$DATE.* | awk '{total+=$1} END {print total}')
echo "📦 Backup size: $BACKUP_SIZE"

# Upload to S3 if specified
if [ "$BACKUP_DEST" = "s3" ] && [ -n "$BACKUP_S3_BUCKET" ]; then
    echo "☁️  Uploading to S3..."
    
    if command -v aws &> /dev/null; then
        aws s3 cp $BACKUP_DIR/database_$DATE.sql.gz s3://$BACKUP_S3_BUCKET/backups/
        aws s3 cp $BACKUP_DIR/redis_$DATE.rdb.gz s3://$BACKUP_S3_BUCKET/backups/
        aws s3 cp $BACKUP_DIR/uploads_$DATE.tar.gz s3://$BACKUP_S3_BUCKET/backups/
        aws s3 cp $BACKUP_DIR/manifest_$DATE.txt s3://$BACKUP_S3_BUCKET/backups/
        echo "✅ Backup uploaded to S3"
    else
        echo "⚠️  AWS CLI not found, skipping S3 upload"
    fi
fi

# Cleanup old backups (keep last 7 days locally)
if [ -n "$BACKUP_RETENTION_DAYS" ]; then
    RETENTION_DAYS=$BACKUP_RETENTION_DAYS
else
    RETENTION_DAYS=7
fi

echo "🧹 Cleaning up old backups (keeping last $RETENTION_DAYS days)..."
find $BACKUP_DIR -name "*_????????_??????.*" -mtime +$RETENTION_DAYS -delete

# Verify backup integrity
echo "🔍 Verifying backup integrity..."

# Check database backup
if ! gunzip -t $BACKUP_DIR/database_$DATE.sql.gz; then
    echo "❌ Database backup verification failed"
    exit 1
fi

# Check Redis backup
if ! gunzip -t $BACKUP_DIR/redis_$DATE.rdb.gz; then
    echo "❌ Redis backup verification failed"
    exit 1
fi

# Check uploads backup
if [ -f "$BACKUP_DIR/uploads_$DATE.tar.gz" ]; then
    if ! tar -tzf $BACKUP_DIR/uploads_$DATE.tar.gz > /dev/null; then
        echo "❌ Uploads backup verification failed"
        exit 1
    fi
fi

echo "✅ Backup integrity verified"

# Send notification (if configured)
if [ -n "$BACKUP_NOTIFICATION_WEBHOOK" ]; then
    curl -X POST "$BACKUP_NOTIFICATION_WEBHOOK" \
        -H "Content-Type: application/json" \
        -d "{\"text\":\"✅ Backup completed successfully at $DATE\"}" \
        > /dev/null 2>&1 || echo "⚠️  Notification webhook failed"
fi

# Display backup summary
echo ""
echo "📊 Backup Summary:"
echo "   Date: $DATE"
echo "   Location: $BACKUP_DIR"
echo "   Files:"
ls -lh $BACKUP_DIR/*_$DATE.*

echo ""
echo "✨ Backup completed successfully!"

# Display restore instructions
echo ""
echo "🔄 To restore from this backup:"
echo "   Database: gunzip < $BACKUP_DIR/database_$DATE.sql.gz | docker-compose -f $COMPOSE_FILE exec -T postgres psql -U $POSTGRES_USER $POSTGRES_DB"
echo "   Uploads:  tar -xzf $BACKUP_DIR/uploads_$DATE.tar.gz"
echo ""