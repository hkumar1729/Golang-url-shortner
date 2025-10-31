# Golang based URL shortner - BACKEND.

---

**URL shortner** is used for shortening of urls. It is a simple web service used to shorten long URLs into compact, shareable links. It helps in managing, sharing, and tracking links more efficiently.

---

### ✅ Core Features

- **URL Shortening**: Convert long URLs into short, easy-to-share links.
- **HTTPS Validation**: Ensures the provided URL starts with https:// (secure links only).
- **Custom Short Codes**: Generates a unique short key for each URL using SHA256 and Base64 encoding.
- **Redirection**: Redirects users from the short link to the original destination.
- **Scalable Design**: Backend implemented with Go (Gin framework) for speed and performance

---

### 🔑 **How It Works**

  - Client sends a POST request with the long URL

    POST /create
    Content-Type: application/json

    {
      "url": "https://domain.com/very/long/path"
    }

  - Server validates the URL
    - Checks if the URL starts with "https://"
    - Rejects invalid or malformed URLs

  - Server generates a short key
    - Uses SHA256 to hash the original URL
    - Encodes the hash using Base64
    - Example: "abc123"

  - Server stores the mapping
    Original URL  --->  Short Key
    "https://domain.com/very/long/path"  --->  "abc123"

  - Server returns the short URL

    {
      "short-url": "https://domain.com/abc123"
    }

  - When user visits the short URL
    GET /abc123  --->  Redirects to the original long URL

---

### ⚙️ Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin web framework
- **Hashing**: SHA256
- **Encoding**: Base64
- **Database**: Postgres + PrismaORM

---