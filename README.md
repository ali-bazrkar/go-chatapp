
# Go ChatApp

As the name suggests, this is a real-time chat application built in **Go**, using **Gorilla WebSockets** and **session-based authentication** with **CSRF protection**.
I will deep dive into technical implementation, things i learned during working on this project, design decisions i made and how they directly help scaling and future possible goals to be implemented. Let's dive into it!
 
### How did I end up here?

I’ve always been drawn to logic and problem-solving.
My journey began in **game development** as a **pixel artist** with artistic background, who also happend to be very curios about computer systems and mathematics. When I entered university, I initially aimed to study **AI and learning algorithms**, but over time, I realized that my passion leaned more toward understanding how systems work rather than training models.
As I explored different areas of computer science, I became fascinated by **backend and low-level programming**. Understanding how operating systems, databases, and compilers are built, and how concurrency allows systems to scale efficiently. I found myself captivated by the idea of building high-performance systems that could handle thousands of requests gracefully

## Project Motivation
Building a communication system was something I always wanted to try. While Go wasn’t my first language, studying its history and its focus on **concurrency as a first-class citizen** made it feel like the perfect choice to learn and implement such a system in.
This project became my practical introduction to backend engineering — learning about networking, HTTP communication, WebSockets, and API design. It’s not a perfect system, but rather a starting point. It represents a strong foundation for my ongoing journey in backend and concurrent system development.

## Screenshots

| Login Page | Registration Page | Chat Room |
| :---: | :---: | :---: |
|<img width="1910" height="920" alt="login!" src="https://github.com/user-attachments/assets/e8140437-3cd5-428a-8ff7-f8be8732d06e"/>| <img width="1908" height="926" alt="register" src="https://github.com/user-attachments/assets/9eacf1a7-7259-4cde-8e5a-4506dbe92eb1" /> |<img width="1920" height="912" alt="chatapp" src="https://github.com/user-attachments/assets/7ca00649-8295-40bd-818e-ddd8b2e3ec51" />|

## Features

1. **Concurrent Communication:** Each user can send and receive messages simultaneously. The `Hub` (chatroom) synchronizes message order, broadcasts messages, and stores them in the database as they are sent.
2. **Persistent History:** Older messages are retrieved when a user joins the chat, refreshes the page, or reconnects after a temporary disconnection. Only missed messages are fetched during temporary disconnects.
3. **User Join/Leave Notification:** Whenever a user joins or leaves, a message notifies all online users in real-time.
4. **Multi-Hub Support:** Currently, the system uses a single hub for simplicity. The database schema and implementation support multiple hubs with minimal changes.
5. **Cookie-Based Authorization:** Uses `session_cookies` (httpOnly) to protect against XSS and `csrf_token` to validate same-origin requests.
6. **Active Session Extension:** The expiration time of active users is extended on valid HTTP requests and periodically during active WebSocket connections, with **rate limiting** to reduce database load.
7. **Online Users List:** Online users are sent as JSON whenever someone joins or leaves, allowing clients to see who is currently online in real-time.
8. **Auto Reconnect:** The browser attempts to reconnect automatically if the connection is lost, stopping after multiple failed attempts.

## Tech Stack

* **Backend:** Go (net/http, Gorilla WebSocket, GORM, bcrypt)  
* **Frontend:** HTML + CSS + Vanilla JS (`index.html` — single-page app)  
* **Database:** SQLite  
* **Auth:** Cookie-based session + CSRF tokens  
* **Communication Protocol:** WebSocket


## How to Run Locally
### Prerequisites
* [Go](https://go.dev/doc/install) (version **1.24.6** or later recommended)
  
### Installation & Startup
1. **Clone the repository:**
```bash
git clone https://github.com/ali-bazrkar/go-chatapp.git
```
2. **Navigate to the project directory:**
```bash
cd go-chatapp
```

3. **Install dependencies:**
(Go will handle this automatically when you run the app, but you can fetch them explicitly if you want.)
```bash
go mod tidy
```

4. **Run the server:**
```bash
go run .
```
*Alternatively, you can build and run the executable:*
```bash
go build -o go-chatapp
./go-chatapp
```

5. **Open the application:**
Open your web browser and go to `http://localhost:3000`.


## API Overview

| Endpoint | Method | Description | Auth Required |
| ----------------- | ------ | --------------------------------- | ------------- |
| `/api/register` | POST | Handles new user registration. | ❌ |
| `/api/login` | POST | Authenticates a user and creates a new session. | ❌ |
| `/api/logout` | POST | End user session (CSRF required) | ✅ |
| `/api/check-auth` | GET | Verifies active sessions and extends their time | ✅ |
| `/ws` | GET | Opens WebSocket connection | ✅ |


## Project Structure

```
go-chatapp/
|
├─ auth/ ​​​​​​​​    ​        ​# Authentication & session management
| ├─ middleware.go
| ├─ session.go
| ├─ utils.go       ​ # password encryption & session generators
| └─ validators.go
├─ db/              ​ # Database setup & queries
| ├─ init.go        ​ # Database connection
| ├─ models.go      ​ # GORM schema structs
| └─ query.go       ​ # GORM CRUD functions
├─ handlers/        ​ # API handlers and websocket Endpoint
| ├─ auth.go       ​  ​# /api/check-auth
| ├─ login.go
| ├─ logout.go
| ├─ register.go
| ├─ setup.go       ​ # APIs and router setup
| └─ websocket.go
├─ model/           ​ # Global structs
| └─ message.go
├─ templates/       ​ # Frontend files (SPA)
| └─ index.html
├─ main.go           ​# Entry point
└─ go.mod           ​ # Dependencies
```

## Future Improvements (TODO)

1. Scale up to full multi-hub support, allowing users to create hubs themselves.
2. Allow users to delete their account or edit their username.
3. Consider switching to **JWT-based authorization**.
4. Rebuild the front-end from scratch using **React**.
