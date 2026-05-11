#include <WiFi.h>
#include <WiFiClientSecure.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include <time.h>
#include <sntp.h>
#include <esp_sntp.h>

// ============= NET CONFIG =============
bool timeSynced = false;

const char* WIFI_SSID = "Tenzly";
const char* WIFI_PASS = "Katyusha";

const char* MQTT_BROKER = "mqtt.tenzly.codes";
const int MQTT_PORT = 13792;
const char* MQTT_USER = "Kitasan";
const char* MQTT_PASS = "Kitasan1234";

// CLIENT_ID bisa diupdate via MQTT → reconnect otomatis
char CLIENT_ID[64] = "ESP32_SoilMonitor_001";

// ============= SSL CERT =============
const char mqttCert[] PROGMEM = R"EOF(
-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIRAIIQz7DSQONZRGPgu2OCiwAwDQYJKoZIhvcNAQELBQAw
TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMTUwNjA0MTEwNDM4
WhcNMzUwNjA0MTEwNDM4WjBPMQswCQYDVQQGEwJVUzEpMCcGA1UEChMgSW50ZXJu
ZXQgU2VjdXJpdHkgUmVzZWFyY2ggR3JvdXAxFTATBgNVBAMTDElTUkcgUm9vdCBY
MTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAK3oJHP0FDfzm54rVygc
h77ct984kIxuPOZXoHj3dcKi/vVqbvYATyjb3miGbESTtrFj/RQSa78f0uoxmyF+
0TM8ukj13Xnfs7j/EvEhmkvBioZxaUpmZmyPfjxwv60pIgbz5MDmgK7iS4+3mX6U
A5/TR5d8mUgjU+g4rk8Kb4Mu0UlXjIB0ttov0DiNewNwIRt18jA8+o+u3dpjq+sW
T8KOEUt+zwvo/7V3LvSye0rgTBIlDHCNAymg4VMk7BPZ7hm/ELNKjD+Jo2FR3qyH
B5T0Y3HsLuJvW5iB4YlcNHlsdu87kGJ55tukmi8mxdAQ4Q7e2RCOFvu396j3x+UC
B5iPNgiV5+I3lg02dZ77DnKxHZu8A/lJBdiB3QW0KtZB6awBdpUKD9jf1b0SHzUv
KBds0pjBqAlkd25HN7rOrFleaJ1/ctaJxQZBKT5ZPt0m9STJEadao0xAH0ahmbWn
OlFuhjuefXKnEgV4We0+UXgVCwOPjdAvBbI+e0ocS3MFEvzG6uBQE3xDk3SzynTn
jh8BCNAw1FtxNrQHusEwMFxIt4I7mKZ9YIqioymCzLq9gwQbooMDQaHWBfEbwrbw
qHyGO0aoSCqI3Haadr8faqU9GY/rOPNk3sgrDQoo//fb4hVC1CLQJ13hef4Y53CI
rU7m2Ys6xt0nUW7/vGT1M0NPAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNV
HRMBAf8EBTADAQH/MB0GA1UdDgQWBBR5tFnme7bl5AFzgAiIyBpY9umbbjANBgkq
hkiG9w0BAQsFAAOCAgEAVR9YqbyyqFDQDLHYGmkgJykIrGF1XIpu+ILlaS/V9lZL
ubhzEFnTIZd+50xx+7LSYK05qAvqFyFWhfFQDlnrzuBZ6brJFe+GnY+EgPbk6ZGQ
3BebYhtF8GaV0nxvwuo77x/Py9auJ/GpsMiu/X1+mvoiBOv/2X/qkSsisRcOj/KK
NFtY2PwByVS5uCbMiogziUwthDyC3+6WVwW6LLv3xLfHTjuCvjHIInNzktHCgKQ5
ORAzI4JMPJ+GslWYHb4phowim57iaztXOoJwTdwJx4nLCgdNbOhdjsnvzqvHu7Ur
TkXWStAmzOVyyghqpZXjFaH3pO3JLF+l+/+sKAIuvtd7u+Nxe5AW0wdeRlN8NwdC
jNPElpzVmbUq4JUagEiuTDkHzsxHpFKVK7q4+63SM1N95R1NbdWhscdCb+ZAJzVc
oyi3B43njTOQ5yOf+1CceWxG1bQVs5ZufpsMljq4Ui0/1lvh+wjChP4kqKOJ2qxq
4RgqsahDYVvTH9w7jXbyLeiNdd8XM2w9U/t7y0Ff/9yi0GE44Za4rF2LN9d11TPA
mRGunUHBcnWEvgJBQl9nJEiU0Zsnvgc/ubhPgXRR4Xq37Z0j4r7g1SgEEzwxA57d
emyPxgcYxn/eR44/KJ4EBs+lVDR3veyJm+kXQ99b21/+jh5Xos1AnX5iItreGCc=
-----END CERTIFICATE-----
)EOF";

