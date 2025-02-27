# Video Transcoding System

ระบบสำหรับอัปโหลดวิดีโอและแปลงไฟล์วิดีโอเป็นคุณภาพต่างๆ เพื่อการสตรีมมิ่ง ระบบนี้ประกอบด้วย Backend API ที่พัฒนาด้วย Go

## คุณสมบัติหลัก

- ✅ อัปโหลดวิดีโอไปยัง S3 หรือ Minio
- ✅ แปลงวิดีโอเป็นความละเอียด 1080p และ 720p ที่ 24fps
- ✅ ตัดวิดีโอเป็นไฟล์ .ts ขนาด 10 วินาทีต่อไฟล์
- ✅ บันทึกข้อมูลไฟล์และ metadata ลงฐานข้อมูล
- ✅ อินเตอร์เฟซสำหรับอัปโหลดวิดีโอที่สวยงาม
- ✅ การเล่นวิดีโอแบบสตรีมมิ่ง
- ✅ สถาปัตยกรรมแบบ Clean Architecture

## โครงสร้างโปรเจค

โครงสร้างโปรเจคออกแบบตามหลัก Clean Architecture ประกอบด้วย:

```
video-transcoding-system/
├── cmd/                       # จุดเริ่มต้นของแอปพลิเคชัน
│   ├── api/                   # API เซิร์ฟเวอร์
│   └── migrate/               # ไฟล์สำหรับจัดการ Migration
├── internal/                  # โค้ดที่ใช้เฉพาะภายในแอปนี้
│   ├── domain/                # Entities (Domain Layer)
│   ├── usecase/               # Use Cases (Application Layer)
│   ├── adapter/               # Interface Adapters (Controller, Repository)
│   └── infrastructure/        # Frameworks & Drivers (External tools)
├── pkg/                       # โค้ดที่สามารถใช้ข้าม projects
│   ├── logger/                # ระบบบันทึกล็อก
│   └── config/                # ระบบการตั้งค่า
├── migrations/                # ไฟล์ SQL migrations
├── scripts/                   # สคริปต์อัตโนมัติต่างๆ
└── frontend/                  # ส่วน Frontend (React+Vite+Tailwind)
```

## ข้อกำหนดเบื้องต้น

- Go 1.23 หรือสูงกว่า
- FFmpeg และ FFprobe
- PostgreSQL
- Minio หรือ AWS S3

## การติดตั้ง

### การตั้งค่า Backend

1. โคลนโปรเจคนี้:
```bash
git clone https://github.com/yourusername/video-transcoding-system.git
cd video-transcoding-system
```

2. รันสคริปต์ตั้งค่า:
```bash
chmod +x setup.sh
./setup.sh
```

3. ปรับแต่งไฟล์ `.env` ตามความต้องการ:
```bash
nano .env
```

### การตั้งค่า Frontend

1. เข้าสู่ไดเรกทอรี frontend:
```bash
cd frontend
```

2. ติดตั้ง dependencies:
```bash
npm install
```

## การใช้งาน

### การรัน Backend

ด้วย Go โดยตรง:
```bash
./run.sh
```

หรือใช้ Makefile:
```bash
make run
```

### การรัน Frontend

```bash
cd frontend
npm run dev
```

### การรันด้วย Docker

รันทั้งระบบด้วย Docker Compose:
```bash
docker-compose up -d
```

## API Endpoints

### Public Endpoints

- `GET /api/v1/health` - ตรวจสอบสถานะ API

### Protected Endpoints (ต้องการ Authentication)

- `POST /api/v1/videos` - อัปโหลดวิดีโอใหม่
- `GET /api/v1/videos/:id` - ดึงข้อมูลวิดีโอตาม ID
- `GET /api/v1/users/videos` - ดึงรายการวิดีโอของผู้ใช้

## การพัฒนา

### การทดสอบ

รันการทดสอบทั้งหมด:
```bash
make test
```

### การ build

Build แอปพลิเคชัน:
```bash
make build
```

### การใช้ Docker

สร้าง image และรันคอนเทนเนอร์:
```bash
make docker-build
make docker-run
```

## การอัปเดต Database Schema

รัน migrations:
```bash
make migrate
```


## User Routes

API Endpoints
```
# Public Routes
POST /api/v1/auth/register     - สมัครสมาชิกใหม่
POST /api/v1/auth/login        - เข้าสู่ระบบ

# Protected Routes (ต้องการ Authentication)
GET /api/v1/users/me           - ดึงข้อมูลผู้ใช้ปัจจุบัน
PUT /api/v1/users/profile      - อัปเดตโปรไฟล์ผู้ใช้
PUT /api/v1/users/password     - เปลี่ยนรหัสผ่าน

# Admin Routes (ต้องการสิทธิ์ผู้ดูแลระบบ)
GET /api/v1/users              - ดึงรายการผู้ใช้ทั้งหมด
DELETE /api/v1/users/:id       - ลบผู้ใช้
PUT /api/v1/users/:id/role     - อัปเดตสิทธิ์ผู้ใช้
```



## ลิขสิทธิ์

โปรเจคนี้อยู่ภายใต้ [MIT License](LICENSE)
