# Course Enrollment Service 🎓

A robust, production-ready backend API service for managing course catalog and student enrollments for the Sonic University platform. Built with Go, Gin framework, PostgreSQL, Redis caching, and comprehensive security features.

## 🌟 Features

- **🔐 JWT Authentication** with role-based access control (Admin/Student)
- **📚 Course Management** with image upload support via AWS S3
- **👥 Student Enrollment System** with duplicate prevention
- **⚡ Redis Caching** for improved performance
- **🔒 HTTPS/SSL** with Let's Encrypt certificates
- **📖 Interactive API Documentation** with Swagger UI
- **🐳 Docker Support** for containerized deployment
- **☁️ AWS Integration** (RDS, S3, EC2)
- **🧪 Comprehensive Testing** (Unit & Integration tests)
- **🚀 CI/CD Pipeline** with GitHub Actions

## 🌐 Live Demo

- **🔗 API Base URL**: https://course-enrollment.site/api/v1/
- **📋 Swagger Documentation**: https://course-enrollment.site/swagger/index.html
- **💚 Health Check**: https://course-enrollment.site/health
- **🎨 Frontend**: https://app.course-enrollment.site

## 🔑 API Endpoints

### 🔐 Authentication
- `POST /api/v1/auth/login` - Admin login (JWT token)
- `GET /api/v1/auth/profile` - Get admin profile (Protected)

### 📚 Courses (Public Read, Admin Write)
- `GET /api/v1/courses` - Get all courses (Public)
- `GET /api/v1/courses/:id` - Get course by ID (Public)
- `POST /api/v1/courses` - Create course (Admin only)
- `POST /api/v1/courses/upload` - Create course with image (Admin only)
- `PUT /api/v1/courses/:id` - Update course (Admin only)
- `DELETE /api/v1/courses/:id` - Delete course (Admin only)

### 👥 Enrollments (Public)
- `POST /api/v1/enrollments` - Enroll student in course
- `GET /api/v1/students/:email/enrollments` - Get student enrollments

### 🛠️ Admin Management (Admin only)
- `GET /api/v1/admin/students` - Get all students
- `GET /api/v1/admin/enrollments` - Get all enrollments
- `DELETE /api/v1/admin/enrollments/:id` - Delete enrollment

### 📊 System
- `GET /health` - Health check with database & Redis status
- `GET /swagger/*` - Interactive API documentation

## 🚀 Quick Start

### 🔑 Admin Credentials
```
Username: admin
Password: admin!dev
```

### 🌐 Try the API Now

**1. Get all courses (Public)**
```bash
curl https://course-enrollment.site/api/v1/courses
```

**2. Login as admin**
```bash
curl -X POST https://course-enrollment.site/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin!dev"}'
```

**3. Create a new course (Admin only)**
```bash
curl -X POST https://course-enrollment.site/api/v1/courses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Advanced React Development",
    "description": "Master React hooks, context, and performance optimization",
    "difficulty": "Advanced"
  }'
```

**4. Enroll a student**
```bash
curl -X POST https://course-enrollment.site/api/v1/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "student_email": "student@example.com",
    "course_id": "YOUR_COURSE_ID"
  }'
```

### 💻 Local Development Setup

**Prerequisites:**
- Go 1.24+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose

**1. Clone & Setup**
```bash
git clone <repository-url>
cd course-enrollment-service
cp .env.example .env
# Edit .env with your database credentials
```

**2. Run with Docker**
```bash
docker-compose up -d
```

**3. Access locally**
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- Health: http://localhost:8080/health

## ⚙️ Configuration

### 🔐 Environment Variables

