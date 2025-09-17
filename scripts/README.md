# Scripts Directory

This directory contains utility scripts for managing the e-commerce SaaS platform.

## Available Scripts

### 🚀 `deploy.sh`
**Production deployment script with full validation and remote server support**

**Local deployment:**
```bash
./scripts/deploy.sh
```

**Remote server deployment:**
```bash
./scripts/deploy.sh --remote --host 1.2.3.4 --user root --password mypassword
```

**Available options:**
- `--remote`: Enable remote deployment mode
- `--host HOST`: Remote server IP/hostname
- `--user USER`: SSH username
- `--password PASS`: SSH password
- `--path PATH`: Remote deployment path (default: /opt/esass)
- `--help`: Show usage information

**Features:**
- Environment variable validation
- Remote server connection testing
- Automatic Docker installation on remote servers
- Port availability checking
- Docker container building and starting
- Health check verification
- Service endpoint testing
- Complete deployment summary

### 📦 `backup.sh`
**Automated backup script for production data**

```bash
./scripts/backup.sh
```

Features:
- PostgreSQL database backup
- Redis data backup
- Uploads directory backup
- Automatic cleanup (7-day retention)
- Backup size reporting

## Setup Instructions

1. **Make scripts executable:**
   ```bash
   chmod +x scripts/*.sh
   ```

2. **Configure environment:**
   ```bash
   cp .env.production.example .env.production
   nano .env.production
   ```

3. **Set up automated backups:**
   ```bash
   # Add to crontab for daily backups at 2 AM
   crontab -e
   # Add line: 0 2 * * * cd /opt/esass && ./scripts/backup.sh
   ```

## Usage Notes

- All scripts should be run from the project root directory
- Scripts will load environment variables from `.env.production`
- Check logs with `docker-compose -f docker-compose.prod.yml logs`
- Backup location: `/opt/backups/`

## Dependencies

- Docker & Docker Compose
- curl (for health checks)
- netcat (for port checking)
- Standard Unix utilities (tar, gzip, find)
- **For remote deployment:** sshpass (or SSH key authentication)

**Install sshpass:**
```bash
# Ubuntu/Debian
sudo apt-get install sshpass

# macOS
brew install sshpass
```

**Alternative: SSH Key Authentication**
```bash
# Generate SSH key pair
ssh-keygen -t rsa -b 4096

# Copy public key to remote server
ssh-copy-id user@server_ip

# Use deployment without password
./scripts/deploy.sh --remote --host server_ip --user username
```