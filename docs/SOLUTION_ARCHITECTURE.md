# Course Enrollment Service - Solution Architecture

## 🎯 **Intern Test Optimized Architecture**

> **Note**: This architecture is designed for **intern backend test** - impressive yet practical!

## 📋 Tổng quan hệ thống

Đây là tài liệu mô tả chi tiết solution architecture cho hệ thống Course Enrollment Service, được tối ưu cho intern test với sự cân bằng giữa **professional setup** và **cost efficiency**.

## 🏗️ Architecture Overview

### High-Level Architecture
```
┌─────────────────────────────────────────────────────────┐
│                    USER BROWSER                         │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                 CLOUDFRONT CDN                          │ ✅ Global Performance
│              (app.yourdomain.com)                       │ ✅ HTTPS Termination
└─────────────────────┬───────────────────────────────────┘ ✅ Professional Setup
                      │
        ┌─────────────┼─────────────┐
        │                           │
┌───────▼────────┐         ┌────────▼────────┐
│   S3 BUCKET    │         │ SINGLE EC2      │ ✅ Simple & Reliable
│ (React Static) │ ✅ Cost │   (Go API)      │ ✅ Cost Optimized
│                │ Effective│ api.domain.com  │ ✅ Easy to Debug
└────────────────┘         └────────┬────────┘
                                    │
                           ┌────────▼────────┐
                           │ RDS POSTGRESQL  │ ✅ Managed Database
                           │   (Database)    │ ✅ Automatic Backups
                           └─────────────────┘ ✅ SSL Encryption
```

### Data Flow
```
User Browser ──► CloudFront ──► S3 (React App)
     │                              │
     │         API Calls            │
     └──────────────────────────────┘
                  │
Internet ──► EC2 (Go API) ──► RDS PostgreSQL
                  │
                  ▼
            S3 (File Uploads)
```

## 🎯 **Architecture Decisions for Intern Test**

### ✅ **KEEP (Impressive & Necessary)**
```yaml
CloudFront CDN:
  ✅ Shows full-stack knowledge
  ✅ Global performance optimization
  ✅ HTTPS termination
  ✅ Professional frontend setup

CI/CD Pipeline:
  ✅ Shows DevOps skills
  ✅ Automated testing & deployment
  ✅ Professional workflow
  ✅ GitHub Actions integration

CloudWatch Monitoring:
  ✅ Basic application monitoring
  ✅ Error tracking & logging
  ✅ Performance metrics
  ✅ Production awareness

Single EC2 Instance:
  ✅ Simple & reliable
  ✅ Cost optimized (~$15/month)
  ✅ Easy to debug & maintain
  ✅ Sufficient for demo/test

RDS PostgreSQL:
  ✅ Managed database service
  ✅ Automatic backups
  ✅ SSL encryption
  ✅ Production-grade setup
```

### ❌ **REMOVE (Overkill for Intern Test)**
```yaml
Application Load Balancer:
  ❌ Single EC2 instance sufficient
  ❌ No need for load distribution
  ❌ Adds unnecessary complexity
  ❌ Extra cost (~$18/month)

Auto Scaling Groups:
  ❌ No scaling requirements
  ❌ Fixed traffic patterns
  ❌ Adds operational complexity
  ❌ Not needed for demo

Multiple AZ Deployment:
  ❌ High availability not critical
  ❌ Single AZ sufficient for test
  ❌ Reduces cost significantly
  ❌ Simpler to manage

Advanced Health Checks:
  ❌ Basic health endpoint sufficient
  ❌ Complex monitoring overkill
  ❌ Simple logging adequate
  ❌ Manual intervention acceptable

Rollback Mechanisms:
  ❌ Manual rollback acceptable
  ❌ Blue-green deployment overkill
  ❌ Git revert sufficient
  ❌ Reduces complexity
```

## � **Why This Architecture is Perfect for Intern Test**

### **Impressive Features**
```yaml
✅ Full-Stack Knowledge:
  - React/Next.js frontend with modern practices
  - Go backend with clean architecture
  - PostgreSQL database with proper relationships
  - AWS cloud infrastructure

✅ DevOps Skills:
  - CI/CD pipeline with GitHub Actions
  - Docker containerization
  - Infrastructure as Code concepts
  - Environment management

✅ Production Awareness:
  - SSL/HTTPS everywhere
  - Database migrations
  - Error handling & logging
  - API documentation (Swagger)
  - Security best practices

✅ Cost Consciousness:
  - Optimized for budget constraints
  - No over-engineering
  - Scalable when needed
  - Professional yet practical
```

