# Database Seeding

File `database/seed.go` menyediakan fungsi untuk seeding database dengan data inisial.

## Struktur Data Seed

### Users
- **admin** - Password: `password123` (di-hash dengan SHA256 sebelum disimpan)
- **John Doe** - Password: `password123` (di-hash dengan SHA256 sebelum disimpan)

### Sensors
- 4 sensor dengan lokasi berbeda:
  - Sensor Ruang 1 - Ruang Tamu (active)
  - Sensor Ruang 2 - Kamar Tidur (active)
  - Sensor Ruang 3 - Dapur (active)
  - Sensor Ruang 4 - Kamar Mandi (inactive)

### Readings
- Data pembacaan kelembaban dari setiap sensor
- Data mencakup 24 jam terakhir dengan 3 interval per sensor
- Status semua readings adalah active (true)

## Cara Menjalankan Seed

1. **Pastikan database sudah running:**
   ```bash
   docker compose up -d
   ```

2. **Jalankan seed command:**
   ```bash
   go run . seed
   ```

3. **Output yang diharapkan:**
   ```
   Starting database seeding...
   Seeding users...
   Seeded user: admin
   Seeded user: John Doe
   Seeding sensors...
   Seeded sensor: Sensor Ruang 1
   Seeded sensor: Sensor Ruang 2
   Seeded sensor: Sensor Ruang 3
   Seeded sensor: Sensor Ruang 4
   Seeding readings...
   ...
   ✓ Seeding completed successfully!
   ```

## Testing Login

Setelah seed berhasil, Anda bisa test login dengan credentials:

```bash
curl -X POST http://localhost:8089/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nama": "admin", "password": "password123"}'
```

Response yang diharapkan (200 OK):
```json
{
  "result": {
    "access_token": "eyJhbGciO..."
  }
}
```

## Cara Kerja Password Hashing

1. **Saat Seed**: Password plain text (`password123`) di-hash menggunakan SHA256
2. **Saat Disimpan ke DB**: Password yang sudah di-hash disimpan ke database
3. **Saat Login**: Password yang dikirim user juga di-hash dengan algoritma yang sama, lalu dibandingkan dengan password yang tersimpan

## Reset Database

Untuk menghapus semua data dan memulai dari awal:

1. **Drop database (dalam Docker):**
   ```bash
   docker compose down -v
   ```

2. **Jalankan ulang database:**
   ```bash
   docker compose up -d
   ```

3. **Jalankan seed lagi:**
   ```bash
   go run . seed
   ```

## Catatan Penting

- Seed hanya akan menambahkan data jika data tersebut belum ada di database
- Jika Anda menjalankan seed berkali-kali, data tidak akan terduplikasi
- Password di-hash menggunakan SHA256 untuk keamanan
- Timestamp readings menggunakan waktu sekarang saat seed dijalankan