**Database (AWS RDS)**
- `DB_HOST` - RDS endpoint
- `DB_PORT` - Database port (5432)
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_SSLMODE` - SSL mode (require for RDS)

**Authentication**
- `JWT_SECRET` - JWT signing secret
- `ADMIN_USERNAME` - Admin username (default: admin)
- `ADMIN_PASSWORD` - Admin password (default: admin!dev)

**Redis Cache**
- `REDIS_HOST` - Redis host
- `REDIS_PORT` - Redis port (6379)
- `REDIS_PASSWORD` - Redis password
- `REDIS_DB` - Redis database number

**AWS S3 (Image Upload)**
- `AWS_REGION` - AWS region
- `AWS_ACCESS_KEY_ID` - AWS access key
- `AWS_SECRET_ACCESS_KEY` - AWS secret key
- `S3_BUCKET_NAME` - S3 bucket for images

**Server**
- `PORT` - Server port (default: 8080)
- `SKIP_MIGRATION` - Skip database migrations (default: false)

## 🏗️ Architecture

### 🌐 High-Level System Architecture
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                    INTERNET                                         │
└─────────────────────────┬───────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              CLOUDFLARE CDN                                        │
│                         (DNS, DDoS Protection, SSL)                                │
└─────────────────────────┬───────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                               VERCEL EDGE                                          │
│                    Frontend: app.course-enrollment.site                            │
│                         (React, Next.js, Global CDN)                               │
└─────────────────────────┬───────────────────────────────────────────────────────────┘
                          │ HTTPS API Calls
                          ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              AWS CLOUD                                             │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                        VPC (10.0.0.0/16)                                   │   │
│  │                                                                             │   │
│  │  ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │  │                    PUBLIC SUBNET (10.0.1.0/24)                     │   │   │
│  │  │                         AZ: ap-southeast-2a                         │   │   │
│  │  │                                                                     │   │   │
│  │  │  ┌─────────────────┐    ┌─────────────────┐                        │   │   │
│  │  │  │   Internet      │    │   NAT Gateway   │                        │   │   │
│  │  │  │   Gateway       │    │   (Managed)     │                        │   │   │
│  │  │  │                 │    │                 │                        │   │   │
│  │  │  └─────────────────┘    └─────────────────┘                        │   │   │
│  │  │                                                                     │   │   │
│  │  │  ┌─────────────────────────────────────────────────────────────┐   │   │   │
│  │  │  │                    EC2 Instance                             │   │   │   │
│  │  │  │                 t3.small (2 vCPU, 2GB)                      │   │   │   │
│  │  │  │              course-enrollment.site                         │   │   │   │
│  │  │  │                                                             │   │   │   │
│  │  │  │  ┌─────────────────┐  ┌─────────────────┐                  │   │   │   │
│  │  │  │  │  Nginx Proxy    │  │   Docker        │                  │   │   │   │
│  │  │  │  │  (Port 80/443)  │  │   Containers    │                  │   │   │   │
│  │  │  │  │  Let's Encrypt  │  │                 │                  │   │   │   │
│  │  │  │  │  SSL/TLS        │  │  ┌─────────────┐│                  │   │   │   │
│  │  │  │  │                 │  │  │ Go API      ││                  │   │   │   │
│  │  │  │  └─────────────────┘  │  │ (Port 8080) ││                  │   │   │   │
│  │  │  │                       │  │ Gin/GORM    ││                  │   │   │   │
│  │  │  │                       │  └─────────────┘│                  │   │   │   │
│  │  │  │                       │  ┌─────────────┐│                  │   │   │   │
│  │  │  │                       │  │ Redis Cache ││                  │   │   │   │
│  │  │  │                       │  │ (Port 6379) ││                  │   │   │   │
│  │  │  │                       │  └─────────────┘│                  │   │   │   │
│  │  │  │                       └─────────────────┘                  │   │   │   │
│  │  │  └─────────────────────────────────────────────────────────────┘   │   │   │
│  │  └─────────────────────────────────────────────────────────────────────┘   │   │
│  │                                                                             │   │
│  │  ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │  │                   PRIVATE SUBNET (10.0.2.0/24)                     │   │   │
│  │  │                         AZ: ap-southeast-2a                         │   │   │
│  │  │                                                                     │   │   │
│  │  │  ┌─────────────────────────────────────────────────────────────┐   │   │   │
│  │  │  │                    RDS PostgreSQL                           │   │   │   │
│  │  │  │                   db.t3.micro                               │   │   │   │
│  │  │  │                Multi-AZ: Disabled                           │   │   │   │
│  │  │  │              Storage: 20GB gp2                              │   │   │   │
│  │  │  │            course-enrollment-db                             │   │   │   │
│  │  │  └─────────────────────────────────────────────────────────────┘   │   │   │
│  │  └─────────────────────────────────────────────────────────────────────┘   │   │
│  │                                                                             │   │
│  │  ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │  │                   PRIVATE SUBNET (10.0.3.0/24)                     │   │   │
│  │  │                         AZ: ap-southeast-2b                         │   │   │
│  │  │                                                                     │   │   │
│  │  │  ┌─────────────────────────────────────────────────────────────┐   │   │   │
│  │  │  │                      S3 Bucket                              │   │   │   │
│  │  │  │               course-enrollment-images                      │   │   │   │
│  │  │  │                Standard Storage                              │   │   │   │
│  │  │  │              Public Read Access                              │   │   │   │
│  │  │  └─────────────────────────────────────────────────────────────┘   │   │   │
│  │  └─────────────────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 🛡️ Security Architecture
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              SECURITY LAYERS                                       │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  Layer 1: Network Security                                                         │
│  ├─ VPC Isolation (10.0.0.0/16)                                                   │
│  ├─ Public/Private Subnet Segregation                                             │
│  ├─ Internet Gateway (Public access)                                              │
│  ├─ NAT Gateway (Private subnet internet access)                                  │
│  └─ Route Tables (Traffic control)                                                │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  Layer 2: Security Groups (Stateful Firewall)                                     │
│  ├─ Web-SG (EC2): 80/443 from 0.0.0.0/0, 22 from Admin IP                       │
│  ├─ App-SG (EC2): 8080 from Web-SG, 6379 from localhost                          │
│  ├─ DB-SG (RDS): 5432 from App-SG only                                           │
│  └─ S3-SG: HTTPS access from App-SG                                               │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  Layer 3: Application Security                                                     │
│  ├─ HTTPS/TLS 1.3 (Let's Encrypt)                                                 │
│  ├─ JWT Authentication (HS256)                                                    │
│  ├─ Role-based Access Control (Admin/Student)                                     │
│  ├─ Input Validation & Sanitization                                               │
│  ├─ SQL Injection Prevention (GORM)                                               │
│  ├─ CORS Policy (Specific origins)                                                │
│  └─ Rate Limiting (Nginx)                                                         │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  Layer 4: Data Security                                                            │
│  ├─ Database Encryption at Rest                                                   │
│  ├─ SSL/TLS for Database Connections                                              │
│  ├─ S3 Bucket Policies (Least Privilege)                                          │
│  ├─ Environment Variables (Secrets)                                               │
│  └─ Password Hashing (bcrypt)                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 🔄 Data Flow Architecture
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                DATA FLOW                                           │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  1. User Request Flow                                                               │
│     Browser → Cloudflare → Vercel → API Gateway → EC2 → Go API                    │
│                                                                                     │
│  2. Authentication Flow                                                             │
│     Login Request → JWT Generation → Redis Cache → Response with Token            │
│                                                                                     │
│  3. Course Management Flow                                                          │
│     Admin Request → JWT Validation → Business Logic → PostgreSQL → Response       │
│                                                                                     │
│  4. Image Upload Flow                                                               │
│     File Upload → Validation → S3 Upload → URL Generation → Database Save         │
│                                                                                     │
│  5. Caching Strategy                                                                │
│     ├─ Course List: Redis Cache (TTL: 5 minutes)                                  │
│     ├─ User Sessions: Redis Store (TTL: 24 hours)                                 │
│     └─ Database Queries: Application-level caching                                │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  6. Monitoring & Logging                                                           │
│     ├─ Application Logs → CloudWatch Logs                                         │
│     ├─ Nginx Access Logs → Local Storage                                          │
│     ├─ Database Metrics → RDS Monitoring                                          │
│     └─ System Metrics → CloudWatch Metrics                                        │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 🌐 Network Flow & Ports
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              NETWORK TOPOLOGY                                      │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  External Traffic Flow:                                                             │
│  Internet (443) → Cloudflare → Vercel → AWS ALB (443) → EC2 Nginx (443/80)       │
│                                                                                     │
│  Internal Traffic Flow:                                                             │
│  Nginx (80/443) → Go API (8080) → PostgreSQL (5432)                              │
│                 → Redis (6379)                                                     │
│                 → S3 (443/HTTPS)                                                   │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  Port Configuration:                                                                │
│  ├─ Public Ports (Internet-facing):                                               │
│  │   ├─ 80 (HTTP) → Redirect to 443                                              │
│  │   ├─ 443 (HTTPS) → Nginx SSL Termination                                      │
│  │   └─ 22 (SSH) → Admin access only (IP restricted)                             │
│  │                                                                                │
│  ├─ Internal Ports (VPC-only):                                                    │
│  │   ├─ 8080 → Go API Server                                                     │
│  │   ├─ 6379 → Redis Cache                                                       │
│  │   ├─ 5432 → PostgreSQL Database                                               │
│  │   └─ 443 → S3 HTTPS API                                                       │
│  │                                                                                │
│  └─ Health Check Ports:                                                           │
│      ├─ /health → Application health                                              │
│      ├─ /metrics → Prometheus metrics                                             │
│      └─ /status → Nginx status                                                    │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 📊 Infrastructure Specifications
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           INFRASTRUCTURE SPECS                                     │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  🖥️ EC2 Instance (Backend):                                                        │
│     ├─ Type: t3.small (2 vCPU, 2 GB RAM)                                          │
│     ├─ OS: Amazon Linux 2023                                                      │
│     ├─ Storage: 20 GB gp3 SSD                                                     │
│     ├─ Network: Enhanced Networking                                               │
│     └─ Monitoring: CloudWatch Agent                                               │
│                                                                                     │
│  🗄️ RDS PostgreSQL (Database):                                                     │
│     ├─ Engine: PostgreSQL 15.4                                                   │
│     ├─ Instance: db.t3.micro (1 vCPU, 1 GB)                                      │
│     ├─ Storage: 20 GB gp2 (Auto-scaling enabled)                                 │
│     ├─ Backup: 7-day retention                                                   │
│     ├─ Multi-AZ: Disabled (Cost optimization)                                    │
│     └─ Encryption: At rest & in transit                                          │
│                                                                                     │
│  ⚡ Redis Cache:                                                                    │
│     ├─ Deployment: Docker container on EC2                                       │
│     ├─ Memory: 512 MB allocated                                                  │
│     ├─ Persistence: RDB snapshots                                                │
│     ├─ Eviction: allkeys-lru                                                     │
│     └─ Network: localhost only                                                   │
│                                                                                     │
│  📦 S3 Storage (Images):                                                           │
│     ├─ Bucket: course-enrollment-images                                          │
│     ├─ Region: ap-southeast-2                                                    │
│     ├─ Access: Public read, Private write                                        │
│     ├─ Encryption: AES-256                                                       │
│     └─ Lifecycle: 90-day IA transition                                           │
│                                                                                     │
│  🌐 Vercel (Frontend):                                                             │
│     ├─ Framework: Next.js 14                                                     │
│     ├─ Deployment: Global Edge Network                                           │
│     ├─ Domain: app.course-enrollment.site                                        │
│     ├─ SSL: Automatic (Vercel managed)                                           │
│     └─ CDN: Global distribution                                                  │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

## 🧪 Testing

### 🔬 Run Tests
```bash
# All tests
make test