### **Technical Depth**
```yaml
Backend (Go):
  ✅ Clean Architecture (Repository pattern)
  ✅ GORM for database operations
  ✅ Gin framework for REST API
  ✅ Swagger documentation
  ✅ Environment configuration
  ✅ Database migrations
  ✅ Error handling
  ✅ UUID primary keys

Frontend (React):
  ✅ Modern React with hooks
  ✅ API integration
  ✅ Responsive design
  ✅ Global CDN delivery
  ✅ HTTPS termination

Infrastructure:
  ✅ AWS cloud services
  ✅ Managed database (RDS)
  ✅ Content delivery (CloudFront)
  ✅ Automated deployment
  ✅ Basic monitoring
```

### **Interview Talking Points**
```yaml
Architecture Decisions:
  "I chose single EC2 over ALB because..."
  "CloudFront provides global performance..."
  "RDS gives us managed backups and SSL..."
  "CI/CD pipeline ensures code quality..."

Scalability:
  "Current architecture handles 100+ users..."
  "Can easily add ALB when traffic grows..."
  "Database can scale vertically first..."
  "Horizontal scaling planned for future..."

Cost Optimization:
  "Balanced professional features with budget..."
  "Avoided over-engineering for demo..."
  "Can scale up when requirements change..."
  "Monitoring helps optimize resources..."
```

## �🌐 Frontend Layer (React/Next.js)

### S3 Static Website Hosting
```yaml
Bucket: course-enrollment-frontend-prod
Configuration:
  - Static website hosting: Enabled
  - Index document: index.html
  - Error document: error.html
  - Public read access: Enabled
  - CORS: Configured for API calls
```

### CloudFront Distribution
```yaml
Origin: S3 bucket
Behaviors:
  - Default: Cache everything (TTL: 24h)
  - /api/*: No cache, forward to EC2
  - /static/*: Long cache (TTL: 1 year)
SSL Certificate: ACM certificate
Custom Domain: app.yourdomain.com
```

### React Configuration
```javascript
// src/config/api.js
const API_BASE_URL = process.env.REACT_APP_API_URL || 'https://api.yourdomain.com';

const api = {
  getCourses: () => fetch(`${API_BASE_URL}/api/v1/courses`),
  createCourse: (data) => fetch(`${API_BASE_URL}/api/v1/courses`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }),
  enrollStudent: (data) => fetch(`${API_BASE_URL}/api/v1/enrollments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
};
```

## 🖥️ Backend Layer (Go API)

### EC2 Configuration
```yaml
Instance: t3.small (2 vCPU, 2GB RAM)
OS: Amazon Linux 2
Security Group:
  - Port 80: 0.0.0.0/0 (HTTP)
  - Port 443: 0.0.0.0/0 (HTTPS)  
  - Port 22: Your IP (SSH)
Elastic IP: Attached
Domain: api.yourdomain.com
```

### CORS Configuration
```go
// internal/router/router.go
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "https://app.yourdomain.com")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

## 🗄️ Database & Storage Layer

### RDS PostgreSQL
```yaml
Instance: db.t3.micro
Engine: PostgreSQL 14
Multi-AZ: No (cost optimization)
Storage: 20GB GP2
Backup: 7 days retention
Security Group: Port 5432 from EC2 only
```

### S3 Buckets
```yaml
Frontend Bucket: course-enrollment-frontend-prod
  - Static website hosting
  - CloudFront origin
  - Public read access

Assets Bucket: course-enrollment-assets-prod  
  - File uploads (course images, documents)
  - Private access via signed URLs
  - CORS enabled for frontend uploads
```

## 💰 Cost Analysis (Intern Test Optimized)

### **Simplified Architecture Cost**
| Component | Service | Configuration | Monthly Cost | Notes |
|-----------|---------|---------------|--------------|-------|
| **Frontend** | S3 + CloudFront | Static hosting + CDN | $3 | ✅ Global performance |
| **Backend** | EC2 t3.small | Single instance | $15 | ✅ Sufficient for demo |
| **Database** | RDS db.t3.micro | PostgreSQL | $8 | ✅ Managed service |
| **Storage** | S3 | File uploads | $1 | ✅ Pay per use |
| **Domain** | Route 53 | DNS hosting | $1 | ✅ Professional setup |
| **SSL** | ACM | Free certificates | $0 | ✅ HTTPS everywhere |
| **Monitoring** | CloudWatch | Basic metrics | $2 | ✅ Essential monitoring |
| **Total** | | | **$30/month** | 🎯 **Perfect for intern test** |