// ============= TOPIC =============
const char* TOPIC_PUBLISH = "HyTerra/toBE";
const char* TOPIC_SUBSCRIBE = "HyTerra/toESP32";

// ============= SENSOR CONFIG =============
const int SENSOR_PIN = 34;              // GPIO34 (ADC1_CH6)
const int RELAY_PIN = 25;               // GPIO untuk relay


// Kalibrasi — bisa di-update via MQTT
int adcDry = 2924;                      // Udara (kering)
int adcWet = 1170;                      // Air (basah)
unsigned long publishInterval = 10000;  // 10 detik default


// ============= PUMP CONTROL =============
bool pump_active = true;                // true = pompa boleh nyala, false = disabled

// ============= GLOBAL =============
WiFiClientSecure net;
PubSubClient mqttClient(net);

unsigned long mqttReconnectTimer = 0;
const unsigned long MQTT_RECONNECT_INTERVAL = 5000;

unsigned long lastPublish = 0;
bool shouldReconnectMQTT = false;       // Flag kalau CLIENT_ID berubah

// ============= FORWARD DECLARATIONS =============
void connectWiFi();
bool reconnectMQTT();
void mqttCallback(char* topic, byte* payload, unsigned int length);
int readSoilMoisture(int& rawADC);
void publishSoilData();
void publishStatus(const char* message, const char* pub_status, bool success);
void publishAck(const char* refAction, bool success, const char* message, JsonObject* extra = nullptr);
bool isTimeSynced();
void setupTime();
String getTimeString();

// ============= SETUP =============
void setup() {
  Serial.begin(115200);
  delay(1000);
  
  analogReadResolution(12);
  analogSetAttenuation(ADC_11db);

  pinMode(RELAY_PIN, OUTPUT);
  digitalWrite(RELAY_PIN, LOW);
  
  net.setCACert(mqttCert);
  
  mqttClient.setServer(MQTT_BROKER, MQTT_PORT);
  mqttClient.setCallback(mqttCallback);
  
  connectWiFi();
  setupTime();
  publishAck("setup", true, "ESP32 Soil Monitor started", nullptr);
  publishStatus("ESP32 Soil Monitor started", "info", true);
}

// ============= TIME SYNC =============
bool isTimeSynced() {
  return time(nullptr) > 1000000000;
}

String getTimeString() {
  time_t now = time(nullptr);
  struct tm timeinfo;
  localtime_r(&now, &timeinfo);

  char buf[25];
  strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S", &timeinfo);
  return String(buf);
}

void setupTime() {
  setenv("TZ", "WIB-7", 1);
  tzset();

  sntp_set_sync_interval(6 * 60 * 60 * 1000UL);
  configTime(0, 0, "pool.ntp.org", "time.nist.gov");

  Serial.print("Waiting for NTP sync");
  uint32_t ntpStart = millis();
  
  while (!isTimeSynced()) {
    if (millis() - ntpStart > 30000) {
      Serial.println("\n[WARNING] NTP timeout, using millis() fallback");
      timeSynced = false;
      return;
    }
    Serial.print(".");
    delay(100);
  }
  
  timeSynced = true;
  Serial.println("\n[TIME] Synced: " + getTimeString());
}

