# Cinema Ticket Booking System

ระบบจองตั๋วหนังแบบ Real-time ที่ออกแบบมาเพื่อรองรับการแย่งจองที่นั่งพร้อมกันหลายคน โดยเน้นเรื่อง concurrency, distributed lock, realtime update, security และโครงสร้างระบบที่พร้อมรันผ่าน Docker ได้ทั้งระบบ

## How to Run the System

### Prerequisites
- Docker
- Docker Compose

### Start all services
```bash
docker compose up --build
```

## Key Features

### User Side
- Login ด้วย Google OAuth 2.0 หรือ Firebase Auth
- ดูผังที่นั่งของรอบฉาย
- เห็นสถานะที่นั่งแบบ real-time
- เลือกที่นั่งและทำการ lock ได้ 5 นาที
- ยืนยัน booking ได้เมื่อชำระเงินสำเร็จ (mock payment flow)

### Admin Side
- ดูรายการ booking ทั้งหมด
- filter ข้อมูล booking ได้อย่างน้อย 1 เงื่อนไข
- ดู audit logs ของเหตุการณ์สำคัญในระบบ

### System Features
- ป้องกัน double booking
- รองรับ concurrent booking
- ใช้ Redis เป็น distributed lock
- ใช้ WebSocket/SSE สำหรับ realtime seat updates
- ใช้ Message Queue สำหรับ async event processing
- รันได้ทั้งระบบด้วย `docker compose up --build`

## Tech Stack Overview

โปรเจกต์นี้พัฒนาตามข้อกำหนดของ assignment โดยใช้ stack หลักดังนี้

### Backend
- **Go** — ใช้สำหรับพัฒนา REST API และ business logic หลักของระบบ
- **Gin** — ใช้เป็น web framework สำหรับจัดการ routing, middleware และ API endpoints

### Frontend
- **Vue 3** — ใช้พัฒนา user interface สำหรับฝั่ง user และ admin dashboard

### Database
- **MongoDB** — ใช้เก็บข้อมูลหลักของระบบ เช่น bookings, audit logs และข้อมูลที่เกี่ยวข้อง

### Cache / Distributed Lock
- **Redis** — ใช้สำหรับทำ distributed lock ของที่นั่ง เพื่อป้องกันการจองซ้ำระหว่างผู้ใช้หลายคน

### Realtime Communication
- **WebSocket** — ใช้กระจาย event การเปลี่ยนแปลงสถานะที่นั่งแบบ real-time ไปยังผู้ใช้ทุกคนที่กำลังดูรอบฉายเดียวกัน

### Message Queue
- **RabbitMQ** — ใช้สำหรับ asynchronous event handling เช่น booking success event, notification หรือ async logging

### Authentication
- **Firebase Auth** — ใช้ยืนยันตัวตนผู้ใช้ และเชื่อม user_id เข้ากับการจอง

### DevOps / Deployment
- **Docker**
- **Docker Compose**

ระบบทั้งหมดสามารถรันได้ด้วยคำสั่งเดียว:

```bash
docker compose up --build
```


## System Architecture Diagram
ภาพรวมการทำงานของระบบ

```md

           +---------------------+
           |     Frontend        |
           |   Vue 3 (User/Admin)|
           +----------+----------+
                      |
                      | HTTP / WebSocket
                      v
           +--------------------------+
           |     Backend API          |
           |      Go + Gin            |
           +----+-----+-----+---------+
                |     |              |
                |     |              |
                v     v              v
          +--------+ +--------+ +-------------+
          | MongoDB| | Redis  | | Message MQ  |
          |        | | Lock   | | RabbitMQ    |
          +--------+ +--------+ +-------------+
                |
                v
         +------------------+
         | Audit Logs /     |
         | Booking Records  |
         +------------------+
```
## Booking Flow

Booking flow เป็นหัวใจหลักของระบบ โดยออกแบบเพื่อป้องกันการจองซ้ำและรองรับผู้ใช้หลายคนที่กดจองพร้อมกัน

### Step 1: User Login
ผู้ใช้เข้าสู่ระบบผ่าน Google OAuth 2.0 หรือ Firebase Auth  
หลังจาก login สำเร็จ ระบบจะได้ `user_id` สำหรับใช้ผูกกับ booking

