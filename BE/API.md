# API Documentation

## Authentication

### POST /v1/auth/login
Login user and generate JWT session token.

**Request Body:**
```json
{
  "nama": "string",
  "password": "string"
}
```

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciO..."
}
```

**Error Responses:**
- 400 Bad Request: Invalid request body or validation errors
- 404 Not Found: User not found
- 500 Internal Server Error

---

## Sensor Management

All sensor endpoints require Bearer token in Authorization header.

### GET /v1/sensor
Get list of sensors with pagination.

**Request Body:**
```json
{
  "pgn": {
    "page": 1,
    "limit": 10
  }
}
```

**Response (200 OK):**
```json
{
  "sensor_lists": [
    {
      "sensor_id": 1,
      "sensor_name": "Sensor 1",
      "sensor_area": "Area A",
      "sensor_status": "active"
    }
  ]
}
```

### PUT /v1/sensor/update
Update sensor status.

**Request Body:**
```json
{
  "sensor_id": 1,
  "sensor_status": "active"
}
```

**Response (200 OK):**
Success message or updated sensor data.

---

## Readings

All readings endpoints require Bearer token in Authorization header.

### GET /v1/reading
Get sensor readings with pagination.

**Request Body:**
```json
{
  "pagination": {
    "page": 1,
    "limit": 10
  }
}
```

**Response (200 OK):**
```json
{
  "readings": [
    {
      "sensor_id": 1,
      "humidity": 65.5,
      "timestamp": "2023-01-01T12:00:00Z",
      "pump_on": true
    }
  ]
}
```

---

## Notes
- All protected endpoints require `Authorization: Bearer <access_token>` header
- Pagination uses standard pagination structure with `page` and `limit` fields
- Error responses follow standard HTTP status codes with descriptive messages