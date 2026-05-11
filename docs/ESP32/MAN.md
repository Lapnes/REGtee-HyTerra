# MANUAL: ESP32 Soil Monitor

## Ringkasan
Dokumentasi ini menjelaskan:
- file `regtee-hyterra\ESP\REGtee HyTerra\src\main.cpp`
- file `regtee-hyterra\ESP\REGtee HyTerra\platformio.ini`

Fokus utama diberikan pada logika utama, konfigurasi MQTT/WiFi, dan format JSON yang digunakan untuk komunikasi.

---

## 1. `main.cpp` — arsitektur dan fungsi utama

### 1.1. Header dan dependensi
Mengimpor library berikut:
- `WiFi.h`: koneksi WiFi ESP32
- `WiFiClientSecure.h`: koneksi TLS/SSL untuk MQTT
- `PubSubClient.h`: klien MQTT
- `ArduinoJson.h`: parsing dan serialisasi JSON
- `time.h`, `sntp.h`, `esp_sntp.h`: sinkronisasi waktu NTP

### 1.2. Konfigurasi jaringan dan MQTT
Variabel global:
- `WIFI_SSID`, `WIFI_PASS`: kredensial WiFi
- `MQTT_BROKER`, `MQTT_PORT`, `MQTT_USER`, `MQTT_PASS`: parameter broker MQTT
- `CLIENT_ID`: ID perangkat MQTT, dapat diperbarui via perintah MQTT

MQTT topic:
- `TOPIC_PUBLISH = "HyTerra/toBE"`
- `TOPIC_SUBSCRIBE = "HyTerra/toESP32"`

### 1.3. Sertifikat SSL untuk MQTT
- `mqttCert[] PROGMEM`: certificate authority (CA) root certificate yang digunakan oleh `WiFiClientSecure`.
- Ini memungkinkan koneksi TLS ke broker MQTT.

### 1.4. Konfigurasi sensor dan relay
- `SENSOR_PIN = 34` → ADC1_CH6 pada ESP32
- `RELAY_PIN = 25`

Kalibrasi default:
- `adcDry = 2924` (nilai ADC udara/kering)
- `adcWet = 1170` (nilai ADC air/basah)
- `publishInterval = 10000` ms (10 detik)

Kontrol pompa:
- `pump_active = true`

### 1.5. Objek global MQTT/WiFi
- `WiFiClientSecure net`
- `PubSubClient mqttClient(net)`
- Timer reconnect MQTT dan flag `shouldReconnectMQTT`

### 1.6. Deklarasi fungsi
Fungsi yang dideklarasikan di awal:
- `connectWiFi()`
- `reconnectMQTT()`
- `mqttCallback(...)`
- `readSoilMoisture(...)`
- `publishSoilData()`
- `publishStatus(...)`
- `publishAck(...)`
- `isTimeSynced()`
- `setupTime()`
- `getTimeString()`

---

## 2. `setup()`
Langkah awal pada startup:
1. `Serial.begin(115200)`
2. `analogReadResolution(12)` dan `analogSetAttenuation(ADC_11db)`
3. Menyeting `RELAY_PIN` sebagai output dan mematikan relay
4. Mengatur sertifikat CA ke `net`
5. Menyeting server MQTT dan callback
6. Memanggil `connectWiFi()` untuk koneksi WiFi
7. Memanggil `setupTime()` untuk sinkronisasi NTP
8. Mengirim ACK awal dan status awal ke broker MQTT

---

## 3. Waktu dan sinkronisasi NTP
### `setupTime()`
- Mengatur timezone ke `WIB-7`
- `configTime(0, 0, "pool.ntp.org", "time.nist.gov")`
- Menunggu sinkronisasi hingga 30 detik
- Jika berhasil, `timeSynced = true`
- Jika timeout, fallback tanpa timestamp nyata

### `isTimeSynced()`
- Mengembalikan `true` jika `time(nullptr) > 1000000000`

### `getTimeString()`
- Mengembalikan waktu dalam format `YYYY-MM-DD HH:MM:SS`

