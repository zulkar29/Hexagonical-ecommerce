# Security Implementation

Security architecture and implementation details for the e-commerce SaaS platform.

## Security Architecture Overview

### Multi-layered Security Approach
1. **Perimeter Security** - WAF, DDoS protection
2. **Application Security** - Authentication, authorization, input validation
3. **Data Security** - Encryption, access control, audit logging
4. **Infrastructure Security** - Network isolation, monitoring

## Authentication & Authorization

### JWT Authentication
```go
// JWT Token Structure
{
  "sub": "user-uuid",
  "tenant_id": "tenant-uuid", 
  "role": "admin|staff|customer",
  "permissions": ["read:products", "write:orders"],
  "exp": 1640995200,
  "iat": 1640908800
}
```

### Multi-Factor Authentication (MFA)
- **TOTP Support**: Google Authenticator, Authy
- **SMS Backup**: Fallback for TOTP
- **Recovery Codes**: One-time backup codes
- **Enforcement**: Required for admin accounts

### Role-Based Access Control (RBAC)
```yaml
Roles:
  super_admin:
    - manage:platform
    - manage:all_tenants
  
  tenant_admin:
    - manage:tenant_settings
    - manage:users
    - manage:products
    - manage:orders
  
  staff:
    - read:products
    - write:products
    - read:orders
    - write:orders
  
  customer:
    - read:own_profile
    - write:own_profile
    - read:own_orders
    - write:cart
```

## Input Validation & Sanitization

### Request Validation
```go
type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=1,max=255"`
    Price       float64 `json:"price" validate:"required,min=0"`
    Description string  `json:"description" validate:"max=10000"`
    SKU         string  `json:"sku" validate:"alphanum,max=100"`
}
```

### SQL Injection Prevention
- **Parameterized Queries**: All database queries use parameters
- **ORM Protection**: GORM prevents SQL injection by default
- **Input Sanitization**: HTML/SQL content sanitized
- **Query Logging**: Monitor for suspicious patterns

### XSS Protection
- **Content Security Policy**: Strict CSP headers
- **Output Encoding**: All user content escaped
- **Sanitization**: HTML purification for rich content
- **Safe Rendering**: React's built-in XSS protection

## Data Protection

### Encryption at Rest
```yaml
Database:
  - PostgreSQL TDE (Transparent Data Encryption)
  - AES-256 encryption
  - Key rotation every 90 days

File Storage:
  - S3 Server-Side Encryption (SSE-S3)
  - Customer-managed keys (SSE-KMS)
  - Encryption in transit (HTTPS/TLS 1.3)