### Step 2: User Opens Seat Map
ผู้ใช้เลือกดูรอบฉาย ระบบจะโหลด seat map ของ showtime นั้น  
สถานะที่นั่งมี 3 แบบ:

- `AVAILABLE`
- `LOCKED`
- `BOOKED`

### Step 3: User Selects a Seat
เมื่อผู้ใช้เลือกที่นั่ง Backend จะพยายามสร้าง lock ของที่นั่งนั้นใน Redis โดยกำหนด TTL = 5 นาที

ถ้า lock สำเร็จ:
- สถานะที่นั่งจะถูกเปลี่ยนเป็น `LOCKED`
- ผู้ใช้คนอื่นจะไม่สามารถเลือกที่นั่งเดียวกันได้
- ระบบจะ broadcast event ผ่าน WebSocket เพื่ออัปเดตทุก client แบบ real-time

ถ้า lock ไม่สำเร็จ:
- แสดงว่า seat ถูกเลือกโดยผู้ใช้อื่นไปแล้ว
- ระบบจะตอบกลับว่าไม่สามารถ lock ที่นั่งนี้ได้

### Step 4: Payment / Confirmation
ภายในช่วงเวลาที่ seat ถูก lock ผู้ใช้ต้องดำเนินการยืนยัน booking (mock payment flow)

หากยืนยันสำเร็จ:
- ระบบเปลี่ยนสถานะ booking เป็น `BOOKED`
- ปลด lock ออกจาก Redis
- ส่ง event ผ่าน Message Queue
- broadcast สถานะ seat ใหม่ไปยังผู้ใช้ทุกคน

### Step 5: Timeout Handling
หากผู้ใช้ไม่ยืนยัน booking ภายใน 5 นาที:
- lock จะหมดอายุ
- seat จะกลับเป็น `AVAILABLE`
- ระบบจะบันทึก audit log
- ระบบ broadcast การเปลี่ยนแปลงแบบ real-time

## Redis Lock Strategy

Redis ถูกใช้เป็น distributed lock เพื่อป้องกันปัญหา double booking เมื่อมีผู้ใช้หลายคนพยายามจองที่นั่งเดียวกันพร้อมกัน

### Lock Key Design
ตัวอย่าง key:

```text
seat_lock:{showtime_id}:{seat_number}
```

โดยจะมีคีย์เก็บ showtime_id พร้อมกับ seat_number เพื่อให้รู้ว่าเป้นของที่นั่งไหนของโชว์อันไหน โดยจะให้มีอายุ 5 นาที เมื่อหมดเวลาล็อคก้จะถูกปล่อยอัตโนมัติ

## Message Queue Usage

ระบบมีการใช้ Message Queue สำหรับงาน asynchronous ภายในระบบ เพื่อแยกงานที่ไม่จำเป็นต้องทำใน request/response cycle ออกจาก flow หลัก

### Use Case
เมื่อ booking สำเร็จ ระบบจะ publish event เช่น:

- `BOOKING_SUCCESS`

consumer สามารถนำ event นี้ไปใช้ต่อได้ เช่น:
- บันทึก async log
- trigger notification แบบ mock
- ส่งข้อมูลไปยัง service อื่นในอนาคต

## API Endpoint

[API Documentation](https://3o605wo4v4.apidog.io)
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/showtimes` |Get All Showtime | No |
| GET | `/showtimes/:id/seats` | Get Seat for select showtime | No |
| POST | `/seats/:seatNumber/lock` | Locking seat | Yes |
| POST | `/seats/:seatNumber/unlock` | Unlock seat | Yes |
| POST | `/booking` | make booking on selected seat | Yes |
| POST | `/booking/:seat_number/confirm` | Confirm on pending booking | Yes |
| GET | `/admin/bookings` | Get bookings | Admin |
| GET | `/admin/logs` | Get eventlog | Admin |

## All protected routes require a JWT token in the header:
```
Authorization: Bearer <token>
```
Tokens are generated on login and contain user ID and role.

## Admin Role
Admin routes check for role: "admin" in the Header. The middleware admin enforces this.