### **Cost Comparison**
```yaml
❌ Over-engineered (ALB + ASG): $65/month
✅ Intern-optimized (Single EC2): $30/month
💰 Savings: $35/month (54% reduction)

Perfect balance of:
  ✅ Professional architecture
  ✅ Cost efficiency
  ✅ Impressive features
  ✅ Easy to maintain
```

### **Scaling Cost Projections**
```yaml
Current (Demo): $30/month
  - Handles 100+ concurrent users
  - Perfect for intern demonstration
  - All professional features included

Future (Production): $65/month
  - Add ALB + Auto Scaling
  - Multi-AZ deployment
  - Advanced monitoring
  - When actually needed
```

## 🚀 Deployment Strategy

### Frontend Deployment
```bash
# Build React app
npm run build

# Deploy to S3
aws s3 sync build/ s3://course-enrollment-frontend-prod --delete

# Invalidate CloudFront cache
aws cloudfront create-invalidation \
  --distribution-id E1234567890 \
  --paths "/*"
```

### Backend Deployment
```bash
# Build Go binary
GOOS=linux GOARCH=amd64 go build -o course-enrollment cmd/server/main.go

# Deploy to EC2
scp course-enrollment ec2-user@api.yourdomain.com:/opt/
ssh ec2-user@api.yourdomain.com "sudo systemctl restart course-enrollment"
```

## 🔒 Security Configuration

### S3 Bucket Policy
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::course-enrollment-frontend-prod/*"
    }
  ]
}
```

### S3 CORS Policy
```json
[
  {
    "AllowedHeaders": ["*"],
    "AllowedMethods": ["GET", "POST", "PUT", "DELETE"],
    "AllowedOrigins": ["https://app.yourdomain.com"],
    "ExposeHeaders": []
  }
]
```

### EC2 Security Groups
```yaml
Inbound Rules:
  - HTTP (80): 0.0.0.0/0
  - HTTPS (443): 0.0.0.0/0
  - SSH (22): Your IP only

Outbound Rules:
  - All traffic: 0.0.0.0/0
```

### IAM Roles & Policies
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "rds:DescribeDBInstances",
        "rds:Connect"
      ],
      "Resource": "arn:aws:rds:region:account:db:course-enrollment-db"
    },
    {
      "Effect": "Allow", 
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": "arn:aws:s3:::course-enrollment-assets/*"
    }
  ]
}
```

## 📊 Monitoring & Logging

### CloudWatch Monitoring
```yaml
Frontend Monitoring:
  - CloudFront cache hit ratio
  - S3 request metrics
  - Error rates (4xx, 5xx)

Backend Monitoring:
  - EC2 CPU/Memory utilization
  - API response times
  - HTTP status codes
  - Custom application metrics

Database Monitoring:
  - RDS CPU/Memory utilization
  - Database connections
  - Query performance
  - Storage usage
```

### Logging Strategy
```yaml
Application Logs:
  - Go application logs → CloudWatch Logs
  - Access logs → CloudWatch Logs
  - Error tracking with structured logging

Infrastructure Logs:
  - CloudFront access logs → S3
  - ALB access logs → S3
  - VPC Flow Logs (if needed)
```

## 🔧 Setup Scripts & Automation

### Frontend Setup
```bash
# Create React app
npx create-react-app course-enrollment-frontend
cd course-enrollment-frontend

# Install dependencies
npm install axios react-router-dom

# Environment variables
echo "REACT_APP_API_URL=https://api.yourdomain.com" > .env.production

# Build script in package.json
{
  "scripts": {
    "deploy": "npm run build && aws s3 sync build/ s3://course-enrollment-frontend-prod --delete"
  }
}
```

### Backend Setup
```bash
# EC2 User Data Script
#!/bin/bash
yum update -y
yum install -y docker
systemctl start docker
systemctl enable docker

# Pull and run Go application
docker run -d \
  --name course-enrollment \
  -p 80:8080 \
  -e DB_HOST=your-rds-endpoint \
  -e DB_PASSWORD=your-password \
  --restart unless-stopped \
  your-app:latest
```

### Health Check Script
```bash
#!/bin/bash
# /opt/health-check.sh
curl -f http://localhost:8080/health || {
  echo "App unhealthy, restarting..."
  docker restart course-enrollment
}

# Add to crontab
# */5 * * * * /opt/health-check.sh
```

## 📈 Scaling Strategy

### Current Capacity
```yaml
Expected Load:
  - Concurrent Users: 100-500
  - API Requests: 1000/hour
  - Database Connections: 20-50
  - Storage: 10GB/month