# With coverage
make test-coverage

# Integration tests only
go test ./tests/...
```

### 📊 Test Coverage
- **Unit Tests**: Repository, Service, Handler layers
- **Integration Tests**: Full API endpoint testing
- **Database Tests**: CRUD operations and constraints
- **Authentication Tests**: JWT token validation

## 🗄️ Database Schema

### 👤 Users Table
```sql
- id (UUID, Primary Key)
- username (VARCHAR, UNIQUE, NOT NULL)
- password_hash (VARCHAR, NOT NULL)
- role (VARCHAR, CHECK: admin/student)
- created_at (TIMESTAMP)
```

### 📚 Courses Table
```sql
- id (UUID, Primary Key)
- title (VARCHAR, NOT NULL)
- description (TEXT, NOT NULL)
- difficulty (VARCHAR, CHECK: Beginner/Intermediate/Advanced)
- image_url (VARCHAR, NULLABLE) -- S3 image URL
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

### 📝 Enrollments Table
```sql
- id (UUID, Primary Key)
- student_email (VARCHAR, NOT NULL)
- course_id (UUID, Foreign Key → courses.id)
- enrolled_at (TIMESTAMP)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
- UNIQUE(student_email, course_id) -- Prevent duplicates
```

## 🛠️ Development

### 📋 Make Commands
```bash
make build          # Build the application
make run            # Build and run locally
make test           # Run all tests
make test-coverage  # Run tests with coverage
make docker-up      # Start with Docker Compose
make docker-down    # Stop Docker containers
make fmt            # Format Go code
make clean          # Clean build artifacts
```

