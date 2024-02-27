This Go Project encompasses an HTTP API server, engineered with the Gorilla Mux router, facilitating
attendance tracking and student creation. The server elegantly manages both the storage and retrieval of
attendance records, crafting student records, and thoughtfully delineating routes for disparate API endpoints.
Furthermore, the program intricately handles cookie management and logging for all incoming requests.
Further combines an RFID module (MFRC522) with an ESP8266 module to create an RFID-based
attendance tracking system. When a student scans their card, the ESP8266 reads their unique ID from the
card, then allows the user to select a subject using buttons. It then sends the information to a web server
over Wi-Fi. The system provides visual feedback through an LCD screen and buzzes to indicate success.