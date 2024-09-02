# rfidtime
####  chip type Alien H3 9654 
# v1
# PC - Decoder/Encoder - 1 Antenna.
# Almacenar los valores por cada EPC - RSSI- Time [Database or something else]

# flow 
# 0x21 get reader information
# 0x51  stop fast inventory
# 0x7f load/modify reader profile

//
# start inventory
# Data: 0c009a010300010000000032ef  (select command) (send like 3 times)
# Data: 09000104fe00803280be
# Data: 05007fcd9303
# Data: 09000104fe00803280be
# Data: 0700010101001e4b


# Data: 09000104fe0180325ce4
# anatomy Data: 09 00 01 04(Qvalue) fe(Session) 01(MaskMem) 8032(MaskAddres) 5c(MaskLen) e4



#########################################
#########################################
TODO List
 -  Define clear Inventory initialization Pro and cons.
 -  Find a way to send a command to stop realtime inventory on demand like graceful shutdown
 -  Test Realtime inventory with more labels.
 -  Find a way to perform Answer Mode inventory/ equipment no anwser always to the same command.
 -  add file log.
 -  Define Group structure for like RunnerID , Group/Category.
   - runner range [1-10000] 2 byte
   - runner category [1-16] 1 byte
        - Juvenil/Libre/Master/ etc 
        - Subcategoria(Edades)
        - Sexo [1bit]
        - Distance [8-max]
   - experimental/reserved  2 bytes

Concurrency :
Notes :
Go Routines:
 - Just like any other function, you can pass it parameters to initialize its state.However
    any values returned by the function are ignored (bodner)
 - 

-  Register Runners [example: name, cellphone, birth date]
-  Register EventInfo [example: Event Name, Event Logo , distances [], start time , contact]
-  Assign Runner to Event.

- TODO LIST
-  Define format for log [example: EventID, ]
-  send telegram/whatsapp message under the following events
    - 1 day before start.
    - x minutes before start.
    - cross the starting line.
    - cross any intermediate line.
    - cross finish.
    - send picture/video after final line.
- User/Runner  
  - User Profile information category [example: category Expert, Novice, Number of Event as Runner, Position etc.]
  - Add medal for being a participant/ add medal for position.
  - Add Runner Certificate for each event (Download).
- RFID system Runner Information collection
  - RFID log information format definition
  - RFID log information format implementation.
  - Store Runner info Time when RSSI has the max value at the start line.
  - Store Runner info Time when RSSI has the max value at the final Line.
  - Store Runner info Time when RSSI has the max value at any Detection additional line.

- Runner Registration on event 
  - Payment 
  - Registration 
  - Assign Runner Event ID.

- Entregables :
    - Report of first 5 places for each category .

    
3 NE -
    - Encoder/Decoder RFID Systems 
    - WebServer to Store User/Event system
    - 