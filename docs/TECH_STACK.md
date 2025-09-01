# Technical Stack & Technologies

## Backend Stack (Golang)

### Core Framework & Libraries
- **Go 1.21+**: Latest stable Go version
- **Fiber**: Express.js-like HTTP web framework for REST APIs
- **GORM**: ORM for database operations
- **Wire**: Dependency injection framework
- **Viper**: Configuration management
- **Zap/Logrus**: Structured logging
- **Validator**: Request validation
- **JWT-Go**: JSON Web Token implementation

### Database & Storage
- **PostgreSQL 15+**: Primary database with JSONB support, self-hosted on VPS
- **Redis 7+**: Caching and session storage, containerized
- **Local File Storage**: Images/files stored on VPS with Cloudflare Images optimization
- **Cloudflare Images**: Image optimization and CDN delivery

### Message Queue & Events
- **WebSocket**: Real-time bidirectional communication
- **NATS/RabbitMQ**: Event streaming and pub/sub
- **Apache Kafka**: High-throughput event streaming (enterprise)
- **Redis Streams**: Lightweight message queue option

### Authentication & Security
- **OAuth2**: Authorization framework
- **BCrypt**: Password hashing
- **Rate Limiting**: Token bucket algorithm
- **CORS**: Cross-origin resource sharing
- **Helmet**: Security headers middleware

### Testing
- **Testify**: Testing toolkit with assertions
- **GoMock**: Mock generation for testing
- **Ginkgo/Gomega**: BDD testing framework
- **Go-SqlMock**: SQL driver mock for testing
- **Playwright**: End-to-end testing for web applications
- **Cypress**: Alternative E2E testing framework
- **Jest**: Unit testing for JavaScript

## Frontend Stack

### Customer Storefront (Next.js)
- **Next.js 14+**: React framework with App Router
- **JavaScript**: Dynamic language with modern ES6+ features
- **Tailwind CSS**: Utility-first CSS framework
- **Shadcn/UI**: Modern UI component library
- **Jotai**: Atomic state management
- **React Query**: Server state management
- **React Hook Form**: Form handling and validation
- **Framer Motion**: Animation library

### Merchant Dashboard (React.js)
- **React 18+**: UI library with concurrent features
- **JavaScript**: Dynamic language with modern ES6+ features
- **Material-UI (MUI)**: Component library for admin interfaces
- **Jotai**: Atomic state management
- **React Query**: Data fetching and caching
- **React Router**: Client-side routing
- **Recharts**: Data visualization and charts
- **React DnD**: Drag and drop functionality

### Shared Frontend Tools
- **Axios**: HTTP client library
- **Date-fns**: Date manipulation library
- **React-i18next**: Internationalization
- **React-Helmet**: Document head management
- **ESLint/Prettier**: Code linting and formatting

## DevOps & Infrastructure

### Containerization & Deployment
- **Docker**: Application containerization
- **Docker Compose**: Multi-container production deployment
- **Nginx**: Reverse proxy and static file serving
- **Let's Encrypt**: Free SSL certificate automation

### VPS Hosting & Infrastructure
- **Hostinger VPS**: Primary hosting provider (cost-effective)
- **DigitalOcean**: Alternative VPS provider
- **Hetzner**: European alternative for better pricing
- **Cloudflare**: Free CDN, DNS, and DDoS protection
- **Cloudflare Images**: Image optimization and delivery
- **Nginx**: Load balancing and reverse proxy

### Monitoring & Observability
- **Uptime Kuma**: Self-hosted uptime monitoring
- **Grafana**: Lightweight metrics visualization (self-hosted)
- **Prometheus**: Basic metrics collection (containerized)
- **Sentry**: Error tracking and monitoring (free tier)
- **Simple logging**: Docker logs + log rotation
- **Custom health checks**: Built-in application monitoring

### CI/CD Pipeline
- **GitHub Actions**: Continuous integration and deployment
- **Docker Hub**: Container registry (free tier)
- **Simple deployment scripts**: Direct VPS deployment
- **SSH-based deployment**: Direct server updates
- **Blue-Green Deployments**: Using Nginx upstream switching
- **Rolling updates**: Docker Compose rolling updates

## Third-Party Integrations

### Payment Processing
- **Stripe**: Primary payment processor
- **PayPal**: Alternative payment option
- **Square**: In-person payments
- **Razorpay**: International payments (Asia)
- **Adyen**: Enterprise payment processing

### Email Services
- **SendGrid**: Transactional emails
- **Mailgun**: Email delivery service
- **Amazon SES**: AWS email service
- **Postmark**: Transactional email delivery

### Shipping & Logistics
- **ShipStation**: Multi-carrier shipping
- **EasyPost**: Shipping API
- **FedEx/UPS/DHL APIs**: Direct carrier integration
- **Canada Post**: Canadian shipping

### Analytics & Marketing
- **Google Analytics**: Web analytics
- **Facebook Pixel**: Social media tracking
- **Mixpanel**: Product analytics
- **Segment**: Customer data platform
- **Klaviyo**: Email marketing automation
- **Mailchimp**: Email marketing platform

## Development Tools

### Version Control & Collaboration
- **Git**: Version control system
- **GitHub**: Code hosting and collaboration
- **GitHub Issues**: Issue tracking
- **GitHub Projects**: Project management

### Development Environment
- **VS Code**: Primary code editor
- **GoLand**: Go-specific IDE option
- **Docker Desktop**: Local containerization
- **Postman/Insomnia**: API testing
- **TablePlus/DBeaver**: Database GUI

### Code Quality & Security
- **SonarQube**: Code quality analysis
- **CodeClimate**: Automated code review
- **Snyk**: Security vulnerability scanning
- **Dependabot**: Dependency updates
- **Pre-commit hooks**: Git hooks for code quality

## Performance Optimization

### Caching Strategy
- **Redis**: Application-level caching
- **CDN**: Static asset caching
- **Database Query Caching**: ORM-level caching
- **API Response Caching**: HTTP response caching
- **Browser Caching**: Client-side caching

### Database Optimization
- **Connection Pooling**: Efficient database connections
- **Read Replicas**: Database read scaling
- **Indexing Strategy**: Optimized database queries
- **Query Optimization**: SQL query performance tuning
- **Database Partitioning**: Large table management

### Frontend Optimization
- **Code Splitting**: Dynamic imports and lazy loading
- **Image Optimization**: WebP format, responsive images
- **Bundle Optimization**: Webpack/Vite optimizations
- **Service Workers**: Offline functionality and caching
- **Progressive Web App**: Native app-like experience

## Security Technologies

### Application Security
- **HTTPS/TLS 1.3**: Encrypted communication
- **OWASP Security Headers**: Security best practices
- **Input Validation**: SQL injection prevention
- **XSS Protection**: Cross-site scripting prevention
- **CSRF Protection**: Cross-site request forgery prevention

### Infrastructure Security
- **VPC**: Virtual private cloud networking
- **Security Groups**: Network access control
- **IAM Roles**: Identity and access management
- **Secrets Manager**: Secure credential storage
- **WAF**: Web application firewall

### Compliance & Privacy
- **GDPR Compliance**: European data protection
- **PCI DSS**: Payment card industry compliance
- **SOC 2**: Security and availability standards
- **Data Encryption**: At-rest and in-transit encryption
- **Audit Logging**: Comprehensive activity logging