---

## 4. WiFi
### `connectWiFi()`
- Mode `WIFI_STA`
- Memulai koneksi ke SSID dan password
- Menunggu hingga 15 detik
- Mencetak IP jika berhasil
- Jika timeout, hanya mencetak pesan, tidak reboot

---

## 5. MQTT, Listening, dan Callback
### 5.1. MQTT flow umum
- ESP32 mencoba menyambung ke broker MQTT `mqtt.tenzly.codes` di port `13792`.
- Setelah koneksi sukses, ESP32 `subscribe` ke topik `HyTerra/toESP32`.
- Semua output MQTT dikirim ke topik `HyTerra/toBE`.
- `mqttClient.loop()` dipanggil setiap iterasi `loop()` untuk:
  - memproses pesan masuk
  - menjaga koneksi MQTT tetap aktif
  - menjalankan callback bila ada pesan baru

### 5.2. `reconnectMQTT()`
- Menunggu setidaknya `MQTT_RECONNECT_INTERVAL = 5000` ms untuk mencoba ulang.
- Jika `mqttClient.connected()` false, fungsi mencoba terhubung lagi.
- Koneksi dilakukan dengan:
  - `CLIENT_ID`
  - `MQTT_USER`
  - `MQTT_PASS`
- Jika koneksi berhasil:
  - subscribe ke `TOPIC_SUBSCRIBE`
  - kirim `publishStatus("ESP32 Soil Monitor started", "info", true)` sebagai notifikasi awal
- Jika gagal, hanya mencetak kode error dan akan retry di iterasi berikutnya.

### 5.3. `mqttCallback(...)` — detail listening dan parsing
`mqttCallback(char* topic, byte* payload, unsigned int length)` adalah fungsi yang dipanggil otomatis ketika broker mengirim pesan ke topik yang disubscribe.

#### aliran callback:
1. Fungsi menerima `topic`, `payload`, dan `length`.
2. Payload dibaca ke buffer lokal `msgBuffer` dan ditambahkan terminator `\0`.
3. Payload kemudian diparse sebagai JSON melalui `deserializeJson(doc, msgBuffer)`.
4. Jika parse gagal atau field `action` tidak ditemukan, callback mengirim ACK error dan berhenti.
5. Jika JSON valid, callback membaca nilai `action` dan mengeksekusi command yang sesuai.

#### validasi callback:
- `topic` hanya digunakan untuk logging; semua perintah dikirim dari `HyTerra/toESP32`.
- `payload` harus berupa JSON valid.
- `action` wajib ada dan berupa string.

#### perilaku callback:
- Perintah valid ditangani langsung di callback.
- Setelah mengeksekusi perintah, callback memanggil `publishAck(...)` untuk memberitahu backend apakah perintah diproses berhasil.
- Jika `client_id` diubah, callback menetapkan `shouldReconnectMQTT = true` supaya loop utama melakukan reconnect dengan ID baru.

---

## 6. Format JSON input / perintah MQTT
Semua perintah masuk melalui payload MQTT di topic `HyTerra/toESP32`.

### aturan umum JSON command
- Tipe payload harus: objek JSON
- Field wajib:
  - `action`: string
- Field tambahan tergantung action.
- Jika JSON tidak valid, maka `mqttCallback()` akan mengirim ACK error.
- Jika `action` tidak dikenali, callback juga mengirim ACK error.

### struktur command standar
```
{
  "action": "<nama_action>",
  ...field lain sesuai action...
}
```

#### 6.1. `update_config`
Digunakan untuk memperbarui parameter runtime.
- `publishInterval`: integer, periode publish dalam ms
- `adcDry`: integer, kalibrasi nilai kering
- `adcWet`: integer, kalibrasi nilai basah
- `client_id`: string, MQTT client ID baru

Contoh:
```
{
  "action": "update_config",
  "publishInterval": 15000,
  "adcDry": 3000,
  "adcWet": 1100,
  "client_id": "ESP32_SoilMonitor_002"
}
```

