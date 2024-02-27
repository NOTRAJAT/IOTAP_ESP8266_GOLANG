#include <SPI.h> //include the SPI bus library
#include <MFRC522.h>
#include <Wire.h>
#include "PCF8574.h"
#include <LiquidCrystal_I2C.h>
#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266HTTPClient.h>
// constexpr uint8_t RST_PIN = D3;     // Configurable, see typical pin layout above
// constexpr uint8_t SS_PIN = D4;
MFRC522 mfrc522(D4, D3); // SS_PIN, RST_PIN
MFRC522::MIFARE_Key key;
PCF8574 pcf8574(0x38);
LiquidCrystal_I2C lcd(0x27, 16, 2);
int button_status[3];
byte readbackblock[18];
char rollno[10];
char sub1[] = "IOTAP";
char sub2[] = "Blockchain";
char *sel = 0;
int httpCode;
char Upload[64];
WiFiClient client;
HTTPClient http;
void req()
{
    if (WiFi.status() == WL_CONNECTED)
    { // Check WiFi connection status

        http.begin(client, "http://192.168.156.251:80/espurl"); // Specify request destination
        http.addHeader("Content-Type", "application/json");     // Specify content-type header
        sprintf(Upload, "{\"Roll_no\":\"%s\",\"Subject\":\"%s\"}", &rollno, sel);
        httpCode = http.POST(Upload); // Send the request
        //    String payload = http.getString();                  //Get the response payload

        Serial.println(httpCode); // Print HTTP return code
        //    Serial.println(payload);    //Print request response payload

        http.end(); // Close connection
    }
    else
    {

        Serial.println("Error in WiFi connection");
    }

    delay(1000); // Send a request every 30 seconds
}

void goodbuzz()
{
    digitalWrite(D8, HIGH);
    delay(200);
    digitalWrite(D8, LOW);
}

void badbuzz()
{
    digitalWrite(D8, HIGH);
    delay(100);
    digitalWrite(D8, LOW);
    delay(10);
    digitalWrite(D8, HIGH);
    delay(90);
    digitalWrite(D8, LOW);
}

void byteToCharWrite(int SIZE)
{
    //     Serial.print("read block: ");
    for (int j = 0; j < SIZE; j++)
    {
        rollno[j] = readbackblock[j];
    }

    rollno[SIZE] = '\0';
}

void ReadSector(int sec)
{

    delay(100);
    readBlock(sec, readbackblock);
    byteToCharWrite(9);
}

void lcdDisplayAdd()
{
    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print("Roll No :");
    lcd.setCursor(0, 1);

    for (int i = 0; i < 9; i++)
    {

        lcd.write(rollno[i]);
    }

    delay(2500);
    lcd.clear();
}

void sector_no()
{

    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print("Read Sector NO");
    button_status[1] = 0;

    while (true)
    {
        delay(200);
        yield();
        button_status[0] = pcf8574.digitalRead(P0);
        button_status[2] = pcf8574.digitalRead(P2);
        if (button_status[0] == LOW)
        {
            button_status[1]++;
            button_status[1] = button_status[1] % 4;
            lcd.clear();
            lcd.setCursor(0, 0);
            lcd.print("Read Sector NO");
            lcd.setCursor(0, 1);
            lcd.print("Selected: ");
            lcd.write(button_status[1] + 48);
            goodbuzz();
        }
        else if (button_status[2] == LOW)
        {
            if (button_status[1] == 0)
            {
                button_status[1] = 1;
            }
            goodbuzz();
            break;
        }
    }
}