### 📁 Project Structure
```
course-enrollment-service/
├── cmd/server/              # 🚀 Application entry point
├── internal/
│   ├── auth/               # 🔐 JWT authentication
│   ├── config/             # ⚙️ Configuration management
│   ├── database/           # 🗄️ Database connection & migrations
│   ├── handler/            # 🌐 HTTP request handlers
│   ├── middleware/         # 🛡️ Authentication middleware
│   ├── models/             # 📊 Data models & validation
│   ├── repository/         # 💾 Data access layer
│   ├── router/             # 🛣️ HTTP routing & CORS
│   └── service/            # 🧠 Business logic layer
├── tests/                  # 🧪 Integration tests
├── migrations/             # 📝 SQL migration scripts
├── docs/                   # 📖 Swagger documentation
├── scripts/                # 🔧 Setup & deployment scripts
└── certs/                  # 🔒 SSL certificates
```

## 🚨 Error Handling

### 📋 Standardized Error Response
```json
{
  "error": "Validation failed",
  "message": "Title is required",
  "details": {
    "field": "title",
    "code": "REQUIRED"
  }
}
```

### 📊 HTTP Status Codes
- `200` ✅ Success
- `201` ✅ Created
- `400` ❌ Bad Request (validation errors)
- `401` 🔒 Unauthorized (invalid/missing JWT)
- `403` 🚫 Forbidden (insufficient permissions)
- `404` 🔍 Not Found
- `409` ⚠️ Conflict (duplicate enrollment)
- `500` 💥 Internal Server Error