```

### Scaling Triggers
```yaml
Scale Up When:
  - CPU > 70% for 5 minutes
  - Memory > 80% for 5 minutes
  - Response time > 2 seconds
  - Error rate > 5%

Scale Down When:
  - CPU < 30% for 15 minutes
  - Memory < 50% for 15 minutes
```

### Migration Path to High Availability
```yaml
Phase 1 (Current): Single Instance
  - 1x EC2 t3.small
  - 1x RDS db.t3.micro
  - Cost: $39/month

Phase 2 (Growth): Auto Scaling
  - 2-4x EC2 instances
  - Application Load Balancer
  - RDS Multi-AZ
  - Cost: $120/month

Phase 3 (Scale): Multi-Region
  - Multiple regions
  - Read replicas
  - ElastiCache
  - Cost: $300+/month
```

## 🛠️ Infrastructure as Code

### Terraform Structure
```
infrastructure/
├── modules/
│   ├── frontend/
│   │   ├── s3.tf
│   │   ├── cloudfront.tf
│   │   └── variables.tf
│   ├── backend/
│   │   ├── ec2.tf
│   │   ├── security_groups.tf
│   │   └── variables.tf
│   └── database/
│       ├── rds.tf
│       ├── subnet_group.tf
│       └── variables.tf
├── environments/
│   ├── dev/
│   ├── staging/
│   └── prod/
└── main.tf
```

### CI/CD Pipeline
```yaml
# .github/workflows/deploy.yml
name: Deploy Course Enrollment Service

on:
  push:
    branches: [main]

jobs:
  deploy-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Node.js
        uses: actions/setup-node@v2
        with:
          node-version: '16'
      - name: Install dependencies
        run: npm install
      - name: Build
        run: npm run build
      - name: Deploy to S3
        run: aws s3 sync build/ s3://course-enrollment-frontend-prod --delete

  deploy-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Build
        run: GOOS=linux GOARCH=amd64 go build -o course-enrollment cmd/server/main.go
      - name: Deploy to EC2
        run: |
          scp course-enrollment ec2-user@${{ secrets.EC2_HOST }}:/opt/
          ssh ec2-user@${{ secrets.EC2_HOST }} "sudo systemctl restart course-enrollment"
```

## 🎯 Best Practices & Recommendations

### Security Best Practices
- ✅ HTTPS everywhere với ACM certificates
- ✅ Proper CORS configuration
- ✅ IAM roles với least privilege principle
- ✅ Security groups với minimal access
- ✅ Regular security updates
- ✅ Secrets management với AWS Secrets Manager

### Performance Optimization
- ✅ CloudFront CDN cho static assets
- ✅ Gzip compression
- ✅ Database connection pooling
- ✅ Proper caching strategies
- ✅ Image optimization
- ✅ Lazy loading cho frontend

### Cost Optimization
- ✅ Right-sizing EC2 instances
- ✅ S3 lifecycle policies
- ✅ CloudWatch cost monitoring
- ✅ Reserved instances cho production
- ✅ Spot instances cho development

### Disaster Recovery
- ✅ Automated RDS backups
- ✅ S3 versioning
- ✅ Cross-region backup strategy
- ✅ Infrastructure as Code
- ✅ Documented recovery procedures

## 📞 Support & Maintenance

### Monitoring Alerts
```yaml
Critical Alerts:
  - Application down (5xx errors > 50%)
  - Database connection failures
  - High CPU/Memory usage (>90%)
  - Disk space critical (>95%)

Warning Alerts:
  - Response time degradation
  - Error rate increase
  - Resource utilization trends
```

### Maintenance Schedule
```yaml
Daily:
  - Check application logs
  - Monitor key metrics
  - Verify backup completion

Weekly:
  - Review performance trends
  - Check security updates
  - Validate monitoring alerts

Monthly:
  - Cost analysis and optimization
  - Security audit
  - Capacity planning review
```

---

**Document Version**: 1.0
**Last Updated**: 2025-09-19
**Author**: Development Team
**Review Date**: Monthly