int readBlock(int blockNumber, byte arrayAddress[])
{
    //  Serial.println("card selected");

    int largestModulo4Number = blockNumber / 4 * 4;
    int trailerBlock = largestModulo4Number + 3; // determine trailer block for the sector

    // authentication of the desired block for access
    byte status = mfrc522.PCD_Authenticate(MFRC522::PICC_CMD_MF_AUTH_KEY_A, trailerBlock, &key, &(mfrc522.uid));

    if (status != MFRC522::STATUS_OK)
    {
        Serial.print("PCD_Authenticate() failed (read): ");
        Serial.println(mfrc522.GetStatusCodeName((MFRC522::StatusCode)status));
        return 3; // return "3" as error message
    }

    // reading a block
    byte buffersize = 18;                                                 // we need to define a variable with the read buffer size, since the MIFARE_Read method below needs a pointer to the variable that contains the size...
    status = mfrc522.MIFARE_Read(blockNumber, arrayAddress, &buffersize); //&buffersize is a pointer to the buffersize variable; MIFARE_Read requires a pointer instead of just a number
    if (status != MFRC522::STATUS_OK)
    {
        Serial.print("MIFARE_read() failed: ");
        Serial.println(mfrc522.GetStatusCodeName((MFRC522::StatusCode)status));
        return 4; // return "4" as error message
    }
    //  Serial.println("block was read");

    return 0;
}

void conn()
{
    WiFi.begin("NG", "1q2w3e4r");

    while (WiFi.status() != WL_CONNECTED)
    {
        delay(500);
        Serial.print('.');
    }
}
void setup()
{
    // put your setup code here, to run once:
    Serial.begin(9600);
    conn();
    Serial.println(WiFi.localIP());
    pcf8574.pinMode(P0, INPUT_PULLUP);
    pcf8574.pinMode(P1, INPUT_PULLUP);
    pcf8574.pinMode(P2, INPUT_PULLUP);
    pinMode(D8, OUTPUT);
    pcf8574.begin();
    lcd.init();      // Initialize the LCD
    lcd.backlight(); // Turn on the backlight
    lcd.clear();
    SPI.begin(); // Init SPI bus
    mfrc522.PCD_Init();
    delay(500);
    for (byte i = 0; i < 6; i++)
    {
        key.keyByte[i] = 0xFF; // keyByte is defined in the "MIFARE_Key" 'struct' definition in the .h file of the library
    }

    delay(150);

    goodbuzz();
    sector_no();

    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print("Scan Your Card");

    goodbuzz();
}

void lcdbanner()
{
    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print("   Attendance  ");
    lcd.setCursor(0, 1);
    lcd.print("     system      ");

    goodbuzz();
}
void loop()
{

    if (!mfrc522.PICC_IsNewCardPresent())
    {

        return;
    }
    if (!mfrc522.PICC_ReadCardSerial())
    {
        return;
    }

    if (WiFi.status() != WL_CONNECTED)
    {
        lcd.clear();
        lcd.setCursor(0, 0);
        lcd.print("Connection");
        lcd.setCursor(0, 1);
        lcd.print("Error");
        conn();
        return;
    }

    ReadSector(button_status[1] + 3);
    goodbuzz();
    Serial.println(rollno);

    lcdDisplayAdd();
    delay(1000);

    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print("Choose Subject:");
    while (true)
    {
        button_status[0] = pcf8574.digitalRead(P0);
        button_status[1] = pcf8574.digitalRead(P1);
        button_status[2] = pcf8574.digitalRead(P2);

        if (button_status[0] == LOW)
        {

            lcd.clear();
            lcd.setCursor(0, 0);
            lcd.print("Choose Subject:");
            lcd.setCursor(0, 1);
            lcd.print(sub1);
            sel = sub1;
            goodbuzz();
        }
        else if (button_status[1] == LOW)
        {
            lcd.clear();
            lcd.setCursor(0, 0);
            lcd.print("Choose Subject:");
            lcd.setCursor(0, 1);
            lcd.print(sub2);
            sel = sub2;
            goodbuzz();
        }
        else if (button_status[2] == LOW)
        {
            goodbuzz();
            break;
        }
        delay(100);
    }

    req();

    if (httpCode == 200)
    {
        lcd.clear();
        lcd.setCursor(0, 0);
        lcd.print("Success Full");
        goodbuzz();
    }
    else
    {
        lcd.clear();
        lcd.setCursor(0, 0);
        lcd.print("Unsuccess Full");
        badbuzz();
    }
    delay(1000);
    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print("Scan Your Card");

    goodbuzz();

    httpCode = 0;
    mfrc522.PICC_HaltA();
    mfrc522.PCD_StopCrypto1();
    // put your main code here, to run repeatedly:
}