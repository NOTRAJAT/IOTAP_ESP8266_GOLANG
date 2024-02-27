#include <SPI.h>     //include the SPI bus library
#include <MFRC522.h> //include the RFID reader library

constexpr uint8_t RST_PIN = D3; // Configurable, see typical pin layout above
constexpr uint8_t SS_PIN = D4;

MFRC522 mfrc522(SS_PIN, RST_PIN); // instatiate a MFRC522 reader object.
MFRC522::MIFARE_Key key;          // create a MIFARE_Key struct named 'key', which will hold the card information

// this is the block number we will write into and then read.

byte readbackblock[18];
byte address[] = "2021ci29f";
int number = 5;
byte block1[16];

int writeBlock(int blockNumber, byte arrayAddress[])
{
    // this makes sure that we only write into data blocks. Every 4th block is a trailer block for the access/security info.
    int largestModulo4Number = blockNumber / 4 * 4;
    int trailerBlock = largestModulo4Number + 3; // determine trailer block for the sector
    if (blockNumber > 2 && (blockNumber + 1) % 4 == 0)
    {
        Serial.print(blockNumber);
        Serial.println(" is a trailer block:");
        return 2;
    }
    Serial.println(blockNumber);
    Serial.println(" is a data block:");

    // authentication of the desired block for access
    byte status = mfrc522.PCD_Authenticate(MFRC522::PICC_CMD_MF_AUTH_KEY_A, trailerBlock, &key, &(mfrc522.uid));
    if (status != MFRC522::STATUS_OK)
    {
        Serial.print("PCD_Authenticate() failed: ");
        Serial.println(mfrc522.GetStatusCodeName((MFRC522::StatusCode)status));
        return 3; // return "3" as error message
    }

    // writing the block
    status = mfrc522.MIFARE_Write(blockNumber, arrayAddress, 16);
    // status = mfrc522.MIFARE_Write(9, value1Block, 16);
    if (status != MFRC522::STATUS_OK)
    {
        Serial.print("MIFARE_Write() failed: ");
        Serial.println(mfrc522.GetStatusCodeName((MFRC522::StatusCode)status));
        return 4; // return "4" as error message
    }
    Serial.println("block was written");
    return 0;
}
void setup()
{
    Serial.begin(9600); // Initialize serial communications with the PC
    SPI.begin();        // Init SPI bus
    mfrc522.PCD_Init(); // Init MFRC522 card (in case you wonder what PCD means: proximity coupling device)
    delay(500);
    Serial.println("Scan a MIFARE Classic card");

    // Prepare the security key for the read and write functions.
    for (byte i = 0; i < 6; i++)
    {
        key.keyByte[i] = 0xFF; // keyByte is defined in the "MIFARE_Key" 'struct' definition in the .h file of the library
    }

    fill(9, block1, 0);

    for (int j = 0; j < 9; j++)
    {
        Serial.write(block1[j]);
    }
}

int readBlock(int blockNumber, byte arrayAddress[])
{
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
    Serial.println("block was read");
    for (int j = 0; j < 16; j++)
    {
        Serial.write(readbackblock[j]);
    }
    Serial.println("");
    return 0;
}

void fill(int end_, byte *BUFFER, int start_)
{
    for (int i = 0; i < end_ - start_; i++)
    {
        BUFFER[i] = address[i + start_];
    }
}
void writeSector(int row)
{

    writeBlock(row, block1);
    delay(250);
    readBlock(row, readbackblock);

    return;
}

// void readsector(int sector){
//
//     Serial.print("read block: ");
//    for (int j=0 ; j<16 ; j++)
//    {
//      Serial.write (readbackblock[j]);
//    }
//
// }

void loop()
{
    // Look for new cards
    if (!mfrc522.PICC_IsNewCardPresent())
    {
        return;
    }

    // Select one of the cards
    if (!mfrc522.PICC_ReadCardSerial())
    {
        return;
    }
    Serial.println("card selected");

    writeSector(number);
    delay(100);
    mfrc522.PICC_DumpToSerial(&(mfrc522.uid));

    mfrc522.PICC_HaltA();
    mfrc522.PCD_StopCrypto1();
}

// Write specific block

// Read specific block