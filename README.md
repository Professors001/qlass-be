# Qlass-be

Backend ของระบบ Qlass — RESTful API สำหรับจัดการห้องเรียนออนไลน์ที่รองรับการสร้างคลาส, สื่อการสอน, แบบทดสอบ และมินิเกมแบบ real-time

## Tech Stack

- **Go 1.25**
- **Gin** Web Framework
- **PostgreSQL** (ผ่าน Supabase)
- **Redis** สำหรับ cache และ real-time features
- **MinIO/Supabase S3** สำหรับ file storage
- **GORM** เป็น ORM
- **JWT** สำหรับ authentication

## Features

- Authentication พร้อม role-based access (student / teacher)
- จัดการคลาสเรียน — สร้าง, เข้าร่วมด้วยโค้ด, ตั้งค่าคลาส
- อัปโหลดวัสดุการสอน (lecture, assignment) ผ่าน MinIO/Supabase
- สร้างและทำแบบทดสอบ (quiz) แบบ real-time
- มินิเกม interactive พร้อม game state management
- ระบบจัดการผู้ใช้และ profile

## Prerequisites

- [Go](https://golang.org) >= 1.25
- PostgreSQL database (Supabase แนะนำ)
- Redis instance
- MinIO หรือ Supabase S3 สำหรับ file storage

## Getting Started

**1. ติดตั้ง dependencies**

```bash
go mod download
```

**2. ตั้งค่า environment variables**

คัดลอก `.env.example` และแก้ไขค่าให้ตรงกับ environment ของคุณ:

```bash
cp .env.example .env
```

```env
APP_PORT=:8080
APP_ENV=development

SUPABASE_URL=postgresql://postgres:password@localhost:5432/postgres
JWT_SECRET=YOUR_JWT_SECRET
REDIS_URL=redis://localhost:6379

MINIO_ENDPOINT=localhost
MINIO_ACCESS_KEY=YOUR_ACCESS_KEY
MINIO_SECRET_KEY=YOUR_SECRET_KEY
MINIO_BUCKET_NAME=qlass-bucket
MINIO_USE_SSL=false

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=YOUR_EMAIL
SMTP_PASS=YOUR_PASSWORD
```

**3. รันด้วย Docker (แนะนำ)**

```bash
docker-compose up -d
```

**4. รัน development server**

```bash
# ติดตั้ง Air สำหรับ hot reload
go install github.com/air-verse/air@latest

# รันด้วย Air
air
```

หรือรันแบบธรรมดา:

```bash
go run cmd/main.go
```

เปิด [http://localhost:8080](http://localhost:8080) สำหรับ API

## Scripts

| Command | Description |
|---|---|
| `go run cmd/main.go` | รัน development server |
| `air` | รันด้วย hot reload |
| `go test ./...` | รัน unit tests |
| `go build -o qlass-be cmd/main.go` | Build สำหรับ production |

## Project Structure

```
qlass-be/
├── cmd/                    # Application entry point
├── adapters/              # External adapters
│   ├── api/              # REST API handlers
│   ├── cache/            # Redis cache
│   └── storage/          # File storage (MinIO/Supabase)
├── domain/                # Business logic
│   ├── entities/         # Database models
│   └── repositories/     # Repository interfaces
├── usecases/              # Application use cases
├── dtos/                  # Data transfer objects
├── router/                # HTTP routing
├── middleware/            # Custom middleware
├── config/                # Configuration management
├── transforms/            # Data transformations
└── utils/                 # Utility functions
```

## API Endpoints

### Authentication
- `POST /auth/register` - สมัครสมาชิก
- `POST /auth/login` - เข้าสู่ระบบ
- `POST /auth/logout` - ออกจากระบบ

### Classes
- `GET /classes` - ดูรายการคลาส
- `POST /classes` - สร้างคลาสใหม่
- `GET /classes/:id` - ดูข้อมูลคลาส
- `PUT /classes/:id` - แก้ไขคลาส
- `POST /classes/join` - เข้าร่วมคลาสด้วยโค้ด

### Quizzes & Games
- `GET /quizzes` - ดูรายการควิซ
- `POST /quizzes` - สร้างควิซใหม่
- `POST /quizzes/:id/start` - เริ่มเกม
- `POST /quizzes/:id/submit` - ส่งคำตอบ
- `GET /games/:id/state` - ดูสถานะเกมแบบ real-time

### Materials
- `GET /materials` - ดูรายการวัสดุการสอน
- `POST /materials` - อัปโหลดวัสดุการสอน
- `GET /materials/:id` - ดาวน์โหลดวัสดุการสอน

## Database Schema

Entities หลัก:
- **User** - ข้อมูลผู้ใช้ (student/teacher)
- **Class** - ข้อมูลคลาสเรียน
- **ClassEnrollment** - การลงทะเบียนเรียน
- **Quiz** - ข้อมูลควิซ
- **QuizQuestion** - คำถามในควิซ
- **QuizOption** - ตัวเลือกคำตอบ
- **ClassMaterial** - วัสดุการสอน
- **Attachment** - ไฟล์แนบ
- **QuizGameLog** - ประวัติการเล่นเกม

## Deployment

### สำหรับ Production

1. ตั้งค่า `APP_ENV=production`
2. ใช้ PostgreSQL แทน SQLite
3. ตั้งค่า SSL สำหรับ database connection
4. ใช้ Redis cluster ถ้าจำเป็น
5. ตั้งค่า CORS ให้ถูกต้อง

### ด้วย Docker

```bash
docker build -t qlass-be .
docker run -p 8080:8080 --env-file .env qlass-be
```