#### 6.2. `set_pump`
- `value`: boolean

Contoh:
```
{
  "action": "set_pump",
  "value": true
}
```

#### 6.3. `restart`
Tanpa field tambahan.

Contoh:
```
{
  "action": "restart"
}
```

#### 6.4. `calibrate_dry`
- Membaca ADC saat sensor berada di udara/kering.

Contoh:
```
{
  "action": "calibrate_dry"
}
```

#### 6.5. `calibrate_wet`
- Membaca ADC saat sensor berada di air/basah.

Contoh:
```
{
  "action": "calibrate_wet"
}
```

#### 6.6. `set_interval`
Alias lama untuk menetapkan interval publish.
- `value`: integer ms

Contoh:
```
{
  "action": "set_interval",
  "value": 20000
}
```

#### 6.7. `read_now`
Menandai ESP32 agar mempublish data pada loop berikutnya.

Contoh:
```
{
  "action": "read_now"
}
```

---

## 7. Format JSON output / publish
Semua payload keluar dipublish ke topik `HyTerra/toBE`.
Semua publish hanya dijalankan jika MQTT sudah terkoneksi.

### 7.1. `publishAck(...)`
Digunakan sebagai respon langsung terhadap perintah yang diterima.

#### Struktur JSON ACK
```
{
  "device": "ESP32_SoilMonitor_001",
  "type": "ack",
  "ref_action": "set_pump",
  "success": true,
  "message": "Pump enabled",
  "pump_active": true,
  "uptime_sec": 125,
  "timestamp": "2026-05-11 12:34:56"
}
```

Penjelasan field:
- `device`: client ID MQTT saat ini
- `type`: selalu `"ack"`
- `ref_action`: action yang diproses
- `success`: boolean hasil eksekusi
- `message`: deskripsi hasil
- `pump_active`: status pompa saat ini
- `uptime_sec`: waktu hidup ESP dalam detik
- `timestamp`: hanya ada bila NTP sinkron
- `data`: objek tambahan bila ada informasi extra

### 7.2. `publishStatus(...)`
Digunakan untuk informasi status umum, bukan respons spesifik action.

#### Struktur JSON status
```
{
  "device": "ESP32_SoilMonitor_001",
  "type": "info",
  "message": "ESP32 Soil Monitor started",
  "status": "success",
  "pump_active": true,
  "uptime_sec": 5
}
```

Penjelasan field:
- `type`: tipe pesan status, misalnya `"info"`
- `status`: `"success"` atau `"error"`

### 7.3. `publishSoilData()`
Dipanggil secara berkala setiap `publishInterval`.

#### Struktur JSON sensor
```
{
  "device": "ESP32_SoilMonitor_001",
  "moisture_percent": 42,
  "raw_adc": 2100,
  "adc_dry": 2924,
  "adc_wet": 1170,
  "interval_ms": 10000,
  "pump_active": true,
  "timestamp": "2026-05-11 12:34:56"
}
```

Jika NTP belum sinkron, `timestamp` tidak disertakan dan diganti `uptime_sec`.

Penjelasan field:
- `moisture_percent`: kelembapan tanah dalam persen
- `raw_adc`: nilai ADC sensor mentah
- `adc_dry`, `adc_wet`: nilai kalibrasi kering/basah
- `interval_ms`: interval publish yang sedang digunakan
- `pump_active`: status pompa saat ini

---

---

## 8. Logika sensor dan publish
### `readSoilMoisture(int& rawADC)`
- Membaca `analogRead(SENSOR_PIN)`
- Mengkonversi ke persentase dengan `map(rawADC, adcDry, adcWet, 0, 100)`
- `constrain(...)` ke rentang 0-100

### `publishSoilData()`
- Membuat JSON sensor
- Mencetak ke Serial
- Mempublish ke MQTT topik `HyTerra/toBE`

---

## 9. Loop utama
### `loop()`
1. Cek koneksi WiFi, panggil `connectWiFi()` jika terputus
2. Jika `shouldReconnectMQTT` true, disconnect dan reset timer
3. Jika MQTT belum terkoneksi, `reconnectMQTT()`
4. Jika terkoneksi, panggil `mqttClient.loop()`
5. Publish data setiap `publishInterval`