// ============= WIFI =============
void connectWiFi() {
  Serial.print("Connecting to WiFi");
  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASS);  
  uint32_t start = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - start < 15000) {
    delay(500); Serial.print(".");
  }
  
  if (WiFi.status() == WL_CONNECTED) {
    Serial.printf("\n[WiFi] IP: %s\n", WiFi.localIP().toString().c_str());
  } else {
    Serial.println("\n[WiFi] Timeout! Retry next loop...");
  }
}

// ============= MQTT FUNCTIONS =============
bool reconnectMQTT() {
  if (millis() - mqttReconnectTimer < MQTT_RECONNECT_INTERVAL) {
    return mqttClient.connected();
  }
  mqttReconnectTimer = millis();
  
  if (mqttClient.connected()) return true;
  
  Serial.print("[MQTT] Connecting...");
  
  if (mqttClient.connect(CLIENT_ID, MQTT_USER, MQTT_PASS)) {
    Serial.printf(" connected to %s:%d as %s\n", MQTT_BROKER, MQTT_PORT, CLIENT_ID);
    mqttClient.subscribe(TOPIC_SUBSCRIBE);
    Serial.printf("[MQTT] Subscribed to %s\n", TOPIC_SUBSCRIBE);
    publishStatus("ESP32 Soil Monitor started", "info", true);
    return true;
  } else {
    Serial.printf(" failed, rc=%d. Retry in %ds...\n", 
                  mqttClient.state(), MQTT_RECONNECT_INTERVAL/1000);
    return false;
  }
}

// ============= ACK FEEDBACK =============
// WAJIB dipanggil setiap kali terima command agar BE tidak retry
void publishAck(const char* refAction, bool success, const char* message, JsonObject* extra) {
  if (!mqttClient.connected()) return;

  StaticJsonDocument<1024> doc;
  doc["device"] = CLIENT_ID;
  doc["type"] = "ack";
  doc["ref_action"] = refAction;
  doc["success"] = success;
  doc["message"] = message;
  doc["pump_active"] = pump_active;
  doc["uptime_sec"] = millis() / 1000;

  if (timeSynced) {
    doc["timestamp"] = getTimeString();
  }

  if (extra != nullptr) {
    doc["data"] = *extra;
  }

  char jsonBuffer[1024];
  serializeJson(doc, jsonBuffer);
  
  Serial.printf("[ACK] %s\n", jsonBuffer);
  mqttClient.publish(TOPIC_PUBLISH, jsonBuffer);
}

