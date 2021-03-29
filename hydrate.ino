#include <ArduinoWebsockets.h>
#include <WiFi.h>
#include <ArduinoJson.h>

#define UPDATE_INTERVAL (1000 * 10)
int last_update_sent = 0;

const char *ssid = "<wifi name>";
const char *password = "<wifi password>";
const char *websockets_server = "wss://example.com/ws";

using namespace websockets;

void onMessageCallback(WebsocketsMessage message)
{
    Serial.print("Got Message: ");
    Serial.println(message.data());
    StaticJsonDocument<200> doc;

    deserializeJson(doc, message.data());

    if (doc["event"] == "hydrate")
    {
        digitalWrite(2, LOW);
        delay(3000);
        digitalWrite(2, HIGH);
    }
}

void onEventsCallback(WebsocketsEvent event, String data)
{
    if (event == WebsocketsEvent::ConnectionOpened)
    {
        Serial.println("Connnection Opened");
    }
    else if (event == WebsocketsEvent::ConnectionClosed)
    {
        Serial.println("Connnection Closed");
    }
    else if (event == WebsocketsEvent::GotPing)
    {
        Serial.println("Got a Ping!");
    }
    else if (event == WebsocketsEvent::GotPong)
    {
        Serial.println("Got a Pong!");
    }
}

WebsocketsClient client;
void setup()
{
    Serial.begin(115200);

    pinMode(2, OUTPUT);
    digitalWrite(2, HIGH);
    
    WiFi.begin(ssid, password);

    for (int i = 0; i < 10 && WiFi.status() != WL_CONNECTED; i++)
    {
        Serial.print(".");
        delay(1000);
    }

    client.onMessage(onMessageCallback);
    client.onEvent(onEventsCallback);

    client.connect(websockets_server);

    client.send("Hi Server!");

    client.ping();
}

void loop()
{
    if (millis() - last_update_sent > UPDATE_INTERVAL)
    {
        client.send("KeepAlive");
        last_update_sent = millis();
    }

    client.poll();
}