---

## 10. `platformio.ini` — konfigurasi build
File `regtee-hyterra\ESP\REGtee HyTerra\platformio.ini` memuat:

```ini
[env:esp32doit-devkit-v1]
platform = espressif32
board = esp32doit-devkit-v1
framework = arduino
monitor_speed = 115200
build_flags = 
    -DMQTT_MAX_PACKET_SIZE=1024
    -DCORE_DEBUG_LEVEL=3

lib_deps =
;MQTT and JSON Libraries
    knolleary/PubSubClient@^2.8
    bblanchon/ArduinoJson@^6.21.2
```

### 10.1. Pilihan board dan framework
- `platform = espressif32`
- `board = esp32doit-devkit-v1`
- `framework = arduino`

### 10.2. Serial monitor
- `monitor_speed = 115200`

### 10.3. Build flags
- `-DMQTT_MAX_PACKET_SIZE=1024`: memperbesar ukuran paket MQTT maksimum agar memungkinkan payload JSON yang cukup besar.
- `-DCORE_DEBUG_LEVEL=3`: mengaktifkan debug level tinggi pada ESP32 core.

### 10.4. Library dependencies
- `knolleary/PubSubClient@^2.8`
- `bblanchon/ArduinoJson@^6.21.2`

---

## 11. konfigurasi

### 11.1. Konfigurasi utama
Konfigurasi di `main.cpp` dibangun sebagai variabel global yang bisa dimodifikasi lewat MQTT:
- `WIFI_SSID`, `WIFI_PASS` — kredensial WiFi
- `MQTT_BROKER`, `MQTT_PORT`, `MQTT_USER`, `MQTT_PASS` — koneksi broker MQTT
- `CLIENT_ID` — client ID MQTT yang bisa diupdate lewat perintah `update_config`
- `TOPIC_PUBLISH`, `TOPIC_SUBSCRIBE` — topic MQTT statis
- `SENSOR_PIN`, `RELAY_PIN` — pin sensor dan relay
- `adcDry`, `adcWet` — nilai kalibrasi sensor tanah kering/ basah
- `publishInterval` — interval publish data dalam milidetik
- `pump_active` — flag logika pompa aktif/disable

### 11.2. Konfigurasi build / project
`platformio.ini` mengatur konfigurasi build dan dependensi:
- `platform = espressif32`
- `board = esp32doit-devkit-v1`
- `framework = arduino`
- `monitor_speed = 115200`
- `build_flags`:
  - `-DMQTT_MAX_PACKET_SIZE=1024` memperbesar paket MQTT khusus JSON
  - `-DCORE_DEBUG_LEVEL=3` mengaktifkan debug ESP32
- `lib_deps`:
  - `knolleary/PubSubClient@^2.8`
  - `bblanchon/ArduinoJson@^6.21.2`

---

## 12. Rangkuman JSON penting
### JSON output `ack`
- `device`, `type`, `ref_action`, `success`, `message`, `pump_active`, `uptime_sec`, `timestamp`, `data`

### JSON output `status`
- `device`, `type`, `message`, `status`, `pump_active`, `uptime_sec`

### JSON output sensor
- `device`, `moisture_percent`, `raw_adc`, `adc_dry`, `adc_wet`, `interval_ms`, `pump_active`, `timestamp`/`uptime_sec`

### JSON input command
- `action`: wajib
- `publishInterval`, `adcDry`, `adcWet`, `client_id`, `value`: optional sesuai action

---

## 13. Cara pakai singkat
1. Upload kode ke ESP32 dengan PlatformIO.
2. Pastikan WiFi berhasil tersambung.
3. Pastikan broker MQTT `mqtt.tenzly.codes` dapat diakses.
4. Kirim payload JSON ke topik `HyTerra/toESP32`.
5. Terima feedback dan data publikasi di `HyTerra/toBE`.