// ============= CALLBACK =============
void mqttCallback(char* topic, byte* payload, unsigned int length) {
  Serial.printf("[MQTT] Message received [%s]\n", topic);

  char msgBuffer[1024];
  if (length >= sizeof(msgBuffer) - 1) {
    Serial.println("[ERROR] Payload too large");
    publishAck("unknown", false, "Payload too large", nullptr);
    return;
  }

  memcpy(msgBuffer, payload, length);
  msgBuffer[length] = '\0';
  Serial.printf("[MQTT] Payload: %s\n", msgBuffer);

  StaticJsonDocument<1024> doc;
  DeserializationError error = deserializeJson(doc, msgBuffer);
  if (error) {
    Serial.printf("[ERROR] JSON parse failed: %s\n", error.c_str());
    publishAck("unknown", false, "Invalid JSON", nullptr);
    return;
  }

  if (!doc.containsKey("action")) {
    Serial.println("[ERROR] Missing 'action' field");
    publishAck("unknown", false, "Missing action field", nullptr);
    return;
  }

  const char* action = doc["action"];

  // -------------------------------------------------
  // 1. UPDATE CONFIG (publishInterval, adcDry, adcWet, client_id)
  // -------------------------------------------------
  if (strcmp(action, "update_config") == 0) {
    bool changed = false;
    StaticJsonDocument<256> ackData;
    JsonObject dataObj = ackData.to<JsonObject>();

    if (doc.containsKey("publishInterval")) {
      unsigned long newInterval = doc["publishInterval"];
      publishInterval = newInterval;
      changed = true;
      dataObj["publishInterval"] = publishInterval;
      Serial.printf("[CONFIG] publishInterval → %lu ms\n", publishInterval);
    }

    if (doc.containsKey("adcDry")) {
      int newDry = doc["adcDry"];
      adcDry = newDry;
      changed = true;
      dataObj["adcDry"] = adcDry;
      Serial.printf("[CONFIG] adcDry → %d\n", adcDry);
    }

    if (doc.containsKey("adcWet")) {
      int newWet = doc["adcWet"];
      adcWet = newWet;
      changed = true;
      dataObj["adcWet"] = adcWet;
      Serial.printf("[CONFIG] adcWet → %d\n", adcWet);
    }

    if (doc.containsKey("client_id")) {
      const char* newId = doc["client_id"];
      if (strcmp(newId, CLIENT_ID) != 0) {
        strlcpy(CLIENT_ID, newId, sizeof(CLIENT_ID));
        changed = true;
        shouldReconnectMQTT = true;  // Reconnect dengan ID baru di loop()
        dataObj["client_id"] = CLIENT_ID;
        Serial.printf("[CONFIG] client_id → %s (will reconnect)\n", CLIENT_ID);
      }
    }

    if (changed) {
      publishAck("update_config", true, "Config updated successfully", &dataObj);
    } else {
      publishAck("update_config", false, "No valid config fields provided", nullptr);
    }
  }

  // -------------------------------------------------
  // 2. SET PUMP (enable/disable)
  // -------------------------------------------------
  else if (strcmp(action, "set_pump") == 0) {
    if (!doc.containsKey("value")) {
      publishAck("set_pump", false, "Missing 'value' field (true/false)", nullptr);
      return;
    }

    bool newState = doc["value"];
    pump_active = newState;

    Serial.printf("[PUMP] Pump %s\n", pump_active ? "ENABLED" : "DISABLED");

    StaticJsonDocument<64> ackData;
    JsonObject dataObj = ackData.to<JsonObject>();
    dataObj["pump_active"] = pump_active;

    publishAck("set_pump", true, pump_active ? "Pump enabled" : "Pump disabled", &dataObj);

    // TODO: nanti tambahkan digitalWrite(RELAY_PIN, pump_active ? HIGH : LOW);
  }

  // -------------------------------------------------
  // 3. RESTART SYSTEM
  // -------------------------------------------------
  else if (strcmp(action, "restart") == 0) {
    publishAck("restart", true, "Restarting ESP32...", nullptr);
    Serial.println("[SYSTEM] Restarting in 500ms...");
    delay(500);  // Beri waktu ACK terkirim
    ESP.restart();
  }

  // -------------------------------------------------
  // 4. KALIBRASI (existing)
  // -------------------------------------------------
  else if (strcmp(action, "calibrate_dry") == 0) {
    adcDry = analogRead(SENSOR_PIN);
    Serial.printf("[ACTION] Calibrate DRY: ADC=%d\n", adcDry);
    
    StaticJsonDocument<64> ackData;
    JsonObject dataObj = ackData.to<JsonObject>();
    dataObj["adcDry"] = adcDry;
    publishAck("calibrate_dry", true, "Dry calibration updated", &dataObj);
  }
  else if (strcmp(action, "calibrate_wet") == 0) {
    adcWet = analogRead(SENSOR_PIN);
    Serial.printf("[ACTION] Calibrate WET: ADC=%d\n", adcWet);
    
    StaticJsonDocument<64> ackData;
    JsonObject dataObj = ackData.to<JsonObject>();
    dataObj["adcWet"] = adcWet;
    publishAck("calibrate_wet", true, "Wet calibration updated", &dataObj);
  }

  // -------------------------------------------------
  // 5. SET INTERVAL (existing, legacy)
  // -------------------------------------------------
  else if (strcmp(action, "set_interval") == 0) {
    unsigned long newInterval = doc["value"] | 10000;
    publishInterval = newInterval;
    Serial.printf("[ACTION] Interval set to %lu ms\n", publishInterval);
    
    StaticJsonDocument<64> ackData;
    JsonObject dataObj = ackData.to<JsonObject>();
    dataObj["publishInterval"] = publishInterval;
    publishAck("set_interval", true, "Interval updated", &dataObj);
  }

  // -------------------------------------------------
  // 6. FORCE READ (existing)
  // -------------------------------------------------
  else if (strcmp(action, "read_now") == 0) {
    Serial.println("[ACTION] Force read triggered");
    publishAck("read_now", true, "Force read executed", nullptr);
    lastPublish = 0;  // Trigger publish di loop berikutnya
  }

  // -------------------------------------------------
  // UNKNOWN COMMAND
  // -------------------------------------------------
  else {
    Serial.printf("[WARN] Unknown command: %s\n", action);
    publishAck(action, false, "Unknown command", nullptr);
  }
}

