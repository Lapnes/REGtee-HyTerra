# REGtee HyTerra

**Smart Agriculture Edge Node for Automated Drip Hydroponics**

**REGtee HyTerra** adalah modul otomasi pertanian cerdas berbasis IoT yang dirancang untuk mengelola sistem hidroponik tetes secara mandiri. Sebagai bagian dari ekosistem **REGtee Cloud**, HyTerra berfungsi sebagai *edge node* yang melakukan pemantauan kelembapan media tanam secara *real-time* dan mengeksekusi siklus penyiraman nutrisi secara presisi.

## 🚀 Fitur Utama

* **Automated Drip System:** Mengatur sirkulasi nutrisi berdasarkan ambang batas (*threshold*) kelembapan tanah atau jadwal tertentu.
* **Dual Communication Protocol:** Mendukung pengiriman data via **MQTT** (untuk kontrol *real-time* dua arah) dan **REST API** (untuk *logging* data terpusat).
* **Centralized Telemetry:** Mengirimkan metrik kesehatan tanaman dan status perangkat ke *backend* REGtee secara berkala.
* **Fail-safe Logic:** Logika internal di ESP32 memastikan tanaman tetap tersiram meskipun koneksi internet terputus (mode *standalone*).

## 🛠️ Arsitektur Sistem

Sistem ini terdiri dari tiga lapisan utama yang terintegrasi:

1. **Hardware Layer (The Node):** ESP32 sebagai otak utama yang terhubung dengan modul relay, pompa, dan sensor kelembapan tanah.
2. **Communication Layer:** Jalur data menggunakan MQTT Broker (seperti Mosquitto) untuk latensi rendah.
3. **Service Layer (REGtee Cloud):** *Backend* (Go) yang mengolah data telemetri, menyimpan log ke dalam basis data, dan menyediakan antarmuka pemantauan.

## 📂 Struktur Repositori

```text
├── firmware/           # Kode C++/Arduino untuk ESP32 (PlatformIO)
├── backend/            # Layanan API berbasis Go/MQTT
├── schematics/         # Diagram rangkaian perangkat keras
└── docs/               # Dokumentasi tambahan dan aset gambar

```

## 🔌 Kebutuhan Perangkat Keras

* **Microcontroller:** ESP32 (DevKit V4)
* **Sensor:** Capacitive Soil Moisture Sensor v1.2
* **Actuator:** 12V Water Pump & Relay Module
* **Power:** External Power Supply 12V

## ⚙️ Cara Memulai

1. Kloning repositori ini:
```bash
git clone https://github.com/lapnes/regtee-hyterra.git

```


2. Konfigurasi kredensial Wi-Fi dan alamat MQTT Broker pada file `config.h` di folder `firmware/`.
3. *Flash* kode ke ESP32 menggunakan Arduino IDE atau PlatformIO.
4. Jalankan *backend* service di *environment* server kamu.

## ⚖️ License & Copyright

Proyek ini dilisensikan di bawah **[Creative Commons Attribution-NonCommercial 4.0 International (CC BY-NC 4.0)](LICENSE)**.

Anda bebas untuk menggunakan, memodifikasi, dan mendistribusikan proyek ini untuk tujuan personal, pendidikan, dan non-komersial, dengan kewajiban memberikan atribusi yang jelas kepada penulis asli.

### 🏢 Commercial Use

Penggunaan perangkat lunak, arsitektur *backend*, dan desain perangkat keras ini untuk tujuan **komersial sangat dilarang** di bawah lisensi standar. Hal ini mencakup (namun tidak terbatas pada) penggabungan kode ini ke dalam produk berbayar atau penggunaan untuk operasional pertanian skala perusahaan.

Jika Anda ingin menggunakan **REGtee HyTerra** untuk tujuan komersial, Anda wajib mendapatkan **Letter of Agreement (LOA) / Lisensi Komersial** terpisah. Silakan hubungi penulis secara langsung melalui profil GitHub ini untuk mendiskusikan persyaratan lisensi.

---

**Created with ❤️ by [Tenzly](https://github.com/lapnes) and [Etherift](https://github.com/etherift) | Part of REGtee Ecosystem**