```

### Encryption in Transit
- **TLS 1.3**: All HTTP communications
- **Certificate Pinning**: Mobile app security
- **HSTS**: HTTP Strict Transport Security
- **Database Connections**: SSL/TLS for all DB connections

### Sensitive Data Handling
```go
// Password Hashing
func HashPassword(password string) (string, error) {
    return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// PII Encryption
func EncryptPII(data string) (string, error) {
    // AES-256-GCM encryption
    // Key from AWS KMS or HashiCorp Vault
}
```

## API Security

### Rate Limiting
```yaml
Limits:
  Anonymous: 100 requests/hour
  Authenticated: 1000 requests/hour
  Admin: 5000 requests/hour
  
Per-Endpoint:
  POST /api/auth/login: 5 requests/minute
  POST /api/orders: 100 requests/hour
  GET /api/products: 1000 requests/hour
```

### API Security Headers
```http
# Security Headers
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
```

### Request Signing (Webhooks)
```go
func VerifyWebhookSignature(payload []byte, signature string) bool {
    mac := hmac.New(sha256.New, webhookSecret)
    mac.Write(payload)
    expectedSignature := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

## Infrastructure Security

### Network Security
```yaml
VPC Configuration:
  - Private subnets for application servers
  - Public subnets only for load balancers
  - NAT Gateway for outbound internet access
  - Security groups with minimal required ports

Security Groups:
  Load Balancer:
    - Inbound: 80, 443 (from 0.0.0.0/0)
    - Outbound: 8080 (to app servers)
  
  App Servers:
    - Inbound: 8080 (from load balancer only)
    - Outbound: 5432 (to database)
  
  Database:
    - Inbound: 5432 (from app servers only)
    - No outbound internet access
```

### Container Security
```dockerfile
# Security-hardened container
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache ca-certificates git

# Non-root user
RUN adduser -D -g '' appuser

# Runtime container
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Run as non-root
USER appuser
```

## Compliance & Privacy

### GDPR Compliance
```go
// Data Subject Rights Implementation
type GDPRService struct{}

func (g *GDPRService) ExportUserData(userID string) (*UserDataExport, error) {
    // Export all user data in JSON format
}

func (g *GDPRService) DeleteUserData(userID string) error {
    // Anonymize or delete user data
    // Maintain referential integrity
}

func (g *GDPRService) ProcessDataRequest(request *DataRequest) error {
    // Handle data access, portability, deletion requests
}
```

### PCI DSS Compliance
- **No Card Storage**: Use Stripe/payment processor vaults
- **Secure Transmission**: TLS 1.3 for all card data
- **Access Logging**: All payment-related access logged
- **Regular Testing**: Quarterly security assessments

### SOC 2 Compliance
- **Security Controls**: Documented security procedures
- **Access Reviews**: Quarterly access reviews
- **Change Management**: Controlled deployment processes
- **Incident Response**: 24-hour incident response plan

## Monitoring & Incident Response

### Security Monitoring
```yaml
Log Sources:
  - Application logs (authentication, authorization)
  - Database audit logs
  - Infrastructure logs (VPC flow logs)
  - WAF logs (blocked requests)

Alerts:
  - Multiple failed login attempts
  - Privilege escalation attempts
  - Unusual API access patterns
  - Database schema changes
  - Large data exports
```

### Incident Response Plan
1. **Detection** - Automated monitoring alerts
2. **Analysis** - Security team investigation
3. **Containment** - Isolate affected systems
4. **Eradication** - Remove threat, patch vulnerabilities
5. **Recovery** - Restore services, monitor for reoccurrence
6. **Lessons Learned** - Post-incident review and improvements

### Security Metrics
- Mean Time to Detection (MTTD)
- Mean Time to Response (MTTR)
- Failed authentication attempts
- API abuse incidents
- Data breach incidents (target: 0)

## Vulnerability Management

### Security Testing
```yaml
Testing Types:
  - Static Application Security Testing (SAST)
  - Dynamic Application Security Testing (DAST)
  - Interactive Application Security Testing (IAST)
  - Software Composition Analysis (SCA)
  - Penetration testing (quarterly)
```

### Dependency Management
```bash
# Go security scanning
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Node.js security auditing
npm audit
npm audit fix
```

### Security Updates
- **Critical**: Applied within 24 hours
- **High**: Applied within 1 week
- **Medium**: Applied within 1 month
- **Low**: Applied during regular maintenance

## Backup & Disaster Recovery

### Secure Backups
```yaml
Backup Security:
  - Encrypted at rest (AES-256)
  - Encrypted in transit (TLS 1.3)
  - Access-controlled (IAM policies)
  - Geographically distributed
  - Regular restore testing

Retention:
  - Daily backups: 30 days
  - Weekly backups: 12 weeks
  - Monthly backups: 12 months
  - Yearly backups: 7 years (compliance)
```

### Business Continuity
- **RTO**: 4 hours (Recovery Time Objective)
- **RPO**: 1 hour (Recovery Point Objective)
- **Multi-region**: Active-passive disaster recovery
- **Failover**: Automated failover procedures

## Security Training & Awareness

### Developer Security Training
- Secure coding practices
- OWASP Top 10 awareness
- Threat modeling
- Security testing
- Incident response procedures

### Security Reviews
- **Code Reviews**: Security-focused peer reviews
- **Architecture Reviews**: Security architecture validation
- **Third-party Reviews**: Vendor security assessments
- **Regular Audits**: Annual security audits

## Third-party Integration Security

### API Integration Security
```go
// Secure API client configuration
client := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS13,
        },
    },
}
```

### Vendor Assessment
- Security questionnaires
- SOC 2 reports review
- Penetration testing results
- Compliance certifications
- Data processing agreements