// ============= SENSOR & PUBLISH =============
int readSoilMoisture(int& rawADC) {
  rawADC = analogRead(SENSOR_PIN);
  int percent = map(rawADC, adcDry, adcWet, 0, 100);
  return constrain(percent, 0, 100);
}

void publishSoilData() {
  if (!mqttClient.connected()) return;

  int rawADC;
  int moisture = readSoilMoisture(rawADC);

  StaticJsonDocument<1024> doc;
  doc["device"] = CLIENT_ID;
  doc["moisture_percent"] = moisture;
  doc["raw_adc"] = rawADC;
  doc["adc_dry"] = adcDry;
  doc["adc_wet"] = adcWet;
  doc["interval_ms"] = publishInterval;
  doc["pump_active"] = pump_active;   // Kirim status pump ke BE

  if (timeSynced) {
    doc["timestamp"] = getTimeString();
  } else {
    doc["uptime_sec"] = millis() / 1000;
  }

  char jsonBuffer[1024];
  serializeJson(doc, jsonBuffer);

  Serial.printf("[PUB] %s\n", jsonBuffer);
  
  if (mqttClient.publish(TOPIC_PUBLISH, jsonBuffer)) {
    Serial.println("[PUB] Success");
  } else {
    Serial.println("[PUB] Failed!");
  }
}

void publishStatus(const char* message, const char* pub_status, bool success) {
  if (!mqttClient.connected()) return;

  StaticJsonDocument<256> doc;
  doc["device"] = CLIENT_ID;
  doc["type"] = pub_status;
  doc["message"] = message;
  doc["status"] = success ? "success" : "error";
  doc["pump_active"] = pump_active;
  doc["uptime_sec"] = millis() / 1000;

  char jsonBuffer[256];
  serializeJson(doc, jsonBuffer);
  mqttClient.publish(TOPIC_PUBLISH, jsonBuffer);
}

// ============= LOOP =============
void loop() {
  // Handle WiFi reconnect
  if (WiFi.status() != WL_CONNECTED) {
    connectWiFi();
  }

  // Handle MQTT reconnect (normal atau karena client_id berubah)
  if (shouldReconnectMQTT) {
    mqttClient.disconnect();
    shouldReconnectMQTT = false;
    mqttReconnectTimer = 0;  // Force reconnect sekarang
  }

  if (!mqttClient.connected()) {
    reconnectMQTT();
  } else {
    mqttClient.loop();
  }

  // Publish data sesuai interval
  unsigned long now = millis();
  if (now - lastPublish >= publishInterval) {
    lastPublish = now;
    publishSoilData();
  }
}