## 🚀 Deployment

### 🌐 Production Infrastructure
- **Frontend**: Vercel (https://app.course-enrollment.site)
- **Backend**: AWS EC2 with Docker
- **Database**: AWS RDS PostgreSQL
- **Cache**: Redis on EC2
- **Storage**: AWS S3 for images
- **SSL**: Let's Encrypt with Nginx reverse proxy
- **CI/CD**: GitHub Actions

### 📈 Monitoring & Health
- Health check endpoint: `/health`
- Database connection monitoring
- Redis connection status
- Automatic SSL certificate renewal

## 🔧 Site Reliability Engineering (SRE)

### 📊 Service Level Objectives (SLOs)
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                 SLO TARGETS                                        │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  🎯 Availability SLO: 99.9% (8.77 hours downtime/year)                            │
│     ├─ API Endpoints: 99.95% uptime                                               │
│     ├─ Database: 99.9% availability                                               │
│     ├─ Cache Layer: 99.5% availability                                            │
│     └─ Image Storage: 99.99% availability                                         │
│                                                                                     │
│  ⚡ Performance SLO: P95 < 500ms, P99 < 1000ms                                    │
│     ├─ API Response Time: P95 < 300ms                                             │
│     ├─ Database Query Time: P95 < 100ms                                           │
│     ├─ Cache Hit Ratio: > 85%                                                     │
│     └─ Image Load Time: P95 < 200ms                                               │
│                                                                                     │
│  🔄 Throughput SLO: 1000 RPS sustained                                            │
│     ├─ Peak Traffic: 2000 RPS (2x capacity)                                       │
│     ├─ Concurrent Users: 500 active sessions                                      │
│     └─ Database Connections: < 80% pool usage                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 🚨 Incident Response & Alerting
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                            ALERTING STRATEGY                                       │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  🔴 Critical Alerts (P0 - Immediate Response):                                     │
│     ├─ API Availability < 99% (5-minute window)                                   │
│     ├─ Database Connection Failures > 5%                                          │
│     ├─ SSL Certificate Expiry < 7 days                                            │
│     ├─ Disk Usage > 90%                                                           │
│     └─ Memory Usage > 95%                                                         │
│                                                                                     │
│  🟡 Warning Alerts (P1 - 30-minute Response):                                     │
│     ├─ API Response Time P95 > 500ms                                              │
│     ├─ Error Rate > 1%                                                            │
│     ├─ Cache Hit Ratio < 80%                                                      │
│     ├─ CPU Usage > 80%                                                            │
│     └─ Failed Login Attempts > 100/hour                                           │
│                                                                                     │
│  🔵 Info Alerts (P2 - 2-hour Response):                                           │
│     ├─ Deployment Notifications                                                   │
│     ├─ Backup Status Reports                                                      │
│     ├─ Security Scan Results                                                      │
│     └─ Performance Trend Analysis                                                 │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  📞 Escalation Matrix:                                                             │
│     ├─ P0: Immediate → On-call Engineer → Team Lead → Manager                     │
│     ├─ P1: 30min → Primary Engineer → Backup Engineer                             │
│     └─ P2: 2hr → Team Slack → Daily Standup                                       │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 📈 Monitoring & Observability Stack
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                          OBSERVABILITY STACK                                       │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  📊 Metrics (Prometheus + Grafana):                                               │
│     ├─ Application Metrics: /metrics endpoint                                     │
│     ├─ System Metrics: Node Exporter                                              │
│     ├─ Database Metrics: PostgreSQL Exporter                                      │
│     ├─ Redis Metrics: Redis Exporter                                              │
│     └─ Nginx Metrics: Nginx Prometheus Module                                     │
│                                                                                     │
│  📝 Logging (ELK Stack):                                                          │
│     ├─ Application Logs: Structured JSON logging                                  │
│     ├─ Access Logs: Nginx combined format                                         │
│     ├─ Error Logs: Go error tracking                                              │
│     ├─ Audit Logs: Authentication & authorization                                 │
│     └─ Security Logs: Failed login attempts                                       │
│                                                                                     │
│  🔍 Tracing (Jaeger):                                                             │
│     ├─ Request Tracing: End-to-end latency                                        │
│     ├─ Database Tracing: Query performance                                        │
│     ├─ Cache Tracing: Redis operations                                            │
│     └─ External API Tracing: S3 operations                                        │
│                                                                                     │
│  🚨 Alerting (AlertManager):                                                      │
│     ├─ Slack Integration: Team notifications                                      │
│     ├─ Email Alerts: Critical incidents                                           │
│     ├─ PagerDuty: On-call escalation                                              │
│     └─ Webhook: Custom integrations                                               │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 🔄 Deployment & Release Strategy
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         DEPLOYMENT STRATEGY                                        │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  🚀 Blue-Green Deployment:                                                        │
│     ├─ Zero-downtime deployments                                                  │
│     ├─ Instant rollback capability                                                │
│     ├─ Health check validation                                                    │
│     └─ Traffic switching via Load Balancer                                        │
│                                                                                     │
│  🎯 Canary Releases:                                                              │
│     ├─ 5% traffic → New version                                                   │
│     ├─ Monitor metrics for 10 minutes                                             │
│     ├─ 25% → 50% → 100% gradual rollout                                          │
│     └─ Automatic rollback on error threshold                                      │
│                                                                                     │
│  🔧 Feature Flags:                                                                │
│     ├─ Runtime feature toggling                                                   │
│     ├─ A/B testing capabilities                                                   │
│     ├─ Gradual feature rollout                                                    │
│     └─ Emergency feature disable                                                  │
│                                                                                     │
│  📋 Release Checklist:                                                            │
│     ├─ ✅ All tests pass (Unit + Integration)                                     │
│     ├─ ✅ Security scan completed                                                 │
│     ├─ ✅ Performance benchmarks met                                              │
│     ├─ ✅ Database migrations tested                                              │
│     ├─ ✅ Rollback plan documented                                                │
│     └─ ✅ Monitoring dashboards updated                                           │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 🛡️ Disaster Recovery & Business Continuity
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                        DISASTER RECOVERY PLAN                                      │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  💾 Backup Strategy:                                                              │
│     ├─ Database: Automated daily backups (7-day retention)                        │
│     ├─ Application: Git repository + Docker images                                │
│     ├─ Configuration: Infrastructure as Code (Terraform)                          │
│     ├─ SSL Certificates: Automated renewal + backup                               │
│     └─ Logs: 30-day retention in S3                                               │
│                                                                                     │
│  🔄 Recovery Procedures:                                                          │
│     ├─ RTO (Recovery Time Objective): < 1 hour                                    │
│     ├─ RPO (Recovery Point Objective): < 15 minutes                               │
│     ├─ Database Restore: Point-in-time recovery                                   │
│     ├─ Application Restore: Docker container deployment                           │
│     └─ DNS Failover: Automated health check switching                             │
│                                                                                     │
│  🌍 Multi-Region Strategy:                                                        │
│     ├─ Primary: ap-southeast-2 (Sydney)                                           │
│     ├─ Secondary: ap-southeast-1 (Singapore)                                      │
│     ├─ Database: Cross-region read replicas                                       │
│     ├─ Storage: S3 cross-region replication                                       │
│     └─ CDN: Global edge locations                                                 │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

## 📝 Git Conventional Commits

### 🎯 Commit Message Format
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### 📋 Commit Types
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           CONVENTIONAL COMMIT TYPES                                │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  🚀 feat: A new feature for the user                                              │
│     Example: feat(auth): add JWT token refresh mechanism                          │
│                                                                                     │
│  🐛 fix: A bug fix for the user                                                   │
│     Example: fix(api): resolve course enrollment duplicate issue                  │
│                                                                                     │
│  📚 docs: Documentation only changes                                              │
│     Example: docs(readme): update API endpoint documentation                      │
│                                                                                     │
│  💄 style: Changes that do not affect the meaning of the code                     │
│     Example: style(handler): format code according to gofmt                       │
│                                                                                     │
│  ♻️ refactor: A code change that neither fixes a bug nor adds a feature          │
│     Example: refactor(service): extract course validation logic                   │
│                                                                                     │
│  ⚡ perf: A code change that improves performance                                  │
│     Example: perf(db): add database indexes for course queries                    │
│                                                                                     │
│  � test: Adding missing tests or correcting existing tests                       │
│     Example: test(integration): add enrollment API test cases                     │
│                                                                                     │
│  🔧 build: Changes that affect the build system or external dependencies          │
│     Example: build(docker): update Go version to 1.24                             │
│                                                                                     │
│  👷 ci: Changes to CI configuration files and scripts                             │
│     Example: ci(github): add automated security scanning                          │
│                                                                                     │
│  🔨 chore: Other changes that don't modify src or test files                      │
│     Example: chore(deps): update dependencies to latest versions                  │
│                                                                                     │
│  ⏪ revert: Reverts a previous commit                                              │
│     Example: revert: feat(auth): add JWT token refresh mechanism                  │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 🏷️ Scopes & Examples
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              COMMIT SCOPES                                         │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  🔐 auth: Authentication & authorization                                           │
│     ├─ feat(auth): implement role-based access control                            │
│     ├─ fix(auth): resolve JWT token expiration issue                              │
│     └─ perf(auth): optimize password hashing performance                          │
│                                                                                     │
│  🌐 api: API endpoints & handlers                                                 │
│     ├─ feat(api): add course image upload endpoint                                │
│     ├─ fix(api): handle malformed JSON requests gracefully                        │
│     └─ docs(api): update Swagger documentation                                    │
│                                                                                     │
│  🗄️ db: Database operations & migrations                                          │
│     ├─ feat(db): add course enrollment constraints                                │
│     ├─ fix(db): resolve connection pool exhaustion                                │
│     └─ perf(db): optimize course listing query                                    │
│                                                                                     │
│  ⚡ cache: Redis caching layer                                                     │
│     ├─ feat(cache): implement course list caching                                 │
│     ├─ fix(cache): handle Redis connection failures                               │
│     └─ perf(cache): optimize cache key structure                                  │
│                                                                                     │
│  🐳 docker: Container & deployment                                                │
│     ├─ feat(docker): add multi-stage build optimization                          │
│     ├─ fix(docker): resolve container startup issues                             │
│     └─ chore(docker): update base image to Alpine 3.19                           │
│                                                                                     │
│  🧪 test: Testing & quality assurance                                             │
│     ├─ test(unit): add course service test coverage                               │
│     ├─ test(integration): add enrollment API tests                                │
│     └─ test(e2e): add complete user journey tests                                 │
│                                                                                     │
│  🔧 config: Configuration & environment                                           │
│     ├─ feat(config): add environment-specific configs                             │
│     ├─ fix(config): resolve SSL certificate path issues                          │
│     └─ chore(config): update default port configuration                          │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 🎯 Breaking Changes & Versioning
```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                          BREAKING CHANGES & SEMVER                                 │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  💥 Breaking Changes:                                                              │
│     feat(api)!: change course ID format to UUID                                   │
│                                                                                     │
│     BREAKING CHANGE: Course IDs are now UUIDs instead of integers.                │
│     This affects all API endpoints that use course IDs.                           │
│                                                                                     │
│     Migration guide:                                                               │
│     - Update client code to handle UUID format                                    │
│     - Run database migration script                                               │
│     - Update API documentation                                                     │
│                                                                                     │
│  📦 Semantic Versioning:                                                          │
│     ├─ MAJOR (X.0.0): Breaking changes                                            │
│     ├─ MINOR (0.X.0): New features (backward compatible)                          │
│     └─ PATCH (0.0.X): Bug fixes (backward compatible)                             │
│                                                                                     │
│  🏷️ Release Tags:                                                                  │
│     ├─ v1.0.0: Initial production release                                         │
│     ├─ v1.1.0: Add authentication features                                        │
│     ├─ v1.1.1: Fix authentication bug                                             │
│     └─ v2.0.0: API redesign (breaking changes)                                    │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

## �🤝 Contributing

### 📋 Development Workflow
1. 🍴 **Fork** the repository
2. 🌿 **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. ✨ **Make** your changes following coding standards
4. 🧪 **Add** tests for new functionality
5. ✅ **Ensure** all tests pass (`make test`)
6. 📝 **Commit** using conventional commits (`git commit -m 'feat(api): add amazing feature'`)
7. 🚀 **Push** to the branch (`git push origin feature/amazing-feature`)
8. 🔄 **Open** a Pull Request with detailed description

### 🎯 Code Quality Standards
```bash
# Run before committing
make fmt           # Format code
make lint          # Run linter
make test          # Run all tests
make security      # Security scan
make docs          # Update documentation
```

### 📋 Pull Request Checklist
- [ ] ✅ Follows conventional commit format
- [ ] 🧪 Tests added/updated and passing
- [ ] 📚 Documentation updated
- [ ] 🔒 Security considerations addressed
- [ ] ⚡ Performance impact assessed
- [ ] 🔄 Backward compatibility maintained
- [ ] 📊 Monitoring/alerting updated if needed

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**🎓 Built with ❤️ for Sonic University Platform**

For questions or support, please open an issue or contact the development team.
